package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"go.uber.org/zap"
	"net/http"
)

func (h *handler) GetLendingData(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	authKey := r.Header.Get("ExternalServiceKey")

	data, err := h.clientService.GetLendingData(ctx, authKey)
	if err != nil {
		if err == errors.ErrBadRequest {
			h.logger.Info("cmd.client-api.handlers.GetLendingData h.clientService.GetLendingData bad request")
			response = responses.BadRequest
			return
		}

		h.logger.Error("cmd.client-api.handlers.GetLendingData h.clientService.GetLendingData", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = data
}
