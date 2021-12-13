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

package kafka

// Config is the configuration options for configuring the Producer.
// This is required if you're running the GitHub bot.
//
// Default: nil | Prefix: TSUBAKI_KAFKA
type Config struct {
	// If the producer should create the topic or not.
	//
	// Default: true | Variable: TSUBAKI_KAFKA_AUTO_CREATE_TOPICS
	AutoCreateTopics bool `yaml:"auto_create_topics"`

	// A list of brokers to connect to. This is a List of `host:port` strings.
	//
	// Default: []string{"localhost:9092"} | Variable: TSUBAKI_KAFKA_BROKERS
	Brokers []string `yaml:"brokers"`

	// Returns the topic to use when sending messages towards the GitHub bot (vice-versa).
	//
	// Warning: This must be the same topic you set from the GitHub bot configuration
	// or it will not receive messages!
	//
	// Default: "arisu:tsubaki" | Variable: TSUBAKI_KAFKA_TOPIC
	Topic string `yaml:"topic"`
}
