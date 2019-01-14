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
	ex, err := os.Executable()
	if err != nil {
		return errors.WithStack(err)
	}
	c.AddConfigPath(homeDir)
	c.SetConfigName(fileName)
	c.AutomaticEnv()

	c.SetDefault("home", homeDir)
	c.SetDefault("file_name", fileName)
	c.SetDefault("executable", ex)
	gr, err := os.Getgroups()
	if err != nil {
		return errors.WithStack(err)
	}
	c.SetDefault("groups", gr)
	host, err := os.Hostname()
	if err != nil {
		return errors.WithStack(err)
	}
	c.SetDefault("env", os.Environ())
	c.SetDefault("uid", os.Getuid())
	c.SetDefault("args", os.Args)
	c.SetDefault("host_name", host)
	c.SetDefault("pid", os.Getpid())
	c.SetDefault("goarch", runtime.GOARCH)
	c.SetDefault("compiler", runtime.Compiler)
	c.SetDefault("runtime_version", runtime.Version())
	c.SetDefault("goos", runtime.GOOS)
	usr, _ := user.Current()
	c.SetDefault("user", usr)
	return nil
}

func (c *Chronic) Annotate() map[string]string {
	settings := c.AllSettings()
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
