/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats"
)

// Index is the root route for this miroservice
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "These aren't the droids you're loo king for.\n")
}

// CreateService calls the FSM through NATS
//
// * It receives an InputPayload
// * Validates all input data and environment definition and makes sure it
//   implements the current spec.
// * Builds an FSMMessage with this data
// * Sends the message to the FSM trhough NATS
func CreateService(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := r.FormValue("service")

	// Get the payload from services-data-microservice
	byt := []byte(body)
	var payload input.Payload
	if err := json.Unmarshal(byt, &payload); err != nil {
		message := fmt.Sprintf("Failed to parse payload. Error: \n %s", err)
		parseErr := errors.New(message)
		http.Error(w, parseErr.Error(), 400)
		return
	}
	prev, _ := getPreviousService(payload.PrevID)

	m, err := BuildFSMMessage(payload, prev)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	msg, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Publish("service.create", msg)
}

// PatchService will patch an errored service
func PatchService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	msg, _ := Request("service.get.mapping", []byte(`{"id":"`+id+`"}`))
	b := strings.Replace(string(msg), "\"service.create\"", "\"service.patch\"", -1)

	Publish("service.patch", []byte(b))
}

// GetService will patch an errored service
func GetService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	msg, err := Request("service.get.mapping", []byte(`{"id":"`+id+`"}`))
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(msg))
}

// DeleteService will remove every part of a service
func DeleteService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload, err := getPreviousService(ps.ByName("id"))
	if err != nil {
		log.Println("Failed to parse payload")
		message := fmt.Sprintf("Failed to parse payload. Error: \n %s", err)
		parseErr := errors.New(message)
		http.Error(w, parseErr.Error(), 400)
		return
	}

	if payload == nil {
		log.Println("Service not found")
		http.Error(w, "Service not found", 400)
		return
	}

	m, err := BuildDeleteMessage(*payload)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	msg, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Publish("service.delete", msg)
}

func getPreviousService(id string) (*output.FSMMessage, error) {
	body, err := Request("service.get.mapping", []byte(`{"id":"`+id+`"}`))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var payload output.FSMMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	if payload.ID == "" {
		return nil, nil
	}

	return &payload, nil
}

var dbURL string
var natsURI string

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/service", CreateService)
	router.PATCH("/service/:id", PatchService)
	router.DELETE("/service/:id", DeleteService)
	router.GET("/service/:id", GetService)

	natsURI = os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = nats.DefaultURL
	}

	dbURL = getDBURL()
	go Subscribe()

	// TODO Depreciate HTTP in favor of NATS
	log.Fatal(http.ListenAndServe(":21000", router))
}
