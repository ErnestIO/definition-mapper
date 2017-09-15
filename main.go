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
	aws "github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	azure "github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
	"gopkg.in/r3labs/graph.v2"
	yaml "gopkg.in/yaml.v2"
)

var n *nats.Conn

type build struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Previous   string `json:"previous_id"`
	Datacenter struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"datacenter"`
	Definition struct {
		Name    string `json:"name"`
		Project string `json:"project"`
	} `json:"build"`
}

func getInputDetails(body []byte) (string, string, string, string, string) {
	var s build
	if err := json.Unmarshal(body, &s); err != nil {
		log.Panic(err)
	}

	return s.ID, s.Name, s.Datacenter.Type, s.Previous, s.Definition.Name
}

func getGraphDetails(body []byte) (string, string, string) {
	var gg map[string]interface{}
	err := json.Unmarshal(body, &gg)
	if err != nil {
		log.Println("could not process graph")
		return "", "", ""
	}

	gx := graph.New()
	err = gx.Load(gg)
	if err != nil {
		log.Println("could not load graph")
		return "", "", ""
	}

	credentials := gx.GetComponents().ByType("credentials")

	return gx.ID, gx.Name, credentials[0].GetProvider()
}

func getDefinition(id string) (map[string]interface{}, error) {
	var d map[string]interface{}

	resp, err := n.Request("build.get.definition", []byte(`{"id":"`+id+`"}`), time.Second)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Data, &d)

	return d, err
}

func getDefinitionDetails(d map[string]interface{}) (string, string) {
	var name string
	var project string

	name, _ = d["name"].(string)
	project, _ = d["project"].(string)

	return name, project
}

func getImportFilters(m map[string]interface{}, name string, provider string) []string {
	var filters []string

	switch provider {
	case "azure", "azure-fake":
		d, ok := m["build"].(map[string]interface{})
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

func getCredentials(gd map[string]interface{}) (map[string]interface{}, error) {
	project, ok := gd["datacenter"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not find project credentials")
	}

	credentials, ok := project["credentials"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not find project credentials")
	}

	credentials["name"] = project["name"]
	credentials["type"] = project["type"]

	return credentials, nil
}

func definitionToGraph(m libmapper.Mapper, body []byte) (*graph.Graph, error) {
	var gd map[string]interface{}
	err := json.Unmarshal(body, &gd)
	if err != nil {
		return nil, err
	}

	definition, ok := gd["build"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("could not convert definition")
	}

	credentials, err := getCredentials(gd)
	if err != nil {
		return nil, err
	}

	sid, ok := gd["id"].(string)
	if ok != true {
		return nil, errors.New("could not find build id")
	}

	name, ok := definition["name"].(string)
	if ok != true {
		return nil, errors.New("could not find build name")
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

// SubscribeCreateBuild : definition.map.creation subscriber
// For a given definition, it will generate the valid build
// and necessary workflow to create the environment on the
// provider
func SubscribeCreateBuild(body []byte) ([]byte, error) {
	id, _, t, p, name := getInputDetails(body)

	m := providers.NewMapper(t)
	if m == nil {
		return body, fmt.Errorf("Unconfigured provider type : '%s'", t)
	}

	g, err := definitionToGraph(m, body)
	if err != nil {
		return body, err
	}

	// If there is a previous build
	if p != "" {
		oMsg, rerr := n.Request("build.get.mapping", []byte(`{"id":"`+p+`"}`), time.Second)
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

	for i := range g.Changes {
		g.Changes[i].SetDefaultVariables()
	}

	g.ID = id
	g.Name = name

	return g.ToJSON()
}

// SubscribeImportBuild : definition.map.import subscriber
// For a given filters it will generate a workflow to fully
// import a provider build.
func SubscribeImportBuild(body []byte) ([]byte, error) {
	var err error

	var gd map[string]interface{}
	err = json.Unmarshal(body, &gd)
	if err != nil {
		return nil, err
	}

	credentials, err := getCredentials(gd)
	if err != nil {
		return nil, err
	}

	id, name, t, _, n := getInputDetails(body)

	m := providers.NewMapper(t)

	filters := getImportFilters(gd, n, t)

	g := m.CreateImportGraph(filters)
	if g, err = g.Diff(graph.New()); err != nil {
		return body, err
	}

	g.ID = id
	g.Name = name
	err = g.AddComponent(m.ProviderCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return g.ToJSON()
}

// SubscribeImportComplete : build.create.done subscriber
// Converts a completed import graph to an inpurt definition
func SubscribeImportComplete(body []byte) error {
	var build struct {
		ID         string       `json:"id"`
		Definition string       `json:"definition"`
		Mapping    *graph.Graph `json:"mapping"`
	}

	id, _, provider := getGraphDetails(body)

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
		def.Name, def.Project = getDefinitionDetails(pd)
	case "azure":
		def := d.(*azure.Definition)
		def.Name, def.Project = getDefinitionDetails(pd)
	}

	data, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	build.ID = id
	build.Definition = string(data)
	build.Mapping = g

	sdata, err := json.Marshal(build)
	if err != nil {
		return err
	}

	_, err = n.Request("build.set.mapping", sdata, time.Second)
	if err != nil {
		return err
	}

	_, err = n.Request("build.set.definition", sdata, time.Second)
	if err != nil {
		return err
	}

	return err
}

// SubscribeDeleteBuild : definition.map.deletion subscriber
// For a given existing build will generate a valid internal
// build with a workflow to delete all its components
func SubscribeDeleteBuild(body []byte) ([]byte, error) {
	var b build
	if err := json.Unmarshal(body, &b); err != nil {
		log.Panic(err)
	}

	t := b.Datacenter.Type
	p := b.Previous

	m := providers.NewMapper(t)

	oMsg, err := n.Request("build.get.mapping", []byte(`{"id":"`+p+`"}`), time.Second)
	if err != nil {
		return body, err
	}

	original, err := mappingToGraph(m, oMsg.Data)
	if err != nil {
		return body, err
	}

	var gd map[string]interface{}
	if err := json.Unmarshal(body, &gd); err != nil {
		return body, err
	}

	credentials, err := getCredentials(gd)
	if err != nil {
		return body, err
	}

	creds := m.ProviderCredentials(credentials)
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
	g.Name = original.Name

	return json.Marshal(g)
}

// SubscribeMapBuild : definition.map.build subscriber
// For a given full build will generate the relative
// definition
func SubscribeMapBuild(body []byte) ([]byte, error) {
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
		if body, err := SubscribeCreateBuild(m.Data); err == nil {
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
		if body, err := SubscribeImportBuild(m.Data); err == nil {
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
		if body, err := SubscribeDeleteBuild(m.Data); err == nil {
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

	if _, err := n.Subscribe("definition.map.build", func(m *nats.Msg) {
		if body, err := SubscribeMapBuild(m.Data); err == nil {
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

	if _, err := n.Subscribe("build.import.done", func(m *nats.Msg) {
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
