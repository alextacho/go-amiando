package amiando

import (
	"encoding/json"
	// "fmt"
	"strings"
	"net/url"
)



func (self *Api) CreateTicketCategory(event, name, price, available, salestart, saleend string) (id string, err error) {
	type Result struct {
		ResultBase
		Id 	ID `json:"id"`
	}
	var result Result
	
	posturl := "https://www.amiando.com/api/event/"+event+"/ticketCategory/create"

	values := url.Values{
				"name": {name}, 
				"price": {price},
				"available": {available},
			}

	if salestart!="" {
		salestart = strings.Replace(salestart," ", "T",-1)
		values.Add("saleStart", salestart)
	}
	if saleend!="" {
		saleend = strings.Replace(saleend," ", "T",-1)
		values.Add("saleEnd", saleend)
	}
	
	// fmt.Println(values)

	j, err := self.httpPost(posturl, values)
	if err != nil {
		return "", err
	}

	// fmt.Println("Result:\n", PrettifyJSON(j))

	err = json.Unmarshal(j, &result)
	if err != nil {
		return "", err
	}

	if result.ResultBase.Success || len(result.ResultBase.Errors) == 0 {
		err = nil
	} else {
		err = &Error{result.Errors}
	}

	// fmt.Printf("%#v\n", result.Errors)

	if err != nil {
		return "", err
	}

	// fmt.Print("ticketcategory id: ", result.Id)
 	
 	return result.Id.String(), nil
 		
}