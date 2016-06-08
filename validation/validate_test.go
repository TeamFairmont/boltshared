// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package validate

import (
	"log"
	"testing"

	"github.com/TeamFairmont/gabs"
	"github.com/stretchr/testify/assert"
)

func TestMatches(tst *testing.T) {

	var err error

	// Create an required parameters to test against
	requiredParams := map[string]string{
		"someGlobalOption1": "string",
		"someGlobalOption2": "bool",
		"someGlobalOption3": "float64",
		"someGlobalOption4": "int64",
	}

	// Create a valid payload to test
	payload, err := gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This is a string",
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": 123,
			"someGlobalOption5": "This is an extra, non-required param"
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result := CheckPayloadReqParams(requiredParams, payload)
	assert.Nil(tst, result, "When payload contains all required params + extras, err should = nil")

	// Test edge cases for int64 values, should throw an error since it's a float64:
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This is a string",
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": 1.0000000002
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Parameter:someGlobalOption4, Expected:int64, Received:float64", result.Error(), "Edge case float64 shouldn't convert to int64.  Err should be thrown.")

	// Test a very big int
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This is a string",
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": 999999999999999999
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Nil(tst, result, "bigint 999999999999999999 is valid, err should = nil")

	//Huge int64 value of 9223372036854775807 is too large.  It will throw an error
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This is a string",
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": 9223372036854775807
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Parameter:someGlobalOption4, Expected:int64, Received:float64", result.Error(), "Numeric value 9223372036854775807 is too large, err will be thrown.")

	// Create an incomplete payload to test
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This payload skips someGlobalOption3",
			"someGlobalOption2": false,
			"someGlobalOption4": 0
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Missing parameter: someGlobalOption3", result.Error(), "Payload missing required params should error.")

	// Create a payload with an invalid field
	//lowest int64 value = -9223372036854775808
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "12345678.90",
			"someGlobalOption2": "This should be bool, but it's string",
			"someGlobalOption3": 54.321,
			"someGlobalOption4": -123
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Parameter:someGlobalOption2, Expected:bool, Received:string", result.Error(), "Payload with invalid type in a required param should error.")

	// Ensure an error is thrown when an int64 is expected but a string is received:
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": "This is a string",
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": "int64 required, but I'm a string!"
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Parameter:someGlobalOption4, Expected:int64, Received:string", result.Error(), "Payload with string value when int64 is expected should error.")

	// Ensure an error is thrown when a string is expected but a float64 is received:
	payload, err = gabs.ParseJSON([]byte(`{
		"initial_input" : {
			"someGlobalOption1": 12345678.90,
			"someGlobalOption2": true,
			"someGlobalOption3": 123.45,
			"someGlobalOption4": 1234567890
		}
	}`))
	if err != nil {
		log.Fatal("Error parsing json:\n", err, "\n")
	}
	// Test it!
	result = CheckPayloadReqParams(requiredParams, payload)
	assert.Equal(tst, "Parameter:someGlobalOption1, Expected:string, Received:float64", result.Error(), "Payload with float64 value when string is expected should error.")
}

func TestCheckPayloadStructure(tst *testing.T) {
	p1, _ := gabs.ParseJSON([]byte(`{
			"initial_input":{},
			"return_value":{},
			"data":{},
			"trace":[],
			"debug":{},
			"nextCommand":"",
			"error":{},
			"config":{},
			"params":{}
		}`))
	err := CheckPayloadStructure(p1)
	assert.Nil(tst, err, "p1 should not have errors")

	p2, _ := gabs.ParseJSON([]byte(`{
				"return_value":{},
				"data":{},
				"trace":[],
				"debug":{},
				"nextCommand":"",
				"error":{},
				"config":{},
				"params":{}
			}`))
	err = CheckPayloadStructure(p2)
	assert.NotNil(tst, err, "Error should exist")
	if err != nil {
		assert.Equal(tst, err.Error(), "Payload missing initial_input", "p2 should be missing initial_input and thus get an error")
	}

	p3, _ := gabs.ParseJSON([]byte(`{
					"initial_input":{},
					"return_value":{},
					"data":{},
					"trace":[],
					"debug":{},
					"error":{},
					"config":{},
					"params":{}
				}`))
	err = CheckPayloadStructure(p3)
	assert.NotNil(tst, err, "Error should exist")
	if err != nil {
		assert.Equal(tst, err.Error(), "Payload missing nextCommand", "p3 should be missing nextCommand and thus get an error")
	}

}
