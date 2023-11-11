package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rl/pkg/config"
	"rl/pkg/messages"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CSI_MAX_PER_MIN   = 5
	CSI_MODEL_PER_MIN = 5
)

func CreateRateLimiterMiddleware(cfg *config.Config, redisAddr string) func(http.Handler) http.Handler {

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			req := messages.Request{}
			err = json.Unmarshal(reqBytes, &req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			csiModelKey := fmt.Sprintf("%s-%s-%d", req.CSI, req.Model, time.Now().Minute())
			csiKey := fmt.Sprintf("%s-%d", req.CSI, time.Now().Minute())
			incrResult := rdb.Incr(r.Context(), csiKey)
			rdb.Expire(r.Context(), csiKey, 1*time.Minute)
			csiMax := calculateCSIMax(cfg, req.CSI)
			if incrResult.Val() > csiMax {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			incrModelResult := rdb.Incr(r.Context(), csiModelKey)
			rdb.Expire(r.Context(), csiModelKey, 1*time.Minute)
			modelMax := calculateModelMax(cfg, req.CSI, req.Model)
			if incrModelResult.Val() > modelMax {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func calculateCSIMax(cfg *config.Config, csi string) int64 {
	val := cfg.Defaults.CSI
	for _, rule := range cfg.Rules {
		if rule.CSI == csi {
			val = rule.Limit
		}
	}

	return val
}

func calculateModelMax(cfg *config.Config, csi string, model string) int64 {
	val := int64(0)
	for _, m := range cfg.Defaults.Models {
		if model == m.Name {
			val = m.Value
		}
	}

	for _, rule := range cfg.Rules {
		if rule.CSI == csi {
			for _, m := range rule.Models {
				if m.Name == model {
					val = m.Value
				}
			}
		}
	}

	return val
}
