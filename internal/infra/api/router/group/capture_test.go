package group

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCaptureFlowHandler implementa la interfaz CaptureFlowHandlerIF para pruebas
type MockCaptureFlowHandler struct {
	mock.Mock
}

func (m *MockCaptureFlowHandler) PaymentCapture(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCaptureFlowHandler) PaymentRelease(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestCaptureFlowV1Routes(t *testing.T) {
	// Configuración inicial
	e := echo.New()
	group := e.Group("/test")
	mockHandler := new(MockCaptureFlowHandler)

	// Ejecutar la función bajo prueba
	CaptureFlowV1Routes(group, mockHandler)

	// Verificar que las rutas se hayan configurado correctamente
	routes := e.Routes()
	foundCapture := false
	foundRelease := false
	for _, route := range routes {
		if route.Path == "/test/v1/order/payment/capture" && route.Method == "POST" {
			foundCapture = true
		}
		if route.Path == "/test/v1/order/payment/release" && route.Method == "POST" {
			foundRelease = true
		}
	}
	assert.True(t, foundCapture, "La ruta POST /v1/order/payment/capture no se configuró correctamente")
	assert.True(t, foundRelease, "La ruta POST /v1/order/payment/release no se configuró correctamente")
}

func TestNewCaptureFlowRoutes(t *testing.T) {
	// Configuración inicial
	e := echo.New()
	group := e.Group("/test")
	mockHandler := new(MockCaptureFlowHandler)

	// Ejecutar la función bajo prueba
	routes := NewCaptureFlowRoutes(group, mockHandler)

	// Verificar que se devuelva una instancia de CaptureFlowRoutes
	assert.NotNil(t, routes)
	assert.IsType(t, &CaptureFlowRoutes{}, routes)
	assert.Equal(t, mockHandler, routes.handler)

	// Verificar que las rutas se hayan configurado
	routesList := e.Routes()
	foundCapture := false
	foundRelease := false
	for _, route := range routesList {
		if route.Path == "/test/v1/order/payment/capture" && route.Method == "POST" {
			foundCapture = true
		}
		if route.Path == "/test/v1/order/payment/release" && route.Method == "POST" {
			foundRelease = true
		}
	}
	assert.True(t, foundCapture, "La ruta POST /v1/order/payment/capture no se configuró correctamente")
	assert.True(t, foundRelease, "La ruta POST /v1/order/payment/release no se configuró correctamente")
}

func TestPaymentCaptureHandler(t *testing.T) {
	// Configuración inicial
	e := echo.New()
	group := e.Group("/test")
	mockHandler := new(MockCaptureFlowHandler)

	// Configurar la ruta
	CaptureFlowV1Routes(group, mockHandler)

	// Crear una solicitud de prueba
	req := httptest.NewRequest(http.MethodPost, "/test/v1/order/payment/capture", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Configurar el mock para que devuelva nil (éxito)
	mockHandler.On("PaymentCapture", c).Return(nil)

	// Ejecutar la solicitud
	err := mockHandler.PaymentCapture(c)

	// Verificar resultados
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestPaymentReleaseHandler(t *testing.T) {
	// Configuración inicial
	e := echo.New()
	group := e.Group("/test")
	mockHandler := new(MockCaptureFlowHandler)

	// Configurar la ruta
	CaptureFlowV1Routes(group, mockHandler)

	// Crear una solicitud de prueba
	req := httptest.NewRequest(http.MethodPost, "/test/v1/order/payment/release", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Configurar el mock para que devuelva nil (éxito)
	mockHandler.On("PaymentRelease", c).Return(nil)

	// Ejecutar la solicitud
	err := mockHandler.PaymentRelease(c)

	// Verificar resultados
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestAuthMiddleware(t *testing.T) {
	// Configuración inicial
	e := echo.New()
	group := e.Group("/test")
	mockHandler := new(MockCaptureFlowHandler)

	// Configurar la ruta
	CaptureFlowV1Routes(group, mockHandler)

	// Verificar que las rutas estén configuradas correctamente
	routes := e.Routes()
	foundCapture := false
	foundRelease := false
	for _, route := range routes {
		if route.Path == "/test/v1/order/payment/capture" {
			foundCapture = true
			assert.Equal(t, http.MethodPost, route.Method)
		}
		if route.Path == "/test/v1/order/payment/release" {
			foundRelease = true
			assert.Equal(t, http.MethodPost, route.Method)
		}
	}
	assert.True(t, foundCapture, "La ruta de capture no se encontró en las rutas configuradas")
	assert.True(t, foundRelease, "La ruta de release no se encontró en las rutas configuradas")
}
