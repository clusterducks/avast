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
    "flag"
    "net/http"
    "text/template"

    "github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var indexTpl = template.Must(template.ParseFiles("index.html"))

func webHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    indexTpl.Execute(w, r.Host)
}

func main() {
    flag.Parse()
    registerClient()
    registerConsul()
    go wsHub.run()

    r := mux.NewRouter()
    r.HandleFunc("/", webHandler)
    r.HandleFunc("/ws", wsHandler)
    r.HandleFunc("/containers", containersHandler)
    r.HandleFunc("/container/{name}/inspect", containerHandler)
    r.HandleFunc("/images", imagesHandler)
    r.HandleFunc("/history/{id}", historyHandler)
    r.HandleFunc("/info", infoHandler)
    r.HandleFunc("/swarm/datacenters", swarmDatacentersHandler)
    r.HandleFunc("/swarm/nodes", swarmNodesHandler)
    r.HandleFunc("/swarm/health/{name}", swarmHealthHandler)
    // @TODO: wrap all routes in a method similar to this:
    // - https://github.com/hashicorp/consul/blob/ae7e96afea84f55c694dfa29877976adfb7ebabf/command/agent/http.go#L292
    http.Handle("/", r)

    panic(http.ListenAndServe(*addr, nil))
}
