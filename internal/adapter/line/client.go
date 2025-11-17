package line

import (
	"line-to-kanban-be/internal/platform/config"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Client struct {
	bot *linebot.Client
}

func NewClient(cfg *config.LineConfig) (*Client, error) {
	bot, err := linebot.New(cfg.ChannelSecret, cfg.ChannelAccessToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot: bot,
	}, nil
}

func (c *Client) GetBot() *linebot.Client {
	return c.bot
}
