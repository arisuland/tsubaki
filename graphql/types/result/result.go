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

// Result is a object that is represented as a result
// of an action. This comes with a `success` and `errors`
// properties. If the `success=true`, the action's result
// was a success. If the `success` value is `false`, then
// something happened (using `result.Err()`, `result.Errs`, or `result.ErrWithMessage`) with detailed
// errors on what happened.
type Result struct {
	// Success represents if the result of this action
	// was successful or not.
	Success bool `json:"success"`

	// Errors represents a list of Error objects
	// if `success=false`.
	Errors []Error `json:"errors"`
}

// Error is a object that is represented as a error
// that went wrong with a Result.
type Error struct {
	// Message is the underlying message if a Error is provided.
	Message string `json:"message"`

	// Code is the underlying error code, if this is a generic
	// `error`, this will be `-1`. Otherwise, read all
	// the error codes: https://docs.arisu.land/graphql/types/result#error-codes
	Code int32 `json:"code"`
}

// Ok returns a successful Result.
func Ok() Result {
	return Result{
		Success: true,
		Errors:  []Error{},
	}
}

// Err takes a underlying error and transforms it to a Result object.
// Use `ErrWithMessage` to use your own error message.
func Err(err []error) Result {
	errors := make([]Error, 0)
	for _, e := range err {
		errors = append(errors, Error{
			Message: e.Error(),
			Code:    -1,
		})
	}

	return Result{
		Success: false,
		Errors:  errors,
	}
}

// ErrWithMessage returns a underlying error with your own message component.
// This is only for singular errors. Use `Errs` to add multiple errors.
func ErrWithMessage(message string, code int32) Result {
	errors := []Error{
		{
			Message: message,
			Code:    code,
		},
	}

	return Result{
		Success: false,
		Errors:  errors,
	}
}

// Errs returns a Result object with x amount of `errors` listed.
func Errs(errors []Error) Result {
	return Result{
		Success: false,
		Errors:  errors,
	}
}
