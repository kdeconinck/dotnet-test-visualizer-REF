// =====================================================================================================================
// = LICENSE:       Copyright (c) 2023 Kevin De Coninck
// =
// =                Permission is hereby granted, free of charge, to any person
// =                obtaining a copy of this software and associated documentation
// =                files (the "Software"), to deal in the Software without
// =                restriction, including without limitation the rights to use,
// =                copy, modify, merge, publish, distribute, sublicense, and/or sell
// =                copies of the Software, and to permit persons to whom the
// =                Software is furnished to do so, subject to the following
// =                conditions:
// =
// =                The above copyright notice and this permission notice shall be
// =                included in all copies or substantial portions of the Software.
// =
// =                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// =                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// =                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// =                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// =                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// =                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// =                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// =                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Package args defines functions for working with the arguments of the application.
package args

import (
	"fmt"
	"os"
)

// FindNamedSingle returns the value of a "named" argument if it's found.
// It returns a NON <nil> error if either the "named" argument hasn't been found, or when the "named" argument doesn't
// have a value.
func FindNamedSingle(key string) (string, error) {
	var args = os.Args[1:]

	for idx, arg := range args {
		if arg == key && idx >= len(args)-1 {
			return "", fmt.Errorf("no value found for arg '%s'", key)
		}

		if arg == key && len(args) >= idx {
			valuePos := idx + 1

			return args[valuePos], nil
		}
	}

	return "", fmt.Errorf("arg '%s' not found", key)
}

// FindNamed returns the value(s) of a "named" argument if it's found.
// It returns a NON <nil> error if either the "named" argument hasn't been found, or when any of the "named" arguments
// doesn't have a value.
func FindNamed(key string) ([]string, error) {
	var args = os.Args[1:]
	var result = make([]string, 0)

	for idx, arg := range args {
		if arg == key && idx >= len(args)-1 {
			return result, fmt.Errorf("no value found for arg '%s'", key)
		}

		if arg == key && len(args) >= idx {
			valuePos := idx + 1

			result = append(result, args[valuePos])
		}
	}

	if len(result) == 0 {
		return result, fmt.Errorf("arg '%s' not found", key)
	}

	return result, nil
}
