package routers

import (
	"fdsp4/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/load", &controllers.MainController{}, "*:Load")
	beego.Router("/list", &controllers.MainController{}, "*:List")
	beego.Router("/clear", &controllers.MainController{}, "*:Clear")
	beego.Router("/json", &controllers.MainController{}, "*:Json")
	beego.Router("/json2", &controllers.MainController{}, "*:GetAdJson2")
	beego.Router("/proto", &controllers.MainController{}, "*:Proto")
}
