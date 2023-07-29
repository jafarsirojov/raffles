package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func (h *handler) GetLeadList(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	status := r.URL.Query().Get("status")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	list, err := h.adminService.GetLeadList(ctx, offset, limit, status)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetLeadList h.adminService.GetLeadList bad request")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.GetLeadList h.adminService.GetLeadList", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = list
}

func (h *handler) GetLeadListXlSX(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	status := r.URL.Query().Get("status")
	url, err := h.adminService.GetLeadListXLSX(ctx, status)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetLeadListXlSX h.adminService.GetLeadList bad request")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.GetLeadListXlSX h.adminService.GetLeadList", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = struct {
		URL string `json:"url"`
	}{URL: url}
}

func (h *handler) UpdateLeadStatus(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	status := r.URL.Query().Get("status")

	if len(strings.TrimSpace(status)) == 0 {
		response = responses.BadRequest
		return
	}

	err := h.adminService.UpdateLeadStatus(ctx, id, status)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UpdateLeadStatus h.adminService.UpdateLeadStatus bad request",
				zap.Int("id", id),
				zap.String("status", status))
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.UpdateLeadStatus h.adminService.UpdateLeadStatus", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) GetLeadAndCommentsByID(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	lc, err := h.adminService.GetLeadAndCommentsByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetLeadAndCommentsByID h.adminService.GetLeadAndCommentsByID not found",
				zap.Int("id", id))
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.admin-api.handlers.GetLeadAndCommentsByID h.adminService.GetLeadAndCommentsByID", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = lc
}

func (h *handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	leadIDStr := mux.Vars(r)["id"]
	leadID, _ := strconv.Atoi(leadIDStr)

	comment := r.URL.Query().Get("comment")

	adminLogin := ctx.Value("login").(string)
	adminID := ctx.Value("id").(int)

	err := h.adminService.AddComment(ctx, structs.Comment{
		LeadID:     leadID,
		AdminID:    adminID,
		AdminLogin: adminLogin,
		Text:       comment,
	})
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddComment h.adminService.AddComment", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
