// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import "github.com/TeamFairmont/gabs"

// FilterPayload recieves a gabs.Container and checks it against []keys, an array of strings.
// Any top-level children of the payload that aren't in []keys are removed. Returns the modified gabs.Container
func FilterPayload(p *gabs.Container, keys []string) (*gabs.Container, error) {
	//copy payload, so origional is not altered
	var payload = gabs.New()
	var err error
	if p != nil {
		payload, err = gabs.ParseJSON([]byte(p.String()))
	} else {
		payload, err = gabs.ParseJSON([]byte(`{}`))
		return nil, err
	}

	//uses gabs children function and checks for an error.
	children, err := payload.ChildrenMap()
	if err != nil {
		return payload, err
	}
	//Loops through every child in children
	//test every key on every child to see if they match in utils.StringInSlice
	//if no match is found, the child is deleted
	if len(keys) >= 1 {
		for child := range children {
			if !StringInSlice(child, keys) {
				payload.Delete(child)
			}
		}
	}
	return payload, err
}
