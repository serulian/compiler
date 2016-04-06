// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sourcemap

import "strings"

// From: https://github.com/mozilla/source-map/blob/master/dist/source-map.js
// A single base 64 digit can contain 6 bits of data. For the base 64 variable
// length quantities we use in the source map spec, the first bit is the sign,
// the next four bits are the actual value, and the 6th bit is the
// continuation bit. The continuation bit tells us whether there are more
// digits in this value following this digit.
//
//   Continuation
//   |    Sign (only applied on the very last bit when continuation is 0)
//   |    |
//   V    V
//   101011

const BASE64_CHARACTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const VLQ_CONTINUATION_BIT_MASK = 32
const VLQ_SIGN_BIT_MASK = 1
const VLQ_VALUE_MASK = VLQ_CONTINUATION_BIT_MASK - 1
const VLQ_BIT_DATA_SHIFT = 5 // 5 bits of data in each character (except the last, which is 4)

// vlqEncode performs full VLQ encoding of the given int slice, returning the string.
func vlqEncode(data []int) string {
	var encoded = ""
	for _, value := range data {
		encoded += vlqEncodeValue(value)
	}
	return encoded
}

// vlqEncodeValue encodes a single integer value into a VLQ string.
func vlqEncodeValue(data int) string {
	// Move the sign bit (0 for positive, 1 for negative) to the end of the value.
	var isNegativeInt uint = 0
	if data < 0 {
		isNegativeInt = 1
		data = -data
	}

	var encoded = ""
	var value uint = (uint(data) << 1) | isNegativeInt
	for {
		// Apply the value mask to the current value to get the necessary five bits.
		valueBits := value & VLQ_VALUE_MASK

		// Shift the value over to account for the bits removed.
		value = value >> VLQ_BIT_DATA_SHIFT

		// Determine whether we need a continuation bit by whether the value still
		// has data.
		if value > 0 {
			valueBits = valueBits | VLQ_CONTINUATION_BIT_MASK
		}

		encoded = encoded + string(BASE64_CHARACTERS[valueBits])
		if value <= 0 {
			return encoded
		}
	}
}

// vlqDecode performs full VLQ decoding of the given string, returning the integer values
// found and whether the decode succeeded.
func vlqDecode(data string) ([]int, bool) {
	var values = make([]int, 0)
	var current = data

	for {
		if len(current) == 0 {
			return values, true
		}

		value, remaining, ok := vlqDecodeValue(current, 0)
		if !ok {
			return values, false
		}

		values = append(values, value)
		current = remaining
	}
}

// vlqDecodeValue attempts to decode the VLQ Base64-encoded integer found starting at the given start
// index in the data string. Returns the integer value, the remaining string data and whether
// the decode was a success.
func vlqDecodeValue(data string, startIndex int) (int, string, bool) {
	if startIndex >= len(data) {
		return -1, "", false
	}

	suffix := data[startIndex:]

	var current = 0
	var shift uint = 0
	for index, currentRune := range suffix {
		// Find the numeric value for the base64 character.
		base64Index := strings.IndexRune(BASE64_CHARACTERS, currentRune)
		if base64Index < 0 {
			return -1, "", false
		}

		// The value of the VLQ is found in the last five bits.
		value := base64Index & VLQ_VALUE_MASK

		// The continuation bit determines whether the current value is a continuation of the
		// previous value.
		isContinuation := (base64Index & VLQ_CONTINUATION_BIT_MASK) == VLQ_CONTINUATION_BIT_MASK

		// Apply the found value to the current value, with proper shifting.
		current += (value << shift)
		shift += VLQ_BIT_DATA_SHIFT

		if !isContinuation {
			// The last bit of the value is a sign bit.
			isNegative := (current & VLQ_SIGN_BIT_MASK) == 1
			current = current >> 1
			if isNegative {
				current = -current
			}

			return current, data[startIndex+index+1:], true
		}
	}

	return -1, "", false
}
