package http

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/utils"
)

type DeUnaAuthHTTPClient struct {
	Config deunaConfig.DeUnaApiConfig
	instrument.Client
}

func NewDeUnaAuthHTTPClient(config *deunaConfig.DeUnaApiConfig, tracer apm.Tracer) resources.DeunaAuthResourceIF {
	headerApiKey := deunaConfig.DeunaApiKeyHeader
	return &DeUnaAuthHTTPClient{
		Config: *config,
		Client: instrument.NewInstrumentedClient(
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithBaseUrl(config.URL),
			instrument.WithRequestTimeout(deUnaTimeout),
			instrument.WithHeaders(map[string]string{
				headerApiKey:   config.ApiKey,
				"Content-Type": "application/json",
			}),
		),
	}
}

// para integrar auth user
func (d DeUnaAuthHTTPClient) AuthUser(
	ctx context.Context,
	request request.DeunaAuthUserRequestDTO,
) (response.DeunaAuthResponseDTO, error) {
	publicKeyPEM := utils.Base64Decode(d.Config.PublicKeyPEM)
	authMessage, err := generateAuthMessage(request.Email, publicKeyPEM)

	if err != nil {
		return response.DeunaAuthResponseDTO{}, fmt.Errorf("error generating auth message: %w", err)
	}

	var res response.DeunaAuthResponseDTO

	httpResponse, err := d.Client.NewRequestWithOptions(
		instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
	).
		SetResult(&res).
		SetBody(request).
		SetContext(ctx).
		SetHeader("X-Auth-Message", authMessage).
		SetQueryParam("load_profile_data", "true").
		Post("/users/external-authorize")

	if err != nil {
		return response.DeunaAuthResponseDTO{}, fmt.Errorf("failed to authenticate user: %w", err)
	}

	if httpResponse.StatusCode() == http.StatusNotFound {
		return response.DeunaAuthResponseDTO{}, fmt.Errorf("user not found with status: %d", httpResponse.StatusCode())
	}

	if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
		return response.DeunaAuthResponseDTO{}, fmt.Errorf("authenticate user failed with status: %d", httpResponse.StatusCode())
	}

	return res, nil
}

func generateAuthMessage(email, publicKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse DER encoded public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("key type is not RSA")
	}

	data := []byte(fmt.Sprintf(`{"email": "%s"}`, email))
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, data, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}
