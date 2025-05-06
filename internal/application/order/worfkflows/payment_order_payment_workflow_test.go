package worfkflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

func TestProcessOrderPaymentWorkflow(t *testing.T) {
	t.Run("must success when executes post processing workflow and call create order activity", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "23442-45sct3-5dt43"
		paymentID := "16363-f4545-1234"

		defer ts.AfterTest()

		createPaymentOrderActivity := new(activities.CreatePaymentOrderActivity)

		env.RegisterActivity(createPaymentOrderActivity)
		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		env.OnActivity(
			createPaymentOrderActivity.CreatePaymentOrder,
			mock.Anything,
			mock.IsType(*new(command.CreatePaymentOrderCommand))).
			Return(response.PaymentOrderResponseDTO{}, nil).Maybe()

		cmd := command.CreatePaymentOrderCommand{
			ReferenceOrderID: orderID,
			Payment:          entities.PaymentOrder{ID: paymentID},
		}

		ts.Env.OnWorkflow(
			PostProcessingPaymentWorkflow, mock.Anything, orderID, paymentID,
		).Return(nil).Once()

		env.ExecuteWorkflow(ProcessOrderPaymentWorkflow, cmd)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("must fail when create activity fails", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "23442-45sct3-5dt43"
		paymentID := "16363-f4545-1234"

		defer ts.AfterTest()

		createPaymentOrderActivity := new(activities.CreatePaymentOrderActivity)

		env.RegisterActivity(createPaymentOrderActivity)
		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		env.OnActivity(
			createPaymentOrderActivity.CreatePaymentOrder,
			mock.Anything,
			mock.IsType(*new(command.CreatePaymentOrderCommand))).
			Return(response.PaymentOrderResponseDTO{}, assert.AnError).Once()

		cmd := command.CreatePaymentOrderCommand{
			ReferenceOrderID: orderID,
			Payment:          entities.PaymentOrder{ID: paymentID},
		}

		env.ExecuteWorkflow(ProcessOrderPaymentWorkflow, cmd)

		ts.Error(env.GetWorkflowError())
	})
}
