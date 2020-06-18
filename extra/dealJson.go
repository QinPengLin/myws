package extra

import (
	"encoding/json"
	"myws/models"
)

func JsonToFromMessageStruct(jsonStr string) (*models.FromMessage,error)  {
	var s models.FromMessage

	err:=json.Unmarshal([]byte(jsonStr), &s)
	if err!=nil{
		return &s,err
	}
	return &s,err
}

func StructToFromMessageJson(structs *models.SendMessage) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}

func StructToSendMessageAbnormalJson(structs *models.SendMessageAbnormal) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}

func StructToSendMessageAndDataAbnormalJson(structs *models.SendMessageAndDataAbnormal) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}