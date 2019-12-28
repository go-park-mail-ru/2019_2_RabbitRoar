package sd

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var log = logging.MustGetLogger("gameLimiter")

type GameLimiter struct {
	maxGames int
	KV       *consulapi.KV
}

func NewGameLimiter(consul *consulapi.Client) *GameLimiter {
	return &GameLimiter{
		maxGames: viper.GetInt("game.max_online"),
		KV:       consul.KV(),
	}
}

func (gl *GameLimiter) GetMaxGames() int {
	return gl.maxGames
}

func (gl *GameLimiter) RunPolling() {
	ticker := time.Tick(5 * time.Second)

	for range ticker {
		kv, _, err := gl.KV.Get("max_games", nil)
		if err != nil {
			log.Error("Error getting max_games keep previous value")
			continue
		}
		if kv == nil {
			log.Error("Error not max_games set keep previous value")
			continue
		}
		maxGames, err := strconv.Atoi(string(kv.Value))
		if err != nil {
			log.Error("Error interpreting max_games keep previous value")
			continue
		}
		log.Info("Updating max_games:", maxGames)
		gl.maxGames = maxGames
	}
}
