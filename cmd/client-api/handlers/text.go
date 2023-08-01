package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"go.uber.org/zap"
	"net/http"
)

func (h *handler) GetTexts(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	ctx := r.Context()

	texts, err := h.clientService.GetTexts(ctx)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetTexts h.clientService.GetTexts not found")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.GetTexts h.clientService.GetTexts", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = texts
}
