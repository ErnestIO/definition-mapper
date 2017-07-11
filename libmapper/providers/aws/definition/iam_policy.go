/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// IamPolicy ...
type IamPolicy struct {
	Name              string                 `json:"name" yaml:"name"`
	PolicyDocument    map[string]interface{} `json:"policy_document" yaml:"policy_document"`
	PolicyDocumentRaw string                 `json:"policy_document_raw,omitempty" yaml:"policy_document_raw,omitempty"`
	Description       string                 `json:"description" yaml:"description"`
	Path              string                 `json:"path" yaml:"path"`
}
