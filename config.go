package sola

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("sola")
}

// LoadConfig by viper
func (s *Sola) LoadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	s.devMode = viper.GetBool(cfgDev)
}

const (
	cfgDev = "sola.dev"
)
