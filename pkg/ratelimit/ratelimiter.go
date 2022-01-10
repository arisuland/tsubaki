// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2022 Noelware
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ratelimit

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
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
		ResetTime: time.Now().Add(1 * time.Hour),
		Remaining: 1200,
		Limit:     1200,
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
	return r.ResetTime.UnixNano() < time.Now().UnixNano()
}

type Ratelimiter struct {
	ratelimits map[string]*Ratelimit
	redis      *redis.Client
}

func NewRatelimiter(redis *redis.Client) Ratelimiter {
	s := time.Now()
	rl := Ratelimiter{
		ratelimits: make(map[string]*Ratelimit),
		redis:      redis,
	}

	count := redis.HLen(context.TODO(), "tsubaki:ratelimits").Val()
	logrus.Debugf("Took %s to get %d ratelimits", time.Since(s).String(), count)

	s = time.Now()
	result, err := redis.HGetAll(context.TODO(), "tsubaki:ratelimits").Result()
	if err != nil {
		logrus.Warnf("Unable to retrieve all ratelimits: %v", err)
	} else {
		for key, value := range result {
			r := &Ratelimit{}
			err := json.Unmarshal([]byte(value), r)
			if err != nil {
				logrus.Warnf("Unable to decode ratelimit for ip %s:\n%v", key, err)
				continue
			}

			rl.ratelimits[key] = r
		}
	}

	logrus.Debugf(
		"Took %s to re-implement all ratelimits (%d/%d ratelimits)",
		time.Now().Sub(s).String(),
		len(rl.ratelimits),
		count,
	)

	// Get all the expired ratelimits
	expired := make([]string, 0)
	for ip, value := range rl.ratelimits {
		if value.Expired() {
			expired = append(expired, ip)
		}
	}

	logrus.Debugf("Found %d ratelimits to expire!", len(expired))
	for _, r := range expired {
		_, err := redis.HDel(context.TODO(), "tsubaki:sessions", r).Result()
		if err != nil {
			logrus.Warnf("Unable to delete expired ratelimit:\n%v", err)
			continue
		}

		logrus.Infof("Deleted expired ratelimit for IP %s", r)
	}

	go rl.resetExpired()
	return rl
}

func (rl Ratelimiter) resetExpired() {
	for ip, value := range rl.ratelimits {
		select {
		case <-time.After(time.Duration(value.ResetTime.UnixNano())):
			{
				logrus.Warnf("Ratelimit for IP %s has been expired.", ip)
				_, err := rl.redis.HDel(context.TODO(), "tsubaki:timeouts", ip).Result()
				if err != nil {
					logrus.Warnf("Unable to remove ratelimit packet for %s:\n%v", ip, err)
					continue
				}
			}
		}
	}
}

func (rl Ratelimiter) cache(ip string, ratelimit *Ratelimit) {
	data, _ := json.Marshal(&ratelimit)
	rl.redis.HMSet(context.TODO(), "tsubaki:ratelimits", ip, string(data))
}

func (rl Ratelimiter) Get(ip string) *Ratelimit {
	res, err := rl.redis.HGet(context.TODO(), "tsubaki:ratelimits", ip).Result()
	if err != nil {
		if err == redis.Nil {
			r := newRatelimit()
			rl.cache(ip, &r)
			return &r
		}

		return nil
	}

	var ratelimit *Ratelimit
	err = json.Unmarshal([]byte(res), &ratelimit)
	if err != nil {
		logrus.Warnf("Unable to unmarshal ratelimit for IP %s, using new ratelimit packet", ip)

		l := newRatelimit()
		ratelimit = &l
	}

	newRl := ratelimit.Consume()
	rl.cache(ip, newRl)
	rl.ratelimits[ip] = newRl

	return newRl
}

func (rl Ratelimiter) Close() error {
	logrus.Info("Storing cached ratelimits in Redis!")

	for ip, value := range rl.ratelimits {
		// Check if it exists
		b, err := rl.redis.HExists(context.TODO(), "tsubaki:ratelimits", ip).Result()
		if err != nil {
			logrus.Fatalf("Unable to check if IP '%s' existed in ratelimit map: %v", ip, err)
			continue
		}

		if !b {
			rl.cache(ip, value)
		}
	}

	logrus.Info("Cached in-memory ratelimits in Redis!")
	return nil
}
