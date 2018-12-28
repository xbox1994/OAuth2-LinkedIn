package controllers

import (
	"facebooksharejs/util"
	"fmt"
	"github.com/astaxie/beego"
)

type AuthorizationController struct {
	beego.Controller
}

func (c *AuthorizationController) Get() {
	host := beego.AppConfig.String("host")
	state := util.Generate(10)
	clientId := beego.AppConfig.String("linkedin_client_id")
	linkedinOauth2AuthorizationUrl := "https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s/auth/linkedin/callback&state=%s"

	uri := fmt.Sprintf(linkedinOauth2AuthorizationUrl, clientId, host, state)

	c.SetSession("_xsrf_Token", state)
	c.Redirect(uri, 302)
}
