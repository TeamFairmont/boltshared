// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package security

import (
	"strconv"
	"testing"
	"time"

	config "github.com/TeamFairmont/boltshared/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateGroup(tst *testing.T) {
	// Config has its own unit tests, but we need it here to verify group structs match
	// and authentication performs as expected.
	cfg, err := config.DefaultConfig()
	assert.Nil(tst, err, "No error")

	// Attempt to authenticate when no groups have been added
	authenticated := AuthenticateGroup("test01", "01234567890~!@#$%^&*-_=+ABCabc", &cfg.Security.Groups)
	assert.False(tst, authenticated, "The user and hmackey should not yet exist in the customized config")

	// Add a group to be authenticated.
	cfg, err = config.CustomizeConfig(cfg, `{
		"security" : {
			"groups": [
				{
					"name": "test01",
					"hmackey": "01234567890~!@#$%^&*-_=+ABCabc"
				}
			]
		}
	}`)
	assert.Nil(tst, err, "No error")
	assert.Equal(tst, "01234567890~!@#$%^&*-_=+ABCabc", cfg.Security.Groups[0].Hmackey, "The security group hmackey should match")

	// Attempt to authenticate the new group successfully.
	authenticated = AuthenticateGroup("test01", "01234567890~!@#$%^&*-_=+ABCabc", &cfg.Security.Groups)
	assert.True(tst, authenticated, "The user and hmackey should exist in the customized config")

	// Attempt to authenticate bad username or password.
	authenticated = AuthenticateGroup("failuser", "01234567890~!@#$%^&*-_=+ABCabc", &cfg.Security.Groups)
	assert.False(tst, authenticated, "The user and hmackey should NOT exist in the customized config")
	authenticated = AuthenticateGroup("test01", "badpw567890~!@#$%^&*-_=+ABCabc", &cfg.Security.Groups)
	assert.False(tst, authenticated, "The user and hmackey should NOT exist in the customized config")
}

func TestEncodeDecodeHMAC(tst *testing.T) {
	// Create a custom config with an HMAC key for encoding
	cfg, err := config.DefaultConfig()
	cfg, err = config.CustomizeConfig(cfg, `{
		"security" : {
			"groups": [
				{
					"name": "test01",
					"hmackey": "01234567890~!@#$%^&*-_=+ABCabc",
					"aesEncryptionKey": "0123456789012key"
				}
			]
		}
	}`)

	// Encode a string and verify it matches the expected output
	// The encoded message expires after 1 second
	stringToEncode := `Bolt is a { "noun": "winner!!1"}`
	hmacSigned, err := EncodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", stringToEncode, strconv.FormatInt(time.Now().Unix(), 10))
	assert.Nil(tst, err, "Should not return an error when encoding")

	assert.NotEqual(tst, stringToEncode, hmacSigned, "The value returned from EncodeHMAC changes every time the function is called.  Make sure the original string doesn't match the returned value.")

	assert.NotEqual(tst, stringToEncode, string(hmacSigned), "The value returned from EncodeHMAC changes every time the function is called.  Make sure the original string doesn't match the returned value if encoded as a string.")

	// Decode the HMAC and verify it matches the original string
	hmacDecoded, err := DecodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", hmacSigned, 30)
	assert.Nil(tst, err, "Should not return an error when decoding")
	assert.Equal(tst, stringToEncode, hmacDecoded, "String prior to encoding should match the output from being encrypted then decrypted.")

	// Test HMAC timestamp expiration.  Payloads should only be decoded if their encoding timestamp was within +- 30 seconds.
	timeNow := time.Now()
	timePlus := timeNow.Add(31 * time.Second)
	timeMinus := timeNow.Add((-31 * time.Second))

	hmacPlusSigned, err := EncodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", stringToEncode, strconv.FormatInt(timePlus.Unix(), 10))
	hmacMinusSigned, err := EncodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", stringToEncode, strconv.FormatInt(timeMinus.Unix(), 10))

	hmacPlusDecoded, err := DecodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", hmacPlusSigned, 30)
	assert.NotNil(tst, err, "Should return an error when decoding an expired payload")
	hmacMinusDecoded, err := DecodeHMAC("01234567890~!@#$%^&*-_=+ABCabc", hmacMinusSigned, 30)
	assert.NotNil(tst, err, "Should return an error when decoding an expired payload")

	assert.NotEqual(tst, stringToEncode, hmacPlusDecoded, "The original string should NOT match the decrypted, expired payload.")
	assert.NotEqual(tst, stringToEncode, hmacMinusDecoded, "The original string should NOT match the decrypted, expired payload.")
	assert.Equal(tst, hmacPlusDecoded, "Error verifying time", `The returned value from an expired payload should be a string that says "Error verifying time"`)
	assert.Equal(tst, hmacMinusDecoded, "Error verifying time", `The returned value from an expired payload should be a string that says "Error verifying time"`)
}
