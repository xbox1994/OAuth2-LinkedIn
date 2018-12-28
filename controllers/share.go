package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type ShareController struct {
	beego.Controller
}

type ShareResponse struct {
	UpdateKey string `json:"updateKey"`
	UpdateUrl string `json:"updateUrl"`
	ErrorCode int32  `json:"errorCode"`
	Message   string `json:"message"`
}

func (c *ShareController) Get() {
	requestBytes := []byte(`{
  "comment": "WPS share comment",
  "content": {
    "title": "WPS share title",
    "description": "WPS share description",
    "submitted-url": "http://thyrsi.com/t6/642/1545816258x2890202791.jpg",
    "submitted-image-url": "http://thyrsi.com/t6/642/1545816258x2890202791.jpg"
  },
  "visibility": {
    "code": "anyone"
  }
}`)
	req, err := http.NewRequest("POST", "https://api.linkedin.com/v1/people/~/shares?format=json",
		bytes.NewBuffer(requestBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetString("access_token"))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("share response body: ", string(body))

	var shareResponse ShareResponse
	err = json.Unmarshal(body, &shareResponse)
	if err != nil {
		panic(err)
	}
	c.Data["json"] = &shareResponse
	if err != nil {
		panic(err)
	}
	c.ServeJSON()
}
