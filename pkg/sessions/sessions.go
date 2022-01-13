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

package sessions

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/prisma/db"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

// sessionType is a type to determine the user's session.
// At the moment, this is just "web" but might come soon!
type sessionType string

// web refers to the web app located at arisu.land!
var web sessionType = "web"

// String stringifies a sessionType value.
func (s sessionType) String() string {
	switch s {
	case "web":
		return "web"

	default:
		return "unknown"
	}
}

// Session represents the current session a user has. This is
// cached in Redis.
//
// IPs are not stored here, only for ratelimiting to determine
// the root network.
type Session struct {
	// ExpiresIn refers the time until this Session will expire.
	ExpiresIn time.Time `json:"expires_in"`

	// Token is the JWT token for this Session.
	Token string `json:"token"`

	// User is the current user of this Session is attached to.
	User types.User `json:"user"`

	// Type represents the current sessionType.
	Type sessionType `json:"session_type"`
}

func NewSession(user types.User, token string) *Session {
	days := 24 * time.Hour
	return &Session{
		ExpiresIn: time.Now().Add(2 * days),
		User:      user,
		Token:     token,
		Type:      web,
	}
}

func (s *Session) Expired() bool {
	return s.ExpiresIn.UnixNano() < time.Now().UnixNano()
}

// SessionManager is the manager for handling all user sessions.
type SessionManager struct {
	sessions map[string]*Session
	prisma   *db.PrismaClient
	redis    *redis.Client
}

func NewSessionManager(redis *redis.Client, prisma *db.PrismaClient) SessionManager {
	s := time.Now()
	m := SessionManager{
		sessions: make(map[string]*Session),
		prisma:   prisma,
		redis:    redis,
	}

	count := redis.HLen(context.TODO(), "tsubaki:sessions").Val()
	logrus.Debugf("Retrieved %d documents in %s", count, time.Since(s).String())

	s = time.Now()
	result, err := redis.HGetAll(context.TODO(), "tsubaki:sessions").Result()
	if err != nil {
		logrus.Warnf("Unable to retrieve all sessions!")
	} else {
		for key, value := range result {
			var session Session
			err := json.Unmarshal([]byte(value), &session)
			if err != nil {
				logrus.Warnf("Unable to unmarshal packet for %s:\n%v", key, err)
				continue
			}

			m.sessions[key] = &session
		}
	}

	logrus.Debugf("Took %s to re-implement all sessions (%d/%d sessions)", time.Since(s).String(), len(m.sessions), count)

	// Get all the expired ratelimits
	expired := make([]string, 0)
	for uid, value := range m.sessions {
		if value.Expired() {
			expired = append(expired, uid)
		}
	}

	logrus.Debugf("Found %d ratelimits to expire!", len(expired))
	for _, r := range expired {
		_, err := redis.HDel(context.TODO(), "tsubaki:sessions", r).Result()
		if err != nil {
			logrus.Warnf("Unable to delete expired ratelimit:\n%v", err)
		}

		logrus.Infof("Deleted expired session for uid %s", r)
	}

	go m.resetExpired()
	return m
}

func (m SessionManager) resetExpired() {
	for id, value := range m.sessions {
		select {
		case <-time.After(time.Duration(value.ExpiresIn.UnixNano())):
			{
				logrus.Warnf("Session for user %s has expired.", id)
				_, err := m.redis.HDel(context.TODO(), "tsubaki:sessions", id).Result()
				if err != nil {
					logrus.Errorf("Unable to delete sesion for user %s:\n%v", id, err)
					continue
				}
			}
		}
	}
}

func (m SessionManager) cache(uid string, session *Session) {
	data, _ := json.Marshal(&session)
	m.redis.HMSet(context.TODO(), "tsubaki:sessions", uid, string(data))
}

func (m SessionManager) Get(uid string) *Session {
	res, err := m.redis.HGet(context.TODO(), "tsubaki:sessions", uid).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		logrus.Errorf("Unable to fetch session for uid %s.", uid)
		return nil
	}

	var session *Session
	err = json.Unmarshal([]byte(res), &session)
	if err != nil {
		logrus.Warnf("Unable to unmarshal session packet for uid %s:\n%v", uid, err)
		return nil
	}

	// find the user and check if it exists
	_, err = m.prisma.User.FindUnique(
		db.User.ID.Equals(uid),
	).Exec(context.TODO())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logrus.Warnf("User %s doesn't exist in the database anymore but sessions still exists!", uid)

			_, err := m.redis.HDel(context.TODO(), "tsubaki:sessions", uid).Result()
			if err != nil {
				logrus.Fatalf("I was unable to delete the session packet for you. Maybe try using `redis-cli`?\n> HDEL tsubaki:sessions %s", uid)
				return nil
			}

			return nil
		}

		logrus.Errorf("Unable to find user %s in database:\n%v", uid, err)
		return nil
	}

	return session
}

func (m SessionManager) New(uid string) *Session {
	logrus.Infof("Creating session for user %s...", uid)
	token, err := pkg.NewToken(uid)
	if err != nil {
		logrus.Errorf("Unable to create JWT token for uid %s:\n%v", uid, err)
		return nil
	}

	// find the user in the database
	user, err := m.prisma.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.TODO())
	if err != nil {
		logrus.Errorf("Unable to find user %s from database:\n%v", uid, err)
		return nil
	}

	u := types.FromDbModel(user)

	sess := NewSession(u, token)
	m.cache(uid, sess)
	m.sessions[uid] = sess

	return nil
}

func (m SessionManager) Delete(uid string) {
	logrus.Warnf("Deleting session for user %s...", uid)
	_, err := m.redis.HDel(context.TODO(), "tsubaki:sessions", uid).Result()
	if err != nil {
		if err == redis.Nil {
			return
		}
		logrus.Errorf("Unable to delete session from user %s:\n%v", uid, err)
	}
}

func (m SessionManager) Close() error {
	logrus.Info("Storing cached sessions in Redis!")

	for ip, value := range m.sessions {
		// Check if it exists
		b, err := m.redis.HExists(context.TODO(), "tsubaki:ratelimits", ip).Result()
		if err != nil {
			logrus.Fatalf("Unable to check if IP '%s' existed in ratelimit map: %v", ip, err)
			continue
		}

		if !b {
			m.cache(ip, value)
		}
	}

	logrus.Info("Cached in-memory sessions in Redis!")
	return nil
}
