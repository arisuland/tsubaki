// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2021 Noelware
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
	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/pkg/managers"
	"arisu.land/tsubaki/prisma/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var Sessions *Manager = nil

type response struct {
	Message string `json:"message"`
}

// sessionType is a type to determine the session's type
// At the moment, all sessions (for now) are just "web"
// but subject is to change.
type sessionType string

var web sessionType = "desktop"

// Session represents the current session that is cached in Redis.
// Your IP is not stored here, only in ratelimiting.
type Session struct {
	// ExpiresIn returns the time that this will expire.
	ExpiresIn time.Time `json:"expires_in"`

	// Token represents the JWT token used for sessions.
	Token string `json:"token"`

	// User represents the current user that this Session is attached.
	//
	// using types.User rather than db.UserModel will guarantee that
	// passwords and/or emails are not leaked in Redis.
	User types.User `json:"user"`

	// Type represents the current sessionType.
	Type sessionType `json:"session_type"`
}

func NewSession(user types.User, token string) *Session {
	// go doesn't have `days` or `weeks`, so this is a substitute for now.
	days := 24 * time.Hour

	return &Session{
		ExpiresIn: time.Now().Add(2 * days),
		Token:     token,
		User:      user,
		Type:      web,
	}
}

func (s *Session) Expired() bool {
	return s.ExpiresIn.UnixNano() < time.Now().UnixNano()
}

// Manager is the manager for handling all user sessions.
// This is mapped in-memory but cached in Redis.
type Manager struct {
	Sessions map[string]*Session
	redis    *managers.RedisManager
	db       managers.Prisma
}

func NewSessionManager(redis *managers.RedisManager, prisma managers.Prisma) Manager {
	if Sessions != nil {
		panic("tried to create new session manager.")
	}

	s := time.Now()
	sessions := Manager{
		Sessions: make(map[string]*Session),
		redis:    redis,
		db:       prisma,
	}

	count := redis.Connection.HLen(context.TODO(), "tsubaki:sessions").Val()
	logrus.Infof("Retrieved %d documents in %s", count, time.Since(s).String())

	s = time.Now()
	result, err := redis.Connection.HGetAll(context.TODO(), "tsubaki:sessions").Result()
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

			sessions.Sessions[key] = &session
		}
	}

	logrus.Infof("Took %s to re-implement all sessions (%d/%d sessions)", time.Since(s).String(), len(sessions.Sessions), count)

	Sessions = &sessions
	// Get all the expired ratelimits
	expired := make([]string, 0)
	for uid, value := range sessions.Sessions {
		if value.Expired() {
			expired = append(expired, uid)
		}
	}

	logrus.Infof("Found %d ratelimits to expire!", len(expired))
	for _, r := range expired {
		err := sessions.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", r).Err()
		if err != nil {
			logrus.Warnf("Unable to delete expired ratelimit:\n%v", err)
		}

		logrus.Infof("Deleted expired session for uid %s", r)
	}

	go sessions.resetExpired()
	return sessions
}

func (m Manager) resetExpired() {
	for id, value := range m.Sessions {
		select {
		case <-time.After(time.Duration(value.ExpiresIn.UnixNano())):
			{
				logrus.Warnf("Session for user %s has expired.", id)
				err := m.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", id).Err()
				if err != nil {
					logrus.Errorf("Unable to delete sesion for user %s:\n%v", id, err)
					continue
				}
			}
		}
	}
}

func (m Manager) cache(uid string, session *Session) {
	data, _ := json.Marshal(&session)
	m.redis.Connection.HMSet(context.TODO(), "tsubaki:sessions", uid, string(data))
}

func (m Manager) Get(uid string) *Session {
	res, err := m.redis.Connection.HGet(context.TODO(), "tsubaki:sessions", uid).Result()
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
	_, err = m.db.Client.User.FindUnique(
		db.User.ID.Equals(uid),
	).Exec(context.TODO())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logrus.Warnf("User %s doesn't exist in the database anymore but sessions still exists!")

			err := m.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", uid).Err()
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

func (m Manager) New(uid string) *Session {
	logrus.Infof("Creating session for user %s...", uid)
	token, err := NewToken(uid)
	if err != nil {
		logrus.Errorf("Unable to create JWT token for uid %s:\n%v", uid, err)
		return nil
	}

	// find the user in the database
	user, err := m.db.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.TODO())
	if err != nil {
		logrus.Errorf("Unable to find user %s from database:\n%v", uid, err)
		return nil
	}

	u := types.FromUserModel(user)

	sess := NewSession(*u, token)
	m.cache(uid, sess)
	m.Sessions[uid] = sess

	return sess
}

func (m Manager) Delete(uid string) {
	logrus.Warnf("Deleting session for user %s...", uid)
	err := m.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", uid).Err()
	if err != nil {
		if err == redis.Nil {
			return
		}
		logrus.Errorf("Unable to delete session from user %s:\n%v", uid, err)
	}
}

func (m Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// if there is no authorization header, let's just skip it.
		if req.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, req)
			return
		}

		auth := req.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer") {
			// TODO: access tokens
			next.ServeHTTP(w, req)
		} else if strings.HasPrefix(auth, "Session") {
			decoded, err := DecodeToken(req.Header.Get("Authorization"))
			if err != nil {
				w.WriteHeader(406)
				_ = json.NewEncoder(w).Encode(&response{
					Message: fmt.Sprintf("Invalid token: %s", req.Header.Get("Authorization")),
				})

				return
			}

			// get user id from MapClaims
			uid, ok := decoded["uid"].(string)
			if !ok {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&response{
					Message: "Unable to cast `uid` ~> string.",
				})

				return
			}

			// add it to the request context
			user, err := m.db.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.TODO())
			if err != nil {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&response{
					Message: "Unable to retrieve user from database.",
				})

				return
			}

			// blep!
			ctx := context.WithValue(req.Context(), "user_id", user.ID)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(406)
			_ = json.NewEncoder(w).Encode(&response{
				Message: "Missing `Bearer` or `Session` prefix.",
			})
		}
	})
}
