package router

import (
	"context"
	"crm/cmd/admin-api/handlers"
	"crm/internal/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Invoke(NewRouter)

type RouterParams struct {
	fx.In
	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
	Handler   handlers.Handler

	//AppConfig *config.AppConfig
}

func NewRouter(params RouterParams) {
	router := mux.NewRouter()

	baseUrl := "/api/crm-admin/v1"

	router.HandleFunc(baseUrl+"/signIn", middleware.ApplyMiddleware(params.Handler.SignIn)).Methods("POST")
	router.HandleFunc(baseUrl+"/leads", middleware.ApplyMiddleware(params.Handler.GetLeadList, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/leads/xlsx", middleware.ApplyMiddleware(params.Handler.GetLeadListXlSX, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/lead/{id:[0-9]+}", middleware.ApplyMiddleware(params.Handler.GetLeadAndCommentsByID, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/lead/{id:[0-9]+}/comment", middleware.ApplyMiddleware(params.Handler.AddComment, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/lead/{id:[0-9]+}/status", middleware.ApplyMiddleware(params.Handler.UpdateLeadStatus, params.Handler.MwCheckAdminAuthToken)).Methods("PUT")

	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}", middleware.ApplyMiddleware(params.Handler.GetEstateByID, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/estates", middleware.ApplyMiddleware(params.Handler.GetEstates, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/estate", middleware.ApplyMiddleware(params.Handler.AddEstate, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/estate", middleware.ApplyMiddleware(params.Handler.UpdateEstate, params.Handler.MwCheckAdminAuthToken)).Methods("PUT")
	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}", middleware.ApplyMiddleware(params.Handler.DeleteEstate, params.Handler.MwCheckAdminAuthToken)).Methods("DELETE")
	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}/approved", middleware.ApplyMiddleware(params.Handler.ApprovedEstate, params.Handler.MwCheckAdminAuthToken)).Methods("PUT")
	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}/uploadImages/{count:[0-9]+}", middleware.ApplyMiddleware(params.Handler.UploadEstateImages, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}/removeImage/{imageName}", middleware.ApplyMiddleware(params.Handler.DeleteEstateImages, params.Handler.MwCheckAdminAuthToken)).Methods("DELETE")
	router.HandleFunc(baseUrl+"/estate/imageBaseURL", middleware.ApplyMiddleware(params.Handler.GetImageBaseURL, params.Handler.MwCheckAdminAuthToken)).Methods("GET")

	router.HandleFunc(baseUrl+"/clients/favorites", middleware.ApplyMiddleware(params.Handler.GetClientsFavorites, params.Handler.MwCheckAdminAuthToken)).Methods("GET")

	router.HandleFunc(baseUrl+"/landing", middleware.ApplyMiddleware(params.Handler.AddLendingPage, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/landing", middleware.ApplyMiddleware(params.Handler.UpdateLendingPage, params.Handler.MwCheckAdminAuthToken)).Methods("PUT")
	router.HandleFunc(baseUrl+"/landing/{id:[0-9]+}", middleware.ApplyMiddleware(params.Handler.GetLendingData, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/landings", middleware.ApplyMiddleware(params.Handler.GetLendingList, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/landing/{id:[0-9]+}/uploadImages/{count:[0-9]+}", middleware.ApplyMiddleware(params.Handler.UploadLendingImages, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/landing/{id:[0-9]+}/removeImage/{imageName}", middleware.ApplyMiddleware(params.Handler.DeleteLendingImages, params.Handler.MwCheckAdminAuthToken)).Methods("DELETE")
	router.HandleFunc(baseUrl+"/landing/featuresAndAmenities", middleware.ApplyMiddleware(params.Handler.GetFeaturesAndAmenities, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/landing/{id:[0-9]+}/uploadBackgroundImage", middleware.ApplyMiddleware(params.Handler.UploadBackgroundImage, params.Handler.MwCheckAdminAuthToken)).Methods("POST")
	router.HandleFunc(baseUrl+"/landing/availability/{id:[0-9]+}/uploadPaymentPlan", middleware.ApplyMiddleware(params.Handler.UploadPaymentPlan, params.Handler.MwCheckAdminAuthToken)).Methods("POST")

	router.HandleFunc(baseUrl+"/texts", middleware.ApplyMiddleware(params.Handler.GetTexts, params.Handler.MwCheckAdminAuthToken)).Methods("GET")
	router.HandleFunc(baseUrl+"/text", middleware.ApplyMiddleware(params.Handler.UpdateText, params.Handler.MwCheckAdminAuthToken)).Methods("PUT")

	handler := cors.AllowAll().Handler(router)

	server := http.Server{
		Addr:    "0.0.0.0:7002", //params.AppConfig.ClientServerHost + params.AppConfig.ClientServerPort,
		Handler: handler,
	}

	params.Lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Admin application started")
				params.Logger.Info("start admin-api")
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("Admin application stopped")
				params.Logger.Info("stop admin-api")
				return server.Shutdown(ctx)
			},
		},
	)
}
