package workspace

import (
	"encoding/json"
	"errors"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	Validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities/workspace"
	"io"
	"time"
)

type Rules struct{}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}

func (r *Rules) Migrate(connection *dynamodb.Client) error {
	return r.createTable(connection)
}

func (r *Rules) createTable(connection *dynamodb.Client) error {

	table := &workspace.Workspace{}

	_, err := connection.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(table.TableName()),
	})

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},

		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table.TableName()),
	}

	if err != nil {
		var rne *types.ResourceNotFoundException
		if errors.As(err, &rne) {
			_, err = connection.CreateTable(context.TODO(), input)
			if err != nil {
				return err
			}		
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Rules) GetMock() interface{} {
	return workspace.Workspace{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: uuid.New().String(),
	}
}

func (r *Rules) Validate(model interface{}) error {
	workspaceModel, err := workspace.InterfaceToModel(model)
	if err != nil {
		return err
	}
	
	return Validation.ValidateStruct(workspaceModel,
		Validation.Field(&workspaceModel.ID, Validation.Required, is.UUIDv4),
		Validation.Field(&workspaceModel.Name, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&workspaceModel.AccountName, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&workspaceModel.Region, Validation.Required, Validation.Length(3, 20)),
		Validation.Field(&workspaceModel.StackName, Validation.Required, Validation.Length(3, 50)),
	)
}