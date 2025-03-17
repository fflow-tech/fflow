package pdf

import (
	"context"
	"encoding/json"
	"io"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/google/go-tika/tika"
)

// TikaPdfParseResult is the result of parsing a PDF file with Tika.
type TikaPdfParseResult struct {
	XTIKAContent string `json:"X-TIKA:content"`
}

// Extractor extracts text from PDF files.
type Extractor struct {
	url string
}

// NewExtractor creates a new Extractor.
func NewExtractor(url string) *Extractor {
	return &Extractor{url: url}
}

// Extract extracts text from PDF files.
func (e *Extractor) Extract(ctx context.Context, input io.Reader) (string, error) {
	c := tika.NewClient(nil, e.url)
	got, err := c.Parse(ctx, input)
	if err != nil {
		return "", err
	}
	var parseResult TikaPdfParseResult
	if err := json.Unmarshal([]byte(got), &parseResult); err != nil {
		return "", err
	}

	return md.NewConverter("", true, nil).ConvertString(parseResult.XTIKAContent)
}
