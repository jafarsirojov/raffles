package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/reply"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.SignIn
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.SignIn json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	token, err := h.authService.SignIn(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.SignIn h.authService.SignIn", zap.Error(err))
		response = responses.Unauthorized
		return
	}

	payload := structs.AuthPayload{Token: token}

	response = responses.Success
	response.Payload = payload
}
