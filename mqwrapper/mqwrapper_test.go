// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mqwrapper

import (
	"testing"

	"github.com/TeamFairmont/amqp"
	"github.com/TeamFairmont/gabs"
	"github.com/stretchr/testify/assert"
)

func TestConnectMQ(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ConnectMQ in short mode (Connectivity required).")
	}

	_, err := ConnectMQ("amqp://guest:guest@localhost:5672/")
	assert.Nil(t, err, "Err is nil, we're connected")
}

func TestCreateConsumeTempQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CreateConsumeTempQueue in short mode (Connectivity required).")
	}

	c, err := ConnectMQ("amqp://guest:guest@localhost:5672/")
	assert.Nil(t, err, "Err is nil, we're connected")

	if !t.Failed() {
		_, _, err = CreateConsumeTempQueue(c.Channel)
		assert.Nil(t, err, "Err is nil, we have a queue")
	}
}

func TestCreateConsumeNamedQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CreateConsumeNamedQueue in short mode (Connectivity required).")
	}

	c, err := ConnectMQ("amqp://guest:guest@localhost:5672/")
	assert.Nil(t, err, "Err is nil, we're connected")

	if !t.Failed() {
		_, _, err = CreateConsumeNamedQueue("unitTestQueue", c.Channel)
		assert.Nil(t, err, "Err is nil, we have a queue")
	}
}

func TestPublishCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PublishCommand in short mode (Connectivity required).")
	}

	c, err := ConnectMQ("amqp://guest:guest@localhost:5672/")
	assert.Nil(t, err, "Err is nil, we're connected")

	var q *amqp.Queue
	if !t.Failed() {
		q, _, err = CreateConsumeTempQueue(c.Channel)
		assert.Nil(t, err, "Err is nil, we have a queue")
	}

	if !t.Failed() {
		err = PublishCommand(c.Channel, "test", "test", &gabs.Container{}, q.Name)
		assert.Nil(t, err, "Err is nil, we have a queue")
	}
}
