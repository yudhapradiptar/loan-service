package server

import (
	"fmt"
	"net/http"

	"loan-service/internal/client"
	"loan-service/internal/config"
	"loan-service/internal/database"
	"loan-service/internal/handlers"
	"loan-service/internal/middleware"
	"loan-service/internal/repository"
	"loan-service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config        *config.Config
	logger        *logrus.Logger
	db            *database.Database
	loanHandler   handlers.LoanHandlerInterface
	healthHandler handlers.HealthHandlerInterface
}

func New(cfg *config.Config) *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	db, err := database.New(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}

	loanRepo := repository.NewLoanRepository(db.DB)
	notificationClient := client.NewNotificationClient(&cfg.Notification)
	loanService := service.NewLoanService(loanRepo, notificationClient)
	validator := validator.New()
	loanHandler := handlers.NewLoanHandler(loanService, validator)
	healthHandler := handlers.NewHealthHandler()

	return &Server{
		config:        cfg,
		logger:        logger,
		db:            db,
		loanHandler:   loanHandler,
		healthHandler: healthHandler,
	}
}

func (s *Server) setupRoutes() http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ErrorMiddleware)

	router.HandleFunc("/health", s.healthHandler.HealthCheck).Methods(http.MethodGet)

	api := router.PathPrefix("/v1").Subrouter()

	api.HandleFunc("/loans", s.loanHandler.GetAllLoans).Methods(http.MethodGet)
	api.HandleFunc("/loans", s.loanHandler.CreateLoan).Methods(http.MethodPost)
	api.HandleFunc("/loans/{uuid}", s.loanHandler.GetLoanByUUID).Methods(http.MethodGet)

	api.HandleFunc("/loans/{uuid}/approve", s.loanHandler.ApproveLoan).Methods(http.MethodPost)

	api.HandleFunc("/loans/{uuid}/invest", s.loanHandler.InvestLoan).Methods(http.MethodPost)
	api.HandleFunc("/loans/{uuid}/disburse", s.loanHandler.DisburseLoan).Methods(http.MethodPost)

	return router
}

func (s *Server) Start() error {
	handler := s.setupRoutes()

	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.logger.Infof("Starting server on %s", addr)

	return http.ListenAndServe(addr, handler)
}
