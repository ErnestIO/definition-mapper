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
	"strconv"
	"time"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/ernestio/definition-mapper/libmapper/providers"
	aws "github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	azure "github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
	"gopkg.in/r3labs/graph.v2"
	yaml "gopkg.in/yaml.v2"
)

var n *nats.Conn

type service struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Previous   string `json:"previous_id"`
	Datacenter struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"datacenter"`
	Definition struct {
		Name       string `json:"name"`
		Datacenter string `json:"datacenter"`
	} `json:"service"`
}

func getInputDetails(body []byte) (string, string, string, string, string) {
	var s service
	if err := json.Unmarshal(body, &s); err != nil {
		log.Panic(err)
	}

	return s.ID, s.Name, s.Datacenter.Type, s.Previous, s.Definition.Name
}

func getGraphDetails(body []byte) (string, string) {
	var gg map[string]interface{}
	err := json.Unmarshal(body, &gg)
	if err != nil {
		log.Println("could not process graph")
		return "", ""
	}

	gx := graph.New()
	err = gx.Load(gg)
	if err != nil {
		log.Println("could not load graph")
		return "", ""
	}

	credentials := gx.GetComponents().ByType("credentials")

	return gx.ID, credentials[0].GetProvider()
}

func getDefinition(id string) (map[string]interface{}, error) {
	var d map[string]interface{}

	resp, err := n.Request("service.get.definition", []byte(`{"id":"`+id+`"}`), time.Second)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Data, &d)

	return d, err
}

func getDefinitionDetails(d map[string]interface{}) (string, string) {
	var name string
	var datacenter string

	name, _ = d["name"].(string)
	datacenter, _ = d["datacenter"].(string)

	return name, datacenter
}

func getImportFilters(m map[string]interface{}, name string, provider string) []string {
	var filters []string

	switch provider {
	case "azure", "azure-fake":
		d, ok := m["service"].(map[string]interface{})
		if !ok {
			return filters
		}

		f, ok := d["import_filters"].([]interface{})
		if !ok {
			return filters
		}
		for _, filter := range f {
			fil := filter.(string)
			filters = append(filters, fil)
		}
	default:
		filters = append(filters, name)
	}

	return filters
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	cm := make(map[string]interface{})

	for k, v := range m {
		cm[k] = v
	}

	return cm
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

	name, ok := definition["name"].(string)
	if ok != true {
		return nil, errors.New("could not find service name")
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
	g.Name = name
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
	id, _, t, p, _ := getInputDetails(body)

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

		for _, c := range g.Components {
			oc := og.Component(c.GetID())
			if oc != nil {
				c.Update(oc)
			}
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

	var gd map[string]interface{}
	err = json.Unmarshal(body, &gd)
	if err != nil {
		return nil, err
	}

	credentials, ok := gd["datacenter"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not find datacenter credentials")
	}

	id, _, t, _, n := getInputDetails(body)

	m := providers.NewMapper(t)

	filters := getImportFilters(gd, n, t)

	g := m.CreateImportGraph(filters)
	if g, err = g.Diff(graph.New()); err != nil {
		return body, err
	}

	g.ID = id
	err = g.AddComponent(m.ProviderCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return g.ToJSON()
}

// SubscribeImportComplete : service.create.done subscriber
// Converts a completed import graph to an inpurt definition
func SubscribeImportComplete(body []byte) error {
	var service struct {
		ID         string `json:"id"`
		Definition string `json:"definition"`
		Mapping    string `json:"mapping"`
	}

	id, provider := getGraphDetails(body)

	pd, err := getDefinition(id)
	if err != nil {
		return err
	}

	var gg map[string]interface{}
	err = json.Unmarshal(body, &gg)
	if err != nil {
		return err
	}

	m := providers.NewMapper(provider)

	g, err := m.LoadGraph(gg)
	if err != nil {
		return err
	}

	d, err := m.ConvertGraph(g)
	if err != nil {
		return err
	}

	switch provider {
	case "aws":
		def := d.(*aws.Definition)
		def.Name, def.Datacenter = getDefinitionDetails(pd)
	case "azure":
		def := d.(*azure.Definition)
		def.Name, def.Datacenter = getDefinitionDetails(pd)
	}

	data, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	gdata, err := g.ToJSON()
	if err != nil {
		return err
	}

	service.ID = id
	service.Definition = string(data)
	service.Mapping = string(gdata)

	sdata, err := json.Marshal(service)
	if err != nil {
		return err
	}

	_, err = n.Request("service.set.mapping", sdata, time.Second)
	if err != nil {
		return err
	}

	_, err = n.Request("service.set.definition", sdata, time.Second)
	if err != nil {
		return err
	}

	return err
}

// SubscribeDeleteService : definition.map.deletion subscriber
// For a given existing service will generate a valid internal
// service with a workflow to delete all its components
func SubscribeDeleteService(body []byte) ([]byte, error) {
	var s service
	if err := json.Unmarshal(body, &s); err != nil {
		log.Panic(err)
	}

	t := s.Datacenter.Type
	dID := strconv.Itoa(s.Datacenter.ID)
	p := s.Previous

	m := providers.NewMapper(t)

	oMsg, rerr := n.Request("service.get.mapping", []byte(`{"id":"`+p+`"}`), time.Second)
	if rerr != nil {
		return body, rerr
	}

	original, merr := mappingToGraph(m, oMsg.Data)
	if merr != nil {
		return body, merr
	}

	oMsg, rerr = n.Request("datacenter.get", []byte(`{"id":`+dID+`}`), time.Second)
	if rerr != nil {
		return body, rerr
	}
	var datacenterDetails map[string]interface{}
	if err := json.Unmarshal(oMsg.Data, &datacenterDetails); err != nil {
		return body, err
	}

	creds := m.ProviderCredentials(datacenterDetails)
	original.UpdateComponent(creds)
	for i := range original.Components {
		original.Components[i].Rebuild(original)
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

	_, _, t, _, _ := getInputDetails(body)
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
			var errorMessage struct {
				Message string `json:"error"`
			}
			errorMessage.Message = err.Error()
			body, err := json.Marshal(errorMessage)
			if err != nil {
				log.Println(err.Error())
			}
			if err = n.Publish(m.Reply, body); err != nil {
				log.Println("Error trying to respond through nats : " + err.Error())
			}
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
			if err = n.Publish(m.Reply, []byte(`{"error":"`+err.Error()+`"}`)); err != nil {
				log.Println("Error trying to respond through nats : " + err.Error())
			}
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
			if err = n.Publish(m.Reply, []byte(`{"error":"`+err.Error()+`"}`)); err != nil {
				log.Println("Error trying to respond through nats : " + err.Error())
			}
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
			if err = n.Publish(m.Reply, []byte(`{"error":"`+err.Error()+`"}`)); err != nil {
				log.Println("Error trying to respond through nats : " + err.Error())
			}
		}
	}); err != nil {
		log.Panic(err)
	}

	if _, err := n.Subscribe("service.import.done", func(m *nats.Msg) {
		if err := SubscribeImportComplete(m.Data); err != nil {
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
