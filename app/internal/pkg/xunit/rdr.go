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

// Package xunit contains functions for parsing XML files contains .NET test result(s) in xUnit's v2+ XML format.
// More information regarding this format can be found @ https://xunit.net/docs/format-xml-v2.
package xunit

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/kdeconinck/dotnet-test-visualizer/internal/pkg/camelcase"
)

// The .NET test result(s) in xUnit's v2+ XML format.
type Result struct {
	XMLName       xml.Name   `xml:"assemblies"`
	Computer      string     `xml:"computer,attr"`
	FinishRTF     string     `xml:"finish-rtf,attr"`
	Id            string     `xml:"id,attr"`
	SchemaVersion string     `xml:"schema-version,attr"`
	StartRTF      string     `xml:"start-rtf,attr"`
	Timestamp     string     `xml:"timestamp,attr"`
	User          string     `xml:"user,attr"`
	Assemblies    []Assembly `xml:"assembly"`
}

// A "test" assembly.
type Assembly struct {
	ConfigFile      string       `xml:"config-file,attr"`
	Environment     string       `xml:"environment,attr"`
	ErrorCount      int          `xml:"errors,attr"`
	FailedCount     int          `xml:"failed,attr"`
	FinishRTF       string       `xml:"finish-rtf,attr"`
	Id              string       `xml:"id,attr"`
	FullName        string       `xml:"name,attr"`
	NotRunCount     int          `xml:"not-run,attr"`
	PassedCount     int          `xml:"passed,attr"`
	RunDate         string       `xml:"run-date,attr"`
	RunTime         string       `xml:"run-time,attr"`
	SkippedCount    int          `xml:"skipped,attr"`
	StartRTF        string       `xml:"start-rtf,attr"`
	TargetFramework string       `xml:"target-framework,attr"`
	TestFramework   string       `xml:"test-framework,attr"`
	Time            float32      `xml:"time,attr"`
	TimeRTF         string       `xml:"time-rtf,attr"`
	Total           int          `xml:"total,attr"`
	Collections     []Collection `xml:"collection"`
	ErrorSet        ErrorSet     `xml:"errors"`
}

// Contains information about the run of a single test collection.
type Collection struct {
	Id           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FailedCount  int    `xml:"failed,attr"`
	NotRunCount  int    `xml:"not-run,attr"`
	PassedCount  int    `xml:"passed,attr"`
	SkippedCount int    `xml:"skipped,attr"`
	Time         string `xml:"time,attr"`
	TimeRTF      string `xml:"time-rtf,attr"`
	TotalCount   int    `xml:"total,attr"`
	Tests        []Test `xml:"test"`
}

// Contains information about the run of a single test.
type Test struct {
	Id         string     `xml:"id,attr"`
	Method     string     `xml:"method,attr"`
	Name       string     `xml:"name,attr"`
	Result     string     `xml:"result,attr"`
	SourceFile string     `xml:"source-file,attr"`
	SourceLine string     `xml:"source-line,attr"`
	Time       float32    `xml:"time,attr"`
	TimeRTF    string     `xml:"time-rtf,attr"`
	Type       string     `xml:"type,attr"`
	Failure    Failure    `xml:"failure"`
	Output     string     `xml:"output"`
	Reason     string     `xml:"reason"`
	TraitSet   TraitSet   `xml:"traits"`
	WarningSet WarningSet `xml:"warnings"`
}

// Contains information a test failure.
type Failure struct {
	ExceptionType string `xml:"exception-type,attr"`
	Message       string `xml:"message"`
	StackTrace    string `xml:"stack-trace"`
}

// A container for 1 or more "trait" elements.
type TraitSet struct {
	Traits []Trait `xml:"trait"`
}

// Contains a single trait name/value pair
type Trait struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// A container for 1 or more "warning" elements.
type WarningSet struct {
	Warnings []string `xml:"warning"`
}

// A container for 0 or more "error" elements.
type ErrorSet struct {
	Errors []Error `xml:"error"`
}

// Contains information about an environment failure that happened outside the scope of running a single unit test.
type Error struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Node struct {
	Name     string
	Tests    []Test
	Children []*Node
}

// Returns all the tests of an assembly.
func (assembly *Assembly) getTests() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		resultSet = append(resultSet, collection.Tests...)
	}

	return resultSet
}

