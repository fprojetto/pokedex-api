package translationapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fprojetto/pokedex-api/pkg/client"
	"github.com/fprojetto/pokedex-api/service"
)

// SuccessInfo represents the success details of a translation API response.
type SuccessInfo struct {
	Total int `json:"total"`
}

// TranslationContent represents the translated text and original text from the API.
type TranslationContent struct {
	Translated  string `json:"translated"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// TranslationResponse represents the overall structure of the translation API response.
type TranslationResponse struct {
	Success  SuccessInfo        `json:"success"`
	Contents TranslationContent `json:"contents"`
}

type TranslationRequest struct {
	Text string `json:"text"`
}

type TranslationClient struct {
	translationAPIURL string
	client            *http.Client
}

func NewClient(translationAPIURL string) (*TranslationClient, error) {
	if translationAPIURL == "" {
		return nil, errors.New("translationAPIURL empty string")
	}

	return &TranslationClient{
		translationAPIURL: translationAPIURL,
		client:            client.HttpClient(),
	}, nil
}

func (c *TranslationClient) Translate(ctx context.Context, style service.TranslationStyle, text string) (string, error) {
	var translatorServiceName string
	switch style {
	case service.Yoda:
		translatorServiceName = "yodish"
	case service.Shakespeare:
		translatorServiceName = "shakespeare-english"
	default:
		return "", errors.New("unsupported translation style")
	}

	getTranslationURL := fmt.Sprintf("%s/translate/%s", c.translationAPIURL, translatorServiceName)
	translationRequest := TranslationRequest{Text: text}
	jsonBody, err := json.Marshal(translationRequest)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", getTranslationURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", service.ErrServiceUnavailable
	}

	var translation TranslationResponse
	if err := json.NewDecoder(res.Body).Decode(&translation); err != nil {
		return "", errors.Join(err, service.ErrServiceUnavailable)
	}

	return translation.Contents.Translated, nil
}
