// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// sourcemap is a package which defines types and methods for working with V3 of the SourceMap
// spec: https://goo.gl/Bn7iTo
package sourcemap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// The supported sourcemap version.s
const VERSION = 3

// sourceMapData represents the data that will be marshalled/unmarshalled to/from JSON
// when forming a source map.
type sourceMapData struct {
	//  File version (always the first entry in the object) and must be a positive integer.
	Version int `json:"version"`

	// An optional name of the generated code that this source map is associated with.
	File string `json:"file,omitempty"`

	// An optional source root, useful for relocating source files on a server or
	// removing repeated values in the "sources" entry.  This value is prepended
	// to the individual entries in the "source" field.
	SourceRoot string `json:"sourceRoot,omitempty"`

	// A list of original sources used by the "mappings" entry.
	Sources []string `json:"sources"`

	// A list of symbol names used by the "mappings" entry.
	Names []string `json:"names"`

	// A string with the encoded mapping data.
	Mappings string `json:"mappings"`
}

// SourceMap is an in-memory representation of a source map being constructed.
type SourceMap struct {
	// The path of the generated source file.
	generatedFilePath string

	// The source root for the source map.
	sourceRoot string

	// lineMappings is a map from 0-indexed line number in the generated source
	// to a map of SourceMapping's, indexed by the 0-indexed column position
	// in the generated source.
	lineMappings map[int]map[int]SourceMapping

	// Set of names encountered.
	namesMap *orderedStringSet

	// Set of source paths encountered.
	sourcesMap *orderedStringSet
}

// ParsedSourceMap is an in-memory immutable representation of an existing source map.
type ParsedSourceMap struct {
	generatedFilePath string
	sourceRoot        string
	sources           []string
	names             []string
	lineMappings      []string

	computedMappings map[int]map[int]SourceMapping
}

// SourceMapping represents a mapping from a position in the input source file
// to an output source file.
type SourceMapping struct {
	// SourcePath is the path of the original source file for this position.
	SourcePath string

	// LineNumber is the 0-indexed line number in the original file.
	LineNumber int

	// ColumnPosition is the 0-indexed column position in the original file.
	ColumnPosition int

	// Name is the original name of the token, if any.
	Name string
}

// Parse parses a sourcemap as contained in the byte data into an in-memory source map.
func Parse(encoded []byte) (*ParsedSourceMap, error) {
	rep := sourceMapData{}

	// Unmarshal the sourcemap JSON into its data struct.
	err := json.Unmarshal(encoded, &rep)
	if err != nil {
		return &ParsedSourceMap{}, fmt.Errorf("Could not parse sourcemap: %v", err)
	}

	// Ensure we have a recognized version.
	if rep.Version != VERSION {
		return &ParsedSourceMap{}, fmt.Errorf("Unsupported sourcemap version: %v", rep.Version)
	}

	lines := strings.Split(rep.Mappings, ";")
	return &ParsedSourceMap{rep.File, rep.SourceRoot, rep.Sources, rep.Names, lines, nil}, nil
}

// NewSourceMap returns a new source map for construction.
func NewSourceMap(generatedFilePath string, sourceRoot string) *SourceMap {
	return &SourceMap{generatedFilePath, sourceRoot, map[int]map[int]SourceMapping{}, newOrderedStringSet(), newOrderedStringSet()}
}

// AddMapping adds a new mapping for the given generated line number and column position to this map.
func (sm *SourceMap) AddMapping(lineNumber int, colPosition int, mapping SourceMapping) {
	if _, ok := sm.lineMappings[lineNumber]; !ok {
		sm.lineMappings[lineNumber] = map[int]SourceMapping{}
	}

	if _, ok := sm.lineMappings[lineNumber][colPosition]; ok {
		return
	}

	sm.lineMappings[lineNumber][colPosition] = mapping

	if mapping.Name != "" {
		sm.namesMap.Add(mapping.Name)
	}

	if mapping.SourcePath != "" {
		sm.sourcesMap.Add(mapping.SourcePath)
	}
}

// OffsetBy returns a source map containing the mappings of this source map offset by the given string.
func (sm *SourceMap) OffsetBy(value string) *SourceMap {
	lines := strings.Split(value, "\n")

	om := NewSourceMap(sm.generatedFilePath, sm.sourceRoot)
	for lineNumber, mappings := range sm.lineMappings {
		for colPosition, mapping := range mappings {
			var updatedLineNumber = lineNumber + len(lines) - 1
			var updatedColPosition = colPosition
			if lineNumber == 0 {
				updatedColPosition = colPosition + len(lines[len(lines)-1])
			}

			om.AddMapping(updatedLineNumber, updatedColPosition, mapping)
		}
	}

	return om
}

// AppendMap adds all the mappings found in the other source map to this source map.
func (sm *SourceMap) AppendMap(other *SourceMap) {
	for lineNumber, mappings := range other.lineMappings {
		for colPosition, mapping := range mappings {
			sm.AddMapping(lineNumber, colPosition, mapping)
		}
	}
}

