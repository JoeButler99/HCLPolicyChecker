package lookup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "strings"
	"encoding/json"
	"io"
)

type API struct {
	EndPoint *url.URL
	Client   *http.Client
}

type Word struct {
	Name  string   `json:"word"`
	score uint64   `json:"score"`
	tags  interface{}  `json:"tags"`
}

var dictionaryAPIRUL = "https://api.datamuse.com/"

func (a *API) Request(req *http.Request) ([]byte, error) {
	// req.Header.Add("Accept", "application/json, */*")
	log.Println("here")
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("or here?")
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("or hereeee?")
	defer resp.Body.Close()
	// log.Println(string(res))
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		return res, nil
	case http.StatusNoContent, http.StatusResetContent:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed")
	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("service is not available: %s", resp.Status)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error: %s", resp.Status)
	case http.StatusConflict:
		return nil, fmt.Errorf("conflict: %s", resp.Status)
	}

	return nil, fmt.Errorf("unknown response status: %s", resp.Status)
}

func (a *API) performLookupRequest(method string) ([]map[string]interface{}, error ) {
	var body io.Reader
	var content []map[string]interface{}
	log.Println(a.EndPoint.String())
	req, err := http.NewRequest(method, a.EndPoint.String(), body)
	if err != nil {
		log.Println(err)
		return content, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	log.Println(fmt.Sprintf("%+v",req))
	res, err := http.Get(a.EndPoint.String())
	if err != nil {
		return content, err
	}
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return content, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(resp, &content)
	if err != nil {
		log.Println(err)
		return content, err
	}
	return content, nil
}

func createAPI(endPoint *url.URL) *API {
	a := new(API)
	a.EndPoint = endPoint
	return a
}

func GetWord(word string) (map[string]interface{}, error) {
	queryParams := make(url.Values)
	queryParams.Add("sp", word)
	queryParams.Add("md", "p")
	var invalidWord map[string]interface{}
	endPoint, err := url.ParseRequestURI(dictionaryAPIRUL + "words")
	if err != nil {
		return invalidWord, err
	}
	endPoint.RawQuery = queryParams.Encode()
	a := createAPI(endPoint)
	content, err := a.performLookupRequest("GET")
	if err != nil {
		return invalidWord, err
	}
	if len(content) == 0 {
		return invalidWord, errors.New("word not found from lookup")
	}
	apiWord := content[0]
	if apiWord["word"].(string) == word {
		return apiWord, nil
	}
	return invalidWord, errors.New("words not found from lookup")
}
