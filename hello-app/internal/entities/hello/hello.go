package hello

import (
	"encoding/json"
	"errors"
	"time"
	"github.com/junydania/pppro/hello-app/internal/entities"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"fmt"
)

type Hello struct {
	entities.Base
	Content string `dynamodbav:"content" json:"content"`
}

func InterfaceToModel(data interface{}) (instance *Hello, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}
	return instance, json.Unmarshal(bytes, &instance)
}

func (h *Hello) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": h.ID.String()}
}

func (h *Hello) TableName() string {
	return "Hellos"
}

func (h *Hello) Bytes() ([]byte, error) {
	return json.Marshal(h)
}

func (h *Hello) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       h.ID.String(),
		"content":   h.Content,
	}
}

func ParseDynamoAtributeToStruct(response map[string] types.AttributeValue) (h Hello, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return h, errors.New("item not found")
	}

	for key, value := range response {
		if key == "_id" {
			tv, _ := value.(*types.AttributeValueMemberS)
			h.ID, err = uuid.Parse(tv.Value)
			if h.ID == uuid.Nil {
				err = errors.New("item not found")
			}
		}

		if key == "created_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			h.CreatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}

		if key == "updated_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			h.UpdatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}
	}
	attributevalue.UnmarshalMap(response, &h)
	return h, nil
}