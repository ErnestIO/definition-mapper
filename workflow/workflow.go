/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package workflow

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Workflow : Is a representation of the workflow to be passed to
// fsm, that represents the endpoints to be called by the fsm
type Workflow struct {
	Arcs []Arc `json:"arcs"`
}

// Arc : Represents a transition on a fsm, that means an arrow from point A
// to point B, and the event that will fire this transition
type Arc struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Event string `json:"event"`
}

// LoadArcs : Load the arcs contained on a file to current workflow
func LoadArcs(file string) ([]Arc, error) {
	var arcs = []Arc{}

	path, _ := filepath.Abs(file)
	f, err := os.Open(path)
	if err != nil {
		return arcs, err
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&arcs)

	return arcs, err
}

// New: Creates a new workflow
func New(path string) *Workflow {
	w := Workflow{}
	w.Arcs = make([]Arc, 0)
	initialArcs, _ := LoadArcs(path)

	w.Add(initialArcs)
	return &w
}

// Add : Adds a group of arcs to a current workflow
func (w *Workflow) Add(a []Arc) {
	// Check for duplicate subject to prevent workflow loops
	for _, arc := range a {
		if w.ContainsSubject(arc.To) {
			return
		}
	}

	// Update the new arc with the last from value
	if len(w.Arcs) > 0 {
		a[0].From = w.Arcs[len(w.Arcs)-1].To
	}

	// Append the arcs to the workflow
	for _, arc := range a {
		w.Arcs = append(w.Arcs, arc)
	}
}

// Finish : adds the service.create.done arcs to current workflow
func (w *Workflow) Finish(path string) error {
	if w.Valid() {
		finalArcs, _ := LoadArcs(path)
		w.Add(finalArcs)
		return nil
	}
	return errors.New("Could not finish workflow!")
}

// Valid : Checks if all arcs on the workflow are valid
func (w *Workflow) Valid() bool {
	for i, arc := range w.Arcs {
		if arc.From == "" || arc.To == "" || arc.Event == "" {
			return false
		}

		// Check the previous Arc's To Value matches the next Arc's From value
		if i != 0 {
			if arc.From != w.Arcs[i-1].To {
				return false
			}
		}
	}
	return true
}

// ContainsEvent : Check if the workflow contains a specific arc
func (w *Workflow) ContainsEvent(arc string) bool {
	for _, a := range w.Arcs {
		if a.Event == arc {
			return true
		}
	}
	return false
}

// ContainsSubject : Check if the workflow contains a specific status on the
// From or To fields
func (w *Workflow) ContainsSubject(subject string) bool {
	for _, a := range w.Arcs {
		if a.To == subject || a.From == subject {
			return true
		}
	}
	return false
}
