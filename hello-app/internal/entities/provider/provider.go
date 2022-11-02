package provider

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

type Provider struct {
	entities.Base
	Alias string `dynamodbav:"alias" json:"alias"`
	RoleArn string `dynamodbav:"role_arn" json:"role_arn"`
	SessionName string `dynamodbav:"session_name" json:"session_name"`
}

func InterfaceToModel(data interface{}) (provider *Provider, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return provider, err
	}

	return provider, json.Unmarshal(bytes, &provider)
}

func (p *Provider) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": p.ID.String()}
}

func (p *Provider) GetFilterAlias() map[string]interface{} {
	return map[string]interface{}{"name": p.Alias}
}

func (p *Provider) TableName() string {
	return "Provider"
}

func (p *Provider) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Provider) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       p.ID.String(),
		"alias":      p.Alias,
		"role_arn": p.RoleArn,
		"session_name": p.SessionName,
		"created_at": p.CreatedAt.Format(entities.GetTimeFormat()),
		"updated_at": p.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}


func ParseDynamoAtributeToStruct(response map[string] types.AttributeValue) (p Provider, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("item not found")
	}

	for key, value := range response {
		if key == "_id" {
			tv, _ := value.(*types.AttributeValueMemberS)
			p.ID, err = uuid.Parse(tv.Value)
			if p.ID == uuid.Nil {
				err = errors.New("item not found")
			}
		}

		if key == "created_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			p.CreatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}

		if key == "updated_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			p.UpdatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}
	}
	attributevalue.UnmarshalMap(response, &p)
	return p, nil
}