// Copyright 2016 Brett Fowle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hashicorp/consul/api"
)

type Node struct {
    Address   string    `json:"address"`
    Node      string    `json:"node"`
}

func consulNodesHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    options := &api.QueryOptions{Datacenter: vars["dc"]}

    nodes, _, err := consul.Catalog.Nodes(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    var cnodes []*Node
    for _, n := range nodes {
        cnodes = append(cnodes, &Node{n.Address, n.Node})
    }

    return cnodes, nil
}

type ConsulNode struct {
    *Node
    Services  []*api.AgentService  `json:"services"`
    Checks    []*api.HealthCheck   `json:"checks"`
}

func consulNodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    options := &api.QueryOptions{Datacenter: vars["dc"]}

    node, _, err := consul.Catalog.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    health, _, err := consul.Health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    services := make([]*api.AgentService, 0, len(node.Services))
    for  _, s := range node.Services {
        services = append(services, s)
    }

    return &ConsulNode{
        &Node{node.Node.Address, node.Node.Node},
        services,
        health,
    }, nil
}
