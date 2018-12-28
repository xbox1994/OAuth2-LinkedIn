package controllers

import (
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

func (c *CallbackController) Get() {
	// xsrf validation
	xsrfToken := c.GetSession("_xsrf_Token")
	if xsrfToken != c.GetString("state") {
		c.Data["json"] = map[string]interface{}{"error": "xsrf error"}
		c.ServeJSON()
		return
	}

	// user cancel authorization request
	if c.GetString("error") != "" {
		c.Data["json"] = map[string]interface{}{
			"error":             c.GetString("error"),
			"error_description": c.GetString("error_description"),
			"state":             c.GetString("state")}
		c.ServeJSON()
		return
	}

	// user accept authorization request
	code := c.GetString("code")
	if code != "" {
		resp, err := http.PostForm("https://www.linkedin.com/oauth/v2/accessToken",
			url.Values{
				"grant_type":    {"authorization_code"},
				"code":          {code},
				"redirect_uri":  {beego.AppConfig.String("host") + "/auth/linkedin/callback"},
				"client_id":     {beego.AppConfig.String("linkedin_client_id")},
				"client_secret": {beego.AppConfig.String("linkedin_client_secret")},
			})
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var accessToken AccessToken
		err = json.Unmarshal(body, &accessToken)
		if err != nil {
			panic(err)
		}
		c.Data["json"] = &accessToken
		c.ServeJSON()
	}
}
