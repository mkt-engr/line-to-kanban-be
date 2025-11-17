package config

import "os"

type LineConfig struct {
	ChannelSecret      string
	ChannelAccessToken string
}

func LoadLineConfig() *LineConfig {
	return &LineConfig{
		ChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		ChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	}
}
