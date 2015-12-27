// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/serulian/compiler/compilergraph"
)

// Value string example:
//
// [9242167c-2d57-4212-8aea-fedf32bd708e]-T00010000G000039[04f97d44-d8fc-4a7d-9c46-955d1bd5add6]F
//  ^ Type ID of SomeType                ^ Special
//										  ^ Nullable
//                                         ^ 0001 = 1 generic
//                                             ^ 0000 = 0 parameters
//                                                 ^ G000039 = generic with 39 chars in length for subreference
// represents
//
// SomeType<SomeGeneric>?
//
// where SomeType is 9242167c-2d57-4212-8aea-fedf32bd708e
// and SomeGeneric is 04f97d44-d8fc-4a7d-9c46-955d1bd5add6
type typeReferenceHeaderSlot struct {
	index     int  // The slot index.
	length    int  // The length of this slot.
	writeable bool // Whether this slot is writeable.
}

var (
	trhSlotStartTypeId        = typeReferenceHeaderSlot{0, 1, false}
	trhSlotTypeId             = typeReferenceHeaderSlot{1, compilergraph.NodeIDLength, true}
	trhSlotEndTypeId          = typeReferenceHeaderSlot{2, 1, false}
	trhSlotFlagSpecial        = typeReferenceHeaderSlot{3, 1, false}
	trhSlotFlagNullable       = typeReferenceHeaderSlot{4, 1, true}
	trhSlotGenericCount       = typeReferenceHeaderSlot{5, 4, true}
	trhSlotParameterCount     = typeReferenceHeaderSlot{6, 4, true}
	trhSlotSubReferenceMarker = typeReferenceHeaderSlot{7, 0, false}
)

var typeReferenceHeaderSlots = [...]typeReferenceHeaderSlot{
	trhSlotStartTypeId,
	trhSlotTypeId,
	trhSlotEndTypeId,
	trhSlotFlagSpecial,
	trhSlotFlagNullable,
	trhSlotGenericCount,
	trhSlotParameterCount,
	trhSlotSubReferenceMarker,
}

const (
	subReferenceGenericChar   = 'G' // The prefix for generic sub references.
	subReferenceParameterChar = 'P' // The prefix for parameter sub references.

	nullableFlagTrue  = 'T' // The value of the trhSlotFlagNullable when the typeref is nullable.
	nullableFlagFalse = 'F' // The value of the trhSlotFlagNullable when the typeref is nullable.

	specialFlagNormal = '-' // The value of the trhSlotFlagSpecial for normal typerefs.
	specialFlagAny    = 'A' // The value of the trhSlotFlagSpecial for "any" type refs.
	specialFlagVoid   = 'V' // The value of the trhSlotFlagSpecial for "void" tyoe refs.
	specialFlagLocal  = 'L' // The value of the trhSlotFlagSpecial for localized type refs.
	specialFlagNull   = 'N' // The value of the trhSlotFlagSpecial for "null" type refs.
)

// typeReferenceHeaderSlotCacheMap holds a cache for looking up the offset of a TRH.
var typeReferenceHeaderSlotCacheMap = map[typeReferenceHeaderSlot]int{}

// The size of the length prefix for subreferences.
const typeRefValueSubReferenceLength = 6

// subReferenceKind represents the kinds of supported subreferences
type subReferenceKind int

const (
	subReferenceGeneric subReferenceKind = iota
	subReferenceParameter
)

// withFlag returns a copy of this type reference with the flag at the given slot replaced with the
// specified rune.
func (tr *TypeReference) withFlag(flagSlot typeReferenceHeaderSlot, value rune) TypeReference {
	return tr.replaceSlot(flagSlot, string(value))
}

// replaceSlot replaces the given slot with the specified value.
func (tr *TypeReference) replaceSlot(slot typeReferenceHeaderSlot, value string) TypeReference {
	if !slot.writeable {
		panic(fmt.Sprintf("Cannot write to slot %v", slot))
	}

	location := getSlotLocation(slot)
	size := slot.length

	if len(value) < size {
		value = value + strings.Repeat(" ", size-len(value))
	}

	return TypeReference{
		tdg:   tr.tdg,
		value: tr.value[:location] + value + tr.value[location+size:],
	}
}

