package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"mime/multipart"
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

	estates, totalCount, err := h.adminService.GetEstates(ctx, offset, limit, status)
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
	response.Payload = structs.EstatesResponse{
		Estates: estates,
		Total:   totalCount,
	}
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

func (h *handler) UploadEstateImages(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	countImagesStr := mux.Vars(r)["count"]
	countImages, _ := strconv.Atoi(countImagesStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadEstateImages r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	var files []multipart.File
	for i := 1; i <= countImages; i++ {
		file, _, err := r.FormFile(strconv.Itoa(i))
		if err != nil {
			h.logger.Error("cmd.admin-api.handlers.UploadEstateImages r.FormFile - Error Retrieving the File", zap.Error(err))
			response = responses.BadRequest
			return
		}

		files = append(files, file)
		file.Close()
	}

	defer func() {
		for i, _ := range files {
			files[i].Close()
		}
	}()

	err = h.adminService.UploadEstateImages(ctx, id, files)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadEstateImages h.adminService.UploadEstateImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadEstateImages h.adminService.UploadEstateImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteEstateImages(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	imageName := mux.Vars(r)["imageName"]

	err := h.adminService.DeleteEstateImages(ctx, id, imageName)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.DeleteEstateImages h.adminService.DeleteEstateImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.DeleteEstateImages h.adminService.DeleteEstateImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) GetImageBaseURL(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	baseURL := h.adminService.GetImageBaseURL()

	response = responses.Success
	response.Payload = struct{ URL string }{URL: baseURL}
}
