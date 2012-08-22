package amiando

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ApiKey struct {
	ID			ID `json:"id"`
	Enabled		bool `json:"enabled"`
	Key 		string `json:"key"`
	Identifier 	string `json:"identifier"`
	Name		string `json:"name"`
}

func CreateAPIKey(name string) (key string, err error) {
	type Result struct {
		ResultBase
		Id 		ID `json:"id"`
		ApiKey 	ApiKey `json:"apiKey"`
	}
	var result Result
	

	posturl := "http://www.amiando.com/api/apiKey/create"

	values := url.Values{ "name": {name}}

	url := posturl + "?version=1&format=json"
	r, err := http.PostForm(url, values)
	if err != nil {
		fmt.Printf("error posting values: %s", err)
		return
	}
	j, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(j, &result)
	if err != nil {
		return "", err
	}

	if result.Success || len(result.Errors) == 0 {
		err = nil
	} else {
		err = &Error{result.Errors}
	}

	if err != nil {
		return "", err
	}
 	
 	return result.ApiKey.Key, nil
 		
}