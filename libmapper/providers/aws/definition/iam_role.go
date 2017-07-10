/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// IamRole ...
type IamRole struct {
	Name                    string                 `json:"name"`
	AssumePolicyDocument    map[string]interface{} `json:"assume_policy_document"`
	AssumePolicyDocumentRaw string                 `json:"assume_policy_document_raw"`
	Policies                []string               `json:"policies"`
	Description             string                 `json:"description"`
	Path                    string                 `json:"path"`
}
