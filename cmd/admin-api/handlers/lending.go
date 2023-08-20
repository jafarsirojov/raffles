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

func (h *handler) AddLendingPage(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Lending
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddLendingPage json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.SaveLending(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddLendingPage h.adminService.SaveLending", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UpdateLendingPage(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Lending
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateLendingPage json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.UpdateLending(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateLendingPage h.adminService.UpdateLending", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) GetLendingList(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	estate, err := h.adminService.GetLendingList(ctx)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetLendingList h.adminService.GetLendingList not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.GetLendingList h.adminService.GetLendingList", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estate
}

func (h *handler) GetLendingData(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	lending, err := h.adminService.GetLendingData(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.GetLendingData h.adminService.GetLendingData not found",
				zap.Int("id", id))
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.GetLendingData h.adminService.GetLendingData", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = lending
}

func (h *handler) UploadLendingImages(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	countImagesStr := mux.Vars(r)["count"]
	countImages, _ := strconv.Atoi(countImagesStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadLendingImages r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	var files []multipart.File
	for i := 1; i <= countImages; i++ {
		file, _, err := r.FormFile(strconv.Itoa(i))
		if err != nil {
			h.logger.Error("cmd.admin-api.handlers.UploadLendingImages r.FormFile - Error Retrieving the File", zap.Error(err))
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

	err = h.adminService.UploadLendingImages(ctx, id, files)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadLendingImages h.adminService.UploadLendingImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadLendingImages h.adminService.UploadLendingImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteLendingImages(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	imageName := mux.Vars(r)["imageName"]

	err := h.adminService.DeleteLendingImages(ctx, id, imageName)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.DeleteLendingImages h.adminService.DeleteLendingImages not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.DeleteLendingImages h.adminService.DeleteLendingImages", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) GetFeaturesAndAmenities(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	estate, err := h.adminService.GetFeaturesAndAmenities(ctx)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info(
				"cmd.admin-api.handlers.GetFeaturesAndAmenities h.adminService.GetFeaturesAndAmenities not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.GetFeaturesAndAmenities h.adminService.GetFeaturesAndAmenities",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estate
}

func (h *handler) UploadBackgroundImage(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundImage r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, _, err := r.FormFile("1")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundImage r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadBackgroundImage(ctx, id, file)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadBackgroundImage h.adminService.UploadBackgroundImage not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundImage h.adminService.UploadBackgroundImage",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadPaymentPlan(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadPaymentPlan r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, _, err := r.FormFile("1")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadPaymentPlan r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadPaymentPlan(ctx, id, file)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadPaymentPlan h.adminService.UploadPaymentPlan not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadPaymentPlan h.adminService.UploadPaymentPlan",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
