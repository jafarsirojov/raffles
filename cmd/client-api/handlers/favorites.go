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

func (h *handler) SaveFavorite(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	var request structs.Favorite
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SaveFavorite json.NewDecoder", zap.Error(err))
		response = responses.BadRequest
		return
	}

	err = h.clientService.SaveFavorite(ctx, 7, request.EstateID)
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.SaveFavorite json.NewDecoder", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	estateIDStr := mux.Vars(r)["estateID"]
	estateID, _ := strconv.Atoi(estateIDStr)

	userID := ctx.Value("id").(int)

	err := h.clientService.DeleteFavorite(ctx, structs.Favorite{
		UserID:   userID,
		EstateID: estateID,
	})
	if err != nil {
		h.logger.Error("cmd.client-api.handlers.DeleteFavorite h.clientService.DeleteFavorite", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
}

func (h *handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	favorites, err := h.clientService.GetEstateFavorites(ctx, 7)
	if err != nil {
		if err == errors.ErrNotFound {
			response = responses.NotFound
			return
		}
		h.logger.Error("cmd.client-api.handlers.GetFavorites h.clientService.GetFavorites",
			zap.Error(err), zap.Int("userID", 0))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = favorites
}
