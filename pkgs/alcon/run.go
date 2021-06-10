package alcon

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Triggers.AddCommand(processCmd)
}

var processCmd = &cobra.Command{
	Use:   "run",
	Short: "Start alerts controller",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		StartAlcon()
	},
}

// StartAlcon starts the alerts controller
func StartAlcon() {
	alcon, err := NewAlcon()
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = alcon.Init(); err != nil {
		logrus.Warn(err)
	}

	if err = alcon.Run(ctx); err != nil {
		logrus.Warn(err)
	}
	logrus.Debug("Starting Alerts controller...")
}
