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
    "github.com/hashicorp/consul/api"
)

var consul *Consul

type Consul struct {
    Client  *api.Client
    Agent   *api.Agent
    Catalog *api.Catalog
    Health  *api.Health
}

func registerConsul() {
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
    }

    consul = &Consul{
        client,
        client.Agent(),
        client.Catalog(),
        client.Health(),
    }
}
