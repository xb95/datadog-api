/*
 * Datadog API for Go
 *
 * Please see the included LICENSE file for licensing information.
 *
 * Copyright 2017 by authors and contributors.
 */

package datadog_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zorkian/go-datadog-api"
)

func TestHelperGetBoolSet(t *testing.T) {
	// Assert that we were able to get the boolean from a pointer field
	m := getTestMonitor()

	if attr, ok := datadog.GetBool(m.Options.NotifyNoData); ok {
		assert.Equal(t, true, attr)
	}
}

func TestHelperGetBoolNotSet(t *testing.T) {
	// Assert GetBool returned false for an unset value
	m := getTestMonitor()

	_, ok := datadog.GetBool(m.Options.NotifyAudit)
	assert.Equal(t, false, ok)
}

func TestHelperStringSet(t *testing.T) {
	// Assert that we were able to get the string from a pointer field
	m := getTestMonitor()

	if attr, ok := datadog.GetStringOk(m.Name); ok {
		assert.Equal(t, "Test monitor", attr)
	}
}

func TestHelperStringNotSet(t *testing.T) {
	// Assert GetString returned false for an unset value
	m := getTestMonitor()

	_, ok := datadog.GetStringOk(m.Message)
	assert.Equal(t, false, ok)
}

func TestHelperIntSet(t *testing.T) {
	// Assert that we were able to get the integer from a pointer field
	m := getTestMonitor()

	if attr, ok := datadog.GetIntOk(m.Id); ok {
		assert.Equal(t, 1, attr)
	}
}

func TestHelperIntNotSet(t *testing.T) {
	// Assert GetInt returned false for an unset value
	m := getTestMonitor()

	_, ok := datadog.GetIntOk(m.Options.RenotifyInterval)
	assert.Equal(t, false, ok)
}

func TestHelperGetJsonNumberSet(t *testing.T) {
	// Assert that we were able to get a JSON Number from a pointer field
	m := getTestMonitor()

	if attr, ok := datadog.GetJsonNumberOk(m.Options.Thresholds.Ok); ok {
		assert.Equal(t, json.Number("2"), attr)
	}
}

func TestHelperGetJsonNumberNotSet(t *testing.T) {
	// Assert GetJsonNumber returned false for an unset value
	m := getTestMonitor()

	_, ok := datadog.GetJsonNumberOk(m.Options.Thresholds.Warning)

	assert.Equal(t, false, ok)
}

func getTestMonitor() *datadog.Monitor {

	o := &datadog.Options{
		NotifyNoData:    datadog.Bool(true),
		Locked:          datadog.Bool(false),
		NoDataTimeframe: 60,
		Silenced:        map[string]int{},
		Thresholds: &datadog.ThresholdCount{
			Ok: datadog.JsonNumber("2"),
		},
	}

	return &datadog.Monitor{
		Query:   datadog.String("avg(last_15m):avg:system.disk.in_use{*} by {host,device} > 0.8"),
		Name:    datadog.String("Test monitor"),
		Id:      datadog.Int(1),
		Options: o,
		Type:    datadog.String("metric alert"),
		Tags:    make([]string, 0),
	}
}

func TestHelperGetStringId(t *testing.T) {
	// Assert GetStringId returned the id without a change if it is a string
	id, err := datadog.GetStringId("abc-xyz-123")
	assert.Equal(t, err, nil)
	assert.Equal(t, id, "abc-xyz-123")

	// Assert GetStringId returned the id as a string if it is an integer
	id, err = datadog.GetStringId(123)
	assert.Equal(t, err, nil)
	assert.Equal(t, id, "123")

	// Assert GetStringId returned an error if the id type is boolean
	_, err = datadog.GetStringId(true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported id type")

	// Assert GetStringId returned an error if the id type is float64
	_, err = datadog.GetStringId(5.2)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported id type")
}

func TestGetFloatFromInterface(t *testing.T) {
	var input interface{}

	input = nil
	val, auto, err := datadog.GetFloatFromInterface(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, auto)
	assert.Nil(t, val)

	input = 0.0
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, false, auto)
	assert.Equal(t, 0.0, *val)

	input = 12.3
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, false, auto)
	assert.Equal(t, 12.3, *val)

	input = 123
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, false, auto)
	assert.Equal(t, 123.0, *val)

	input = int64(1234567890123456789.0)
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, false, auto)
	assert.Equal(t, 1234567890123456789.0, *val)

	input = "auto"
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, true, auto)
	assert.Nil(t, val)

	input = "wrong!"
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.NotNil(t, err)

	input = false
	val, auto, err = datadog.GetFloatFromInterface(&input)
	assert.NotNil(t, err)
}
