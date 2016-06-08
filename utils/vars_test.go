// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilString(tst *testing.T) {
	assert.Equal(tst, "yesnil", NilString(nil, "yesnil"), "Should be yesnil")
	assert.Equal(tst, "original string", NilString("original string", "yesnil"), "Should be original string")
}

func TestStringInSlice(tst *testing.T) {
	slice := []string{"abc", "123"}
	assert.Equal(tst, true, StringInSlice("123", slice), "Should contain 123")
	assert.Equal(tst, false, StringInSlice("test", slice), "Should not contain test")
}
