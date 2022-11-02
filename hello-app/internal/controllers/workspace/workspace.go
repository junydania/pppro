package workspace

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/entities/workspace"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/repository/adapter"
	"time"
)

type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	ListOne(ID uuid.UUID) (entity workspace.Workspace, err error)
	ListAll() (entities []workspace.Workspace, err error)
	Create(entity *workspace.Workspace) (uuid.UUID, error)
	Update(ID uuid.UUID, entity *workspace.Workspace) error
	Remove(ID uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) ListOne(id uuid.UUID) (entity workspace.Workspace, err error) {
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return workspace.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []workspace.Workspace, err error) {
	entities = []workspace.Workspace{}
	var entity workspace.Workspace

	filter := expression.Name("name").NotEqual(expression.Value(""))
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
			entity, err := workspace.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	
	return entities, nil
}

func (c *Controller) Create(entity *workspace.Workspace) (uuid.UUID, error) {
	entity.CreatedAt = time.Now()
	entity.DeleteDate = time.Now().Add(time.Hour * 24)
	entity.UptimeHours = 24
	
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}

func (c *Controller) Update(id uuid.UUID, entity *workspace.Workspace) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}
	found.ID = id
	found.Name = entity.Name
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