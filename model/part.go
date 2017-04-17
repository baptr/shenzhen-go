// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"html/template"
	"net/http"

	"github.com/google/shenzhen-go/model/parts"
	"github.com/google/shenzhen-go/model/pin"
)

// Part abstracts the implementation of a node. Concrete implementations should be
// able to be marshalled to and unmarshalled from JSON sensibly.
type Part interface {
	// Clone returns a copy of this part.
	Clone() interface{}

	// Impl returns Go source code implementing the part.
	// The head is executed, then the body is executed (# Multiplicity
	// instances of the body concurrently), then the tail (once the body/bodies
	// are finished).
	//
	// This allows cleanly closing channels for nodes with Multiplicity > 1.
	// The tail is deferred so that the body can use "return" and it is still
	// executed.
	Impl() (head, body, tail string)

	// Imports returns any extra import lines needed for the Part.
	Imports() []string

	// Pins returns any pins - "channel arguments" - to the part.
	// inputs and outputs map argument names to types (the "<-chan" /
	// "chan<-" part of the type is implied).
	Pins() []pin.Definition

	// TypeKey returns the "type" of part.
	TypeKey() string

	// Update sets fields in the part based on info in the given Request.
	Update(*http.Request) error
}

// PartType has metadata common to this type of part, and is also a part factory.
// The HTML is loaded with the editor.
type PartType struct {
	New    func() Part
	Panels []struct {
		Name   string
		Editor template.HTML
	}
}

// PartTypes translates part type strings into useful information.
var PartTypes = map[string]*PartType{
	/*	"Aggregator":     func() Part { return new(parts.Aggregator) },
		"Broadcast":      func() Part { return new(parts.Broadcast) }, */
	"Code": {
		New:    func() Part { return new(parts.Code) },
		Panels: parts.CodePanels,
	},
	/*	"Filter":         func() Part { return new(parts.Filter) },
		"HTTPServer":     func() Part { return new(parts.HTTPServer) },
		"StaticSend":     func() Part { return new(parts.StaticSend) },
		"TextFileReader": func() Part { return new(parts.TextFileReader) },
		"Unslicer":       func() Part { return new(parts.Unslicer) },*/
}
