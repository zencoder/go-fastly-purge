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
	VALID_PURGE_URL               string = "http://www.example.com/test/sample.jpg"
	VALID_PURGE_STATUS_OK         string = "ok"
	VALID_PURGE_ID                string = "154-1434616760-1946753"
	VALID_PURGE_SERVICE           string = "3h16jblNfHsGtnGcUxA32F"
	VALID_PURGE_API_KEY           string = "bd16a4bbcf66be5fdeb955a19cb76a32"
	VALID_PURGE_KEY               string = "test/key"
	VALID_PURGE_HEADER_SOFT_PURGE string = "Fastly-Soft-Purge"
	VALID_PURGE_HEADER_KEY        string = "Fastly-Key"
	VALID_PURGE_API_ENDPOINT      string = "https://api.fastly.com"
)

var (
	VALID_PURGE_ALL_URL string = fmt.Sprintf("%s/service/%s/purge_all", VALID_PURGE_API_ENDPOINT, VALID_PURGE_SERVICE)
	VALID_PURGE_KEY_URL string = fmt.Sprintf("%s/service/%s/purge/%s", VALID_PURGE_API_ENDPOINT, VALID_PURGE_SERVICE, VALID_PURGE_KEY)
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
		assert.Equal(s.T(), VALID_PURGE_URL, string(b))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
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
		assert.Equal(s.T(), VALID_PURGE_URL, string(b))
		assert.Equal(s.T(), "1", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(VALID_PURGE_URL, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), VALID_PURGE_ID, id)
}

func (s *FastlyPurgeSuite) TestPurgeAllInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_ALL_URL, string(b))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll(VALID_PURGE_SERVICE, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeAllSoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_ALL_URL, string(b))
		assert.Equal(s.T(), "1", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll(VALID_PURGE_SERVICE, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeyInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_KEY_URL, string(b))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, VALID_PURGE_KEY, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeySoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(s.T(), VALID_PURGE_KEY_URL, string(b))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithOverrideURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, VALID_PURGE_KEY, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
}
