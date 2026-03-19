package getall

import (
	"errors"
	"log/slog"
	"net/http"
	"rest_API/internal/lib/api/response"
	"rest_API/internal/lib/logger/sl"
	"rest_API/internal/storage"
	"rest_API/internal/storage/interfaces"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	URLs map[string]string `json:"urls"`
}

func New(logger *slog.Logger, urlGetter interfaces.URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getall.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request id", middleware.GetReqID(r.Context())),
		)

		urls, err := urlGetter.GetAllURLs()
		if errors.Is(err, storage.ErrURLNotExists) {
			logger.Info("urls does not exists")

			render.JSON(w, r, response.Error("does not exists"))

			return
		}

		if err != nil {
			logger.Info("failed to get urls", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		logger.Info("urls successfully retrieved", slog.Int("count", len(urls)))

		render.JSON(w, r, Response{
			Response: response.Ok(),
			URLs:     urls,
		})
	}
}
