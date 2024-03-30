package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Address               string
	ReadBufferSize        int
	PrimeNodes            []string
	FindNodeSpeed         time.Duration
	ExpirationTime        time.Duration
	TransactionKeepTime   time.Duration
	TransactionTickerTime time.Duration
	AcquirerMaxSize       int
	AcquirerIntervalTime  time.Duration
	DatabaseType          string
	DatabaseFileName      string
}

var DefaultConfig = defaultConfig()

func defaultConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return &Config{
		Address:               viper.GetString("address"),
		ReadBufferSize:        viper.GetInt("read-buffer-size"),
		PrimeNodes:            viper.GetStringSlice("prime-nodes"),
		FindNodeSpeed:         time.Millisecond * time.Duration(viper.GetInt("find-node-speed")),
		ExpirationTime:        time.Minute * time.Duration(viper.GetInt("expiration-time")),
		TransactionKeepTime:   time.Second * time.Duration(viper.GetInt("transaction-keep-time")),
		TransactionTickerTime: time.Second * time.Duration(viper.GetInt("transaction-ticker-time")),
		AcquirerMaxSize:       viper.GetInt("acquirer-max-size"),
		AcquirerIntervalTime:  time.Second * time.Duration(viper.GetInt("acquirer-transaction-keep-time")),
		DatabaseType:          viper.GetString("database-type"),
		DatabaseFileName:      viper.GetString("database-filename"),
	}
}
