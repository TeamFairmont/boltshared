// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"testing"

	"github.com/TeamFairmont/gabs"
	"github.com/stretchr/testify/assert"
)

//calls FilterPayload and passes in a sample JSON code.
func TestFilterPayload(t *testing.T) {
	testJSON, err := gabs.ParseJSON([]byte(`{
		"someThing":{
			"v1/stuff": 111
	},
		"!@#$%&*()_+":{
			"blcokThing": 111
},
"apploTree":{
	"v1/stuff": 111
},
"a1pploTree":{
	"v1/stuff": 111
},
"1apploTree":{
	"v1/stuff": 111
},
"crappleTree":{
	"v1/stuff": 111
},
		"apiCalls": {
			"v1/addProduct": {
				"resultTimeoutMs": 100,
				"cache": {
					"enabled": false,
					"expirationTimeSec": 600,
					"allowOverride": false
				},
				"requiredParams": {
					"someGlobalOption1": "string",
					"someGlobalOption2": "bool"
				}
			}
		}
	}`))
	//if ther is an error fail
	assert.Equal(t, err, nil)
	keys := []string{"apiCalls",
		"zapiCalls",
		"otherThing",
		"apploTree",
		"!@#$%&*()_+",
	}
	runFilter(t, testJSON, keys)
	runFilter(t, nil, keys)
	keys = nil
	runFilter(t, testJSON, keys)
	testJSON = nil
	runFilter(t, testJSON, keys)
}

func runFilter(t *testing.T, testJSON *gabs.Container, keys []string) {
	//rund FilterPayload
	filtered, err := FilterPayload(testJSON, keys)
	assert.Nil(t, err, "FilterPayload returned an error")
	//Get the childredn for testing
	if testJSON != nil {
		children, err := filtered.ChildrenMap()
		assert.Nil(t, err, "gabs ChildrenMap returned an error")
		if keys != nil {
			//check to see that all children returned by OutputRequest are in the filter keys list
			for child := range children {
				assert.True(t, StringInSlice(child, keys))
			} //check that the old payload does not match the new one
			assert.False(t, testJSON.String() == filtered.String(), "should NOT be equal")
		} else { //if no filter keys are sent, the old payload should match the new one
			assert.Equal(t, testJSON.String(), filtered.String(), "should be true")
		}
	} else {
		//if gabs container passed to FilterPayload was nil
		assert.Nil(t, testJSON)
		assert.Equal(t, testJSON, filtered)
	}
}
