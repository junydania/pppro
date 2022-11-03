package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/junydania/pppro/hello-app/config"
	"github.com/junydania/pppro/hello-app/internal/repository/adapter"
	"github.com/junydania/pppro/hello-app/internal/repository/instance"
	"github.com/junydania/pppro/hello-app/internal/routes"
	"github.com/junydania/pppro/hello-app/internal/rules"
	RulesHello "github.com/junydania/pppro/hello-app/internal/rules/hello"
	"github.com/junydania/pppro/hello-app/utils/logger"
	"log"
	"net/http"
)

func main() {
	configs := config.GetConfig()
	
	connection := instance.GetConnection()
	repository := adapter.NewAdapter(connection)

	logger.INFO("Waiting service starting.... ", nil)

	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error on migrate: ", err)
		}
	}
	logger.PANIC("", checkTables(connection))

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("Service running on port ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.Client) []error {
	var errors []error
	callMigrateAndAppendError(&errors, connection, &RulesHello.Rules{})
	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.Client, rule rules.Interface) {
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.Client) error {
	response, err := connection.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found: ", nil)
		}
		for _, tableName := range response.TableNames {
			logger.INFO("Table found: ", tableName)
		}
	}
	return err
}