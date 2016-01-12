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
    //"bufio"
    "encoding/json"
    "fmt"
    "io"

    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/events"
)

func (c *connection) echoEvents() {
    fmt.Printf(" --> Connected to client (%v)...\n", c.ws.RemoteAddr())

    options := types.EventsOptions{}
    r, err := cli.Events(options)
    if err != nil {
        panic(err)
    }
    defer r.Close()

    d := json.NewDecoder(r)
    messages := make(chan events.Message)

    go func() {
        for {
            var event events.Message
            if err := d.Decode(&event); err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Println(err)
            }
            messages <- event
        }
    }()

    for {
        select {
        case msg := <-messages:
            fmt.Println(" --> Received docker event...")
            fmt.Println(msg)

            data, err := json.Marshal(msg)
            if err != nil {
                fmt.Println(err)
            }
            wsHub.broadcast <- data
        }
    }
}
