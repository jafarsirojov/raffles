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

func (h *handler) MwCheckAuthToken(middleware middleware.Handler) middleware.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		var response structs.Response
		if len(strings.TrimSpace(authToken)) == 0 {
			h.logger.Warn("MwCheckAuthToken, token is empty", zap.Any("url", r.URL))
			defer reply.Json(w, http.StatusOK, &response)
			response = responses.Unauthorized
			return
		}

		userID, login, err := h.authService.CheckAuthToken(r.Context(), authToken)
		if err != nil {
			h.logger.Error("MwCheckAuthToken h.authService.CheckAuthToken", zap.Error(err), zap.Any("url", r.URL))
			defer reply.Json(w, http.StatusOK, &response)
			response = responses.Unauthorized
			return
		}

		if len(login) != 0 && userID != 0 {
			r = r.WithContext(context.WithValue(r.Context(), "login", login))
			r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
			middleware(w, r)
			return
		}

		h.logger.Debug("MwCheckAuthToken - admin sent a request", zap.String("login", login), zap.Any("url", r.URL))
		defer reply.Json(w, http.StatusOK, &response)
		response = responses.Forbidden
	}
}
