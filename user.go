package amiando

import (
	"encoding/json"
	"fmt"
	
	"net/url"
)

type AmiandoUser struct {
	Id			ID `json:"id"`
	FirstName	string `json:"firstName"`
	LastName 	string `json:"lastName"`
	Username 	string `json:"username"`
	Password	string `json:"password"`
	Language	string `json:"language"`
}



func (self *Api) CreateAmiandoUser(firstname, lastname, username, password, language string) (id string, err error) {
	type Result struct {
		ResultBase
		Id 		ID `json:"id"`
	}
	var result Result
	

	posturl := "https://www.amiando.com/api/user/create"

	values := url.Values{"firstName": {firstname}, 
				"lastName": {lastname},
				"username": {username},
				"password": {password},
				"language": {language},
			}

	j, err := self.httpPost(posturl, values)
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

	
 	
 	return result.Id.String(), nil
 		
}

func (self *Api) AddPaymentToUser(userid, accountHolder, bankName, swift, iban, country string) error {
	type Result struct {
		Success bool `json:"success"`
		Errors []string `json:"errors"`
	}
	var result Result
	
	posturl := "https://www.amiando.com/api/user/"+userid+"/bankAccount"

	values := url.Values{"accountHolder": {accountHolder}, 
				"bankName": {bankName},
				"swift": {swift},
				"iban": {iban},
				"country": {country},
			}

	j, err := self.httpPost(posturl, values)
	if err != nil {
		return err
	}

	fmt.Println("Result:\n", PrettifyJSON(j))

	err = json.Unmarshal(j, &result)
	if err != nil {
		return err
	}

	if result.Success || len(result.Errors) == 0 {
		return nil
	}
	return &Error{result.Errors}

	

}

func (self *Api) AddBillingInfoToUser(userid, firstName, lastName, company, street, zipCode, city, country string) error {
	type Result struct {
		Success bool `json:"success"`
		Errors []string `json:"errors"`
	}
	var result Result
	
	posturl := "https://www.amiando.com/api/user/"+userid+"/address/billing"

	values := url.Values{"firstName": {firstName}, 
				"lastName": {lastName},
				"company": {company},
				"street": {street},
				"zipCode": {zipCode},
				"city": {city},
				"country": {country},
			}

	j, err := self.httpPost(posturl, values)
	if err != nil {
		return err
	}

	fmt.Println("Result:\n", PrettifyJSON(j))

	err = json.Unmarshal(j, &result)
	if err != nil {
		return err
	}

	if result.Success || len(result.Errors) == 0 {
		return nil
	}
	return &Error{result.Errors}

	

}