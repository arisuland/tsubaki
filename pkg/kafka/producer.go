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

package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

// Producer is the main Kafka writer to send out messages to its consumers.
type Producer struct {
	// Writer is the kafka.Writer instance to send out messages to.
	Writer *kafka.Writer

	// Config is the configuration passed down from NewProducer
	Config Config
}

// NewProducer creates a new Producer with the provided Config object.
func NewProducer(config Config) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		Config: config,
		Writer: writer,
	}
}

// Write produces a new message as the form of JSON to all consumers, in which,
// it will return an error if anything occurs.
func (prod *Producer) Write(key []byte, data interface{}) error {
	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	err = prod.Writer.WriteMessages(context.TODO(), kafka.Message{
		// TODO: figure out what partition to do?
		// for now, i think 0 is a good way to do this iunno D:
		Partition: 0,
		Key:       key,
		Value:     bytes,
	})

	if err != nil {
		return err
	}

	return nil
}
