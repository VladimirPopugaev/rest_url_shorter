package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "rest_url_shorter/internal/lib/api/response"
	"rest_url_shorter/internal/lib/logger/sl"
	"rest_url_shorter/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.redirect.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("url is required field")

			render.JSON(w, r, resp.Error("url is required field"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("got url", slog.String("url", url))

		//TODO: add redirect to url
		http.Redirect(w, r, url, http.StatusFound)
	}
}
