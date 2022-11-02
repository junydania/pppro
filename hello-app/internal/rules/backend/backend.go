package backend

import (
	"encoding/json"
	"errors"
	"io"
	"time"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities/backend"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	Validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
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

	table := &backend.Backend{}

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
			ReadCapacityUnits:  aws.Int64(4),
			WriteCapacityUnits: aws.Int64(4),
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
	return backend.Backend{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Bucket: uuid.New().String(),
	}
}

func (r *Rules) Validate(model interface{}) error {
	backendModel, err := backend.InterfaceToModel(model)
	if err != nil {
		return err
	}

	return Validation.ValidateStruct(backendModel,
		Validation.Field(&backendModel.ID, Validation.Required, is.UUIDv4),

	)
}