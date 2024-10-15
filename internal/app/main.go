package app

import (
	"fmt"
	"net/http"
	"os"
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
		fmt.Println(err.Error())
		return
	}

	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filePath := "./notes.json"

	// check is file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		fmt.Println("Файл успешно создан:", filePath)

	} else {
		fmt.Println("Файл уже существует:", filePath)
	}

	log := zapLog.ZapLogger

	stor := storage.New("./notes.json")

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
