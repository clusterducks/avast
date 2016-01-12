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
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hashicorp/consul/api"
)

var consul  *api.Client
var agent   *api.Agent
var catalog *api.Catalog
var health  *api.Health

func registerConsul() {
    config := api.DefaultConfig()
    consul, _ = api.NewClient(config)
    agent, catalog, health = consul.Agent(),
        consul.Catalog(),
        consul.Health()
}

func swarmDatacentersHandler(w http.ResponseWriter, r *http.Request) {
    dc, _ := catalog.Datacenters()

    data, err := json.Marshal(&dc)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

type NodesMeta struct {
    Nodes   []*api.Node     `json:"nodes"`
    Meta    *api.QueryMeta  `json:"meta"`
}

func swarmNodesHandler(w http.ResponseWriter, r *http.Request) {
    // @TODO: query by dc
    nodes, meta, _ := catalog.Nodes(nil)
    //members, _ := agent.Members(false)

    data, err := json.Marshal(NodesMeta{nodes, meta})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

type HealthCheckMeta struct {
    HealthCheck []*api.HealthCheck  `json:"healthCheck"`
    Meta        *api.QueryMeta      `json:"meta"`
}

func swarmHealthHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    // @TODO: query by dc
    check, meta, _ := health.Node(vars["name"], nil)

    data, err := json.Marshal(HealthCheckMeta{check, meta})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
