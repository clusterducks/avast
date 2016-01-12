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
    "time"

    "github.com/docker/engine-api/types"
    "github.com/dustin/go-humanize"
)

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
    root.Add("", nodes)

    data, err := json.Marshal(&root)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
