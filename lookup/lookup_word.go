package lookup

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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


func (a *API) performLookupRequest(method string) ([]map[string]interface{}, error ) {
	var body io.Reader
	var content []map[string]interface{}
	req, err := http.NewRequest(method, a.EndPoint.String(), body)
	if err != nil {
		return content, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
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


// GetWord takes a string and searches and returns the best scoring match from:
// "https://api.datamuse.com/words/"
// If the string does not exactly match it will return an empty struct 
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
