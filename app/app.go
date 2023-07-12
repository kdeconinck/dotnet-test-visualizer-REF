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

// Package main implements Visualizo, a CLI application for visualizing .NET test result(s).
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/kdeconinck/dotnet-test-visualizer/internal/pkg/camelcase"
	"github.com/kdeconinck/dotnet-test-visualizer/internal/pkg/xunit"
)

// // First, we loop over all the test(s) that doesn't contain any nesting.
// for _, test := range resultTree.Tests {
// 	status := "\033[1;32m✓\033[0m"
// 	if test.Result != "Pass" {
// 		status = "\033[1;31m⛌\033[0m"
// 	}

// 	if test.Time <= tresholdFast {
// 		fmt.Printf("    🚀 %s %s. (%v seconds)\r\n", status, test.TestName(), test.Time)
// 	} else if test.Time <= tresholdNormal {
// 		fmt.Printf("    🕐 %s %s. (%v seconds)\r\n", status, test.TestName(), test.Time)
// 	} else {
// 		fmt.Printf("    🐌 %s %s. (%v seconds)\r\n", status, test.TestName(), test.Time)
// 	}
// }

// These are constants that should be passed using a configuration file but they are hard-coded right now.
const (
	tresholdFast   float32 = 0.05
	tresholdNormal float32 = 0.1
)

