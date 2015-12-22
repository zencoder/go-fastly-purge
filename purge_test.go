package fastlypurge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FastlyPurgeSuite struct {
	suite.Suite
	validResponseJSON     string
	validResponseNoIDJSON string
}

const (
	VALID_PURGE_PATH              string = "/test/sample.jpg"
	VALID_PURGE_STATUS_OK         string = "ok"
	VALID_PURGE_ID                string = "154-1434616760-1946753"
	VALID_PURGE_SERVICE           string = "3h16jblNfHsGtnGcUxA32F"
	VALID_PURGE_API_KEY           string = "bd16a4bbcf66be5fdeb955a19cb76a32"
	VALID_PURGE_KEY               string = "test/key"
	VALID_PURGE_HEADER_SOFT_PURGE string = "Fastly-Soft-Purge"
	VALID_PURGE_HEADER_KEY        string = "Fastly-Key"
	VALID_PURGE_API_ENDPOINT      string = "https://api.fastly.com"
)

const (
	INVALID_PURGE_MODE int64 = 5
)

var (
	VALID_PURGE_ALL_PATH string = fmt.Sprintf("/service/%s/purge_all", VALID_PURGE_SERVICE)
	VALID_PURGE_KEY_PATH string = fmt.Sprintf("/service/%s/purge/%s", VALID_PURGE_SERVICE, VALID_PURGE_KEY)
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
	jsonBytesNoID, _ := json.Marshal(&PurgeResponse{
		Status: &status,
	})
	s.validResponseNoIDJSON = string(jsonBytesNoID)
}

func (s *FastlyPurgeSuite) SetupSuite() {

}

func (s *FastlyPurgeSuite) TestNewPurge() {
	p := NewPurge()
	assert.NotNil(s.T(), p)
	assert.Equal(s.T(), "", p.APIKey)
	assert.Equal(s.T(), "", p.FastlyURL)
}

func (s *FastlyPurgeSuite) TestNewPurgeWithAPIKey() {
	p := NewPurgeWithAPIKey(VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)
	assert.Equal(s.T(), VALID_PURGE_API_KEY, p.APIKey)
	assert.Equal(s.T(), VALID_PURGE_API_ENDPOINT, p.FastlyURL)
}

func (s *FastlyPurgeSuite) TestPurgeRequestErrorInvalidPurgeMode() {
	p := NewPurge()
	id, err := p.purgeRequest(VALID_PURGE_PATH, "PURGE", 5, true)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Invalid Purge Mode"), err)
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeRequestErrorInvalidURL() {
	p := NewPurge()
	id, err := p.PurgeURL("", PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Failed to parse URL")
	assert.Contains(s.T(), err.Error(), "empty url")
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeURLInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Host)
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), VALID_PURGE_ID, id)
}

func (s *FastlyPurgeSuite) TestPurgeURLInstantErrorJSONDecode() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		type InvalidJSON struct {
			Invalid chan error `json:"invalid"`
		}
		i := &InvalidJSON{}
		jsonBytes, _ := json.Marshal(i)
		fmt.Fprintf(w, string(jsonBytes))
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Failed to decode JSON with error: EOF"), err)
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeURLInstantError500ResponseCode() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Invalid response code, expected 200, got 500"), err)
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeURLInstantErrorInvalidStatus() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		status := "error"
		jsonBytes, _ := json.Marshal(&PurgeResponse{
			Status: &status,
		})
		fmt.Fprintf(w, string(jsonBytes))
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Purge failed with Status, error"), err)
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeURLInstantErrorNoID() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("No ID returned for Purge"), err)
	assert.Equal(s.T(), "", id)
}

func (s *FastlyPurgeSuite) TestPurgeURLSoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "PURGE", r.Method)
		assert.Equal(s.T(), VALID_PURGE_PATH, r.URL.Path)
		assert.Equal(s.T(), "1", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseJSON)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	id, err := p.PurgeURL(ts.URL+VALID_PURGE_PATH, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), VALID_PURGE_ID, id)
}

func (s *FastlyPurgeSuite) TestPurgeAllInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_ALL_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll(VALID_PURGE_SERVICE, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeAllInstantErrorNoAPIKey() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_ALL_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll(VALID_PURGE_SERVICE, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("API Key is required for Purge All"), err)
}

func (s *FastlyPurgeSuite) TestPurgeAllInstantErrorNoService() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_ALL_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll("", PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Service is required for Purge All"), err)
}

func (s *FastlyPurgeSuite) TestPurgeAllSoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_ALL_PATH, r.URL.Path)
		assert.Equal(s.T(), "1", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeAll(VALID_PURGE_SERVICE, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeyInstant() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_KEY_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, VALID_PURGE_KEY, PURGE_MODE_INSTANT)
	assert.Nil(s.T(), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeyInstantErrorNoAPIKey() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_KEY_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := newPurgeWithFastlyURL(ts.URL)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, VALID_PURGE_KEY, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("API Key is required for Purge By Key"), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeyInstantErrorNoService() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_KEY_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey("", VALID_PURGE_KEY, PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Service is required for Purge By Key"), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeyInstantErrorNoKey() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_KEY_PATH, r.URL.Path)
		assert.Equal(s.T(), "", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, "", PURGE_MODE_INSTANT)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Key is required for Purge By Key"), err)
}

func (s *FastlyPurgeSuite) TestPurgeKeySoft() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), VALID_PURGE_KEY_PATH, r.URL.Path)
		assert.Equal(s.T(), "1", r.Header.Get(VALID_PURGE_HEADER_SOFT_PURGE))
		assert.Equal(s.T(), VALID_PURGE_API_KEY, r.Header.Get(VALID_PURGE_HEADER_KEY))
		fmt.Fprintf(w, s.validResponseNoIDJSON)
	}))
	defer ts.Close()

	p := NewPurgeWithFastlyURLAndAPIKey(ts.URL, VALID_PURGE_API_KEY)
	assert.NotNil(s.T(), p)

	err := p.PurgeKey(VALID_PURGE_SERVICE, VALID_PURGE_KEY, PURGE_MODE_SOFT)
	assert.Nil(s.T(), err)
}
