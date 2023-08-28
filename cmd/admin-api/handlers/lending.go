package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"crm/pkg/util"
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

	err := r.ParseMultipartForm(30 << 20)
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

const (
	BackgroundImage     = "BackgroundImage"
	BackgroundForMobile = "BackgroundForMobile"
	FilePlan            = "FilePlan"
	MainLogo            = "MainLogo"
	PartnerLogo         = "PartnerLogo"
	OurLogo             = "OurLogo"
)

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.Upload r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	keys := []string{BackgroundImage, BackgroundForMobile, FilePlan, MainLogo, PartnerLogo, OurLogo}
	var file multipart.File
	var info *multipart.FileHeader
	var methodKey string
	for _, key := range keys {
		file, info, err = r.FormFile(key)
		if err == nil {
			methodKey = key
			break
		}
	}
	defer file.Close()

	if len(methodKey) == 0 {
		h.logger.Error("cmd.admin-api.handlers.Upload r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}

	switch methodKey {
	case BackgroundImage:
		err = h.adminService.UploadBackgroundImage(ctx, id, file)
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadBackgroundImage not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadBackgroundImage",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	case BackgroundForMobile:
		err = h.adminService.UploadBackgroundForMobile(ctx, id, file)
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadBackgroundForMobile not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadBackgroundForMobile",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	case FilePlan:
		err = h.adminService.UploadFilePlan(ctx, id, file)
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadFilePlan not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadFilePlan",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	case MainLogo:
		err = h.adminService.UploadMainLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadMainLogo not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadMainLogo",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	case PartnerLogo:
		err = h.adminService.UploadPartnerLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadPartnerLogo not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadPartnerLogo",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	case OurLogo:
		err = h.adminService.UploadOurLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
		if err != nil {
			if err == errors.ErrNotFound {
				h.logger.Info("cmd.admin-api.handlers.Upload h.adminService.UploadOurLogo not found")
				response = responses.NotFound
				return
			}
			h.logger.Error("cmd.admin-api.handlers.Upload h.adminService.UploadOurLogo",
				zap.Error(err))
			response = responses.InternalErr
			return
		}

		response = responses.Success
		return

	default:
		h.logger.Error("cmd.admin-api.handlers.Upload not found file")
		response = responses.BadRequest
		return
	}
}

func (h *handler) GetFileURL(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)
	var ctx = r.Context()

	url := h.adminService.GetFileURL(ctx)

	response = responses.Success
	response.Payload = struct{ URL string }{URL: url}
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

	file, _, err := r.FormFile("BackgroundImage")
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

func (h *handler) UploadBackgroundForMobile(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundForMobile r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, _, err := r.FormFile("BackgroundForMobile")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundForMobile r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadBackgroundForMobile(ctx, id, file)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadBackgroundForMobile h.adminService.UploadBackgroundForMobile not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadBackgroundForMobile h.adminService.UploadBackgroundForMobile",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) AddAvailability(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Availability
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddAvailability json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.SaveAvailability(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddAvailability h.adminService.SaveAvailability", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Availability
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateAvailability json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.UpdateAvailability(ctx, request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UpdateAvailability h.adminService.UpdateAvailability", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) RemoveAvailability(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	var request structs.Availability
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.RemoveAvailability json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.adminService.RemoveAvailability(ctx, id)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.RemoveAvailability h.adminService.RemoveAvailability", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadFilePlan(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadFilePlan r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, _, err := r.FormFile("FilePlan")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadFilePlan r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadFilePlan(ctx, id, file)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadFilePlan h.adminService.UploadFilePlan not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadFilePlan h.adminService.UploadFilePlan",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadMainLogo(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadMainLogo r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, info, err := r.FormFile("MainLogo")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadMainLogo r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadMainLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadMainLogo h.adminService.UploadMainLogo not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadMainLogo h.adminService.UploadMainLogo",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadPartnerLogo(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadPartnerLogo r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, info, err := r.FormFile("PartnerLogo")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadPartnerLogo r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadPartnerLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadPartnerLogo h.adminService.UploadPartnerLogo not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadPartnerLogo h.adminService.UploadPartnerLogo",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) UploadOurLogo(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadOurLogo r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, info, err := r.FormFile("OurLogo")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.UploadOurLogo r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.UploadOurLogo(ctx, id, file, util.GetFileTypeByFilename(info.Filename))
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.UploadOurLogo h.adminService.UploadOurLogo not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.UploadOurLogo h.adminService.UploadOurLogo",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) AddFeatureAndAmenity(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	featureName := r.URL.Query().Get("featureName")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddFeatureAndAmenity r.ParseMultipartForm", zap.Error(err))
		response = responses.BadRequest
		return
	}

	file, info, err := r.FormFile("FeatureAndAmenity")
	if err != nil {
		h.logger.Error("cmd.admin-api.handlers.AddFeatureAndAmenity r.FormFile - Error Retrieving the File", zap.Error(err))
		response = responses.BadRequest
		return
	}
	defer file.Close()

	err = h.adminService.AddFeatureAndAmenity(ctx, file, util.GetFileTypeByFilename(info.Filename), featureName)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.admin-api.handlers.AddFeatureAndAmenity h.adminService.AddFeatureAndAmenity not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.AddFeatureAndAmenity h.adminService.AddFeatureAndAmenity",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteFeatureAndAmenity(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	err := h.adminService.DeleteFeatureAndAmenity(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info(
				"cmd.admin-api.handlers.DeleteFeatureAndAmenity h.adminService.DeleteFeatureAndAmenity not found")
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.admin-api.handlers.DeleteFeatureAndAmenity h.adminService.DeleteFeatureAndAmenity",
			zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}
