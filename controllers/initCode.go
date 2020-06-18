package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"encoding/json"
)

func initStatusCode(){
	var statusCode  map[string]map[string]string
	codeConf:=beego.AppConfig.String("codeConf")
	err := json.Unmarshal([]byte(codeConf), &statusCode)
    if err != nil {
        logs.Info("初始化状态码失败！请检查，程序将无法正常运行")
	}
	successValueCode,successOkCode := statusCode["success"]["code"]
	if successOkCode {
		success["code"]=successValueCode
	}else{
		success["code"]="c_1"
	}
	successValueMsg,successOkMsg := statusCode["success"]["msg"]
	if successOkMsg {
		success["msg"]=successValueMsg
	}else{
		success["code"]="OK"
	}

	moranoValueCode,moranoOkCode := statusCode["morano"]["code"]
	if moranoOkCode {
		morano["code"]=moranoValueCode
	}else{
		morano["code"]="c_1001"
	}
	moranoValueMsg,moranoOkMsg := statusCode["morano"]["msg"]
	if moranoOkMsg {
		morano["msg"]=moranoValueMsg
	}else{
		morano["code"]="无需修改"
	}

	morazyValueCode,morazyOkCode := statusCode["morazy"]["code"]
	if morazyOkCode {
		morazy["code"]=morazyValueCode
	}else{
		morazy["code"]="c_1002"
	}
	morazyValueMsg,morazyOkMsg := statusCode["morazy"]["msg"]
	if morazyOkMsg {
		morazy["msg"]=morazyValueMsg
	}else{
		morazy["code"]="uid已经被占用"
	}

	morayesValueCode,morayesOkCode := statusCode["morayes"]["code"]
	if morayesOkCode {
		morayes["code"]=morayesValueCode
	}else{
		morayes["code"]="c_1003"
	}
	morayesValueMsg,morayesOkMsg := statusCode["morayes"]["msg"]
	if morayesOkMsg {
		morayes["msg"]=morayesValueMsg
	}else{
		morayes["code"]="uid修改或者新增成功"
	}

	usernoValueCode,usernoOkCode := statusCode["userno"]["code"]
	if usernoOkCode {
		userno["code"]=usernoValueCode
	}else{
		userno["code"]="c_102"
	}
	usernoValueMsg,usernoOkMsg := statusCode["userno"]["msg"]
	if usernoOkMsg {
		userno["msg"]=usernoValueMsg
	}else{
		userno["code"]="玩家不存在或者离线了"
	}

	useryesValueCode,useryesOkCode := statusCode["useryes"]["code"]
	if useryesOkCode {
		useryes["code"]=useryesValueCode
	}else{
		useryes["code"]="c_103"
	}
	useryesValueMsg,useryesOkMsg := statusCode["useryes"]["msg"]
	if useryesOkMsg {
		useryes["msg"]=useryesValueMsg
	}else{
		useryes["code"]="发送给了玩家"
	}

	verifyValueCode,verifyOkCode := statusCode["verify"]["code"]
	if verifyOkCode {
		verify["code"]=verifyValueCode
	}else{
		verify["code"]="c_104"
	}
	verifyValueMsg,verifyOkMsg := statusCode["verify"]["msg"]
	if verifyOkMsg {
		verify["msg"]=verifyValueMsg
	}else{
		verify["code"]="你还未设置uid"
	}

	illegalValueCode,illegalOkCode := statusCode["illegal"]["code"]
	if illegalOkCode {
		illegal["code"]=illegalValueCode
	}else{
		illegal["code"]="c_105"
	}
	illegalValueMsg,illegalOkMsg := statusCode["illegal"]["msg"]
	if illegalOkMsg {
		illegal["msg"]=illegalValueMsg
	}else{
		illegal["code"]="uid不合法"
	}

	gnerrValueCode,gnerrOkCode := statusCode["gnerr"]["code"]
	if gnerrOkCode {
		gnerr["code"]=gnerrValueCode
	}else{
		gnerr["code"]="c_106"
	}
	gnerrValueMsg,gnerrOkMsg := statusCode["gnerr"]["msg"]
	if gnerrOkMsg {
		gnerr["msg"]=gnerrValueMsg
	}else{
		gnerr["code"]="群名不能为空"
	}

	lenerrValueCode,lenerrOkCode := statusCode["lenerr"]["code"]
	if lenerrOkCode {
		lenerr["code"]=lenerrValueCode
	}else{
		lenerr["code"]="c_107"
	}
	lenerrValueMsg,lenerrOkMsg := statusCode["lenerr"]["msg"]
	if lenerrOkMsg {
		lenerr["msg"]=lenerrValueMsg
	}else{
		lenerr["code"]="群名人数超过设定"
	}

	gadderrValueCode,gadderrOkCode := statusCode["gadderr"]["code"]
	if gadderrOkCode {
		gadderr["code"]=gadderrValueCode
	}else{
		gadderr["code"]="c_108"
	}
	gadderrValueMsg,gadderrOkMsg := statusCode["gadderr"]["msg"]
	if gadderrOkMsg {
		gadderr["msg"]=gadderrValueMsg
	}else{
		gadderr["code"]="新增群失败"
	}

	grouperrValueCode,grouperrOkCode := statusCode["grouperr"]["code"]
	if grouperrOkCode {
		grouperr["code"]=grouperrValueCode
	}else{
		grouperr["code"]="c_109"
	}
	grouperrValueMsg,grouperrOkMsg := statusCode["grouperr"]["msg"]
	if grouperrOkMsg {
		grouperr["msg"]=grouperrValueMsg
	}else{
		grouperr["code"]="群id错误"
	}

	noauthValueCode,noauthOkCode := statusCode["noauth"]["code"]
	if noauthOkCode {
		noauth["code"]=noauthValueCode
	}else{
		noauth["code"]="c_110"
	}
	noauthValueMsg,noauthOkMsg := statusCode["noauth"]["msg"]
	if noauthOkMsg {
		noauth["msg"]=noauthValueMsg
	}else{
		noauth["code"]="无权限"
	}

	notomeValueCode,notomeOkCode := statusCode["notome"]["code"]
	if notomeOkCode {
		notome["code"]=notomeValueCode
	}else{
		notome["code"]="c_111"
	}
	notomeValueMsg,notomeOkMsg := statusCode["notome"]["msg"]
	if notomeOkMsg {
		notome["msg"]=notomeValueMsg
	}else{
		notome["code"]="不能给自己发送消息"
	}
}