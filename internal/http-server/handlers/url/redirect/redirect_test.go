package redirect_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"rest_url_shorter/internal/http-server/handlers/url/redirect"
	"rest_url_shorter/internal/http-server/handlers/url/redirect/mocks"
	"rest_url_shorter/internal/lib/api"
	"rest_url_shorter/internal/lib/logger/handlers/slogdiscard"
	"testing"
)

func TestGetURLHandler(t *testing.T) {
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
			url:   "https://www.google.com",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				urlGetterMock.On("GetURL", tt.alias).
					Return(tt.url, tt.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			redirectURL, err := api.GetRedirectedURL(testServer.URL + "/" + tt.alias)
			require.NoError(t, err)

			assert.Equal(t, tt.url, redirectURL)
		})
	}
}
