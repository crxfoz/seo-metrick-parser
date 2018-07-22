package worker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var reporterTest Reporter

func init() {
	reporterTest = make(map[string]map[string]interface{})
}

func TestReporter_AddUrl(t *testing.T) {
	u := "http://google.com"

	reporterTest.AddUrl(u)
	_, ok := reporterTest[u]
	assert.Equal(t, true, ok)
}

func TestReporter_AddForUrl(t *testing.T) {
	u := "http://yandex.ru"

	res := reporterTest.AddUrl(u)
	reporterTest.AddForUrl(res, "tyc", map[string]string{"value": "100"})

	v, ok := reporterTest[u]["tyc"]
	assert.Equal(t, true, ok)
	vv, ok := v.(map[string]string)
	assert.Equal(t, true, ok)

	assert.Equal(t, "100", vv["value"])
}
