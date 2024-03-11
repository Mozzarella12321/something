package deleteUrl

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/mozzarella12321/orders-api/internal/lib/api/response"
	"github.com/mozzarella12321/orders-api/internal/lib/logger/sl"
	"github.com/mozzarella12321/orders-api/internal/postgresql"
)

type URLDeleter interface {
	DeleteUrl(alias string) (string, error)
}

type Request struct {
	URL   string `json:"url,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
		}

		log.Info("request body decoded", slog.Any("request", req))

		url, err := urlDeleter.DeleteUrl(req.Alias)

		if errors.Is(err, postgresql.ErrURLNotFound) {
			log.Info("url does not exist", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}
		log.Info("deleted url", slog.String("url", url))

		ResponseOK(w, r, req.Alias)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
