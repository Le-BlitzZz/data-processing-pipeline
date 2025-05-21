package server

import (
	"Le-BlitzZz/streaming-etl-app/internal/uploader/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, conf *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	log.Info("server: started")

	router := gin.Default()

	tcpSocket := fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort())

	server := &http.Server{
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       -1,
		WriteTimeout:      -1,
		Handler:           router,
		Addr:              tcpSocket,
	}

	log.Infof("server: listening on %s", server.Addr)

	go StartHttp(server)

	<-ctx.Done()

	log.Info("server: shutting down")

	err := server.Close()
	if err != nil {
		log.Errorf("server: shutdown failed (%s)", err)
	}
}

// StartHttp starts the Web server in http mode.
func StartHttp(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("server: shutdown complete")
		} else {
			log.Errorf("server: %s", err)
		}
	}
}
