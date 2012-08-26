package amiando

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// Api

func NewApi(key string) *Api {
	return &Api{Key: key}
}

type Api struct {
	Key  string
	http http.Client
}

func (self *Api) httpGet(url string) (body []byte, err error) {
	response, err := self.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func (self *Api) httpPost(url string, values url.Values) (body []byte, err error) {

	url = url + "?apikey=" + self.Key + "&version=1&format=json"
	r, err := http.PostForm(url, values)
	if err != nil {
		fmt.Printf("error posting values: %s", err)
		return
	}
	body, err = ioutil.ReadAll(r.Body)
	r.Body.Close()
	
	return body, err
}

func (self *Api) Call(resourceFormat string, resourceArg interface{}, result ErrorReporter) (err error) {
	result.Reset()

	sep := "?"
	if strings.Contains(resourceFormat, "?") {
		sep = "&"
	}
	resourceFormat = "http://www.amiando.com/api/" + resourceFormat + sep + "apikey=%s&version=1&format=json"
	url := fmt.Sprintf(resourceFormat, resourceArg, self.Key)

	j, err := self.httpGet(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, result)
	if err != nil {
		return err
	}

	return result.Err()
}

func (self *Api) Call_debug(resourceFormat string, resourceArg interface{}, result ErrorReporter) (err error) {
	result.Reset()

	sep := "?"
	if strings.Contains(resourceFormat, "?") {
		sep = "&"
	}
	resourceFormat = "http://www.amiando.com/api/" + resourceFormat + sep + "apikey=%s&version=1&format=json"
	url := fmt.Sprintf(resourceFormat, resourceArg, self.Key)

	fmt.Println("URL: ", url)
	j, err := self.httpGet(url)
	if err != nil {
		return err
	}
	fmt.Println("Result:\n", PrettifyJSON(j))

	// Catch nasty problems in json.Unmarshal
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			fmt.Println(err)
			debug.PrintStack()
			os.Exit(-1)
		}
	}()

	err = json.Unmarshal(j, result)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", result)

	return result.Err()
}

func (self *Api) Payment(id ID, out interface{}) (err error) {
	type Result struct {
		ResultBase
		Payment interface{} `json:"payment"`
	}
	result := Result{Payment: out}
	return self.Call("payment/%v", id, &result)
}

func (self *Api) TicketIDsOfPayment(paymentID ID) (ids []ID, err error) {
	type Result struct {
		ResultBase
		Tickets []ID `json:"tickets"`
	}
	var result Result
	err = self.Call("payment/%v/tickets", paymentID, &result)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

func (self *Api) Ticket(id ID, out interface{}) (err error) {
	type Result struct {
		ResultBase
		Ticket interface{} `json:"ticket"`
	}
	result := Result{Ticket: out}
	return self.Call("ticket/%v", id, &result)
}

func (self *Api) User(id ID, out interface{}) (err error) {
	type Result struct {
		ResultBase
		User interface{} `json:"user"`
	}
	result := Result{User: out}
	return self.Call("user/%v", id, &result)
}

func (self *Api) HostId(username string) (ids []ID, err error) {
	type Result struct {
		ResultBase
		Ids []ID `json:"ids"`
	}
	var result Result
	err = self.Call_debug("user/find?username=%v", username, &result)
	if err != nil {
		return nil, err
	}
	return result.Ids, nil
}

func (self *Api) CreateEvent(hostId string, title string, country string, date string, identifier string) (id string, err error) {
	type Result struct {
		Success bool `json:"success"`
		Errors []string `json:"errors"`
		Id ID `json:"id"`
	}
	var result Result
	

	posturl := "http://www.amiando.com/api/event/create"

	date = strings.Replace(date," ", "T",-1)

	// fmt.Println(hostId + "," + title + "," + country + "," + identifier + "," + date)
	values := url.Values{ "hostId": {hostId}, 
				"title": {title},
				"country": {country},
				"selectedDate": {date},
				"identifier": {identifier},
			}

	j, err := self.httpPost(posturl, values)
	if err != nil {
		return "", err
	}
	// fmt.Println("Result:\n", PrettifyJSON(j))

	err = json.Unmarshal(j, &result)
	if err != nil {
		return "", err
	}

	if result.Success || len(result.Errors) == 0 {
		err = nil
	} else {
		err = &Error{result.Errors}
	}
	
	// fmt.Printf("%#v\n", result.Errors)

	if err != nil {
		return "", err
	}
 	// fmt.Print("event id: %v", result.Id)
 	
 	return result.Id.String(), nil
 	
	
 	
}

func (self *Api) ActivateEvent(id string) error {
	type Result struct {
		Success bool `json:"success"`
		Errors []string `json:"errors"`
	}
	var result Result
	
	posturl := "https://www.amiando.com/api/event/"+id+"/activate"

	values := url.Values{}


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

func TestAmiandoWebHook(posturl string) error {

	values := url.Values{
		"eventId" : {"123"},
		"eventIdentifier" : {"test_startuplivevienna7"},
		"numberOfTickets" : {"2"},
		"ticketFirstName0" : {"Alex"},
		"ticketLastName0" : {"Tascha"},
		"ticketEmail0" : {"hemmshoe@gmail.com"},
		"ticketFirstName1" : {"Xandl"},
		"ticketLastName1" : {"Mann"},
		"ticketEmail1" : {"hemmshoe@gmail.com"},
	}

	fmt.Printf("posting to: %s", posturl)

	_, err := http.PostForm(posturl, values)
	if err != nil {
		fmt.Printf("error posting values: %s", err)
		return err
	}

	fmt.Printf("posting done")

	return err
}
