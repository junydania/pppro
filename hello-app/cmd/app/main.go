package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"bitbucket.org/codapayments/coda-stack-management-app/config"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/repository/adapter"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/repository/instance"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/routes"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/rules"
	RulesWorkspace "bitbucket.org/codapayments/coda-stack-management-app/internal/rules/workspace"
	RulesProvider "bitbucket.org/codapayments/coda-stack-management-app/internal/rules/provider"
	RulesBackend "bitbucket.org/codapayments/coda-stack-management-app/internal/rules/backend"
	"bitbucket.org/codapayments/coda-stack-management-app/utils/logger"
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
	callMigrateAndAppendError(&errors, connection, &RulesWorkspace.Rules{})
	callMigrateAndAppendError(&errors, connection, &RulesProvider.Rules{})
	callMigrateAndAppendError(&errors, connection, &RulesBackend.Rules{})
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