// The `main` entry point for the application.
func main() {
	// Prints the ASCII header.
	fmt.Println("    _  _ ___ _____   _____       _    __   ___              _ _            ")
	fmt.Println("   | \\| | __|_   _| |_   _|__ __| |_  \\ \\ / (_)____  _ __ _| (_)______ _ _ ")
	fmt.Println("  _| .` | _|  | |     | |/ -_|_-<  _|  \\ V /| (_-< || / _` | | |_ / -_) '_|")
	fmt.Println(" (_)_|\\_|___| |_|     |_|\\___/__/\\__|   \\_/ |_/__/\\_,_\\__,_|_|_/__\\___|_|  ")

	// Load the file that contains the xUnit test result(s).
	xUnitResults := loadXunitTestResults("results2.xml")

	fmt.Println("")

	// Overal information.
	fmt.Printf("Amount of assemblies: %v\r\n", len(xUnitResults.Assemblies))

	if xUnitResults.Computer != "" {
		fmt.Printf("Computer:             %s\r\n", xUnitResults.Computer)
	}

	if xUnitResults.User != "" {
		fmt.Printf("User:                 %s\r\n", xUnitResults.User)
	}

	if xUnitResults.StartRTF != "" {
		fmt.Printf("Start time:           %s\r\n", xUnitResults.StartRTF)
	}

	if xUnitResults.FinishRTF != "" {
		fmt.Printf("End time:             %s\r\n", xUnitResults.FinishRTF)
	} else if xUnitResults.Timestamp != "" {
		fmt.Printf("End time:             %s\r\n", xUnitResults.Timestamp)
	}

	// Loop over the assemblies and print the results accordingly.
	for _, assembly := range xUnitResults.Assemblies {
		fmt.Println("")
		fmt.Printf("Assembly:         %s", assembly.Name())

		if assembly.FailedCount != 0 {
			fmt.Printf(" - \033[1;31m⛌ Failed (%v of %v failed).\033[0m\r\n", assembly.FailedCount, assembly.Total)
		} else {
			fmt.Printf(" - \033[1;32m✓ Passed (%v of %v passed).\033[0m \r\n", assembly.PassedCount, assembly.Total)
		}

		fmt.Printf("Date / time:      %s %s\r\n", assembly.RunDate, assembly.RunTime)
		fmt.Printf("Total time:       %v seconds.\r\n", assembly.Time)

		// Print information about the assembly.
		fmt.Println("")
		fmt.Printf("  # tests:        %v\r\n", assembly.Total)
		fmt.Printf("  # Passed tests: %v\r\n", assembly.PassedCount)
		fmt.Printf("  # Failed tests: %v\r\n", assembly.FailedCount)
		fmt.Printf("  # Errors:       %v\r\n", assembly.ErrorCount)
		fmt.Println("")

		// Build the tree which contains all the test(s).
		resultTree := assembly.BuildTree()

		// Sort the tree because Go doesn't guarantee the order of the elements in a map when iterating.
		keys := make([]string, 0)
		for k := range resultTree {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Loop over all elements in the map (in order).
		for _, key := range keys {
			for _, node := range resultTree[key] {
				if key != "" {
					fmt.Printf("  Trait: %s\r\n", key)
				}

				// for _, node := range nodes {
				for _, test := range node.Tests {
					status := "\033[1;32m✓\033[0m"
					if test.Result != "Pass" {
						status = "\033[1;31m⛌\033[0m"
					}

					suffix := ""

					if key != "" {
						suffix = "  "
					}

					if test.Time <= tresholdFast {
						fmt.Printf("%s  🚀 %s %s (%v seconds)\r\n", suffix, status, test.TestName(), test.Time)
					} else if test.Time <= tresholdNormal {
						fmt.Printf("%s  🕐 %s %s (%v seconds)\r\n", suffix, status, test.TestName(), test.Time)
					} else {
						fmt.Printf("%s  🐌 %s %s (%v seconds)\r\n", suffix, status, test.TestName(), test.Time)
					}
				}

				// // Travel over all the nested test(s).
				// for _, node := range node.Children {
				// 	fmt.Println("")
				// 	PrintTree(node, "")
				// }
			}
		}
	}

	fmt.Println("")
}

func PrintTree(node *xunit.Node, indent string) {
	fmt.Printf("%s  %s.\r\n", indent, getFriendlyName(node.Name))
	for _, test := range node.Tests {
		status := "\033[1;32m✓\033[0m"
		if test.Result != "Pass" {
			status = "\033[1;31m⛌\033[0m"
		}

		if test.Time <= tresholdFast {
			fmt.Printf("%s     🚀 %s %s (%v seconds)\r\n", indent, status, test.TestName(), test.Time)
		} else if test.Time <= tresholdNormal {
			fmt.Printf("%s     🕐 %s %s (%v seconds)\r\n", indent, status, test.TestName(), test.Time)
		} else {
			fmt.Printf("%s     🐌 %s %s (%v seconds)\r\n", indent, status, test.TestName(), test.Time)
		}
	}

	if len(node.Tests) > 0 {
		fmt.Println("")
	}

	for _, child := range node.Children {
		PrintTree(child, indent+"  ")
	}
}

func getFriendlyName(inputString string) string {
	return strings.Join(camelcase.Split(inputString), " ")
}

// // Displays the results of the given tests.
// func displayTestResults(tests []xunit.Test) {
// 	for _, test := range tests {
// 		if test.Time <= tresholdFast {
// 			fmt.Printf("    🚀 %s\r\n", test.String())
// 		} else if test.Time <= tresholdNormal {
// 			fmt.Printf("    🕐 %s\r\n", test.String())
// 		} else {
// 			fmt.Printf("    🐌 %s\r\n", test.String())
// 		}
// 	}
// }

// // Displays the results of the given tests as a tree.
// func displayTreeTestResults(tests []xunit.Test) {
// 	for _, test := range tests {
// 		fmt.Println(test.Name)
// 		// if test.Time <= tresholdFast {
// 		// 	fmt.Printf("    🚀 %s\r\n", test.String())
// 		// } else if test.Time <= tresholdNormal {
// 		// 	fmt.Printf("    🕐 %s\r\n", test.String())
// 		// } else {
// 		// 	fmt.Printf("    🐌 %s\r\n", test.String())
// 		// }
// 	}
// }

// Load the xUnit test result(s) file from filePath and return into a struct representing these result(s).
func loadXunitTestResults(filePath string) xunit.Result {
	var result xunit.Result

	rdr, err := os.Open(filePath)

	// TODO: Handle the error gracefully.
	if err != nil {
		fmt.Println(err)
	}

	// Close the XML file as the file has been read.
	defer rdr.Close()

	xml.Unmarshal(readXunitTestResults(rdr), &result)

	return result
}

func readXunitTestResults(rdr io.Reader) []byte {
	bytes, err := io.ReadAll(rdr)

	// TODO: Handle the error gracefully.
	if err != nil {
		fmt.Println(err)
	}

	return bytes
}
