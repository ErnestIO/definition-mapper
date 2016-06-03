/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats"
)

// Loads the configuration file config.json
func getDBURL() string {
	var cfg map[string]interface{}

	// Default db Name
	dbName := "services"
	if os.Getenv("DB_NAME") != "" {
		dbName = os.Getenv("DB_NAME")
	}

	// Check for ENV variable
	if os.Getenv("DB_URI") != "" {
		return fmt.Sprintf("%s/%s?sslmode=disable", os.Getenv("DB_URI"), dbName)
	}

	// Get config from conf-store
	nc, err := nats.Connect(natsURI)
	if err != nil {
		log.Panic(err)
	}

	resp, err := nc.Request("config.get.postgres", nil, time.Second)
	if err != nil {
		log.Println("could not load config")
		log.Panic(err)
	}

	err = json.Unmarshal(resp.Data, &cfg)
	if err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf("%s/%s?sslmode=disable", cfg["url"], dbName)
}
