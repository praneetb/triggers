package alcon

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitWebServer initializes the Web Server
func InitWebServer(ctrl *Alcon) (err error) {

	// http handlers
	http.HandleFunc("/", DefaultHandler)

	return nil
}

// RunWebServer runs the Web Server
func RunWebServer(ctx context.Context, ctrl *Alcon) (err error) {
	server := ctrl.RESTServer

	server.Addr = viper.GetString("server.endpoint")

	logrus.WithFields(logrus.Fields{
		"IP:Port": server.Addr,
	}).Info("Starting API Server")

	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {

	case <-ctx.Done():
		logrus.Debugf("REST/GRPC WebServer Terminated")
		err = ctx.Err()
	case err = <-errChan:

	}
	return err
}
