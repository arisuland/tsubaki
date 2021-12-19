// Ratelimiter code is heavily modified from Delly (Discord Extreme List)
// Modified with permission! (Thanks Ice <3)

package ratelimiter

import (
	"arisu.land/tsubaki/pkg/managers"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ratelimitedResponse struct {
	Message    string `json:"message"`
	RetryAfter int64  `json:"retry_after"`
}

type Ratelimit struct {
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
	Limit     int       `json:"limit"`
}

func newRatelimit() Ratelimit {
	return Ratelimit{
		ResetTime: time.Now(),
		Remaining: 500,
		Limit:     500,
	}
}

func (r Ratelimit) Consume() *Ratelimit {
	r.Remaining = r.Remaining - 1
	return &r
}

func (r Ratelimit) Exceeded() bool {
	return !r.Expired() && r.Remaining == 0
}

func (r Ratelimit) Expired() bool {
	return r.ResetTime.UnixNano() < r.ResetTime.UnixNano()
}

type Ratelimiter struct {
	Ratelimits map[string]*Ratelimit
	Redis      *managers.RedisManager
	logger     slog.Logger
}

func NewRatelimiter(redis *managers.RedisManager) Ratelimiter {
	log := slog.Make(sloghuman.Sink(os.Stdout))
	s := time.Now()
	rl := Ratelimiter{
		Ratelimits: make(map[string]*Ratelimit),
		Redis:      redis,
		logger:     log,
	}

	count := redis.Connection.HLen(context.TODO(), "tsubaki:ratelimits").Val()
	log.Info(context.Background(), fmt.Sprintf("Took %s to get %d ratelimits", time.Now().Sub(s).String(), count))

	s = time.Now()
	result, err := redis.Connection.HGetAll(context.TODO(), "tsubaki:ratelimits").Result()
	if err != nil {
		log.Warn(
			context.Background(),
			"Unable to retrieve ratelimits, they have been reset.",
			slog.F("error", xerrors.Errorf("%v", err)),
		)
	} else {
		for key, value := range result {
			r := &Ratelimit{}
			err := json.Unmarshal([]byte(value), r)
			if err != nil {
				log.Warn(context.Background(), fmt.Sprintf("Unable to decode ratelimit for %s, skipping.", key))
				continue
			}

			rl.Ratelimits[key] = r
		}
	}

	log.Info(context.Background(), fmt.Sprintf("Took %s to re-add ratelimits (%d/%d ratelimits)", time.Now().Sub(s).String(), len(rl.Ratelimits), count))

	go rl.resetExpired()
	return rl
}

func (rl Ratelimiter) resetExpired() {
	for ip, value := range rl.Ratelimits {
		select {
		case <-time.After(time.Duration(value.ResetTime.UnixNano())):
			{
				rl.logger.Warn(context.Background(), fmt.Sprintf("Ratelimit for IP %s has been expired.", ip))
				err := rl.Redis.Connection.HDel(context.TODO(), "tsubaki:timeouts", ip).Err()
				if err != nil {
					rl.logger.Warn(
						context.Background(),
						fmt.Sprintf("Unable to delete IP %s:", ip),
						slog.F("error", xerrors.Errorf("%v", err)),
					)

					continue
				}
			}
		}
	}
}

func (rl Ratelimiter) Cache(ip string, ratelimit *Ratelimit) {
	data, _ := json.Marshal(&ratelimit)
	rl.Redis.Connection.HMSet(context.TODO(), "tsubaki:ratelimits", ip, string(data))
}

func (rl Ratelimiter) Get(ip string) *Ratelimit {
	res, err := rl.Redis.Connection.HGet(context.TODO(), "tsubaki:ratelimits", ip).Result()
	if err != nil {
		if err == redis.Nil {
			r := newRatelimit()
			rl.Cache(ip, &r)
			return &r
		}

		return nil
	}

	var ratelimit *Ratelimit
	err = json.Unmarshal([]byte(res), &ratelimit)
	if err != nil {
		rl.logger.Warn(
			context.TODO(),
			fmt.Sprintf("Unable to unmarshal ratelimit for ip %s:", ip),
			slog.F("error", xerrors.Errorf("%v", err)),
		)

		l := newRatelimit()
		ratelimit = &l
	}

	newRl := ratelimit.Consume()
	rl.Cache(ip, newRl)
	rl.Ratelimits[ip] = newRl

	return newRl
}

func (rl Ratelimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ratelimit := rl.Get(realIP(req))
		headers := w.Header()

		if ratelimit.Exceeded() {
			retry := strconv.FormatInt(time.Now().Sub(ratelimit.ResetTime).Milliseconds(), 10)
			headers.Set("Content-Type", "application/json")
			headers.Set("Retry-After", retry)

			w.WriteHeader(429)
			_ = json.NewEncoder(w).Encode(&ratelimitedResponse{
				Message:    fmt.Sprintf("Current IP %s has exceeded all ratelimits, try again later >:3", req.RemoteAddr),
				RetryAfter: time.Now().Sub(ratelimit.ResetTime).Milliseconds() / 1000,
			})

			return
		}

		headers.Set("X-RateLimit-Limit", strconv.Itoa(ratelimit.Limit))
		headers.Set("X-RateLimit-Remaining", strconv.Itoa(ratelimit.Remaining))
		headers.Set("X-RateLimit-Reset", strconv.FormatInt(ratelimit.ResetTime.Unix()*1000, 10))

		next.ServeHTTP(w, req)
	})
}

// https://github.com/go-chi/httprate/blob/master/httprate.go#L25-L47
func realIP(req *http.Request) string {
	var ip string
	if tcip := req.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := req.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		idx := strings.Index(xff, ", ")
		if idx == -1 {
			idx = len(xff)
		}

		// python moment
		ip = xff[:idx]
	} else {
		var err error

		ip, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ip = req.RemoteAddr
		}
	}

	return ip
}
