package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
	"os/user"
	"runtime"
)

func init() {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller, "user", os.Getenv("USER"))
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

var homeDir, _ = homedir.Dir()
var sys = &afero.Afero{
	Fs: afero.NewOsFs(),
}
const fileName = ".chronic"

type Chronic struct {
	Cron *cron.Cron
	*viper.Viper
}

func New() *Chronic {
	return &Chronic{
		Cron:  cron.New(),
		Viper: viper.New(),
	}
}

func (c *Chronic) Init() error {
	c.AddConfigPath(homeDir)
	c.SetConfigName(fileName)
	c.AutomaticEnv()

	c.SetDefault("os.env", os.Environ())

	c.SetDefault("runtime.goarch", runtime.GOARCH)
	c.SetDefault("runtime.compiler", runtime.Compiler)
	c.SetDefault("runtime.version", runtime.Version())
	c.SetDefault("runtime.goos", runtime.GOOS)
	usr, _ := user.Current()
	c.SetDefault("os.user", usr)
	return nil
}

func Annotate(v *viper.Viper) map[string]string {
		settings := v.AllSettings()
		an := make(map[string]string)
		for k, v := range settings {
			if t, ok := v.(string); ok == true {
				an[k] = t
			}
		}
		return an
}

func (c *Chronic) Write() error {
	// If a config file is found, read it in.
	b, err := sys.Exists(homeDir+"/.chronic.yaml")
	if err != nil {
		return errors.WithStack(err)
	}
	if !b {
		f, err := sys.Create(homeDir+"/.chronic.yaml")
		if err != nil {
			return errors.WithStack(err)
		}
		c.SetConfigFile(f.Name())
	}
	if err := c.ReadInConfig(); err != nil {
		log.Println("failed to read config file, writing defaults...")
		if err := c.WriteConfig(); err != nil {
			return errors.Wrap(err, "failed to write config")
		}

	} else {
		log.Println("Using config file-->", c.ConfigFileUsed())
		if err := c.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
