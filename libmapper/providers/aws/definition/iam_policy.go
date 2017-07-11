/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// IamPolicy ...
type IamPolicy struct {
	Name              string                 `json:"name"`
	PolicyDocument    map[string]interface{} `json:"policy_document"`
	PolicyDocumentRaw string                 `json:"policy_document_raw"`
	Description       string                 `json:"description"`
	Path              string                 `json:"path"`
}
