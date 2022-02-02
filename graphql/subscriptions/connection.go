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
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Connection is the current GraphQL subscription connection available
// in this context.
type Connection struct {
	conn *websocket.Conn
	uuid string
}

// Ack is the p
type Ack struct {
	OperationId string                 `json:"id,omitempty"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

func newConnection(conn *websocket.Conn) Connection {
	return Connection{
		conn: conn,
		uuid: uuid.NewV4().String(),
	}
}

func (c Connection) Connect() {
	logrus.Infof("Now processing messages for connection with UUID %s.", c.uuid)
}