// Returns the name (NO extension) of the assembly.
func (assembly *Assembly) getNameWithoutExtension() string {
	assemblyName := assembly.Name()

	return strings.TrimRight(assemblyName, ".dl")
}

func (assembly *Assembly) BuildTree() map[string][]*Node {
	mapSet := make(map[string][]*Node)

	root := &Node{Name: ""}
	resultSet := make([]*Node, 0)

	for _, item := range assembly.testsWithoutTraits() {
		currentNode := root

		if item.HasDisplayName() {
			currentNode.Tests = append(currentNode.Tests, item)
		} else if !item.HasDisplayName() && !item.isNested() {
			currentNode.Tests = append(currentNode.Tests, item)
		} else if item.isNested() {
			parts := make([]string, 0)
			fullNameParts := strings.Split(item.Name, "+")
			parts = append(parts, fullNameParts[0][strings.LastIndex(fullNameParts[0], ".")+1:])
			parts = append(parts, strings.Split(fullNameParts[len(fullNameParts)-1], ".")[0])

			for idx, part := range parts {
				var childNode *Node

				for _, child := range currentNode.Children {
					if child.Name == part {
						childNode = child
						break
					}
				}

				if childNode == nil {
					childNode = &Node{Name: part}
					if idx == len(parts)-1 {
						childNode.Tests = append(childNode.Tests, item)
					}

					currentNode.Children = append(currentNode.Children, childNode)
				}

				currentNode = childNode
			}
		}
	}

	resultSet = append(resultSet, root)
	mapSet[""] = resultSet

	for _, trait := range assembly.UniqueTraits() {
		root := &Node{Name: ""}
		resultSet := make([]*Node, 0)

		for _, item := range assembly.testsForTrait(trait) {
			currentNode := root

			if item.HasDisplayName() {
				currentNode.Tests = append(currentNode.Tests, item)
			} else if !item.HasDisplayName() && !item.isNested() {
				currentNode.Tests = append(currentNode.Tests, item)
			} else if item.isNested() {
				parts := make([]string, 0)
				fullNameParts := strings.Split(item.Name, "+")
				parts = append(parts, fullNameParts[0][strings.LastIndex(fullNameParts[0], ".")+1:])
				parts = append(parts, strings.Split(fullNameParts[len(fullNameParts)-1], ".")[0])

				for idx, part := range parts {
					var childNode *Node

					for _, child := range currentNode.Children {
						if child.Name == part {
							childNode = child
							break
						}
					}

					if childNode == nil {
						childNode = &Node{Name: part}
						if idx == len(parts)-1 {
							childNode.Tests = append(childNode.Tests, item)
						}

						currentNode.Children = append(currentNode.Children, childNode)
					}

					currentNode = childNode
				}
			}
		}

		resultSet = append(resultSet, root)
		mapSet[trait.Name+" - "+trait.Value] = resultSet
	}

	return mapSet
}

func PrintTree(node *Node, indent string) {
	fmt.Println(indent + node.Name + " - ")
	fmt.Println(len(node.Tests))

	for _, child := range node.Children {
		PrintTree(child, indent+"  ")
	}
}

// Name returns the name and the extension of the assembly.
func (assembly *Assembly) Name() string {
	if strings.Contains(assembly.FullName, "/") {
		return assembly.FullName[strings.LastIndex(assembly.FullName, "/")+1:]
	}

	return assembly.FullName[strings.LastIndex(assembly.FullName, "\\")+1:]
}

// Returns `true` if `test` has is a "nested" test, `false` otherwise.
func (test *Test) isNested() bool {
	return strings.Contains(test.Name, "+")
}

// HasDisplayName returns `true` if `test` has a display name, `false` otherwise.
func (test *Test) HasDisplayName() bool {
	return strings.Contains(test.Name, " ")
}

// UniqueTraits returns the unique traits over a single assembly.
func UniqueTraits(tests []Test) []Trait {
	resultSet := make([]Trait, 0)

	for _, test := range tests {
		for _, trait := range test.TraitSet.Traits {
			isTraitAlreadyFound := false
			for _, result := range resultSet {
				if result.Name == trait.Name && result.Value == trait.Value {
					isTraitAlreadyFound = true
					break
				}
			}

			if !isTraitAlreadyFound {
				resultSet = append(resultSet, trait)
			}
		}
	}

	return resultSet
}

