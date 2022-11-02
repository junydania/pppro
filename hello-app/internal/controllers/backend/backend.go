package backend

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities/backend"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/repository/adapter"
	"time"
)
type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	ListOne(ID uuid.UUID) (entity backend.Backend, err error)
	ListAll() (entities []backend.Backend, err error)
	Create(entity *backend.Backend) (uuid.UUID, error)
	Update(ID uuid.UUID, entity *backend.Backend) error
	Remove(ID uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) ListOne(id uuid.UUID) (entity backend.Backend, err error) {
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return backend.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []backend.Backend, err error) {
	entities = []backend.Backend{}
	var entity backend.Backend

	filter := expression.Name("alias").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := backend.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (c *Controller) Create(entity *backend.Backend) (uuid.UUID, error) {
	entity.CreatedAt = time.Now()
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}

func (c *Controller) Update(id uuid.UUID, entity *backend.Backend) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}
	found.ID = id
	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), entity.TableName())
	return err
}

func (c *Controller) Remove(id uuid.UUID) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}
	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}