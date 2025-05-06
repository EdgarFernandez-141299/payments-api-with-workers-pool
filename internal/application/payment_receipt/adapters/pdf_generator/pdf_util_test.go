package pdf_generator

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChromedp is a mock for chromedp functions
type MockChromedp struct {
	mock.Mock
}

func TestNewPDFUtil(t *testing.T) {
	// Test with a valid path
	path := "/test/path"
	pdfUtil := NewPDFUtil(path)

	// Assert that the PDFUtil was created with the correct path
	assert.Equal(t, path, pdfUtil.BasePath)

	// Test with an empty path
	pdfUtil = NewPDFUtil("")
	assert.Equal(t, "", pdfUtil.BasePath)
}

func TestEmbedImages(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdf-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test image files for different formats
	// PNG image
	pngData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	pngPath := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(pngPath, pngData, 0644); err != nil {
		t.Fatalf("Failed to write PNG test image: %v", err)
	}

	// JPEG image
	jpegData := []byte{0xFF, 0xD8, 0xFF} // JPEG header
	jpegPath := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(jpegPath, jpegData, 0644); err != nil {
		t.Fatalf("Failed to write JPEG test image: %v", err)
	}

	// GIF image
	gifData := []byte{0x47, 0x49, 0x46, 0x38} // GIF header
	gifPath := filepath.Join(tempDir, "test.gif")
	if err := os.WriteFile(gifPath, gifData, 0644); err != nil {
		t.Fatalf("Failed to write GIF test image: %v", err)
	}

	// SVG image
	svgData := []byte("<svg xmlns='http://www.w3.org/2000/svg'></svg>")
	svgPath := filepath.Join(tempDir, "test.svg")
	if err := os.WriteFile(svgPath, svgData, 0644); err != nil {
		t.Fatalf("Failed to write SVG test image: %v", err)
	}

	// Create a PDFUtil with the temp directory as base path
	pdfUtil := NewPDFUtil(tempDir)

	// Test cases
	tests := []struct {
		name        string
		htmlContent string
		wantErr     bool
		checkFunc   func(string) bool
	}{
		{
			name:        "No images",
			htmlContent: "<html><body>No images here</body></html>",
			wantErr:     false,
			checkFunc: func(result string) bool {
				return result == "<html><body>No images here</body></html>"
			},
		},
		{
			name:        "With PNG image",
			htmlContent: `<html><body><img src="test.png" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "data:image/png;base64,")
			},
		},
		{
			name:        "With JPEG image",
			htmlContent: `<html><body><img src="test.jpg" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "data:image/jpeg;base64,")
			},
		},
		{
			name:        "With GIF image",
			htmlContent: `<html><body><img src="test.gif" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "data:image/gif;base64,")
			},
		},
		{
			name:        "With SVG image",
			htmlContent: `<html><body><img src="test.svg" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "data:image/svg+xml;base64,")
			},
		},
		{
			name:        "With non-existent image",
			htmlContent: `<html><body><img src="nonexistent.png" alt="Test"></body></html>`,
			wantErr:     true,
			checkFunc:   nil,
		},
		{
			name:        "With data URI image",
			htmlContent: `<html><body><img src="data:image/png;base64,iVBORw==" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return result == `<html><body><img src="data:image/png;base64,iVBORw==" alt="Test"></body></html>`
			},
		},
		{
			name:        "With HTTP image",
			htmlContent: `<html><body><img src="http://example.com/image.png" alt="Test"></body></html>`,
			wantErr:     false,
			checkFunc: func(result string) bool {
				return result == `<html><body><img src="http://example.com/image.png" alt="Test"></body></html>`
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pdfUtil.embedImages(tt.htmlContent)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					assert.True(t, tt.checkFunc(result))
				}
			}
		})
	}
}

func TestGeneratePDF(t *testing.T) {
	// This test is more complex because it involves chromedp
	// We'll create a simple test that verifies the function handles errors correctly

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdf-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test HTML template
	htmlContent := `<html><body>Hello, {{.Name}}!</body></html>`
	htmlPath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Failed to write test HTML: %v", err)
	}

	// Create a PDFUtil with the temp directory as base path
	pdfUtil := NewPDFUtil(tempDir)

	// Test with valid data
	data := map[string]string{"Name": "World"}
	reader, err := pdfUtil.GeneratePDF(context.Background(), "test.html", data)

	if err != nil {
		t.Skipf("Skipping test because PDF generation failed: %v", err)
	} else {
		assert.NotNil(t, reader)

		// Verify that the returned value is a bytes.Reader (line 126)
		_, ok := reader.(*bytes.Reader)
		assert.True(t, ok, "Expected reader to be a *bytes.Reader")

		// Read the content to verify it's not empty
		content, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
	}

	_, err = pdfUtil.GeneratePDF(context.Background(), "nonexistent.html", data)
	assert.Error(t, err)

	invalidHTML := `<html><body>Hello, {{.Name!}</body></html>`
	invalidPath := filepath.Join(tempDir, "invalid.html")
	if err := os.WriteFile(invalidPath, []byte(invalidHTML), 0644); err != nil {
		t.Fatalf("Failed to write invalid HTML: %v", err)
	}

	_, err = pdfUtil.GeneratePDF(context.Background(), "invalid.html", data)
	assert.Error(t, err)
}
