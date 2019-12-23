package wf

import (
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func HelloWorldWorkflow(ctx workflow.Context, value string) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ao := workflow.ActivityOptions{
		TaskList:               "helloworld-worker",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger.Info("helloworld workflow started")

	var helloworldResult string
	err = workflow.ExecuteActivity(ctx, "main.HelloworldActivity", "World ! Hrishi here.....").Get(ctx, &helloworldResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	logger.Info("Workflow completed.", zap.String("Result", helloworldResult))

	return nil
}
