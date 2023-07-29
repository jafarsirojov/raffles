package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *handler) GetEstateByID(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	estate, err := h.adminService.GetEstateByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetEstateByID h.adminService.GetEstateByID not found",
				zap.Int("id", id))
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.GetEstateByID h.adminService.GetEstateByID", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estate
}

func (h *handler) GetEstates(w http.ResponseWriter, r *http.Request) {
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

	estates, err := h.adminService.GetEstates(ctx, offset, limit, status)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetEstates h.adminService.GetEstates not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.GetEstates h.adminService.GetEstates", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estates
}

func (h *handler) AddEstate(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Estate
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddEstate json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.AddEstate(ctx, request)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.AddEstate h.adminService.AddEstate not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.AddEstate h.adminService.AddEstate", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UpdateEstate(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Estate
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateEstate json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.UpdateEstate(ctx, request)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UpdateEstate h.adminService.UpdateEstate not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UpdateEstate h.adminService.UpdateEstate", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteEstate(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := h.adminService.DeleteEstate(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.DeleteEstate h.adminService.DeleteEstate not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.DeleteEstate h.adminService.DeleteEstate", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) ApprovedEstate(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := h.adminService.ApprovedEstate(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.ApprovedEstate h.adminService.ApprovedEstate not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.ApprovedEstate h.adminService.ApprovedEstate", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadImage r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, _, err := r.FormFile("myFile")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadImage r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadImages(ctx, id, &file)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadImage h.adminService.UploadImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadImage h.adminService.UploadImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	imageName := mux.Vars(r)["imageName"]

	err := h.adminService.DeleteImages(ctx, id, imageName)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.DeleteImage h.adminService.DeleteImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.DeleteImage h.adminService.DeleteImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
