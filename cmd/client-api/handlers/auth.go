package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/reply"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.SignUp
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SignUp json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	token, err := h.authService.SignUp(ctx, request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SignUp h.authService.SignIn", zap.Error(err))
		response = responses.InternalErr
		return
	}

	payload := structs.AuthPayload{Token: token}

	response = responses.Success
	response.Payload = payload
}

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.SignIn
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SignIn json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	token, err := h.authService.SignIn(ctx, request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SignIn h.authService.SignIn", zap.Error(err))
		response = responses.Unauthorized
		return
	}

	payload := structs.AuthPayload{Token: token}

	response = responses.Success
	response.Payload = payload
}
