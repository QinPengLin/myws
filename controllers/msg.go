package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"myws/extra"
	"myws/models"
	"myws/common"
	"strconv"
	"time"
	"fmt"
)

var (
	clients   = make(map[string]*websocket.Conn)     //保存链接
	broadcast = make(chan *models.SendMessage)              //消息通道

	uidToKey  =  make(map[string]string)             //保存uid和key的对应关系
	groupIdToUid  =  make(map[string]string)         //保存群ID和UId的对应关系

	success       =  make(map[string]string)
	morano        =  make(map[string]string)
	morazy        =  make(map[string]string)
	morayes       =  make(map[string]string)
	userno        =  make(map[string]string)
	useryes       =  make(map[string]string)
	verify        =  make(map[string]string)
	illegal       =  make(map[string]string)
	gnerr         =  make(map[string]string)
	lenerr        =  make(map[string]string)
	gadderr       =  make(map[string]string)
	grouperr      =  make(map[string]string)
	noauth        =  make(map[string]string)
	notome        =  make(map[string]string)

	general       =  make(map[string]string)
)

func init() {
	initStatusCode()
	go SendWorld()
}

//广播给所有玩家，也可以广播给某个人
func SendWorld() {
	for {
		msg := <-broadcast

		if msg.ToUid!=models.Radio {
			logs.Info("不是世界消息【"+msg.ToUid+"】")
			continue
		}

		msgJson,_:=extra.StructToFromMessageJson(msg)
		sendMsg:=extra.Base64EncodeString(msgJson)
		for k,client := range clients {
			logs.Info("信息需要发送给【"+k+"】=>uid【"+msg.ToUid+"】")
			err := client.WriteJSON(sendMsg)
			if err != nil {
				logs.Info("key为【"+k+"】断掉了 error:",err)
				client.Close()
				delete(clients, k)
			}
		}
	}
}

func  Read(wsLink *websocket.Conn,key string)  {
	nowUid:=""
	verifyUser:=beego.AppConfig.String("verify_user")
	for {
		logs.Info("【"+key+"】进入获取信息模式")
		_, message, err := wsLink.ReadMessage()
		if err != nil {
			logs.Info("key【"+key+"】离开了 error:",err)
			if nowUid!="" {
			    //删除uidToKey
			    delete(uidToKey, nowUid)
			}
			//关闭连接，删除clients
			wsLink.Close()
			delete(clients, key)
			break
		}

		messageStr,_ := extra.Base64DecodeString(string(message))
		logs.Info("读取了链接【"+key+"】的发来的信息:",messageStr)
		fromData,_:=extra.JsonToFromMessageStruct(messageStr)

		//验证玩家是否绑定uid
		if verifyUser==models.VerifyUser {
		    //如果当前的nowUid为空说明还未设置uid必须要先设置uid才能发送消息
		    if nowUid=="" && fromData.Method!=models.AddUid {
				ToMe(key,extra.ReturnMessageAbnormal(verify))
			    continue
		    }
		}

		//异常状态码
		codes:=success["code"]
		abnormalMsg:=success["msg"]

		/**新增或者修改绑定uid=>key*/
		//如果是修改或者新增uid
		if fromData.Method==models.ModifyUid || fromData.Method==models.AddUid {
			if fromData.Uid=="" {
				//uid不能为空
				logs.Info("修改或者新增的uid【" + fromData.Uid + "】不能为空")
				ToMe(key,extra.ReturnMessageAbnormal(illegal))
				continue
			}
			//先查看之前是否有关系
			value,ok := uidToKey[fromData.Uid]
			if ok {//存在
				if value==key {
					logs.Info("uid【"+fromData.Uid+"】无需修改")
					codes=morano["code"]
					abnormalMsg=morano["msg"]
				}else {
					logs.Info("将uid【" + fromData.Uid + "】已经存在")
					codes=morazy["code"]
					abnormalMsg=morazy["msg"]
				}
			}else {//不存在就修改

				/**判断用户是否合法*/
				if verifyUser==models.VerifyUser {
					if !GetUserToRedis(fromData.Uid) {
						logs.Info("修改或者新增的uid【" + fromData.Uid + "】不合法")
						ToMe(key,extra.ReturnMessageAbnormal(illegal))
						continue
					}
				}
				/**判断用户是否合法*/

				//如果不存在，就先查找key是否存在于uidToKey如果存在就删除不存在就新增（这就是修改uid功能）
				for uidToKeyK,uidToKeyV := range uidToKey {
					if uidToKeyV==key {//存在就删除
						delete(uidToKey, uidToKeyK)
						logs.Info("删除了之前的uid【"+uidToKeyK+"】key为【"+key+"】的用户列表")
						break
					}
				}
				//新增或者删除
				uidToKey[fromData.Uid]=key

				//记录当前的uid
				nowUid=fromData.Uid
				logs.Info("将uid【"+fromData.Uid+"】绑定到了key为【"+key+"】")
				codes=morayes["code"]
				abnormalMsg=morayes["msg"]
			}

			general["code"]=codes
			general["msg"]=abnormalMsg
			ToMe(key,extra.ReturnMessageAbnormal(general))
			logs.Info("key为【"+key+"】的用户列表给自己发消息")
			continue
		}
		/**新增或者修改绑定uid=>key*/
		
		//如果是创建群就走创建群的方法
		if fromData.Method==models.CreateGroup {
			if fromData.GroupName=="" {//群名不能为空
			    ToMe(key,extra.ReturnMessageAbnormal(gnerr))
			    continue
			}
			go CreateGroup(key,nowUid,fromData.GroupName,fromData.GroupInitMembers)
			continue
		}

		//如果是群聊就走群聊的方法
		if fromData.Method==models.ToGroup {
			//查找群id是否存在
			groupData:=GetGroupDataToRedis(fromData.ToGroupId)
			if groupData=="" {//群id不能为空
				ToMe(key,extra.ReturnMessageAbnormal(grouperr))
				continue
			}
			go SendToGroup(nowUid,fromData,groupData,key)
			continue
		}

		/**将fromData信息放到SendMessage*/
		strNowTime:=strconv.FormatInt(time.Now().Unix(),10)
		var senMessag models.SendMessage
		senMessag.Time=strNowTime
		senMessag.FromUid=nowUid
		senMessag.ToUid=fromData.ToUid
		senMessag.Msg=fromData.Content
		/**将fromData信息放到SendMessage*/

		if senMessag.ToUid!="" {
			abnormal:=extra.ReturnMessageAbnormal(success)
			senMessag.Abnormal=&abnormal

		    if senMessag.ToUid==models.Radio{
		        //将信息放入通道,世界消息(只有发送给所有人才会把消息放入通道,这个通道中也可以把消息发送给某个人)
		        broadcast <- &senMessag
				continue
		    }else{
		    	//单对单
				//获取toUid对应的key
	            _,ok := uidToKey[fromData.ToUid]
	            if !ok {//不存在就保存记录等到玩家上线发送并且跳过该次循环（后面做）
					logs.Info("内存中无该玩家uid【"+fromData.ToUid+"】OneToOne")
					go ToMe(key,extra.ReturnMessageAbnormal(userno))
					continue
	            }else{
	            	//不能给自己发送消息（给自己发消息会产生websocket并发写入）
					if fromData.ToUid==nowUid {
						logs.Info("ToUid【" + fromData.ToUid + "】不能给自己发消息")
						ToMe(key,extra.ReturnMessageAbnormal(notome))
						continue
					}
					//发送给对应的玩家
					go OneToOne(fromData.ToUid,senMessag)

					//给自己发消息
					var senMessagToMe models.SendMessage
					senMessagToMe.Time=strNowTime
					senMessagToMe.FromUid=nowUid
					senMessagToMe.ToUid=fromData.ToUid
					senMessagToMe.Msg=fromData.Content

					userYesMsg:=extra.ReturnMessageAbnormal(useryes)
					senMessagToMe.Abnormal=&userYesMsg
					go OneToOne(nowUid,senMessagToMe)
					continue
				}
			}
		}
		//ToUid为空的情况
		logs.Info("ToUid【" + fromData.ToUid + "】不能为空")
		ToMe(key,extra.ReturnMessageAbnormal(illegal))
	}
}

