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

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var StatusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a Teapot",
	420: "Blaze It Bitches!", // not in rfc, just there for fun :3
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	509: "Bandwidth Limit Exceeded",
	510: "Not Exceeded",
	511: "Network Authentication Required",
}

func GetStatusCode(code int) string {
	// Get all keys
	keys := make([]int, 0, len(StatusCodes))
	for k := range StatusCodes {
		keys = append(keys, k)
	}

	// Iterate over them to find the value
	for _, key := range keys {
		if code == key {
			return StatusCodes[key]
		}
	}

	// If none was found, just return "Unknown"
	return "Unknown"
}

// WriteJson is a simple utility function to write data into JSON.
func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

// GetJsonBody is a simple utility function to retrieve this http.Request's
// body as a JSON object.
func GetJsonBody(req *http.Request) (int, map[string]interface{}, error) {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return http.StatusUnsupportedMediaType, nil, fmt.Errorf("content type was not application/json, received %s", contentType)
	}

	var data map[string]interface{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return 406, nil, fmt.Errorf("wrong type provided for field '%s'", unmarshalErr.Field)
		} else {
			return 400, nil, err
		}
	}

	return -1, data, nil
}
