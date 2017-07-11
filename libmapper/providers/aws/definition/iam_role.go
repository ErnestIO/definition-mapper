/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// IamRole ...
type IamRole struct {
	Name                    string                 `json:"name" yaml:"name"`
	AssumePolicyDocument    map[string]interface{} `json:"assume_policy_document" yaml:"assume_policy_document"`
	AssumePolicyDocumentRaw string                 `json:"assume_policy_document_raw,omitempty" yaml:"assume_policy_document_raw,omitempty"`
	Policies                []string               `json:"policies" yaml:"policies"`
	Description             string                 `json:"description" yaml:"description"`
	Path                    string                 `json:"path" yaml:"path"`
}
