package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"rest_url_shorter/internal/http-server/handlers/url/save"
	"rest_url_shorter/internal/http-server/handlers/url/save/mocks"
	"rest_url_shorter/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "http://test.url",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "http://test.url",
		},
		{
			name:      "Empty url",
			alias:     "test_alias",
			url:       "",
			respError: "field URL is required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range cases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tt.respError == "" || tt.mockError != nil {
				urlSaverMock.On("SaveURL", tt.url, mock.AnythingOfType("string")).
					Return(tt.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tt.url, tt.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
