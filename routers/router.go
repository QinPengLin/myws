package routers

import (
	"myws/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WsController{})

	beego.Router("/set_user", &controllers.UserApiController{}, "post:SetUser")  //设置用户
	beego.Router("/get_group_message", &controllers.UserApiController{}, "post:GetGroupMessage")  //获取群历史记录

}
