/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
	"github.com/nats-io/nats"
)

var natsClient *nats.Conn
var err error

// MapCreateService builds a valid service from the input and replies with it
func MapCreateService(msg *nats.Msg) {
	var payload input.Payload

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		natsClient.Publish(msg.Reply, []byte(`{"error":"Failed to parse payload."}`))
		return
	}
	prev, _ := getPreviousService(payload.PrevID)

	m, err := BuildFSMMessage(payload, prev)
	if err != nil {
		natsClient.Publish(msg.Reply, []byte(err.Error()))
		return
	}

	body, err := json.Marshal(m)
	if err != nil {
		natsClient.Publish(msg.Reply, []byte(err.Error()))
		return
	}

	natsClient.Publish(msg.Reply, body)
}

// MapDeleteService : builds and responds with a service based on input
func MapDeleteService(msg *nats.Msg) {
	var input input.Payload

	if err := json.Unmarshal(msg.Data, &input); err != nil {
		natsClient.Publish(msg.Reply, []byte(`{"error":"Failed to parse payload."}`))
		return
	}
	payload, err := getPreviousService(input.PrevID)

	if err != nil {
		natsClient.Publish(msg.Reply, []byte(`{"error":"Failed to parse payload."}`))
		return
	}

	if payload == nil {
		natsClient.Publish(msg.Reply, []byte(`{"error":"Service not found."}`))
		return
	}

	m, err := BuildDeleteMessage(*payload)
	if err != nil {
		log.Println(err)
		natsClient.Publish(msg.Reply, []byte(`{"error":"Internal error."}`))
		return
	}

	body, err := json.Marshal(m)
	if err != nil {
		natsClient.Publish(msg.Reply, []byte(`{"error":"Internal error."}`))
		return
	}

	natsClient.Publish(msg.Reply, body)
}

// Get the previous service based on an ID
func getPreviousService(id string) (*output.FSMMessage, error) {
	msg, err := natsClient.Request("service.get.mapping", []byte(`{"id":"`+id+`"}`), time.Second)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var payload output.FSMMessage
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return nil, err
	}

	if payload.ID == "" {
		return nil, nil
	}

	return &payload, nil
}

func main() {
	natsURI := os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = nats.DefaultURL
	}
	natsClient, err = nats.Connect(natsURI)
	if err != nil {
		log.Panic(err)
	}
}

func Subscribe(natsURI *nats.Conn) {
	natsClient.Subscribe("description.map_create", MapCreateService)
	natsClient.Subscribe("description.map_delete", MapDeleteService)
}
