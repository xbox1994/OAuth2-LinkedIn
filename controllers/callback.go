package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CallbackController struct {
	beego.Controller
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

type ShareResponse struct {
	UpdateKey string `json:"updateKey"`
	UpdateUrl string `json:"updateUrl"`
	ErrorCode int32  `json:"errorCode"`
	Message   string `json:"message"`
}

func getAccessToken(code string) AccessToken {
	resp, err := http.PostForm("https://www.linkedin.com/oauth/v2/accessToken",
		url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {beego.AppConfig.String("host") + "/linkedin/auth/callback"},
			"client_id":     {beego.AppConfig.String("linkedin_client_id")},
			"client_secret": {beego.AppConfig.String("linkedin_client_secret")},
		})
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		panic(err)
	}
	return accessToken
}

func share(accessToken string) ShareResponse {

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
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("share response body: ", string(body))

	var shareResponse ShareResponse
	err = json.Unmarshal(body, &shareResponse)
	if err != nil {
		panic(err)
	}
	return shareResponse
}

func (c *CallbackController) Get() {
	// csrf validation
	csrfToken := c.GetSession("_csrf_Token")
	if csrfToken != c.GetString("state") {
		c.Data["json"] = map[string]interface{}{"error": "xsrf error"}
		c.ServeJSON()
		return
	}

	// user cancel authorization request or linkedin server error
	if c.GetString("error") != "" {
		c.Data["json"] = map[string]interface{}{
			"error":             c.GetString("error"),
			"error_description": c.GetString("error_description")}
		c.ServeJSON()
		return
	}

	// user accept authorization request
	code := c.GetString("code")
	if code != "" {
		// get access token by code
		accessToken := getAccessToken(code)
		// share by access token
		shareResponse := share(accessToken.AccessToken)
		c.Data["json"] = &shareResponse
		c.ServeJSON()
		return
	}
}

