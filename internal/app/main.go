package app

import (
	"net/http"
	"web-server/internal/config"
	"web-server/internal/service"
	"web-server/internal/storage"
	"web-server/internal/transport/http/handler"
	"web-server/internal/transport/http/router"
	"web-server/pkg/logger"

	"go.uber.org/zap"
)

func StartServer() {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
	}

	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("logger init error", zap.Error(err))
	}

	log := zapLog.ZapLogger

	dbConn, err := storage.GetConnect(cfg.DBDSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}

	defer dbConn.Close()

	stor := storage.New(dbConn, log)

	serv := service.New(stor)

	hand := handler.New(serv, log)

	r := router.New(&hand)

	log.Info("starting server", zap.String("address", cfg.ServerAddress))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
