package config

import (
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(string) string
	GetInt(string) int
	GetDuration(key string) time.Duration
}

type config struct {
	cfg *viper.Viper
}

func NewConfig() Config {

	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("json")
	cfg.AddConfigPath(RootDir() + "/config")
	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()

	cfg.WatchConfig()

	return &config{cfg: cfg}
}

func (c *config) GetString(key string) string {
	return c.cfg.GetString(key)
}

func (c *config) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.cfg.GetDuration(key)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	sp := strings.Split(b, "post-loader-service")
	return sp[0] + "/post-loader-service/"
}
