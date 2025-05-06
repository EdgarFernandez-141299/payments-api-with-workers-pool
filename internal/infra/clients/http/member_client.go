package http

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	commonResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters/resources"
	response "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

const (
	enterpriseHeader       = "x-enterprise-id"
	userIDHeader           = "x-identity-user-id"
	usernameHeader         = "x-username"
	MemberApiTimeout       = 30 * time.Second
	getMemberByIDPath      = "/api/v1/members/%s"
	getUserProfileInfoPath = "/api/v1/members/profile-info"
)

type MemberHTTPClientImpl struct {
	instrument.Client
}

func NewMemberHTTPClient(tracer apm.Tracer) commonResources.MemberAPIResourceIF {
	return &MemberHTTPClientImpl{
		instrument.NewInstrumentedClient(
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithBaseUrl(config.Config().MembersApi.URL),
			instrument.WithRequestTimeout(MemberApiTimeout),
		),
	}
}

func (m *MemberHTTPClientImpl) GetMemberByID(oldCtx context.Context, id, enterpriseId string) (response.MemberDTO, error) {
	return decorators.TraceDecorator(
		oldCtx,
		"MemberHTTPClient.GetMemberById",
		func(ctx context.Context, decorators decorators.Span) (response.MemberDTO, error) {
			var responseBody response.MemberResponse[response.MemberDTO]
			_, err := m.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true),
			).
				SetContext(ctx).
				SetHeader(enterpriseHeader, enterpriseId).
				SetHeader(userIDHeader, id).
				SetHeader(usernameHeader, id).
				SetResult(&responseBody).
				Get(fmt.Sprintf(getMemberByIDPath, id))

			if responseBody.Data.ID != id {
				return response.MemberDTO{}, fmt.Errorf("member not found")
			}

			return responseBody.Data, err
		},
	)
}

func (m *MemberHTTPClientImpl) GetUserProfileInfo(oldCtx context.Context, userId string, enterpriseId string) (response.UserProfileInfoDTO, error) {
	return decorators.TraceDecorator(
		oldCtx,
		"MemberHTTPClient.GetUserProfileInfo",
		func(ctx context.Context, decorators decorators.Span) (response.UserProfileInfoDTO, error) {
			var responseBody response.UserProfileInfoDTO
			_, err := m.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true),
			).
				SetContext(ctx).
				SetResult(&responseBody).
				SetHeader(enterpriseHeader, enterpriseId).
				SetHeader(userIDHeader, userId).
				SetQueryParam("userID", userId).
				Get(getUserProfileInfoPath)

			if responseBody.UserID != userId {
				return response.UserProfileInfoDTO{}, fmt.Errorf("member not found")
			}

			return responseBody, err
		},
	)
}