// Build returns the built source map.
func (sm *SourceMap) Build() *ParsedSourceMap {
	// Sort both sets to ensure consistent source map production.
	sm.namesMap.Sort()
	sm.sourcesMap.Sort()

	// Find the maximum line number.
	var maxLine = 0
	for lineNumber, _ := range sm.lineMappings {
		if lineNumber > maxLine {
			maxLine = lineNumber
		}
	}

	// Allocate the line mappings slice.
	lineMappings := make([]string, maxLine+1)

	// Indexes represent:
	// 0 - Column position in the generated file
	// 1 - 0-index in the source array
	// 2 - 0-index starting line in the original source
	// 3 - 0-index column position in the original source
	// 4 - 0-index into the names array
	globalValues := [5]int{0, 0, 0, 0, 0}

	// Calculate the line mappings for each slice.
	for lineNumber, _ := range lineMappings {
		mappings, hasMappings := sm.lineMappings[lineNumber]
		if !hasMappings {
			continue
		}

		globalValues[0] = 0 // Special case in source maps: Reset for every line, unlike the other fields.

		// Collect the col positions for the mappings and sort them by column position.
		columnPositions := make([]int, len(mappings))
		var index = 0
		for columnPosition, _ := range mappings {
			columnPositions[index] = columnPosition
			index++
		}

		sort.Ints(columnPositions)

		// Compute the mappings for this line.
		var buffer bytes.Buffer
		for _, columnPosition := range columnPositions {
			mapping := mappings[columnPosition]

			// Build the 5 values for the specific segment.
			unadjustedSegmentValues := [5]int{
				columnPosition,
				sm.sourcesMap.IndexOf(mapping.SourcePath),
				mapping.LineNumber,
				mapping.ColumnPosition,
				sm.namesMap.IndexOf(mapping.Name),
			}

			// Create the adjusted segment values. SourceMap V3 stores each
			// field integer value relative to the *first instance of that field
			// in a segment* (except the column position, which is automatically reset
			// on each line).
			adjustedSegmentValues := make([]int, 0, 5)
			for index, unadjustedSegmentValue := range unadjustedSegmentValues {
				if unadjustedSegmentValue < 0 {
					break
				}

				// Add the value adjusted based on the current global value.
				adjustedSegmentValues = append(adjustedSegmentValues, unadjustedSegmentValue-globalValues[index])

				// Set the global value to the *unadjusted value*.
				globalValues[index] = unadjustedSegmentValue
			}

			// Convert the adjusted segment values into VLQ and append the segment
			// to the buffer.
			buffer.WriteString(vlqEncode(adjustedSegmentValues))
			buffer.WriteRune(',')
		}

		lineMappingsString := buffer.String()
		lineMappings[lineNumber] = lineMappingsString[0 : len(lineMappingsString)-1] // Remove trailing ,
	}

	return &ParsedSourceMap{
		generatedFilePath: sm.generatedFilePath,
		sourceRoot:        sm.sourceRoot,
		sources:           sm.sourcesMap.OrderedItems(),
		names:             sm.namesMap.OrderedItems(),
		lineMappings:      lineMappings,
		computedMappings:  nil,
	}
}

// Marshal turns the parsed source map into its JSON form.
func (psm *ParsedSourceMap) Marshal() ([]byte, error) {
	data := sourceMapData{
		Version:    VERSION,
		File:       psm.generatedFilePath,
		SourceRoot: psm.sourceRoot,
		Sources:    psm.sources,
		Names:      psm.names,
		Mappings:   strings.Join(psm.lineMappings, ";"),
	}

	return json.MarshalIndent(data, "", "  ")
}

// computeMappings computes the full mappings map for all of the mappings found in the parsed
// source map.
func (psm *ParsedSourceMap) computeMappings() {
	// Indexes represent:
	// 0 - Column position in the generated file
	// 1 - 0-index in the source array
	// 2 - 0-index starting line in the original source
	// 3 - 0-index column position in the original source
	// 4 - 0-index into the names array
	mappedValues := [5]int{0, 0, 0, 0, 0}

	computed := map[int]map[int]SourceMapping{}

	for index, lineMapping := range psm.lineMappings {
		if lineMapping == "" {
			continue
		}

		computed[index] = map[int]SourceMapping{}
		mappedValues[0] = 0 // Special: Always reset for each line.
		segments := strings.Split(lineMapping, ",")

		for _, segment := range segments {
			// VLQ decode the segment.
			values, ok := vlqDecode(segment)
			if !ok {
				return
			}

			// Update the various positions and indices.
			for index, value := range values {
				mappedValues[index] += value
			}

			var sourcePath = ""
			var name = ""

			if mappedValues[1] >= 0 && mappedValues[1] < len(psm.sources) {
				sourcePath = psm.sources[mappedValues[1]]
			}

			if mappedValues[4] >= 0 && mappedValues[4] < len(psm.names) {
				name = psm.names[mappedValues[4]]
			}

			mapping := SourceMapping{
				SourcePath:     sourcePath,
				LineNumber:     mappedValues[2],
				ColumnPosition: mappedValues[3],
				Name:           name,
			}

			computed[index][mappedValues[0]] = mapping
		}
	}

	psm.computedMappings = computed
}

// LookupMapping returns the mapping for the given (0-indexed) line number and column position, if any.
func (psm *ParsedSourceMap) LookupMapping(lineNumber int, colPosition int) (SourceMapping, bool) {
	if psm.computedMappings == nil {
		psm.computeMappings()
	}

	lineMappings, ok := psm.computedMappings[lineNumber]
	if !ok {
		return SourceMapping{}, false
	}

	return lineMappings[colPosition], true
}