//将消息发送给某一个人
func OneToOne(toUid string,msg models.SendMessage) {
	//获取toUid对应的key
	key,ok := uidToKey[toUid]
	if !ok {//不存在就保存记录等到玩家上线发送并且跳过该次循环（后面做）
		logs.Info("内存中无该玩家uid【"+toUid+"】OneToOne")
		return
	}
		
	msgJson,_:=extra.StructToFromMessageJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	err := clients[key].WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key+"】断掉了OneToOneerror:",err)
		clients[key].Close()
		delete(clients, key)
	}
	logs.Info("信息需要发送给【"+key+"】OneToOne")
}

//服务器返回给自己的消息
func ToMe(key string,msg models.SendMessageAbnormal)  {
	msgJson,_:=extra.StructToSendMessageAbnormalJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	err := clients[key].WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key+"】断掉了ToMe error:",err)
		clients[key].Close()
		delete(clients, key)
	}
}

//服务器返回带有数据的消息给自己
func ToMeData(key string,msg models.SendMessageAndDataAbnormal)  {
	msgJson,_:=extra.StructToSendMessageAndDataAbnormalJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	err := clients[key].WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key+"】断掉了ToMe error:",err)
		clients[key].Close()
		delete(clients, key)
	}
}

//获取指定用户是否在redis中
func GetUserToRedis(user_id string) bool {
	d := common.RedisConn.Get()
	defer d.Close()
	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,user_id)
	if err != nil {
		logs.Info("GetUserToRedis查询redis错误,Err:%v", err)
		return false
	}

    fmt.Println(ok)
	if ok.(int64)==1 {
		logs.Info("GetUserToRedis查询redis账号存在")
		return true
	}

	return false
}

//获取指群ID是否在redis中
func GetGroupDataToRedis(group_id string) string {
	d := common.RedisConn.Get()
	defer d.Close()
	data, _ := redis.String(d.Do(models.REDIS_HGET,models.GROUP_LIST,group_id))
	return data
}

//获取指uid对应的群id数据字符
func GetUidToGroupDataRedis(uid string) string {
	d := common.RedisConn.Get()
	defer d.Close()
	data, _ := redis.String(d.Do(models.REDIS_HGET,models.USER_TO_GROUPS,uid))
	return data
}