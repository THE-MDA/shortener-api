package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "rest_API/internal/lib/api/response"
	"rest_API/internal/lib/logger/sl"
	"rest_API/internal/lib/random"
	"rest_API/internal/storage"

	"rest_API/internal/storage/interfaces"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URl   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLenght = 5

func New(log *slog.Logger, urlSaver interfaces.URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidatorError(validatorErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		id, err := urlSaver.SaveURL(req.URl, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URl))
			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.Ok(),
		Alias:    alias,
	})
}