// getSlotAsInt returns the int found at the given slot.
func (tr *TypeReference) getSlotAsInt(slot typeReferenceHeaderSlot) int {
	strValue := tr.getSlot(slot)

	// Special common case.
	if strValue == "0000" {
		return 0
	}

	i, err := strconv.Atoi(strValue)
	if err != nil {
		panic(fmt.Sprintf("Expected int value for slot %v, found: %v", slot, strValue))
	}
	return i
}

// getSlot returns the string found at the given slot.
func (tr *TypeReference) getSlot(slot typeReferenceHeaderSlot) string {
	location := getSlotLocation(slot)
	size := slot.length

	return tr.value[location : location+size]
}

// lengthPrefixedValue returns this type reference's value string, prefixed with its length value
// padded to the correct number of digits.
func (tr *TypeReference) lengthPrefixedValue() string {
	return padNumberToString(len(tr.value), typeRefValueSubReferenceLength) + tr.value
}

// withSubReference returns a copy of this type reference with the given subreference added.
func (tr *TypeReference) withSubReference(kind subReferenceKind, subref TypeReference) TypeReference {
	slot, kindRune := getSubReferenceSlotAndChar(kind)
	count := tr.getSlotAsInt(slot)

	// Add to the count and append the subreference.
	updated := tr.replaceSlot(slot, padNumberToString(count+1, slot.length))
	return TypeReference{
		tdg:   tr.tdg,
		value: updated.value + string(kindRune) + subref.lengthPrefixedValue(),
	}
}

// getSubReferences returns all sub references of the given kind on this type reference.
func (tr *TypeReference) getSubReferences(kind subReferenceKind) []TypeReference {
	slot, kindChar := getSubReferenceSlotAndChar(kind)

	// Find the number of applicable subreferences.
	count := tr.getSlotAsInt(slot)
	subrefs := make([]TypeReference, count)

	if count == 0 {
		return subrefs
	}

	// For all subreferences loop and filter.
	var currentIndex int = getSlotLocation(trhSlotSubReferenceMarker)
	var collectedCount int = 0

	for {
		// Retrieve the character for the subreference (G or P)
		subReferenceChar := tr.value[currentIndex]

		// Retrieve the length of the subreference's value string.
		subReferenceLengthStr := tr.value[currentIndex+1 : currentIndex+1+typeRefValueSubReferenceLength]
		subReferenceLength, err := strconv.Atoi(subReferenceLengthStr)
		if err != nil {
			panic(fmt.Sprintf("Expected int value for subreference length, found: %v", subReferenceLengthStr))
		}

		// Move the current index forward to the point at which the subreference value string begins.
		currentIndex = currentIndex + 1 + typeRefValueSubReferenceLength
		subReferenceValue := tr.value[currentIndex : currentIndex+subReferenceLength]

		// If we have found a subreference of the known kind, add it.
		if subReferenceChar == kindChar {
			subrefs[collectedCount] = TypeReference{
				tdg:   tr.tdg,
				value: subReferenceValue,
			}

			collectedCount = collectedCount + 1
		}

		// Once we've found all the subreferences of the expected kind, we're done.
		if collectedCount == count {
			return subrefs
		}

		// Move the current index past the subreference value string.
		currentIndex = currentIndex + subReferenceLength
	}
}

// getSubReferenceSlotAndChar returns the count slot and character for the given kind of subreference.
func getSubReferenceSlotAndChar(kind subReferenceKind) (typeReferenceHeaderSlot, uint8) {
	switch kind {
	case subReferenceGeneric:
		return trhSlotGenericCount, subReferenceGenericChar

	case subReferenceParameter:
		return trhSlotParameterCount, subReferenceParameterChar

	default:
		panic("Unknown kind of subreference")
		return trhSlotGenericCount, '_'
	}
}

