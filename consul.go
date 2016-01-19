package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/mux"
    consulapi "github.com/hashicorp/consul/api"
    "github.com/hashicorp/consul/watch"
)

var consulRegistry *ConsulRegistry

type ClientNode struct {
    Name        string  `json:"name"`
    Address     string  `json:"address"`
}

type ConsulNode struct {
    *ClientNode
    Services  []*consulapi.AgentService  `json:"services"`
    Checks    []*consulapi.HealthCheck   `json:"checks"`
}

type ConsulRegistry struct {
    addr        string
    client      *consulapi.Client
    agent       *consulapi.Agent
    catalog     *consulapi.Catalog
    health      *consulapi.Health
    services    map[string]*consulapi.ServiceEntry
    nodes       []*ConsulNode
    sync.RWMutex
}

type Watcher struct {
    addr        string
    wp          *watch.WatchPlan
    watchers    map[string]*watch.WatchPlan
}

type WatchEvent struct {
    From        string      `json:"from"`
    Type        string      `json:"type"`
    Data        interface{} `json:"data"`
    Timestamp   time.Time   `json:"timestamp"`
}

//---------------------
//- REGISTRY

func newConsulRegistry() {
    config := consulapi.DefaultConfig()
    c, err := consulapi.NewClient(config)
    if err != nil {
        fmt.Println(err)
        return
    }

    consulRegistry = &ConsulRegistry{
        addr:       config.Address,
        client:     c,
        agent:      c.Agent(),
        catalog:    c.Catalog(),
        health:     c.Health(),
        services:   make(map[string]*consulapi.ServiceEntry),
    }
}

//---------------------
//- DATACENTERS

func (cr *ConsulRegistry) DatacentersHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    datacenters, err := cr.catalog.Datacenters()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    return datacenters, nil
}

//---------------------
//- NODES

func (cr *ConsulRegistry) NodesHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}
    nodes, _, err := cr.catalog.Nodes(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    var cnodes []*ClientNode
    for _, n := range nodes {
        cnodes = append(cnodes, &ClientNode{
            Name: n.Node,
            Address: n.Address,
        })
    }

    return cnodes, nil
}

func (cr *ConsulRegistry) NodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}
    node, _, err := cr.catalog.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    health, _, err := cr.health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    services := make([]*consulapi.AgentService, 0, len(node.Services))
    for  _, s := range node.Services {
        services = append(services, s)
    }

    return &ConsulNode{
        ClientNode: &ClientNode{
            Name: node.Node.Node,
            Address: node.Node.Address,
        },
        Services: services,
        Checks: health,
    }, nil
}

//---------------------
//- HEALTH

func (cr *ConsulRegistry) HealthHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}
    check, _, err := cr.health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    return check, nil
}

//---------------------
//- DISCOVERY

func (cr *ConsulRegistry) registerWatcher(watchType string) error {
    w, err := newWatcher(cr.addr, watchType)
    if err != nil {
        fmt.Println(err)
        return err
    }
    go w.Run()

    return nil
}

func (w *Watcher) registerServiceWatcher(service string) error {
    wp, err := watch.Parse(map[string]interface{}{
        "type": "service",
        "service": service,
    })
    if err != nil {
        return err
    }

    wp.Handler = func(idx uint64, data interface{}) {
        switch d := data.(type) {
        case []*consulapi.ServiceEntry:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData("service", &i)

                consulRegistry.Lock()
                consulRegistry.services[i.Service.Service] = i
                consulRegistry.Unlock()
            }
        }
    }

    go wp.Run(w.addr)
    w.watchers[service] = wp

    return nil
}

func newWatcher(addr string, watchType string) (*Watcher, error) {
    wp, err := watch.Parse(map[string]interface{}{"type": watchType})
    if err != nil {
        return nil, err
    }

    w := &Watcher{
        addr,
        wp,
        make(map[string]*watch.WatchPlan),
    }

    wp.Handler = func(idx uint64, data interface{}) {
        // @TODO: type switch seems to convert back to interface{}
        // if applying multiple types on the case (e.g. case []*A, []*B
        // it would be nice to combine these case statements for similar
        // types; try using reflect.TypeOf instead, perhaps
        switch d := data.(type) {
        // nodes
        case []*consulapi.Node:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData(watchType, &i)
            }
        // checks
        case []*consulapi.HealthCheck:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData(watchType, &i)
            }
        // services
        case map[string][]string:
            for i, _ := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                if _, ok := w.watchers[i]; ok {
                    continue
                }
                w.registerServiceWatcher(i)
            }

            consulRegistry.RLock()
            rs := consulRegistry.services
            consulRegistry.RUnlock()

            // remove unknown services from registry
            for s, _ := range rs {
                if _, ok := d[s]; !ok {
                    consulRegistry.Lock()
                    delete(consulRegistry.services, s)
                    consulRegistry.Unlock()
                }
            }

            // remove unknown services from watchers
            for i, svc := range w.watchers {
                if _, ok := d[i]; !ok {
                    svc.Stop()
                    delete(w.watchers, i)
                }
            }
        }
    }

    return w, nil
}

func broadcastData(watchType string, data interface{}) {
    evt := &WatchEvent {
        "consul",
        watchType,
        data,
        time.Now(),
    }

    msg, err := json.Marshal(evt)
    if err != nil {
        fmt.Println(err)
    }
    wsHub.broadcast <- msg
}

func (w *Watcher) Run() {
    w.wp.Run(w.addr)
}

func (w *Watcher) Stop() {
    if w.wp == nil {
        return
    }
    w.wp.Stop()
}

func (cr *ConsulRegistry) EchoDiscovery() {
    cr.registerWatcher("nodes")
    cr.registerWatcher("checks")
    cr.registerWatcher("services")
}
