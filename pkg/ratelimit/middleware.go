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
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://github.com/go-chi/httprate/blob/master/httprate.go#L25-L47

func RealIP(req *http.Request) string {
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

func (rl Ratelimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		limit := rl.Get(RealIP(req))
		headers := w.Header()

		if limit.Exceeded() {
			retry := strconv.FormatInt(time.Now().Sub(limit.ResetTime).Milliseconds(), 10)
			headers.Set("Retry-After", retry)

			w.WriteHeader(429)
			_ = json.NewEncoder(w).Encode(&ratelimitedResponse{
				Message:    fmt.Sprintf("Current IP %s has exceeded all ratelimits, try again later >:3", req.RemoteAddr),
				RetryAfter: time.Now().Sub(limit.ResetTime).Milliseconds() / 1000,
			})

			return
		}

		headers.Set("X-RateLimit-Limit", strconv.Itoa(limit.Limit))
		headers.Set("X-RateLimit-Remaining", strconv.Itoa(limit.Remaining))
		headers.Set("X-RateLimit-Reset", strconv.FormatInt(limit.ResetTime.Unix()*1000, 10))

		next.ServeHTTP(w, req)
	})
}
