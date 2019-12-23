package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
)

var HostPort = "127.0.0.1:7933"
var Domain = "samples-domain"
var TaskListName = "helloworld-worker"
var ClientName = "helloworld-worker"
var CadenceService = "cadence-frontend"

func init() {
	workflow.Register(HelloWorldWorkflow)
}

var logger *zap.Logger

func main() {

	logger, _ = zap.NewDevelopment()

	logger.Info("Starting Workflow")
	startWorkflow()
}

func startWorkflow() {
	logger.Info("Start workflow method")
	workflowOptions := client.StartWorkflowOptions{
		ID:                              fmt.Sprintf("helloworld_%s", uuid.New()),
		TaskList:                        "helloworld-worker",
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	serviceClient := buildCadenceClient()
	workflowClient := client.NewClient(serviceClient, Domain, nil)
	domainClient := client.NewDomainClient(serviceClient, nil)
	_, err := domainClient.Describe(context.Background(), Domain)
	if err != nil {
		logger.Info("Domain doesn't exist", zap.String("Domain", Domain), zap.Error(err))
	} else {
		logger.Info("Domain successfully registered.", zap.String("Domain", Domain))
	}
	workflowClient.StartWorkflow(context.Background(), workflowOptions, HelloWorldWorkflow, "World ! Hrishi Here.... ")
}

func buildCadenceClient() workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService))
}
