package handlers

import (
	"context"
	"crm/internal/middleware"
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/reply"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (h *handler) MwCheckAdminAuthToken(middleware middleware.Handler) middleware.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		var response structs.Response
		if len(strings.TrimSpace(authToken)) == 0 {
			h.logger.Warn("MwCheckAdminAuthToken, token is empty", zap.Any("url", r.URL))
			defer reply.Json(w, http.StatusOK, &response)
			response = responses.Unauthorized
			return
		}

		id, login, err := h.authService.CheckAdminAuthToken(r.Context(), authToken)
		if err != nil {
			h.logger.Error("MwCheckAdminAuthToken h.authService.CheckAdminAuthToken", zap.Error(err), zap.Any("url", r.URL))
			defer reply.Json(w, http.StatusOK, &response)
			response = responses.Unauthorized
			return
		}

		if len(login) != 0 && id != 0 {
			r = r.WithContext(context.WithValue(r.Context(), "login", login))
			r = r.WithContext(context.WithValue(r.Context(), "id", id))
			middleware(w, r)
			return
		}

		h.logger.Debug("MwCheckAdminAuthToken - admin sent a request", zap.String("login", login), zap.Any("url", r.URL))
		defer reply.Json(w, http.StatusOK, &response)
		response = responses.Forbidden
	}
}
