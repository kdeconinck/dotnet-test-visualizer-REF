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

// Quality assurance: Validate and measure the performance of the public API of the "args" package.
package args_test

import (
	"errors"
	"os"
	"testing"

	"github.com/kdeconinck/dotnet-test-visualizer/internal/pkg/args"
)

// UT: Test if you can find a named argument.
//
// NOTE: These test(s) can't be executed in parallel since they modify a "global" variable.
func TestFindNamed(t *testing.T) {
	// Execution.
	for utName, utData := range map[string]struct {
		args           []string
		arg            string
		expectedValues []string
		expectedError  error
	}{
		"When NO arguments are provided, a NON <nil> error is returned.": {
			args: []string{
				"app", // This represents the name of the application.
			},
			arg:            "--port",
			expectedValues: []string{},
			expectedError:  errors.New("arg '--port' not found"),
		},

		"When the named argument is NOT found, an error is returned.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--logger",
			},
			arg:            "--port",
			expectedValues: []string{},
			expectedError:  errors.New("arg '--port' not found"),
		},

		"When NO value for the named argument is found, an error is returned.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--output",
			},
			arg:            "--output",
			expectedValues: []string{},
			expectedError:  errors.New("no value found for arg '--output'"),
		},

		"When a value for the named argument is found, that value is returned.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--port", "8080",
			},
			arg: "--port",
			expectedValues: []string{
				"8080",
			},
			expectedError: nil,
		},

		"When multiple values for the named argument is found, these values are returned.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--port", "8080", "--port", "9090",
			},
			arg: "--port",
			expectedValues: []string{
				"8080",
				"9090",
			},
			expectedError: nil,
		},
	} {
		utName, utData := utName, utData

		t.Run(utName, func(t *testing.T) {
			// There might be other test(s) relying on the arguments, so once the test is complete, we need to restore
			// the values that we havechanged previously.
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Arrange.
			os.Args = utData.args

			// Act.
			result, err := args.FindNamed(utData.arg)

			// Assert.
			if len(result) != len(utData.expectedValues) {
				t.Fatalf("\n\n"+
					"UT Name:  %s\n"+
					"\033[32mExpected: %d argument(s).\033[0m\n"+
					"\033[31mActual:   %d argument(s).\033[0m\n\n", utName, len(utData.expectedValues), len(result))
			} else {
				for idx := 0; idx < len(result); idx++ {
					if result[idx] != utData.expectedValues[idx] {
						t.Fatalf("\n\n"+
							"UT Name:            %s\n"+
							"\033[32mExpected (index %d): %s\033[0m\n"+
							"\033[31mActual (index %d):   %s\033[0m\n\n", utName, idx, utData.expectedValues[idx], idx,
							result[idx])
					}
				}
			}

			if err != nil && utData.expectedError == nil {
				t.Fatalf("\n\n"+
					"UT Name:          %s\n"+
					"\033[32mExpected (error): <nil>\033[0m\n"+
					"\033[31mActual (error):   %s\033[0m\n\n", utName, err)
			} else if err == nil && utData.expectedError != nil {
				t.Fatalf("\n\n"+
					"UT Name:          %s\n"+
					"\033[32mExpected (error): %s\033[0m\n"+
					"\033[31mActual (error):   <nil>\033[0m\n\n", utName, utData.expectedError)
			} else if err != nil && utData.expectedError != nil && err.Error() != utData.expectedError.Error() {
				t.Fatalf("\n\n"+
					"UT Name:          %s\n"+
					"\033[32mExpected (error): %s\033[0m\n"+
					"\033[31mActual (error):   %s\033[0m\n\n", utName, utData.expectedError, err)
			}
		})
	}
}

// Benchmark: Find a named argument.
func BenchmarkFindNamed(b *testing.B) {
	// Execution.
	for utName, utData := range map[string]struct {
		args []string
		arg  string
	}{
		"Find a named argument when NO arguments are provided.": {
			args: []string{
				"app", // This represents the name of the application.
			},
			arg: "--port",
		},

		"Find a named argument that's NOT present.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--logger",
			},
			arg: "--port",
		},

		"Find a named argument that doesn't have a value.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--output",
			},
			arg: "--output",
		},

		"Find a named argument that does have a value.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--output", "disk",
			},
			arg: "--output",
		},

		"Find a named argument that has multiple values.": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--port", "8080", "--port", "9090",
			},
			arg: "--port",
		},

		"Find a named argument that has multiple values (Larger dataset).": {
			args: []string{
				"app", // This represents the name of the application.

				// From this point forward, these are the arguments that are passed to the application.
				"--port", "1010", "--port", "2020", "--port", "3030", "--port", "4040", "--port", "5050", "--port",
				"6060", "--port", "7070", "--port", "8080", "--port", "9090",
			},
			arg: "--port",
		},
	} {
		utName, utData := utName, utData

		b.Run(utName, func(b *testing.B) {
			// There might be other test(s) relying on the arguments, so once the test is complete, we need to restore
			// the values that we havechanged previously.
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Arrange.
			os.Args = utData.args

			// Reset the timer to ensure that the benchmarks are accurate.
			b.ResetTimer()

			// Run.
			for n := 0; n < b.N; n++ {
				_, _ = args.FindNamed(utData.arg)
			}
		})
	}
}
