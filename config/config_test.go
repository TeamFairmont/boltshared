// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(tst *testing.T) {
	cfg, err := DefaultConfig()
	assert.Equal(tst, "{}", string(cfg.WorkerConfig), "Should be empty json object")
	assert.Equal(tst, 0, len(cfg.APICalls), "Should be empty map")
	assert.Equal(tst, 0, len(cfg.CommandMetas), "Should be empty map")
	assert.Nil(tst, err, "No error")

	defcache := "{\"type\":\"\",\"host\":\"localhost:6379\",\"pass\":\"\",\"timeoutMs\":2000}"
	cac, _ := json.Marshal(cfg.Cache)
	assert.Equal(tst, defcache, string(cac), "Cache structs should match")

	defsecurity := "{\"verifyTimeout\":30,\"groups\":[],\"handlerAccess\":null,\"corsDomains\":[],\"corsAutoAddLocal\":true}"
	sec, _ := json.Marshal(cfg.Security)
	assert.Equal(tst, defsecurity, string(sec), "Security structs should match")
}

func TestCustomizeConfig(tst *testing.T) {
	cfg, err := DefaultConfig()
	assert.Nil(tst, err, "No error")

	CustomizeConfig(cfg, `{
            "engine": { "bind": ":8888" },
            "cache": { "type": "memcached" },
			"security" : {
				"groups": [
					{
						"name": "test01",
						"hmackey": "01234567890~!@#$%^&*-_=+ABCabc"
					},
					{
						"name": "test02",
						"hmackey": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcd"
					}
				]
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
						"someGlobalOption2": "boolean"
					},
					"commands": [
						{
							"name": "product/checkDuplicates",
							"resultTimeoutMs": 6000,
							"returnAfter": false,
							"configParams": {
								"caseSensitive": false
							}
						}
					]
				},
				"v2/Junk": {}
			},
			"commandMeta": {
				"product/checkDuplicates": {
					"requiredParams": {
						"sku":"string",
						"name": "multilang-strings"
					}
				},
				"product/saveToDb": {
					"requiredParams": {
						"sku":"string",
						"price":"double",
						"salePrice":"double",
						"name": "multilang-strings",
						"shortDescription": "multilang-strings"
					}
				},
				"updateSearchIndex": {
					"requiredParams": {
						"indexFields": "array",
						"uniqueIdField": "string"
					}
				}
			}
    }`)

	assert.Equal(tst, ":8888", cfg.Engine.Bind, "Customized engine bind (:8888) should override the default (:443)")
	assert.Equal(tst, "memcached", cfg.Cache.Type, "Customized cache type (memcached) should should override the default (redis)")
	assert.Equal(tst, "localhost:6379", cfg.Cache.Host, "No customized cache host name should override the default (localhost)")
	assert.NotEqual(tst, "", cfg.Cache.Host, "The cache host name should contain a non-empty value")
	assert.NotNil(tst, cfg.Cache.Host, "The cache host name should be non-nil")

	assert.Equal(tst, 2, len(cfg.Security.Groups), "Security group count length should be 2")
	assert.Equal(tst, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcd", cfg.Security.Groups[1].Hmackey, "The 2nd security group hmackey should match")

	assert.Equal(tst, "product/checkDuplicates", cfg.APICalls["v1/addProduct"].Commands[0].Name, "Api commands should match")
	assert.Equal(tst, 2, len(cfg.APICalls), "Api command count should match")

	assert.NotNil(tst, cfg.CommandMetas["product/checkDuplicates"], "Command meta loaded successfully")
	assert.Equal(tst, 3, len(cfg.CommandMetas), "Command meta count matches expected")
}

func TestBuildConfig(tst *testing.T) {
	// Build a config and test its schema using the versions in TeamFairmont/boltengine/etc/bolt/
	cfg, err := BuildConfig("../../boltengine/etc/bolt/", "../../boltengine/etc/bolt/config.json")
	// If err is nil, the schema in TeamFairmont/boltengine/etc/bolt/config-schema.json has validated the json in TeamFairmont/boltengine/etc/bolt/config.json
	assert.Nil(tst, err, "err should be nil")

	// Make sure BuildConfig returned a working config:
	assert.Equal(tst, ":8888", cfg.Engine.Bind, "config.json's engine bind (:8888) should override the default (:443)")

	// Sending an invalid path should return an error
	cfg, err = BuildConfig("/bad/path/", "/bad/path/config.json")
	assert.Equal(tst, "open etc/bolt/config.json: no such file or directory", err.Error(), "Bad path to config should not resolve")
	// An invalid path when the api is running will load the config json and schema from etc/bolt/.
	// This path is invalid from the present working directory when running the unit tests.
}
