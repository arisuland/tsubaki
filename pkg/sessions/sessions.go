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
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
	"net/http"
	"os"
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
	weeks := 7 * days

	return &Session{
		ExpiresIn: time.Now().Add(1 * weeks),
		Token:     token,
		User:      user,
		Type:      web,
	}
}

// Manager is the manager for handling all user sessions.
// This is mapped in-memory but cached in Redis.
type Manager struct {
	Sessions map[string]*Session
	redis    *managers.RedisManager
	log      slog.Logger
	db       managers.Prisma
}

func NewSessionManager(redis *managers.RedisManager, prisma managers.Prisma) Manager {
	if Sessions != nil {
		panic("tried to create new session manager.")
	}

	l := slog.Make(sloghuman.Sink(os.Stdout))
	s := time.Now()
	sessions := Manager{
		Sessions: make(map[string]*Session),
		redis:    redis,
		log:      l,
		db:       prisma,
	}

	count := redis.Connection.HLen(context.TODO(), "tsubaki:sessions").Val()
	l.Info(context.Background(), fmt.Sprintf("Took %s to get %d sessions.", time.Since(s).String(), count))

	s = time.Now()
	result, err := redis.Connection.HGetAll(context.TODO(), "tsubaki:sessions").Result()
	if err != nil {
		l.Warn(
			context.Background(),
			"Unable to retrieve sessions, they will be not cached in memory.",
			slog.F("error", xerrors.Errorf("%v", err)),
		)
	} else {
		for key, value := range result {
			var session Session
			err := json.Unmarshal([]byte(value), &session)
			if err != nil {
				l.Warn(context.Background(), fmt.Sprintf("Unable to decode session for user %s, skipping.", key))
				continue
			}

			sessions.Sessions[key] = &session
		}
	}

	l.Info(
		context.Background(),
		fmt.Sprintf(
			"Took %s to re-implement sessions (%d/%d sessions)",
			time.Since(s).String(),
			len(sessions.Sessions),
			count,
		),
	)

	Sessions = &sessions
	go sessions.resetExpired()
	return sessions
}

func (m Manager) resetExpired() {
	for id, value := range m.Sessions {
		select {
		case <-time.After(time.Duration(value.ExpiresIn.UnixNano())):
			{
				m.log.Warn(context.Background(), fmt.Sprintf("Sessions for %s has expired.", id))
				err := m.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", id).Err()
				if err != nil {
					m.log.Warn(
						context.Background(),
						fmt.Sprintf("Unable to delete session for user %s:", id),
						slog.F("error", xerrors.Errorf("%v", err)),
					)

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

func (m Manager) Get(uid string, req *http.Request) *Session {
	res, err := m.redis.Connection.HGet(context.TODO(), "tsubaki:sessions", uid).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		m.log.Warn(context.Background(), fmt.Sprintf("Unable to get session from uid %s.", uid))
		return nil
	}

	var session *Session
	err = json.Unmarshal([]byte(res), &session)
	if err != nil {
		m.log.Warn(
			context.Background(),
			fmt.Sprintf("Unable to unmarshal session packet for uid %s:", uid),
			slog.F("error", xerrors.Errorf("%v", err)))

		// useless to create a new session packet.
		return nil
	}

	// find the user and check if it exists
	_, err = m.db.Client.User.FindUnique(
		db.User.ID.Equals(uid),
	).Exec(context.TODO())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			m.log.Warn(
				context.Background(),
				fmt.Sprintf("User %s doesn't exist anymore (but has sessions still), removing from Redis...", uid),
			)

			err := m.redis.Connection.HDel(context.TODO(), "tsubaki:sessions", uid).Err()
			if err != nil {
				m.log.Error(
					context.Background(),
					fmt.Sprintf("Unfortunately, I was unable to delete the packet for you. Run this command in `redis-cli` to do so:\n\nHDEL tsubaki:sessions %s", uid),
				)

				return nil
			}

			return nil
		}

		m.log.Error(
			context.Background(),
			"Unable to find user from database:",
			slog.F("error", xerrors.Errorf("%v", err)))

		return nil
	}

	return session
}

func (m Manager) New(uid string) *Session {
	m.log.Info(context.Background(), fmt.Sprintf("Creating user session for %s...", uid))
	token, err := NewToken(uid)
	if err != nil {
		m.log.Error(
			context.Background(),
			fmt.Sprintf("Unable to create JWT token for uid %s:", uid),
			slog.F("error", xerrors.Errorf("%v", err)))

		return nil
	}

	// find the user in the database
	user, err := m.db.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.TODO())
	if err != nil {
		m.log.Error(
			context.Background(),
			"Unable to find user from database:",
			slog.F("error", xerrors.Errorf("%v", err)))

		return nil
	}

	u := types.FromUserModel(user)

	sess := NewSession(*u, token)
	m.cache(uid, sess)
	m.Sessions[uid] = sess

	return sess
}

func (m Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// if there is no authorization header, let's just skip it.
		if req.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, req)
			return
		}

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
		ctx := context.WithValue(req.Context(), "user", user)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
