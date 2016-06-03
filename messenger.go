/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/nats-io/nats"
)

// Subscribe : It manages the subscription to service creation results on nats
func Subscribe() {
	nc, err := nats.Connect(natsURI)
	if err != nil {
		log.Panic(err)
	}

	nc.Subscribe("service.create.done", func(m *nats.Msg) {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal(err)
			return
		}
		var Service struct {
			ID       string `json:"id"`
			Endpoint string `json:"endpoint"`
		}
		err = json.Unmarshal(m.Data, &Service)
		if err != nil {
			log.Fatal(err)
		}
		ServiceID := Service.ID
		time.Sleep(2 * time.Second)
		// TODO : Mysql should never, ever be used on this microservice, this
		// needs to be delegated to service-data-microservice
		_, err = db.Query(`UPDATE services set service_result=$1, service_status='done', service_endpoint=$2, service_error=NULL where service_id=$3`, m.Data, Service.Endpoint, ServiceID)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	})

	nc.Subscribe("service.create.error", func(m *nats.Msg) {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			fmt.Printf(err.Error())
		}
		var Service struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(m.Data, &Service)
		if err != nil {
			log.Fatal(err)
		}
		ServiceID := Service.ID
		// TODO : Mysql should never, ever be used on this microservice, this
		// needs to be delegated to service-data-microservice
		time.Sleep(2 * time.Second)
		_, err = db.Query(`UPDATE services set service_error=$1, service_status='errored' where service_id=$2`, m.Data, ServiceID)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	})

	nc.Subscribe("service.delete.done", func(m *nats.Msg) {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			fmt.Printf(err.Error())
		}
		var Service struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		err = json.Unmarshal(m.Data, &Service)
		if err != nil {
			log.Fatal(err)
		}
		// TODO : Mysql should never, ever be used on this microservice, this
		// needs to be delegated to service-data-microservice
		time.Sleep(2 * time.Second)
		_, err = db.Query(`DELETE FROM services WHERE service_name=$1`, Service.Name)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	})

	nc.Subscribe("service.delete.error", func(m *nats.Msg) {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			fmt.Printf(err.Error())
		}
		var Service struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(m.Data, &Service)
		if err != nil {
			log.Fatal(err)
		}
		ServiceID := Service.ID
		// TODO : Mysql should never, ever be used on this microservice, this
		// needs to be delegated to service-data-microservice
		_, err = db.Query(`UPDATE services set service_error=$1, service_status='errored' where service_id=$2`, m.Data, ServiceID)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
	})
}

// Publish : Manages nats message publication
func Publish(subject string, msg []byte) {
	nc, err := nats.Connect(natsURI)
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
	nc.Publish(subject, msg)
}

// Request : does a nats request based on given config
func Request(subject string, msg []byte) ([]byte, error) {
	nc, err := nats.Connect(natsURI)
	if err != nil {
		return nil, err
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := nc.Request(subject, msg, 5*time.Second)
	return res.Data, err
}
