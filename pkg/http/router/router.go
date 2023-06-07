package router

import (
	"event-history/pkg/eventinfo"
	"event-history/pkg/http/internal/handler"
	"event-history/pkg/http/internal/middleware"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(lgr *zap.Logger, eventsService eventinfo.Service) http.Handler {
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())

	eventsHandler := handler.NewEventsHandler(lgr, eventsService)

	router.HandleFunc("/", withMiddlewares(lgr, middleware.WithErrorHandler(lgr, eventsHandler.Create))).Methods(http.MethodPost)
	router.HandleFunc("/latest/{user_id}/{key}", withMiddlewares(lgr, middleware.WithErrorHandler(lgr, eventsHandler.Get))).Methods(http.MethodGet)
	router.HandleFunc("/", withMiddlewares(lgr, middleware.WithErrorHandler(lgr, eventsHandler.Update))).Methods(http.MethodPut)
	router.HandleFunc("/{user_id}/{key}", withMiddlewares(lgr, middleware.WithErrorHandler(lgr, eventsHandler.Delete))).Methods(http.MethodDelete)
	router.HandleFunc("/{user_id}/{key}", withMiddlewares(lgr, middleware.WithErrorHandler(lgr, eventsHandler.GetHistory))).Methods(http.MethodGet)

	return router
}

func withMiddlewares(lgr *zap.Logger, hnd http.HandlerFunc) http.HandlerFunc {
	return middleware.WithSecurityHeaders(middleware.WithReqResLog(lgr, hnd))
}
