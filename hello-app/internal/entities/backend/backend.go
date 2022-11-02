package backend

import (
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities"
	"fmt"
)

type Backend struct {
	entities.Base
	Bucket string `dynamodbav:"bucket" json:"bucket"`
	Key string `dynamodbav:"key" json:"key"`
	Region string `dynamodbav:"region" json:"region"`
	LockTableName string `dynamodbav:"lock_table" json:"lock_table"`
	RoleArn string `dynamodbav:"role_arn" json:"role_arn"`
	SessionName string `dynamodbav:"session_name" json:"session_name"`
}

func InterfaceToModel(data interface{}) (backend *Backend, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return backend, err
	}
	return backend, json.Unmarshal(bytes, &backend)
}

func (b *Backend) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": b.ID.String()}
}

func (b *Backend) TableName() string {
	return "Backend"
}

func (b *Backend) Bytes() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Backend) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       b.ID.String(),
		"bucket": b.Bucket,
		"region":      b.Region,
		"role_arn": b.RoleArn,
		"session_name": b.SessionName,
		"created_at": b.CreatedAt.Format(entities.GetTimeFormat()),
		"updated_at": b.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}

func ParseDynamoAtributeToStruct(response map[string] types.AttributeValue) (b Backend, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return b, errors.New("item not found")
	}

	for key, value := range response {
		if key == "_id" {
			tv, _ := value.(*types.AttributeValueMemberS)
			b.ID, err = uuid.Parse(tv.Value)
			if b.ID == uuid.Nil {
				err = errors.New("item not found")
			}
		}

		if key == "created_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			b.CreatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}

		if key == "updated_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			b.UpdatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}
	}
	attributevalue.UnmarshalMap(response, &b)
	return b, nil
}