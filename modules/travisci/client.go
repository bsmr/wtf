package travisci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var TRAVIS_HOSTS = map[bool]string{
	false: "travis-ci.org",
	true:  "travis-ci.com",
}

func BuildsFor(apiKey string, pro bool) (*Builds, error) {
	builds := &Builds{}

	travisAPIURL.Host = "api." + TRAVIS_HOSTS[pro]

	resp, err := travisRequest(apiKey, "builds")
	if err != nil {
		return builds, err
	}

	parseJson(&builds, resp.Body)

	return builds, nil
}

/* -------------------- Unexported Functions -------------------- */

var (
	travisAPIURL = &url.URL{Scheme: "https", Path: "/"}
)

func travisRequest(apiKey string, path string) (*http.Response, error) {
	params := url.Values{}
	params.Add("limit", "10")

	requestUrl := travisAPIURL.ResolveReference(&url.URL{Path: path, RawQuery: params.Encode()})

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Travis-API-Version", "3")

	bearer := fmt.Sprintf("token %s", apiKey)
	req.Header.Add("Authorization", bearer)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(resp.Status)
	}

	return resp, nil
}

func parseJson(obj interface{}, text io.Reader) {
	jsonStream, err := ioutil.ReadAll(text)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonStream))

	for {
		if err := decoder.Decode(obj); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
}
