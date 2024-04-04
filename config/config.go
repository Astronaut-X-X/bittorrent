package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Address               string
	ReadBufferSize        int
	SendQueueSize         int
	NodeQueueSize         int
	PrimeNodes            []string
	SendMessageSpeed      time.Duration
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
		SendQueueSize:         viper.GetInt("send-queue-size"),
		NodeQueueSize:         viper.GetInt("node-queue-size"),
		PrimeNodes:            viper.GetStringSlice("prime-nodes"),
		SendMessageSpeed:      time.Millisecond * time.Duration(viper.GetInt("send-message-speed")),
		FindNodeSpeed:         time.Millisecond * time.Duration(viper.GetInt("find-node-speed")),
		ExpirationTime:        time.Minute * time.Duration(viper.GetInt("expiration-time")),
		TransactionKeepTime:   time.Second * time.Duration(viper.GetInt("transaction-keep-time")),
		TransactionTickerTime: time.Second * time.Duration(viper.GetInt("transaction-ticker-time")),
		AcquirerMaxSize:       viper.GetInt("acquirer-max-size"),
		AcquirerIntervalTime:  time.Second * time.Duration(viper.GetInt("acquirer-interval-time")),
		DatabaseType:          viper.GetString("database-type"),
		DatabaseFileName:      viper.GetString("database-filename"),
	}
}
