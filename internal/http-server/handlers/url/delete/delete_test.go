package delete_test

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	del "rest_url_shorter/internal/http-server/handlers/url/delete"
	"rest_url_shorter/internal/http-server/handlers/url/delete/mocks"
	"rest_url_shorter/internal/lib/logger/handlers/slogdiscard"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		affectedRows int64
		respError    string
		mockError    error
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			affectedRows: 0,
		},
		{
			name:         " ",
			alias:        "s",
			affectedRows: 0,
		},
	}

	for _, tt := range cases {
		//tt := tt

		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tt.respError == "" || tt.mockError != nil {
				urlDeleterMock.On("DeleteURL", tt.alias).
					Return(tt.affectedRows, tt.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", del.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			urlForTest := testServer.URL + "/" + tt.alias

			client := &http.Client{}

			req, err := http.NewRequest("DELETE", urlForTest, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer func() { _ = resp.Body.Close() }()

			var response del.Response
			err = json.NewDecoder(resp.Body).Decode(&response)

			require.NoError(t, err)

			require.Equal(t, tt.affectedRows, response.AffectedRows)

			/*handler := del.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/"+tt.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp del.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.affectedRows, resp.AffectedRows)*/
		})
	}
}
