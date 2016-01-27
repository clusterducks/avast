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

type ConsulKV struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type ConsulService struct {
	ID      string   `json:"id"`
	Service string   `json:"service"`
	Tags    []string `json:"tags"`
	Port    int      `json:"port"`
	Address string   `json:"address"`
}

type ConsulHealthCheck struct {
	Node        string `json:"node"`
	CheckID     string `json:"checkId"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	Output      string `json:"output"`
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
}

type ConsulNode struct {
	Name     string               `json:"name"`
	Address  string               `json:"address"`
	Services []*ConsulService     `json:"services"`
	Checks   []*ConsulHealthCheck `json:"checks"`
}

type ConsulRegistry struct {
	addr     string
	client   *consulapi.Client
	agent    *consulapi.Agent
	catalog  *consulapi.Catalog
	health   *consulapi.Health
	services map[string]*consulapi.ServiceEntry
	nodes    []*ConsulNode
	sync.RWMutex
}

type Watcher struct {
	addr     string
	wp       *watch.WatchPlan
	watchers map[string]*watch.WatchPlan
}

type WatchEvent struct {
	From      string      `json:"from"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

//---------------------
//- REGISTRY

func newConsulRegistry() {
	// @TODO: put a watch on "barney/docker/swarm/leader" key
	// for current swarm leader
	config := consulapi.DefaultConfig()
	c, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	consulRegistry = &ConsulRegistry{
		addr:     config.Address,
		client:   c,
		agent:    c.Agent(),
		catalog:  c.Catalog(),
		health:   c.Health(),
		services: make(map[string]*consulapi.ServiceEntry),
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

	var cnodes []*ConsulNode
	for _, n := range nodes {
		node, err := cr.fetchNode(n.Node)
		if err != nil {
			fmt.Println(err)
		}
		cnodes = append(cnodes, node)
	}

	return cnodes, nil
}

func (cr *ConsulRegistry) fetchNode(name string) (*ConsulNode, error) {
	options := &consulapi.QueryOptions{}
	node, _, err := cr.catalog.Node(name, options)
	if err != nil {
		return nil, err
	}

	services := make([]*ConsulService, 0, len(node.Services))
	for _, s := range node.Services {
		services = append(services, &ConsulService{
			s.ID,
			s.Service,
			s.Tags,
			s.Port,
			s.Address,
		})
	}

	checks, _, err := cr.health.Node(name, options)
	health := make([]*ConsulHealthCheck, 0, len(checks))
	for _, c := range checks {
		health = append(health, &ConsulHealthCheck{
			c.Node,
			c.CheckID,
			c.Name,
			c.Status,
			c.Notes,
			c.Output,
			c.ServiceID,
			c.ServiceName,
		})
	}

	return &ConsulNode{
		node.Node.Node,
		node.Node.Address,
		[]*ConsulService(services),
		health,
	}, nil
}

func (cr *ConsulRegistry) NodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	node, err := cr.fetchNode(vars["name"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
		return nil, err
	}

	return node, nil
}

//---------------------
//- HEALTH

func (cr *ConsulRegistry) HealthHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	options := &consulapi.QueryOptions{}
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

func (cr *ConsulRegistry) registerWatcher(watchType string, opts map[string]string) error {
	w, err := newWatcher(cr.addr, watchType, opts)
	if err != nil {
		fmt.Println(err)
		return err
	}
	go w.Run()

	return nil
}

func (w *Watcher) registerServiceWatcher(service string) error {
	wp, err := watch.Parse(map[string]interface{}{
		"type":       "service",
		"datacenter": datacenter,
		"service":    service,
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

func newWatcher(addr string, watchType string, opts map[string]string) (*Watcher, error) {
	var options = map[string]interface{}{
		"type":       watchType,
		"datacenter": datacenter,
	}
	for k, v := range opts {
		options[k] = v
	}

	wp, err := watch.Parse(options)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		addr,
		wp,
		make(map[string]*watch.WatchPlan),
	}

	wp.Handler = func(idx uint64, data interface{}) {
		switch d := data.(type) {
		// key
		case *consulapi.KVPair:
			fmt.Printf("[ %v ]\t%v\n", time.Now(), d)
			broadcastData(watchType, &ConsulKV{
				d.Key,
				d.Value,
			})
		// nodes
		case []*consulapi.Node:
			nodes := make([]*ConsulNode, 0, len(d))
			for _, i := range d {
				fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
				nodes = append(nodes, &ConsulNode{
					Name:    i.Node,
					Address: i.Address,
				})
			}
			broadcastData(watchType, &nodes)
		// checks
		case []*consulapi.HealthCheck:
			for _, i := range d {
				fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
				broadcastData(watchType, &ConsulHealthCheck{
					i.Node,
					i.CheckID,
					i.Name,
					i.Status,
					i.Notes,
					i.Output,
					i.ServiceID,
					i.ServiceName,
				})
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
		default:
			fmt.Printf("Could not interpret data type: %v", &d)
		}
	}

	return w, nil
}

func broadcastData(watchType string, data interface{}) {
	evt := &WatchEvent{
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
	cr.registerWatcher("key", map[string]string{"key": "barney/docker/swarm/leader"})
	cr.registerWatcher("nodes", nil)
	cr.registerWatcher("checks", nil)
	cr.registerWatcher("services", nil)
}
