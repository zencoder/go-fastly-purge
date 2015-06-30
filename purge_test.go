package fastly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FastlyPurgeSuite struct {
	suite.Suite
	validResponseJSON string
}

const (
	VALID_PURGE_URL       string = "http://www.example.com/test/sample.jpg"
	VALID_PURGE_STATUS_OK string = "ok"
	VALID_PURGE_ID        string = "154-1434616760-1946753"
)

func TestFastlyPurgeSuite(t *testing.T) {
	suite.Run(t, new(FastlyPurgeSuite))
}

func (s *FastlyPurgeSuite) SetupTest() {
	status := VALID_PURGE_STATUS_OK
	id := VALID_PURGE_ID
	jsonBytes, _ := json.Marshal(&PurgeResponse{
		Status: &status,
		ID:     &id,
	})
	s.validResponseJSON = string(jsonBytes)
}

func (s *FastlyPurgeSuite) SetupSuite() {

}

func (s *FastlyPurgeSuite) TestPurgeURLInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_URL, b)
		assert.Equal(s.T(), "", r.Header.Get("Fastly-Soft-Purge"))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(VALID_PURGE_URL, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), VALID_PURGE_ID, id)
}

func (s *FastlyPurgeSuite) TestPurgeURLSoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_URL, b)
		assert.Equal(s.T(), "1", r.Header.Get("Fastly-Soft-Purge"))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(VALID_PURGE_URL, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), VALID_PURGE_ID, id)
}
