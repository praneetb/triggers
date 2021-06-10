package alcon

import (
	"context"

	httprouter "github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Route endpoints for the controller
type Route struct {
	Name   string
	Method string
	Path   string
	Handle httprouter.Handle
}

var endpoints = []Route{
	{
		"root",
		"GET",
		"/default",
		DefaultHandler,
	},
}

// CreateNewHTTPRouter creates router endpoints
func CreateNewHTTPRouter() *httprouter.Router {
	hrouter := httprouter.New()

	for _, route := range endpoints {
		logrus.Debugf("hello %v", route.Name)
		hrouter.Handle(route.Method, route.Path, route.Handle)
	}
	return hrouter
}

// InitWebServer initializes the Web Server
func InitWebServer(ctrl *Alcon) (err error) {
	// Create Router handlers
	hrouter := CreateNewHTTPRouter()
	ctrl.Router = hrouter
	ctrl.RESTServer.Handler = hrouter

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
