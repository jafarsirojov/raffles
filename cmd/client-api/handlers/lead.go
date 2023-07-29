package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *handler) SaveLead(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Lead
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SaveLead json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.clientService.SaveLead(ctx, request)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.client-api.handlers.SaveLead h.clientService.SaveLead bad request",
				zap.Any("request", request))
			response = responses.BadRequest
			return
		}

		h.logger.Error("cmd.client-api.handlers.SaveLead h.clientService.SaveLead", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