// UniqueTraits returns the unique traits over a single assembly.
func (assembly *Assembly) UniqueTraits() []Trait {
	resultSet := make([]Trait, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			for _, trait := range test.TraitSet.Traits {
				isTraitAlreadyFound := false
				for _, result := range resultSet {
					if result.Name == trait.Name && result.Value == trait.Value {
						isTraitAlreadyFound = true
						break
					}
				}

				if !isTraitAlreadyFound {
					resultSet = append(resultSet, trait)
				}
			}
		}
	}

	return resultSet
}

// NonNestedTests returns all the NON nested test(s) of an assembly.
func (assembly *Assembly) NonNestedTests() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if !test.isNested() {
				resultSet = append(resultSet, test)
			}
		}

		resultSet = append(resultSet, collection.Tests...)
	}

	return resultSet
}

// NestedTests returns all the nested test(s) of an assembly.
func (assembly *Assembly) NestedTests() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if test.isNested() {
				resultSet = append(resultSet, test)
			}
		}

		resultSet = append(resultSet, collection.Tests...)
	}

	return resultSet
}

// NonNestedTestsWithoutTraits returns all the NON nested test(s) of an assembly NOT belonging to any trait.
func (assembly *Assembly) testsWithoutTraits() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if len(test.TraitSet.Traits) == 0 {
				resultSet = append(resultSet, test)
			}
		}
	}

	return resultSet
}

// NonNestedTestsWithoutTraits returns all the NON nested test(s) of an assembly NOT belonging to any trait.
func (assembly *Assembly) NonNestedTestsWithoutTraits() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if !test.isNested() && len(test.TraitSet.Traits) == 0 {
				resultSet = append(resultSet, test)
			}
		}
	}

	return resultSet
}

// NestedTestsWithoutTraits returns all the nested test(s) of an assembly NOT belonging to any trait.
func (assembly *Assembly) NestedTestsWithoutTraits() []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if test.isNested() && len(test.TraitSet.Traits) == 0 {
				resultSet = append(resultSet, test)
			}
		}
	}

	return resultSet
}

// Returns all the test(s) that does belong to the given trait.
func (assembly *Assembly) testsForTrait(trait Trait) []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			for _, testTrait := range test.TraitSet.Traits {
				if testTrait.Name == trait.Name && testTrait.Value == trait.Value {
					resultSet = append(resultSet, test)
				}
			}
		}
	}

	return resultSet
}

// NonNestedTestsForTrait returns all the NON nested test(s) of an assembly belonging to the given trait.
func (assembly *Assembly) NonNestedTestsForTrait(trait Trait) []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if !test.isNested() {
				for _, testTrait := range test.TraitSet.Traits {
					if testTrait.Name == trait.Name && testTrait.Value == trait.Value {
						resultSet = append(resultSet, test)
					}
				}
			}
		}
	}

	return resultSet
}

// NestedTestsForTrait returns all the nested test(s) of an assembly belonging to the given trait.
func (assembly *Assembly) NestedTestsForTrait(trait Trait) []Test {
	resultSet := make([]Test, 0)

	for _, collection := range assembly.Collections {
		for _, test := range collection.Tests {
			if test.isNested() {
				for _, testTrait := range test.TraitSet.Traits {
					if testTrait.Name == trait.Name && testTrait.Value == trait.Value {
						resultSet = append(resultSet, test)
					}
				}
			}
		}
	}

	return resultSet
}

// // String returns the human-readble version of test.
// func (test *Test) String() string {
// 	if test.Result == "Pass" {
// 		return fmt.Sprintf("\033[1;32m ✓ \033[0m %s (%v seconds)", test.FriendlyName(), test.Time)
// 	}

// 	return fmt.Sprintf("\033[1;31m ⛌ \033[0m %s (%v seconds)", test.FriendlyName(), test.Time)
// }

// TestName returns the name of the test.
func (test *Test) TestName() string {
	if test.HasDisplayName() {
		return test.Name
	}

	return getFriendlyName(test.Name[strings.LastIndex(test.Name, ".")+1:])
}

func getFriendlyName(inputString string) string {
	return strings.Join(camelcase.Split(inputString), " ")
}

// FriendlyName returns the human-reable version of the name of the test.
func (test *Test) FriendlyName() string {
	return strings.Join(camelcase.Split(test.TestName()), " ")
}
