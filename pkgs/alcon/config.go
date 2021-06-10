package alcon

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	Triggers.PersistentFlags().StringVarP(&configFile, "config", "c",
		"", "Configuration File")
	viper.AutomaticEnv()
}

// Triggers defines root command.
var Triggers = &cobra.Command{
	Use:   "alcon",
	Short: "Triggers alerts controller command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func initConfig() {
	if configFile == "" {
		configFile = viper.GetString("config")
	}
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Viper cannot read Config, Error: %v", err)
	}

	levelStr := viper.GetString("server.log_level")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Fatalf("Cannot set Log Level, Error: %v", err)
	}
	logrus.SetLevel(level)
}
