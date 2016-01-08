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

    "github.com/docker/engine-api/types"
    "github.com/gorilla/mux"
)

func containersHandler(w http.ResponseWriter, r *http.Request) {
    options := types.ContainerListOptions{All: true}
    containers, err := cli.ContainerList(options)
    if err != nil {
        panic(err)
    }

    data, err := json.Marshal(&containers)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func containerHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    container, err := cli.ContainerInspect(vars["name"])
    if err != nil {
        panic(err)
    }

    data, err := json.Marshal(&container)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
