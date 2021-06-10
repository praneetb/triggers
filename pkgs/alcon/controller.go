package alcon

import (
	"context"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Alcon instance of the alerts controller
type Alcon struct {
	RESTServer *http.Server
	Router     *httprouter.Router
}

// NewAlcon makes instance of alerts controller
func NewAlcon() (*Alcon, error) {
	alcon := &Alcon{
		RESTServer: &http.Server{},
	}

	return alcon, nil
}

// Init initializes the alerts controller.
func (c *Alcon) Init() error {
	err := InitWebServer(c)
	if err != nil {
		logrus.WithError(err).Error("Cannot initialize web-server")
		return err
	}
	return nil
}

// Run initializes the alerts controller.
func (c *Alcon) Run(ctx context.Context) error {
	defer func() {
		if err := c.Close(); err != nil {
			logrus.WithError(err).Info("Closing Controller failed")
		}
	}()

	err := RunWebServer(ctx, c)
	if err != nil {
		logrus.WithError(err).Error("Cannot Run web-server")
		return err
	}

	return nil
}

// Close closes trigger controller resources.
func (c *Alcon) Close() error {
	return nil
}
