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

package result

// Result represents a response of the action that was executed. This is used
// in the database controllers.
type Result struct {
	// Success determines if this Result was a success.
	Success bool `json:"success"`

	// Data returns the underlying data that was successful,
	// this can be empty if Result.Errors are nil.
	Data interface{} `json:"data,omitempty"`

	// StatusCode returns the status code to use for this Result object.
	StatusCode int `json:"-"` // this shouldn't be in the JSON object when sent to the end user.

	// Errors returns the underlying Error object of what happened.
	// This is usually used in the result.Err() or result.Errs()
	// function fields.
	Errors []Error `json:"errors,omitempty"`
}

// Error represents the error that occurred in the resulted action.
type Error struct {
	// Code represents the error code that is used, you can read up
	// on all the error codes here: https://docs.arisu.land/api/reference#error-codes
	Code string `json:"code"`

	// Message is a brief message of what happened.
	Message string `json:"message"`
}

// Ok returns a Result object with the data attached.
func Ok(data interface{}) *Result {
	return OkWithStatus(200, data)
}

// OkWithStatus returns a Result object with a different status code
// rather than 200 OK.
func OkWithStatus(status int, data interface{}) *Result {
	return &Result{
		StatusCode: status,
		Success:    true,
		Data:       data,
	}
}

// NoContent a result object using the 201 status code.
func NoContent() *Result {
	return &Result{
		StatusCode: 204,
	}
}

// Err returns a Result object with any error that occurred.
func Err(status int, code string, message string) *Result {
	return &Result{
		StatusCode: status,
		Success:    false,
		Errors: []Error{
			NewError(code, message),
		},
	}
}

// Errs returns a Result object for multiple errors that might've occurred.
func Errs(status int, errors ...Error) *Result {
	return &Result{
		StatusCode: status,
		Errors:     errors,
		Success:    false,
	}
}

// NewError constructs a new Error object.
func NewError(code string, message string) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}
