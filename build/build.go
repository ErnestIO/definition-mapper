/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package build

import "github.com/r3labs/graph"

// Build ...
type Build struct {
	ID         string       `json:"id"`
	Definition string       `json:"definition"`
	Mapping    *graph.Graph `json:"mapping"`
}
