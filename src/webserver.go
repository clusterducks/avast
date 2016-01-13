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
    "flag"
    "net/http"
    "strings"
    "text/template"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

var router *mux.Router
var addr = flag.String("addr", ":8080", "http service address")
var indexTpl = template.Must(template.ParseFiles("index.html"))

func setHeaders(w http.ResponseWriter, headers map[string]string) {
    for k, v := range headers {
        w.Header().Set(http.CanonicalHeaderKey(k), v)
    }
}

func startWebserver() {
    flag.Parse()

    router = mux.NewRouter()
    router.HandleFunc("/",                          wrap(indexHandler))
    router.HandleFunc("/ws",                        wrap(wsHandler))
    router.HandleFunc("/docker/containers",         wrap(dockerContainersHandler))
    router.HandleFunc("/docker/container/{name}",   wrap(dockerContainerHandler))
    router.HandleFunc("/docker/images",             wrap(dockerImagesHandler))
    router.HandleFunc("/docker/history/{id}",       wrap(dockerHistoryHandler))
    router.HandleFunc("/docker/info",               wrap(dockerInfoHandler))
    router.HandleFunc("/consul/datacenters",        wrap(consulDatacentersHandler))
    router.HandleFunc("/consul/nodes",              wrap(consulNodesHandler))
    router.HandleFunc("/consul/nodes/{dc}",         wrap(consulNodesHandler))
    router.HandleFunc("/consul/node/{name}",        wrap(consulNodeHandler))
    router.HandleFunc("/consul/health/{name}",      wrap(consulHealthHandler))
    router.HandleFunc("/consul/health/{name}/{dc}", wrap(consulHealthHandler))

    http.Handle("/", router)
    panic(http.ListenAndServe(*addr, handlers.CompressHandler(router)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    indexTpl.Execute(w, r.Host)

    return nil, nil
}

func wrap(handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
    f := func(w http.ResponseWriter, r *http.Request) {
        setHeaders(w, map[string]string{
            "Content-Type": "text/html; charset=utf8",
            "Access-Control-Allow-Origin": "*",
        })

        obj, err := handler(w, r)

    HAS_ERR:
        if err != nil {
            code := http.StatusInternalServerError
            errMsg := err.Error()
            if strings.Contains(errMsg, "Permission denied") {
                code = http.StatusForbidden
            }

            w.WriteHeader(code)
            w.Write([]byte(err.Error()))
            return
        }

        if obj != nil {
            buf, err := json.Marshal(obj)
            if err != nil {
                goto HAS_ERR
            }

            w.Header().Set("Content-Type", "application/json")
            w.Write(buf)
        }
    }

    return f
}
