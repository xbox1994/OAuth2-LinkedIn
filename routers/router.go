package routers

import (
	"facebooksharejs/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/linkedin/auth/authorization", &controllers.AuthorizationController{})
    beego.Router("/linkedin/auth/callback", &controllers.CallbackController{})
}
