package translationapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fprojetto/pokedex-api/api-client/translationapi"
	"github.com/fprojetto/pokedex-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslationClient_Translate(t *testing.T) {
	tests := []struct {
		name           string
		style          service.TranslationStyle
		text           string
		mockStatus     int
		mockResponse   any
		expectedResult string
		expectedError  error
	}{
		{
			name:       "successful yoda translation",
			style:      service.Yoda,
			text:       "hello",
			mockStatus: http.StatusOK,
			// Refactored: Use named types SuccessInfo and TranslationContent
			mockResponse: translationapi.TranslationResponse{
				Success: translationapi.SuccessInfo{
					Total: 1,
				},
				Contents: translationapi.TranslationContent{
					Translated: "hello, you must",
					Text:       "hello",
					Translation: "yoda",
				},
			},
			expectedResult: "hello, you must",
			expectedError:  nil,
		},
		{
			name:       "successful shakespeare translation",
			style:      service.Shakespeare,
			text:       "hello",
			mockStatus: http.StatusOK,
			// Refactored: Use named types SuccessInfo and TranslationContent
			mockResponse: translationapi.TranslationResponse{
				Success: translationapi.SuccessInfo{
					Total: 1,
				},
				Contents: translationapi.TranslationContent{
					Translated: "Hark, greetings!",
					Text:       "hello",
					Translation: "shakespeare",
				},
			},
			expectedResult: "Hark, greetings!",
			expectedError:  nil,
		},
		{
			name:          "unsupported translation style",
			style:         service.TranslationStyle("pirate"), // Example of unsupported style
			text:          "hello",
			expectedError: errors.New("unsupported translation style"),
		},
		{
			name:       "api internal server error",
			style:      service.Yoda,
			text:       "hello",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]string{
				"error": "internal server error",
			},
			expectedError: service.ErrServiceUnavailable,
		},
		{
			name:       "api not found error",
			style:      service.Yoda,
			text:       "hello",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]string{
				"error": "not found",
			},
			expectedError: service.ErrServiceUnavailable,
		},
		{
			name:       "invalid json response",
			style:      service.Yoda,
			text:       "hello",
			mockStatus: http.StatusOK,
			mockResponse: "this is not valid json", // Invalid JSON
			expectedError: errors.Join(errors.New("invalid character 'h' in literal true (expecting 'r')"), service.ErrServiceUnavailable),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var expectedPath string
				switch tt.style {
				case service.Yoda:
					expectedPath = "/translate/yodish"
				case service.Shakespeare:
					expectedPath = "/translate/shakespeare-english"
				default:
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}

				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				w.WriteHeader(tt.mockStatus)
				if tt.mockStatus == http.StatusOK {
					if invalidJSON, ok := tt.mockResponse.(string); ok {
						w.Write([]byte(invalidJSON)) // Write invalid JSON for error case
					} else {
						json.NewEncoder(w).Encode(tt.mockResponse)
					}
				} else {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer ts.Close()

			client, err := translationapi.NewClient(ts.URL)
			require.NoError(t, err, "Failed to create translationapi client")

			result, err := client.Translate(context.Background(), tt.style, tt.text)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.name == "invalid json response" {
					assert.Contains(t, err.Error(), "invalid character 'h' in literal true (expecting 'r')")
					assert.Contains(t, err.Error(), service.ErrServiceUnavailable.Error())
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

// Helper to create a pointer to a boolean, used for mock responses if needed, though not strictly required by translationapi.
// This is kept from pokeapi_test.go as a reference, but not directly used in translationapi tests.
func BoolPtr(b bool) *bool {
	return &b
}
