// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var c *Collector

func TestNewStatCollector(test *testing.T) {
	c = NewStatCollector("test")
	assert.Equal(test, c.Name(), "test", "Name should be test")
}

func TestDisableTimes(test *testing.T) {
	c.DisableTimes()
	assert.True(test, c.DisableTime, "Time collection should be disabled")
	assert.Nil(test, c.InitDate, "Init date should be nil")
	assert.Nil(test, c.ChangeDate, "Init date should be nil")

	c.Child("testTimes").V(5)
	assert.True(test, c.Child("testTimes").DisableTime, "Time collection should be disabled")
	assert.Nil(test, c.Child("testTimes").InitDate, "Init date should be nil")
	assert.Nil(test, c.Child("testTimes").ChangeDate, "Init date should be nil")

	c.EnableTimes()
	c.V(1)
	assert.False(test, c.DisableTime, "Time collection should be enabled")
	assert.NotNil(test, c.InitDate, "Init date should be non-nil")
	assert.NotNil(test, c.ChangeDate, "Init date should be non-nil")

	//js, _ := c.JSON()
	//fmt.Println(js)
}

func TestSetV(test *testing.T) {
	c.V(1)
	assert.Equal(test, 1, c.GetV(), "Should be int 1")

	c.V(10.05)
	assert.Equal(test, 10.05, c.GetV(), "Should be float 10.05")

	c.V("some string")
	assert.Equal(test, "some string", c.GetV(), "Should be 'some string'")
}

func TestChild(test *testing.T) {
	has := c.Has("willfail")
	assert.False(test, has, "Shouldn't have child willfail")

	c.Child("newkidontheblock").V(101)
	has = c.Has("newkidontheblock")
	assert.True(test, has, "Should have child newkidontheblock")
}

func TestJSON(test *testing.T) {
	cj := NewStatCollector("jsontest")
	cj.Ch("c1").V(1).Ch("cc1").V(1.1)
	cj.Child("c2").Value(2).Child("cc2").Value(2.1)
	cj.Child("cs").V("string value").Child("ccs").V("child string value")
	cj.Child("cs").Child("ccs2").V("child string value 2")
	cj.Child("cnv") //sets a child with no value

	cj.Ch("avg").AvgLen(10, 3)
	cj.Ch("avg").Avg(20)
	cj.Ch("avg").Avg(25)

	js, err := cj.JSON()
	assert.Nil(test, err, "Shouldn't be an error")

	assert.Contains(test, js, `"value": 1,`)
	assert.Contains(test, js, `"value": 1.1`)
	assert.Contains(test, js, `"value": 2,`)
	assert.Contains(test, js, `"value": 2.1`)
	assert.Contains(test, js, `"value": "child string value"`)
	assert.Contains(test, js, `"value": "child string value 2"`)
	assert.Contains(test, js, `"value": 0`)
	assert.Contains(test, js, `18.333334,`) //avg
}

func TestIncr(test *testing.T) {
	ct := NewStatCollector("float32").V(float32(10.0)).Incr()
	assert.Equal(test, float32(11.0), ct.GetV())

	ct = NewStatCollector("float64").V(float64(10.0)).Incr()
	assert.Equal(test, float64(11.0), ct.GetV())

	ct = NewStatCollector("int").V(int(10)).Incr()
	assert.Equal(test, int(11), ct.GetV())

	ct = NewStatCollector("int32").V(int32(10)).Incr()
	assert.Equal(test, int32(11), ct.GetV())

	ct = NewStatCollector("int64").V(int64(10)).Incr()
	assert.Equal(test, int64(11), ct.GetV())

	ct = NewStatCollector("string").V("somestring").Incr()
	assert.Equal(test, "somestring", ct.GetV())

	ct = NewStatCollector("nil").V(nil).Incr()
	assert.Nil(test, ct.GetV())
}

func TestDecr(test *testing.T) {
	ct := NewStatCollector("float32").V(float32(10.0)).Decr()
	assert.Equal(test, float32(9.0), ct.GetV())

	ct = NewStatCollector("float64").V(float64(10.0)).Decr()
	assert.Equal(test, float64(9.0), ct.GetV())

	ct = NewStatCollector("int").V(int(10)).Decr()
	assert.Equal(test, int(9), ct.GetV())

	ct = NewStatCollector("int32").V(int32(10)).Decr()
	assert.Equal(test, int32(9), ct.GetV())

	ct = NewStatCollector("int64").V(int64(10)).Decr()
	assert.Equal(test, int64(9), ct.GetV())

	ct = NewStatCollector("string").V("somestring").Decr()
	assert.Equal(test, "somestring", ct.GetV())

	ct = NewStatCollector("nil").V(nil).Decr()
	assert.Nil(test, ct.GetV())
}

func TestAvg(test *testing.T) {
	ca := NewStatCollector("avgtest")
	ca.AvgLen(10, 3)
	ca.Avg(20)
	ca.Avg(30)
	v := ca.GetV().([]float32)
	//fmt.Println(v)
	assert.Equal(test, v[0], float32(20))

	ca.Avg(20)
	ca.Avg(20)
	ca.Avg(10)
	ca.Avg(5)
	v = ca.GetV().([]float32)
	//fmt.Println(v)
	assert.Equal(test, v[0], float32(11.666667))

	cb := NewStatCollector("avgtest2")
	cb.Avg(0)
	cb.Avg(100)
	cb.Avg(200)
	cb.Avg(300)
	cb.Avg(400)
	cb.Avg(500)
	cb.Avg(600)
	cb.Avg(700)
	cb.Avg(800)
	cb.Avg(900)
	cb.Avg(1000)
	v = cb.GetV().([]float32)
	//fmt.Println(v)
	assert.Equal(test, v[0], float32(550))
}
