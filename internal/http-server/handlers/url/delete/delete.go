package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "rest_url_shorter/internal/lib/api/response"
	"rest_url_shorter/internal/lib/logger/sl"
)

type Response struct {
	resp.Response
	AffectedRows int64 `json:"affectedRows"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.35.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) (int64, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("url is required field")

			render.JSON(w, r, resp.Error("url is required field"))
			return
		}

		deletedRowsCount, err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error(
				"failed to delete url",
				sl.Err(err),
				slog.String("alias", alias),
			)

			render.JSON(w, r, resp.Error("failed to delete url"))
		}

		log.Info("url was deleted", slog.String("alias", alias))

		responseOK(w, r, deletedRowsCount)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, count int64) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		AffectedRows: count,
	})
}
