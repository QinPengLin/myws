package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"myws/common"
	"myws/extra"
	"myws/models"
	"strings"
	"strconv"
	"time"
	//"fmt"
)

//创建群
func  CreateGroup(key,uid,groupName,groupInitMembers string)  {
	successUid := ""
	failUid := ""
	var createUids []models.GroupMembers
	var groupMember models.GroupMembers

	if groupInitMembers!="" {//需要对传过来的uid做验证
		ds:=beego.AppConfig.String("group_uid_partition_code")
		uidStrArr,l:=extra.StrPartition(groupInitMembers,ds)
		//检测uid个数是否合法
		maxGroupMembers,_:=beego.AppConfig.Int("max_group_members")
		if l>maxGroupMembers {//超过群最大容纳人数
			var senMessagLenerr models.SendMessageAbnormal
			senMessagLenerr.Code=lenerr["code"]
			senMessagLenerr.AbnormalMsg=lenerr["msg"]
			ToMe(key,senMessagLenerr)
			return
		}

		for _,uidV := range uidStrArr {
			if GetUserToRedis(uidV) {
				if uidV!=uid {
					successUid = successUid + ds + uidV
					groupMember.Uid = uidV
					groupMember.Level = "2"
					createUids = append(createUids,groupMember)
				}
			}else {
				failUid = failUid+ds+uidV
			}
		}
	}

	groupMember.Uid = uid
	groupMember.Level = "1"
	createUids = append(createUids,groupMember)

	//创建群id
	groupId:=extra.CreateOnlKey()

	_,addErr:=AddGroupToRdis(groupId,createUids,groupName)
	if addErr!=nil {
		var senMessagGadderr models.SendMessageAbnormal
		senMessagGadderr.Code=gadderr["code"]
		senMessagGadderr.AbnormalMsg=gadderr["msg"]
		ToMe(key,senMessagGadderr)
		return
	}

	//成功
	var senMessagGadd models.SendMessageAndDataAbnormal
	senMessagGadd.Code=success["code"]
	senMessagGadd.AbnormalMsg=success["msg"]
	senMessagGadd.Data=groupId
	ToMeData(key,senMessagGadd)
	return

}

//创建群并且维护uid=>Groupids的关系
func AddGroupToRdis(groupId string,groupData []models.GroupMembers,groupName string) (bool,error) {
	//获取createUids在user_to_groups中的所有信息
	d := common.RedisConn.Get()
	defer d.Close()

	//循环获取
	existUserToGroup:=make(map[string]string)
	for _,groupDataV := range groupData {
		ok, _ := redis.String(d.Do(models.REDIS_HGET,models.USER_TO_GROUPS,groupDataV.Uid))
		if ok=="" {
			existUserToGroup[groupDataV.Uid] = groupId
		}else {
			existUserToGroup[groupDataV.Uid] = ok + "|" + groupId
		}
		existUserToGroup[groupDataV.Uid] = strings.Trim(existUserToGroup[groupDataV.Uid], "|")
	}

	//新增群
	groupMemberDataJson, _ := json.Marshal(groupData)
	groupMemberDataJsonStr:=string(groupMemberDataJson)
	_, hsetErr := d.Do(models.REDIS_HSET,models.GROUP_LIST, groupId, groupMemberDataJsonStr)
	if hsetErr != nil {
		logs.Info("AddGroupToRdis保存账号错误, hsetErr:%v", hsetErr)
		return false,hsetErr
	}
	//新增群id=>群名
	_, hsetToNameErr := d.Do(models.REDIS_HSET,models.GROUPID_TO_GROUPNAME, groupId, groupName)
	if hsetToNameErr != nil {
		logs.Info("AddGroupToRdis保存账号错误, hsetToNameErr:%v", hsetToNameErr)
		return false,hsetToNameErr
	}
	//批量修改或者新增
	_, hmsetErr := d.Do(models.REDIS_HMSET, redis.Args{}.Add(models.USER_TO_GROUPS).AddFlat(existUserToGroup)...)
	if hmsetErr!=nil {
		logs.Info("AddGroupToRdis查询redis错误, hmsetErr:%v", hmsetErr)
		return false,hmsetErr
	}
	return true,nil
}

//向某个群发送消息
func SendToGroup(fromUid string,fromData *models.FromMessage,groupData string,key string) {

	strNowTime:=strconv.FormatInt(time.Now().Unix(),10)
	var senMessagToMe models.SendMessage

	senMessagToMe.Time=strNowTime
	senMessagToMe.FromUid=fromUid 
	senMessagToMe.ToGroupId=fromData.ToGroupId
	senMessagToMe.Msg=fromData.Content

	//将groupData解析成数组
	var groupDataArr []models.GroupMembers

	json.Unmarshal([]byte(groupData), &groupDataArr)
	inGroup:=false
	for _,gV := range groupDataArr {
		if fromUid==gV.Uid {
			inGroup = true
			break
		}
	}
	//需要先判断发消息人是否在该群中
	if inGroup {//存在才能群发
		//将群发消息写入群发记录表
		insertData := &models.GroupChatRecord{
			GroupId:       fromData.ToGroupId,
			FormUid:       fromUid,
			Content:       fromData.Content,
			CreationTime:  strNowTime,
		}
		models.InsertGroupChatRecord(insertData)
	    for _,groupDataArrV := range groupDataArr {
		    go OneToOne(groupDataArrV.Uid,senMessagToMe)
		}
	}else{//不存在就提示无权限
		var senMessagNoauth models.SendMessageAbnormal
		senMessagNoauth.Code=noauth["code"]
		senMessagNoauth.AbnormalMsg=noauth["msg"]
		ToMe(key,senMessagNoauth)
	}
}