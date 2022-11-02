package workspace

import (
	"encoding/json"
	"errors"
	"time"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"fmt"
)

type Workspace struct {
	entities.Base
	Name string `dynamodbav:"name" json:"name"`
	AccountName string `dynamodbav:"account_name" json:"account_name"`
	Region string `dynamodbav:"region" json:"region"`
	StackName string `dynamodbav:"stack_name" json:"stack_name"`
	UptimeHours int `dynamodbav:"uptime_hours" json:"uptime_hours"`	
	DeleteDate time.Time `json:"delete_date"`
	WorkspaceDetails struct {
		InstanceId string `dynamodbav:"instance_id" json:"instance_id"`
		InstanceIp string  `dynamodbav:"instance_ip" json:"instance_ip"`
		Email string `dynamodbav:"email" json:"email"`
		Username string `dynamodbav:"username" json:"username"`
		VpcId string `dynamodbav:"vpc_id" json:"vpc_id"`
		SecurityGroups []string `dynamodbav:"security_groups" json:"security_groups"`
		SubnetIds []string `dynamodbav:"subnet_ids" json:"subnet_ids"`
	} `dynamodbav:"workspace_details" json:"workspace_details"`
}

func InterfaceToModel(data interface{}) (instance *Workspace, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}
	return instance, json.Unmarshal(bytes, &instance)
}

func (w *Workspace) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": w.ID.String()}
}

func (w *Workspace) TableName() string {
	return "Workspaces"
}

func (w *Workspace) Bytes() ([]byte, error) {
	return json.Marshal(w)
}

func (w *Workspace) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       w.ID.String(),
		"name":      w.Name,
		"account_name": w.AccountName,
		"region": w.Region,
		"stack_name": w.StackName,
		"uptime_hours": w.UptimeHours,
		"delete_date": w.DeleteDate.Format(entities.GetTimeFormat()),
		"created_at": w.CreatedAt.Format(entities.GetTimeFormat()),
		"updated_at": w.UpdatedAt.Format(entities.GetTimeFormat()),
		"workspace_details": map[string]interface{}{
			"instance_id": w.WorkspaceDetails.InstanceId, 
			"instance_ip": w.WorkspaceDetails.InstanceIp,
			"email": w.WorkspaceDetails.Email,
			"username": w.WorkspaceDetails.Username,
			"vpc_id": w.WorkspaceDetails.VpcId,
			"security_groups": w.WorkspaceDetails.SecurityGroups,
			"subnet_ids": w.WorkspaceDetails.SubnetIds,
		},
	}
}

func ParseDynamoAtributeToStruct(response map[string] types.AttributeValue) (w Workspace, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return w, errors.New("item not found")
	}

	for key, value := range response {
		if key == "_id" {
			tv, _ := value.(*types.AttributeValueMemberS)
			w.ID, err = uuid.Parse(tv.Value)
			if w.ID == uuid.Nil {
				err = errors.New("item not found")
			}
		}

		if key == "delete_date" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			w.DeleteDate = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}

		if key == "created_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			w.CreatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}

		if key == "updated_at" {
			tv, _ := value.(*types.AttributeValueMemberS)
			timeFormat := entities.GetTimeFormat()
			timeParsed, error := time.Parse(timeFormat, tv.Value)
			w.UpdatedAt = timeParsed
			if error != nil {
				fmt.Println(error)
				return
			}
		}
	}
	attributevalue.UnmarshalMap(response, &w)
	return w, nil
}