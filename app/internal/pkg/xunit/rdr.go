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

import "encoding/xml"

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
	Name            string       `xml:"name,attr"`
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
