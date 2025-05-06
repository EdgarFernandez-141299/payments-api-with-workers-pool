package pdf_generator

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFUtil struct {
	BasePath string
}

// NewPDFUtil creates a new PDFUtil instance.
func NewPDFUtil(basePath string) *PDFUtil {
	return &PDFUtil{BasePath: basePath}
}

// embedImages replaces <img src="…"> references that point to local files
// with data‑URI versions so the HTML is fully self‑contained.
func (p *PDFUtil) embedImages(htmlContent string) (string, error) {
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"[^>]*>`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		imgSrc := match[1]

		if strings.HasPrefix(imgSrc, "data:") || strings.HasPrefix(imgSrc, "http") {
			continue
		}

		imgPath := filepath.Join(p.BasePath, imgSrc)
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			return "", err
		}

		mimeType := "image/png"
		switch {
		case strings.HasSuffix(imgSrc, ".jpg"), strings.HasSuffix(imgSrc, ".jpeg"):
			mimeType = "image/jpeg"
		case strings.HasSuffix(imgSrc, ".gif"):
			mimeType = "image/gif"
		case strings.HasSuffix(imgSrc, ".svg"):
			mimeType = "image/svg+xml"
		}

		dataURI := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imgData)
		htmlContent = strings.ReplaceAll(htmlContent, `src="`+imgSrc+`"`, `src="`+dataURI+`"`)
	}
	return htmlContent, nil
}

// GeneratePDF renders the given template with the provided data and returns
// an io.Reader that streams the resulting PDF.
func (p *PDFUtil) GeneratePDF(
	ctx context.Context,
	htmlPath string,
	data map[string]string,
) (io.Reader, error) {

	tmpl, err := template.ParseFiles(filepath.Join(p.BasePath, htmlPath))
	if err != nil {
		return nil, err
	}

	var htmlBuf bytes.Buffer
	if err = tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, err
	}

	htmlContent, err := p.embedImages(htmlBuf.String())
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
		chromedp.NoSandbox,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pdfBuffer []byte

	printParams := page.PrintToPDF().
		WithPrintBackground(true).
		WithPreferCSSPageSize(true).
		WithPaperWidth(4.17). // 400px at 96dpi
		WithPaperHeight(6.5)  // Estimated height based on content

	err = chromedp.Run(taskCtx,
		chromedp.Navigate("data:text/html;charset=utf-8,"+strings.ReplaceAll(htmlContent, "#", "%23")),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = printParams.Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(pdfBuffer), nil
}
