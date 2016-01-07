package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "bufio"

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

func main() {
    var err error
    cli, err = client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
    if err != nil {
        panic(err)
    }

    /*
    fmt.Println(" --> Listing running containers...")

    options := types.ContainerListOptions{All: true}
    containers, err := cli.ContainerList(options)
    if err != nil {
        panic(err)
    }

    for _, c := range containers {
        fmt.Println(c.ID)
    }
    */

    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/", rootHandler)

    panic(http.ListenAndServe(":8080", nil))
}
