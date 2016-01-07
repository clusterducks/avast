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
    "io/ioutil"
    "bufio"
    "net/http"
    "encoding/json"

    "github.com/gorilla/websocket"
    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
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

func echoEvents(conn *websocket.Conn) {
    fmt.Println(" --> Listening for events...")

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
            panic(err)
        }

        fmt.Println(" --> Received docker event...")
        fmt.Print(evt)

        err = conn.WriteJSON(evt)
        if err != nil {
            fmt.Println(err)
        }
    }
}

type ImageNode struct {
  Id        string
  ParentId  string
  Children  []*ImageNode  `json:"children,omitempty"`
  *types.Image
}

func (node *ImageNode) Add(parent string, nodes []*ImageNode) {
    for _, n := range nodes {
        if n.ParentId == parent {
            node.Children = append(node.Children, n)
            n.Add(n.Id, nodes)
        }
    }
}

func main() {
    var err error
    var containers []types.Container
    var images []types.Image

    cli, err = client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
    if err != nil {
        panic(err)
    }

    fmt.Println(" --> Listing docker containers...")

    copts := types.ContainerListOptions{All: true}
    containers, err = cli.ContainerList(copts)
    if err != nil {
        panic(err)
    }

    for _, c := range containers {
        fmt.Printf("ID: %v\n\tNames: %v\n\tImage: %v\n\tImageID: %v\n", c.ID, c.Names, c.Image, c.ImageID)
    }

    containerJson, _ := json.Marshal(&containers)
    fmt.Println(string(containerJson))

    ///////////////////////////////////////////////////////////////////////////

    fmt.Println(" --> Listing docker images...")

    iopts := types.ImageListOptions{All: true}
    images, err = cli.ImageList(iopts)
    if err != nil {
        panic(err)
    }

    for _, i := range images {
        fmt.Printf("ID: %v\n\tParentID: %v\n\n", i.ID, i.ParentID)
    }

    //inodes := make([]*ImageNode, len(images))
    //for i, img := range images {
    //    fmt.Printf("ID: %v\n\tParentID: %v\n", img.ID, img.ParentID)
    //    inodes[i] = &ImageNode{
    //      string(img.ID),
    //      string(img.ParentID),
    //      nil,
    //      &img,
    //    }
    //}

    //root := &ImageNode{"", "", nil, nil}
    //root.Add("", inodes)

    //imageJson, _ := json.Marshal(&root)
    //fmt.Println(string(imageJson))

    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/", rootHandler)

    panic(http.ListenAndServe(":8080", nil))
}
