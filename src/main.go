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
    "time"
    //"strings"
    "io/ioutil"
    "bufio"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/events"
    "github.com/dustin/go-humanize"
)

var cli *client.Client

var defaultHeaders = map[string]string{"User-Agent": "engine-api-cli-1.0"}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    content, err := ioutil.ReadFile("index.html")
    if err != nil {
        fmt.Println("Could not open file.", err)
    }
    fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Origin") != "http://" + r.Host {
        http.Error(w, "Origin not allowed", 403)
        return
    }

    conn, err := upgrader.Upgrade(w, r, w.Header())
    if err != nil {
        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
        fmt.Println(err)
    }

    go echoEvents(conn)
}

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

type ImageNode struct {
    ID          string              `json:"id"`
    ParentID    string              `json:"parentId"`
    RepoTags    []string            `json:"repoTags"`
    RepoDigests []string            `json:"repoDigests"`
    Created     string              `json:"created"`
    Size        string              `json:"size"`
    VirtualSize string              `json:"virtualSize"`
    Labels      map[string]string   `json:"labels"`
    //types.Image
    Children    []*ImageNode    `json:"children"`
}

func (node *ImageNode) Add(parent string, nodes []*ImageNode) {
    for _, n := range nodes {
        if n.ParentID == parent {
            node.Children = append(node.Children, n)
            n.Add(n.ID, nodes)
        }
    }
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
    options := types.ImageListOptions{All: true}
    images, err := cli.ImageList(options)
    if err != nil {
        panic(err)
    }

    nodes := make([]*ImageNode, len(images))
    for i, img := range images {
        nodes[i] = &ImageNode{
            //img
            img.ID,
            img.ParentID,
            img.RepoTags,
            img.RepoDigests,
            humanize.Time(time.Unix(img.Created, 0)),
            humanize.Bytes(uint64(img.Size)),
            humanize.Bytes(uint64(img.VirtualSize)),
            img.Labels,
            nil,
        }
    }

    root := &ImageNode{
        //types.Image{},
        "", "", nil, nil,
        "", "", "", nil, nil,
    }
    root.Add("", nodes)

    data, err := json.Marshal(&root)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    history, err := cli.ImageHistory(vars["id"])
    if err != nil {
        panic(err)
    }

    data, err := json.Marshal(&history)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}


func echoEvents(conn *websocket.Conn) {
    fmt.Printf(" --> Connected to client (%v)...\n", conn.RemoteAddr())

    evtOpts := types.EventsOptions{}
    evtReader, err := cli.Events(evtOpts)
    if err != nil {
        panic(err)
    }
    defer evtReader.Close()

    rd := bufio.NewReader(evtReader)
    for {
        evt, err := rd.ReadString('\n')
        if err != nil {
            fmt.Println(err)
        }

        var message events.Message
        err = json.Unmarshal([]byte(evt), &message)
        if err != nil {
            fmt.Println(err)
        }

        fmt.Println(" --> Received docker event...")
        fmt.Println(message)

        err = conn.WriteJSON(message)
        if err != nil {
            fmt.Println(err)
        }
    }
}

func main() {
    var err error

    cli, err = client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
    if err != nil {
        panic(err)
    }

    router := mux.NewRouter()
    router.HandleFunc("/", rootHandler)
    router.HandleFunc("/ws", wsHandler)
    router.HandleFunc("/containers/list", containersHandler)
    router.HandleFunc("/container/{name}/inspect", containerHandler)
    router.HandleFunc("/images/list", imagesHandler)
    router.HandleFunc("/history/{id}", historyHandler)
    http.Handle("/", router)

    panic(http.ListenAndServe(":8080", nil))
}
