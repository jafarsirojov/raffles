package router

import (
	"context"
	"crm/cmd/client-api/handlers"
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

	baseUrl := "/api/crm-client/v1"

	router.HandleFunc(baseUrl+"/lead", params.Handler.SaveLead).Methods("POST")

	router.HandleFunc(baseUrl+"/estates", params.Handler.GetEstates).Methods("GET")
	router.HandleFunc(baseUrl+"/estates/luxury", params.Handler.GetLuxuryEstates).Methods("GET")
	router.HandleFunc(baseUrl+"/estate/{id:[0-9]+}", params.Handler.GetEstateByID).Methods("GET")
	router.HandleFunc(baseUrl+"/estate/searchOptions", params.Handler.GetSearchOptions).Methods("GET")
	router.HandleFunc(baseUrl+"/estate/clearCache", params.Handler.ClearCache).Methods("DELETE")
	router.HandleFunc(baseUrl+"/estate/imageBaseURL", params.Handler.GetImageBaseURL).Methods("GET")

	handler := cors.AllowAll().Handler(router)

	server := http.Server{
		Addr:    "0.0.0.0:7001", //params.AppConfig.ClientServerHost + params.AppConfig.ClientServerPort,
		Handler: handler,
	}

	params.Lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Client application started")
				params.Logger.Info("start client-api")
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("Client application stopped")
				params.Logger.Info("stop client-api")
				return server.Shutdown(ctx)
			},
		},
	)
}
