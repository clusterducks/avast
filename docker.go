package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/events"
    "github.com/dustin/go-humanize"
    "github.com/gorilla/mux"
)

var dockerClient *DockerClient

type DockerClient struct {
    *client.Client
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
    Children    []*ImageNode        `json:"children"`
}

type DockerEvent struct {
    From        string          `json:"from"`
    Type        string          `json:"type"`
    Data        *events.Message `json:"data"`
    Timestamp   time.Time       `json:"timestamp"`
}

//---------------------
//- CLIENT

func newDockerClient() {
    cli, err := client.NewEnvClient()
    if err != nil {
        fmt.Println(err)
        return
    }
    dockerClient = &DockerClient{cli}
}

//---------------------
//- CONTAINERS

func (dc *DockerClient) ContainersHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    options := types.ContainerListOptions{All: true}
    containers, err := dc.ContainerList(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return containers, nil
}

func (dc *DockerClient) ContainerHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    container, err := dc.ContainerInspect(vars["name"])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return container, nil
}

//---------------------
//- IMAGES

func (node *ImageNode) add(parent string, nodes []*ImageNode) {
    for _, n := range nodes {
        if n.ParentID == parent {
            node.Children = append(node.Children, n)
            n.add(n.ID, nodes)
        }
    }
}

func (dc *DockerClient) ImagesHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    options := types.ImageListOptions{All: true}
    images, err := dc.ImageList(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    nodes := make([]*ImageNode, len(images))
    for i, img := range images {
        nodes[i] = &ImageNode{
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

    root := &ImageNode{}
    root.add("", nodes)

    return root, nil
}

func (dc *DockerClient) HistoryHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    history, err := dc.ImageHistory(vars["id"])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return history, nil
}

//---------------------
//- INFO

func (dc *DockerClient) InfoHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    info, err := dc.Info()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return info, nil
}

//---------------------
//- EVENTS

func (dc *DockerClient) EchoEvents() {
    options := types.EventsOptions{}
    r, err := dc.Events(options)
    if err != nil {
        fmt.Println(err)
    }
    defer r.Close()

    d := json.NewDecoder(r)
    messages := make(chan *DockerEvent)

    go func() {
        for {
            var data events.Message
            if err := d.Decode(&data); err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Println(err)
            }
            messages <- &DockerEvent{
                "docker",
                "event",
                &data,
                time.Now(),
            }
        }
    }()

    for {
        select {
        case msg := <-messages:
            data, err := json.Marshal(msg)
            if err != nil {
                fmt.Println(err)
            }
            wsHub.broadcast <- data
        }
    }
}
