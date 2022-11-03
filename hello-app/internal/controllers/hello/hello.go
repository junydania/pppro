package hello
import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/junydania/pppro/hello-app/internal/entities/hello"
	"github.com/junydania/pppro/hello-app/internal/repository/adapter"
	"time"
)

type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	ListOne(ID uuid.UUID) (entity hello.Hello, err error)
	ListAll() (entities []hello.Hello, err error)
	Create(entity *hello.Hello) (uuid.UUID, error)
	Remove(ID uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) ListOne(id uuid.UUID) (entity hello.Hello, err error) {
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return hello.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []hello.Hello, err error) {
	entities = []hello.Hello{}
	var entity hello.Hello

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
			entity, err := hello.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

func (c *Controller) Create(entity *hello.Hello) (uuid.UUID, error) {
	entity.CreatedAt = time.Now()

	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}


func (c *Controller) Remove(id uuid.UUID) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}
	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}