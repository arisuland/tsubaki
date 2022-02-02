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

package subscriptions

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// SubscriptionServer is the current server context for handling WebSocket connections.
// This follows Apollo's [graphql-transport-ws](https://github.com/apollographql/subscriptions-transport-ws/blob/master/PROTOCOL.md) protocol,
// since [Fubuki](https://github.com/auguwu/Arisu/tree/master/web) implements the Apollo Client to execute queries, mutations, and subscriptions.
//
// But, why does Tsubaki embed subscriptions? Since we emit the following:
//    - user / project db updates
//    - administrator information (for admin dashboard, if enabled)
type SubscriptionServer struct {
	upgrader websocket.Upgrader
}

// NewSubscriptionServer creates a new SubscriptionServer.
func NewSubscriptionServer() SubscriptionServer {
	return SubscriptionServer{
		upgrader: websocket.Upgrader{
			Subprotocols: []string{"graphql-ws"},
		},
	}
}

// ServeHTTP is the middleware for implementing GraphQL subscriptions.
func (h SubscriptionServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, v := range websocket.Subprotocols(req) {
		if v != "graphql-ws" {
			continue
		}

		conn, err := h.upgrader.Upgrade(w, req, nil)
		if err != nil {
			w.Header().Set("X-WebSocket-Upgrade-Failure", err.Error())
			return
		}

		if conn.Subprotocol() != "graphql-ws" {
			w.Header().Set("X-WebSocket-Upgrade-Failure", fmt.Sprintf("tried to upgrade websocket, but wrong protocol: %s", conn.Subprotocol()))
			_ = conn.Close()

			return
		}

		connection := newConnection(conn)
		go connection.Connect()

		return
	}

	w.Header().Set("X-WebSocket-Upgrade-Failure", "no subprotocols available :(")
}
