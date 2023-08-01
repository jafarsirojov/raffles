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

func (h *handler) GetTexts(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	ctx := r.Context()

	texts, err := h.adminService.GetTexts(ctx)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetTexts h.adminService.GetTexts not found")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.GetTexts h.adminService.GetTexts", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = texts
}

func (h *handler) UpdateText(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Text
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateText json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.UpdateText(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateText h.adminService.UpdateText", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
