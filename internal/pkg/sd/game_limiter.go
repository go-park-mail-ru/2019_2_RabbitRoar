package sd

import (
	"encoding/json"
	consulapi "github.com/hashicorp/consul/api"
	"time"
)

type MaxGames struct {
	max_games int64
	cache_time int
}

type GameLimiter struct {
	mg MaxGames
	kv *consulapi.KV
	lastUpdate int
}

func (gl *GameLimiter) Init() (error) {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := consulapi.NewClient(config)
	if (err != nil) {
		return err
	}
	gl.mg.cache_time = 1000
	gl.mg.max_games = 0
	gl.kv = consul.KV()
	gl.lastUpdate = 0
	return nil
}

func (gl *GameLimiter) GetMaxGames() (int64, error) {
	currentTime := time.Now().Second()
	if currentTime - gl.lastUpdate < gl.mg.cache_time {
		return gl.mg.max_games, nil
	}

	kvp, _, err := gl.kv.Get("max_games", nil)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(kvp.Value, &gl.mg); err != nil {
		return 0, err
	}

	gl.lastUpdate = time.Now().Second()
	return gl.mg.max_games, nil
}
