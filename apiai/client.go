package apiai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL        = "https://api.api.ai/v1"
	defaultVersion = "20150910"
	defaultLang    = "en"
)

type ApiAIClient interface {
	Query(sessionID string, queries ...string) (*QueryResponse, error)
}

type Options struct {
	Version     string
	AccessToken string
	Lang        string
}

func NewClient(options Options) ApiAIClient {
	version := defaultVersion
	if len(options.Version) > 0 {
		version = options.Version
	}

	lang := defaultLang
	if len(options.Lang) > 0 {
		lang = options.Lang
	}

	var client = &http.Client{Timeout: time.Second * 20}
	a := apiAIClient{client: client, accessToken: options.AccessToken, version: version, lang: lang}
	return &a
}

type apiAIClient struct {
	client      *http.Client
	accessToken string
	version     string
	lang        string
}

func (a *apiAIClient) Query(sessionID string, queries ...string) (*QueryResponse, error) {
	var qryReq = QueryRequest{Query: queries, Lang: a.lang, SessionID: sessionID}
	reqBody, err := json.Marshal(qryReq)
	if err != nil {
		log.Printf("Could not encode query request as JSON: %v", err)
		return nil, err
	}

	url := baseURL + "/query?v=" + a.version
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Failed to create http request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.accessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := a.client.Do(req)
	if err != nil {
		log.Printf("Failed to send query request: %v", err)
		return nil, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read query response body: %v", err)
		return nil, err
	}

	var qryRes QueryResponse
	if err := json.Unmarshal(resBody, &qryRes); err != nil {
		log.Printf("Could not decode query response body as JSON: %v", err)
		return nil, err
	}
	return &qryRes, nil
}