// getSlotLocation returns the slot location (0-indexed) in the value string.
func getSlotLocation(slot typeReferenceHeaderSlot) int {
	if location, ok := typeReferenceHeaderSlotCacheMap[slot]; ok {
		return location
	}

	var location int
	for _, currentSlot := range typeReferenceHeaderSlots {
		if currentSlot == slot {
			typeReferenceHeaderSlotCacheMap[slot] = location
			return location
		}

		location = location + currentSlot.length
	}

	panic("Could not retrieve location for slot")
	return -1
}

// buildTypeReferenceValue returns a string value for representing the given type reference data in a single
// string.
func buildTypeReferenceValue(typeNode compilergraph.GraphNode, nullable bool, generics ...TypeReference) string {
	var buffer bytes.Buffer

	// Referenced type ID.
	buffer.WriteByte('[')
	buffer.WriteString(string(typeNode.NodeId))
	buffer.WriteByte(']')

	// Special: This is a normal type reference.
	buffer.WriteByte(specialFlagNormal)

	// Nullable.
	if nullable {
		buffer.WriteByte(nullableFlagTrue)
	} else {
		buffer.WriteByte(nullableFlagFalse)
	}

	// Generic count and parameter count.
	buffer.WriteString(padNumberToString(len(generics), trhSlotGenericCount.length))
	buffer.WriteString(padNumberToString(0, trhSlotParameterCount.length))

	if len(generics) == 0 {
		return buffer.String()
	}

	// For each generic, add the length of its value string and then the value string itself.
	for _, generic := range generics {
		buffer.WriteByte(subReferenceGenericChar)
		buffer.WriteString(padNumberToString(len(generic.value), typeRefValueSubReferenceLength))
		buffer.WriteString(generic.value)
	}

	return buffer.String()
}

// buildLocalizedRefValue returns a string value for representing a localized type reference.
func buildLocalizedRefValue(genericNode compilergraph.GraphNode) string {
	var buffer bytes.Buffer

	genericKind := genericNode.Get(NodePredicateGenericKind)
	genericIndex := genericNode.Get(NodePredicateGenericIndex)

	localId := "local:" + genericKind + ":" + genericIndex

	// Referenced type ID. For localized types, this is the generic index and kind.
	buffer.WriteByte('[')
	buffer.WriteString(localId)
	buffer.WriteString(strings.Repeat(" ", trhSlotTypeId.length-len(localId)))
	buffer.WriteByte(']')

	// Flags: special and nullable.
	buffer.WriteByte(specialFlagLocal)
	buffer.WriteByte(nullableFlagFalse)

	// Generic count and parameter count.
	buffer.WriteString(padNumberToString(0, trhSlotGenericCount.length))
	buffer.WriteString(padNumberToString(0, trhSlotParameterCount.length))

	return buffer.String()
}

// buildSpecialTypeReferenceValue returns a string value for representing the given special type
// reference.
func buildSpecialTypeReferenceValue(specialFlagValue byte) string {
	var buffer bytes.Buffer

	// Referenced type ID. For special types, we leave this empty.
	buffer.WriteByte('[')
	buffer.WriteString(strings.Repeat(" ", trhSlotTypeId.length))
	buffer.WriteByte(']')

	// Flags: special and nullable.
	buffer.WriteByte(specialFlagValue)
	buffer.WriteByte(nullableFlagFalse)

	// Generic count and parameter count.
	buffer.WriteString(padNumberToString(0, trhSlotGenericCount.length))
	buffer.WriteString(padNumberToString(0, trhSlotParameterCount.length))

	return buffer.String()
}

// padNumberToString converts the given number into a string, padding it to ensure it has the
// specified number of digits.
func padNumberToString(number int, digits int) string {
	numberString := strconv.Itoa(number)
	return strings.Repeat("0", digits-len(numberString)) + numberString
}
