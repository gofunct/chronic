package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"
	kitlog "github.com/go-kit/kit/log"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfg string) func() {
	return func() {
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller, "user", os.Getenv("USER"))
		log.SetOutput(kitlog.NewStdlibAdapter(logger))

		if cfg != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfg)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".goscript" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".chronic")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}