package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"rest_API/internal/lib/api/response"
	"rest_API/internal/lib/logger/sl"
	"rest_API/internal/storage"
	"rest_API/internal/storage/interfaces"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func New(logger *slog.Logger, urlDeleter interfaces.URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", "alias", alias)

			render.JSON(w, r, response.Error("not found"))

			return
		}

		if err != nil {
			logger.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		logger.Info("url deleted successfully", slog.String("alias", alias))
		render.JSON(w, r, response.Ok())
	}
}
