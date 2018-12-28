package routers

import (
	"facebooksharejs/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/auth/linkedin", &controllers.AuthorizationController{})
    beego.Router("/auth/linkedin/callback", &controllers.CallbackController{})
    beego.Router("/auth/linkedin/share", &controllers.ShareController{})
}
