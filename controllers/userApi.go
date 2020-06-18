package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"myws/common"
	"myws/models"
	"time"
	"strconv"
        "strings"
)



type UserApiController struct {
	beego.Controller
}

//将用户设置到聊天系统
func (c *UserApiController) SetUser() {
	user_id := c.GetString("user_id")
	item := c.GetString("item")
	if user_id=="" || item=="" {
		logs.Info("参数错误,请检查")
		c.Data["json"] = models.ApiResponse{Code: models.ErrParam, Message: "参数错误,请检查", Data: ""}
		c.ServeJSON()
		return
	}

	d := common.RedisConn.Get()
	defer d.Close()

	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,item+"_"+user_id)
	if err != nil {
		logs.Info("查询redis错误,Err:%v", err)
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "redis保存账号错误1", Data: ""}
		c.ServeJSON()
		return
	}

	if ok.(int64)==1 {
		logs.Info("账号存在")
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "账号存在", Data: ""}
		c.ServeJSON()
		return
	}

	strNowTime:=strconv.FormatInt(time.Now().Unix(),10)
	_, err1 := d.Do(models.REDIS_HSET,models.USER_LIST, item+"_"+user_id, strNowTime)
	if err1 != nil {
		logs.Info("redis保存账号错误, Err:%v", err1)
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "redis保存账号错误", Data: ""}
		c.ServeJSON()
		return
	}

	logs.Info("redis保存账号成功")
	c.Data["json"] = models.ApiResponse{Code: models.SUCCESS, Message: "redis保存账号成功", Data: ""}
	c.ServeJSON()
	return
}
//获取群消息
func (c *UserApiController) GetGroupMessage() {
	userId := c.GetString("user_id")
	groupId := c.GetString("group_id")
	page, _ := c.GetInt("page")
	if groupId=="" || userId=="" {
		logs.Info("参数错误,请检查")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrParam, 
			Message: "参数错误,请检查", 
			Data: "",
		}
		c.ServeJSON()
		return
	}
	d := common.RedisConn.Get()
	defer d.Close()

	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,userId)
	if err != nil {
		logs.Info("查询redis错误,Err:%v", err)
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem, 
			Message: "redis获取失败1", 
			Data: "",
		}
		c.ServeJSON()
		return
	}
	if ok.(int64)!=1 {
		logs.Info("账号不存在")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem, 
			Message: "账号不存在", 
			Data: "",
		}
		c.ServeJSON()
		return
	}

	//查找该用户用有的群
	groupDdat:=GetUidToGroupDataRedis(userId)
	if !strings.Contains(groupDdat,groupId) {
		logs.Info("你没有该群权限")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem,
			Message: "你没有该群权限",
			Data: "",
			}
		c.ServeJSON()
		return
	}

	//获取群聊天历史记录
	limit,_ :=beego.AppConfig.Int("max_history_msg_record")
	data:= models.GetGroupChatRecord(groupId,page,limit)

	c.Data["json"] = models.ApiResponse{
		Code: models.SUCCESS,
		Message: "ok1",
		Data: data,
		Page: strconv.Itoa(page),
		}
	c.ServeJSON()
	return
}