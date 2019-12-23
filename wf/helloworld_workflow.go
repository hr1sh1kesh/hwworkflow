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
		TaskList:               "HelloWorldWorker",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger.Info("helloworld workflow started")

	future := workflow.ExecuteActivity(ctx, "HelloworldActivity", value)
	var result string
	if err := future.Get(ctx, &result); err != nil {
		return err
	}

	logger.Info("Workflow completed.", zap.String("Result", result))

	return nil
}
