/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/ernestio/definition-mapper/libmapper/providers"
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
	"gopkg.in/r3labs/graph.v2"
)

var n *nats.Conn

func getInputDetails(body []byte) (string, string, string, string) {
	var service struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Previous   string `json:"previous_id"`
		Datacenter struct {
			Type string `json:"type"`
		} `json:"datacenter"`
	}

	if err := json.Unmarshal(body, &service); err != nil {
		log.Panic(err)
	}

	return service.ID, service.Name, service.Datacenter.Type, service.Previous
}

func definitionToGraph(m libmapper.Mapper, body []byte) (*graph.Graph, error) {
	var gd map[string]interface{}
	err := json.Unmarshal(body, &gd)
	if err != nil {
		return nil, err
	}

	definition, ok := gd["service"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not convert definition")
	}

	credentials, ok := gd["datacenter"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not find datacenter credentials")
	}

	sid, ok := gd["id"].(string)
	if ok != true {
		return nil, errors.New("could not find service id")
	}

	d, err := m.LoadDefinition(definition)
	if err != nil {
		return nil, err
	}

	g, err := m.ConvertDefinition(d)
	if err != nil {
		return nil, err
	}

	// set graph ID and credentials
	g.ID = sid
	err = g.AddComponent(m.ProviderCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return g, nil
}

func mappingToGraph(m libmapper.Mapper, body []byte) (*graph.Graph, error) {
	var gm map[string]interface{}
	err := json.Unmarshal(body, &gm)
	if err != nil {
		return nil, err
	}

	return m.LoadGraph(gm)
}

// SubscribeCreateService : definition.map.creation subscriber
// For a given definition, it will generate the valid service
// and necessary workflow to create the environment on the
// provider
func SubscribeCreateService(body []byte) ([]byte, error) {
	id, _, t, p := getInputDetails(body)

	m := providers.NewMapper(t)
	if m == nil {
		return body, fmt.Errorf("Unconfigured provider type : '%s'", t)
	}

	g, err := definitionToGraph(m, body)
	if err != nil {
		return body, err
	}

	// If there is a previous service
	if p != "" {
		oMsg, rerr := n.Request("service.get.mapping", []byte(`{"id":"`+p+`"}`), time.Second)
		if rerr != nil {
			return body, rerr
		}
		og, merr := mappingToGraph(m, oMsg.Data)
		if merr != nil {
			return body, merr
		}

		g, err = g.Diff(og)
		if err != nil {
			return body, err
		}
	} else {
		g, err = g.Diff(graph.New())
		if err != nil {
			return body, err
		}
	}

	g.ID = id

	return g.ToJSON()
}

// SubscribeImportService : definition.map.import subscriber
// For a given filters it will generate a workflow to fully
// import a provider service.
func SubscribeImportService(body []byte) ([]byte, error) {
	var err error
	var filters []string

	id, n, t, _ := getInputDetails(body)
	// TODO Allow multi-filters for azure development
	filters = append(filters, n)

	m := providers.NewMapper(t)

	g := m.CreateImportGraph(filters)
	if g, err = g.Diff(graph.New()); err != nil {
		return body, err
	}

	g.ID = id

	return g.ToJSON()
}

// SubscribeDeleteService : definition.map.deletion subscriber
// For a given existing service will generate a valid internal
// service with a workflow to delete all its components
func SubscribeDeleteService(body []byte) ([]byte, error) {
	_, _, t, p := getInputDetails(body)
	m := providers.NewMapper(t)

	oMsg, rerr := n.Request("service.get.mapping", []byte(`{"id":"`+p+`"}`), time.Second)
	if rerr != nil {
		return body, rerr
	}

	original, merr := mappingToGraph(m, oMsg.Data)
	if merr != nil {
		return body, merr
	}

	empty := graph.New()

	g, err := empty.Diff(original)
	if err != nil {
		return body, err
	}

	g.ID = p

	return json.Marshal(g)
}

// SubscribeMapService : definition.map.service subscriber
// For a given full service will generate the relative
// definition
func SubscribeMapService(body []byte) ([]byte, error) {
	var gd map[string]interface{}

	if err := json.Unmarshal(body, &gd); err != nil {
		return body, err
	}

	_, _, t, _ := getInputDetails(body)
	m := providers.NewMapper(t)

	original, err := m.LoadGraph(gd)
	if err != nil {
		return body, err
	}
	definition, err := m.ConvertGraph(original)
	if err != nil {
		return body, err
	}

	return json.Marshal(definition)
}

// ManageDefinitions : Manages all subscriptions
func ManageDefinitions() {
	if _, err := n.Subscribe("definition.map.creation", func(m *nats.Msg) {
		if body, err := SubscribeCreateService(m.Data); err == nil {
			if err = n.Publish(m.Reply, body); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}); err != nil {
		log.Panic(err)
	}

	if _, err := n.Subscribe("definition.map.import", func(m *nats.Msg) {
		if body, err := SubscribeImportService(m.Data); err == nil {
			if err = n.Publish(m.Reply, body); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}); err != nil {
		log.Panic(err)
	}

	if _, err := n.Subscribe("definition.map.deletion", func(m *nats.Msg) {
		if body, err := SubscribeDeleteService(m.Data); err == nil {
			if err = n.Publish(m.Reply, body); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}); err != nil {
		log.Panic(err)
	}

	if _, err := n.Subscribe("definition.map.service", func(m *nats.Msg) {
		if body, err := SubscribeMapService(m.Data); err == nil {
			if err = n.Publish(m.Reply, body); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}); err != nil {
		log.Panic(err)
	}
}

func setup() {
	n = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func main() {
	setup()
	ManageDefinitions()
	runtime.Goexit()
}
