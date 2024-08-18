package start

import (
	"diploma1/internal/app/config"
	"diploma1/internal/app/service/logging"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Good struct {
	Match      string `json:"match"`
	Reward     int    `json:"reward"`
	RewardType string `json:"reward_type"`
}

func AccrualFilling() {
	path := fmt.Sprintf("http://%s/api/goods", config.State().AccrualAddress)
	good := Good{
		Match:      "Bork",
		Reward:     10,
		RewardType: "%",
	}

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(&good).
		Post(path)
	if err != nil {
		logging.Sugar.Error(err)
		return
	}
}
