package fastlypurge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type PurgeResponse struct {
	Status *string `json:"status,omitempty"`
	ID     *string `json:"id,omitempty"`
}

type PurgeMode int64

const (
	PURGE_MODE_INSTANT PurgeMode = iota
	PURGE_MODE_SOFT
)

const (
	PURGE_API_ENDPOINT string = "https://api.fastly.com"
)

const (
	PURGE_HEADER_SOFT_PURGE string = "Fastly-Soft-Purge"
	PURGE_HEADER_KEY        string = "Fastly-Key"
)

type Purge struct {
	APIKey    string
	FastlyURL string
}

func NewPurge() *Purge {
	return &Purge{}
}

func newPurgeWithFastlyURL(fastlyURL string) *Purge {
	return &Purge{
		FastlyURL: fastlyURL,
	}
}

func NewPurgeWithAPIKey(apiKey string) *Purge {
	return &Purge{
		FastlyURL: PURGE_API_ENDPOINT,
		APIKey:    apiKey,
	}
}

func NewPurgeWithFastlyURLAndAPIKey(fastlyURL string, apiKey string) *Purge {
	return &Purge{
		FastlyURL: fastlyURL,
		APIKey:    apiKey,
	}
}

func (p *Purge) purgeRequest(reqURL string, httpMethod string, purgeMode PurgeMode, idExpected bool) (string, error) {
	if purgeMode != PURGE_MODE_INSTANT && purgeMode != PURGE_MODE_SOFT {
		return "", errors.New("Invalid Purge Mode")
	}

	_, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return "", fmt.Errorf("Failed to parse URL, with error: %s", err)
	}

	req, err := http.NewRequest(httpMethod, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create HTTP request, with error: %s", err.Error())
	}

	if purgeMode == PURGE_MODE_SOFT {
		req.Header.Add(PURGE_HEADER_SOFT_PURGE, "1")
	}

	if p.APIKey != "" {
		req.Header.Add(PURGE_HEADER_KEY, p.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to perform HTTP request with error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Invalid response code, expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var pr PurgeResponse
	if err := dec.Decode(&pr); err != nil {
		return "", fmt.Errorf("Failed to decode JSON with error: %s", err.Error())
	}

	if pr.Status == nil || *pr.Status != "ok" {
		return "", fmt.Errorf("Purge failed with Status, %s", *pr.Status)
	}

	if idExpected == true {
		if pr.ID == nil || *pr.ID == "" {
			return "", errors.New("No ID returned for Purge")
		} else {
			return *pr.ID, nil
		}
	} else {
		return "", nil
	}
}

func (p *Purge) PurgeURL(url string, purgeMode PurgeMode) (string, error) {
	return p.purgeRequest(url, "PURGE", purgeMode, true)
}

func (p *Purge) PurgeAll(service string, purgeMode PurgeMode) error {
	if p.APIKey == "" {
		return errors.New("API Key is required for Purge All")
	}
	if service == "" {
		return errors.New("Service is required for Purge All")
	}

	url := fmt.Sprintf("%s/service/%s/purge_all", p.FastlyURL, service)

	_, err := p.purgeRequest(url, "POST", purgeMode, false)

	return err
}

func (p *Purge) PurgeKey(service string, key string, purgeMode PurgeMode) error {
	if p.APIKey == "" {
		return errors.New("API Key is required for Purge By Key")
	}
	if service == "" {
		return errors.New("Service is required for Purge By Key")
	}
	if key == "" {
		return errors.New("Key is required for Purge By Key")
	}

	url := fmt.Sprintf("%s/service/%s/purge/%s", p.FastlyURL, service, key)

	_, err := p.purgeRequest(url, "POST", purgeMode, false)

	return err
}
