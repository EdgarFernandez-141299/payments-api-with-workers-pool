package worfkflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	errors2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestExecuteOrderWorkflow(t *testing.T) {
	orderID := "ref-342"

	t.Run("order created successfully and execute payment orders pipelines", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		defer ts.AfterTest()

		checkOrderExistenceActivity := new(activities.CheckOrderActivity)

		ts.Env.RegisterActivity(checkOrderExistenceActivity)
		ts.Env.RegisterWorkflow(ProcessOrderPaymentWorkflow)
		ts.Env.RegisterWorkflow(CaptureFlowWorkflow)

		ts.Env.OnActivity(
			checkOrderExistenceActivity.CheckOrderExistence,
			mock.Anything,
			orderID,
		).
			Return(true, nil)

		ts.Env.ExecuteWorkflow(ProcessOrderWorkflow, PaymentWorkflowInput{
			OrderID: orderID,
		})

		ts.NoError(ts.Env.GetWorkflowError())
	})

	t.Run("order not exists and execute payment orders pipelines", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)

		defer ts.AfterTest()

		checkOrderActivity := new(activities.CheckOrderActivity)

		ts.Env.RegisterActivity(checkOrderActivity)
		ts.Env.RegisterWorkflow(ProcessOrderPaymentWorkflow)
		ts.Env.RegisterWorkflow(CaptureFlowWorkflow)

		ts.Env.OnActivity(
			checkOrderActivity.CheckOrderExistence,
			mock.Anything,
			orderID,
		).
			Return(false, nil)

		ts.Env.ExecuteWorkflow(ProcessOrderWorkflow, PaymentWorkflowInput{
			OrderID: orderID,
			PaymentCommands: []command.CreatePaymentOrderCommand{
				{
					Payment: entities.PaymentOrder{ID: "red-234"},
				},
			},
		})

		gotErr := ts.Env.GetWorkflowError()

		ts.Error(gotErr)
		assert.True(t, errors.IsBusinessErrorWithCode(gotErr, errors2.OrderNotFoundErrorCode))
	})

	t.Run("order check fails and not execute payment orders pipelines", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)

		defer ts.AfterTest()

		checkOrderActivity := new(activities.CheckOrderActivity)

		ts.Env.RegisterActivity(checkOrderActivity)
		ts.Env.RegisterWorkflow(ProcessOrderPaymentWorkflow)
		ts.Env.RegisterWorkflow(CaptureFlowWorkflow)

		ts.Env.OnActivity(
			checkOrderActivity.CheckOrderExistence,
			mock.Anything,
			orderID,
		).
			Return(false, assert.AnError)

		ts.Env.ExecuteWorkflow(ProcessOrderWorkflow, PaymentWorkflowInput{
			OrderID: orderID,
		})

		ts.NoError(ts.Env.GetWorkflowError())
	})

	t.Run("order created successfully and execute payment orders pipelines and execute ProcessOrderPaymentWorkflow and fails", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)

		defer ts.AfterTest()

		createOrderActivity := new(activities.CheckOrderActivity)

		ts.Env.RegisterActivity(createOrderActivity)
		ts.Env.RegisterWorkflow(ProcessOrderPaymentWorkflow)
		ts.Env.RegisterWorkflow(CaptureFlowWorkflow)

		ts.Env.OnActivity(
			createOrderActivity.CheckOrderExistence,
			mock.Anything,
			orderID,
		).
			Return(true, nil)

		ts.Env.OnWorkflow(
			ProcessOrderPaymentWorkflow, mock.Anything, mock.IsType(*new(command.CreatePaymentOrderCommand)),
		).Return(assert.AnError).Once()

		ts.Env.ExecuteWorkflow(ProcessOrderWorkflow, PaymentWorkflowInput{
			OrderID: orderID,
			PaymentCommands: []command.CreatePaymentOrderCommand{
				{
					ID:      "1",
					Payment: entities.PaymentOrder{ID: "21"},
				},
				{
					ID:      "2",
					Payment: entities.PaymentOrder{ID: "34"},
				},
			},
		})

		got := &PaymentOrderWorkflowOut{}
		gotErr := ts.Env.GetWorkflowResult(got)

		expectedResponse := &PaymentOrderWorkflowOut{
			ReferenceOrderID: orderID,
			Payments: []PaymentStatus{
				NewFailedPaymentStatus("21", ""),
				NewNotProcessedStatus("34"),
			},
		}

		ts.NoError(gotErr)
		assert.Equal(t, expectedResponse, got)
	})

	t.Run("order created successfully and execute payment orders pipelines and execute ProcessOrderPaymentWorkflow", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)

		defer ts.AfterTest()

		createOrderActivity := new(activities.CheckOrderActivity)

		ts.Env.RegisterActivity(createOrderActivity)
		ts.Env.RegisterWorkflow(ProcessOrderPaymentWorkflow)
		ts.Env.RegisterWorkflow(CaptureFlowWorkflow)

		ts.Env.OnActivity(
			createOrderActivity.CheckOrderExistence,
			mock.Anything,
			orderID,
		).
			Return(true, nil)

		ts.Env.OnWorkflow(
			ProcessOrderPaymentWorkflow, mock.Anything, mock.IsType(*new(command.CreatePaymentOrderCommand)),
		).Return(nil).Maybe()

		ts.Env.ExecuteWorkflow(ProcessOrderWorkflow, PaymentWorkflowInput{
			OrderID: orderID,
			PaymentCommands: []command.CreatePaymentOrderCommand{
				{
					ID:      "1",
					Payment: entities.PaymentOrder{ID: "21"},
				},
				{
					ID:      "2",
					Payment: entities.PaymentOrder{ID: "34"},
				},
			},
		})

		got := &PaymentOrderWorkflowOut{}
		gotErr := ts.Env.GetWorkflowResult(got)

		expectedResponse := &PaymentOrderWorkflowOut{
			ReferenceOrderID: orderID,
			Payments: []PaymentStatus{
				NewProcessingPaymentStatus("21"),
				NewProcessingPaymentStatus("34"),
			},
		}

		ts.NoError(gotErr)
		assert.Equal(t, expectedResponse, got)
	})
}
