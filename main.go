package main

import (
  "fmt"
  "bufio"

  "github.com/docker/engine-api/client"
  "github.com/docker/engine-api/types"
)

func main() {
    defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
    cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
    if err != nil {
        panic(err)
    }

    fmt.Println(" --> Listing running containers...")

    options := types.ContainerListOptions{All: true}
    containers, err := cli.ContainerList(options)
    if err != nil {
        panic(err)
    }

    for _, c := range containers {
        fmt.Println(c.ID)
    }

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
      fmt.Print(evt)
    }
}
