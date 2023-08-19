package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (h *handler) GetLendingData(w http.ResponseWriter, r *http.Request) {
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
