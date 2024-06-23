//
// SPDX-FileCopyrightText: Copyright 2023 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.3
//
// Change history:
//    2023-03-26: V1.0.0: Created.
//    2024-01-05: V1.0.1: More efficient generation of MinInt64Text.
//    2024-01-08: V1.0.2: Conversion to digits sped up.
//    2024-01-09: V1.0.3: Conversion to digits sped up a little bit more.
//

// Package numberformat implements a formatter for integer numbers that groups digits
// in blocks of three and separates them by a separator. It uses a static buffer, so it
// must not be called at the same time from multiple callers.
package numberformat

import (
	"math"
)

// ======== Private constants ========

// defaultSeparator is the default group separator.
const defaultSeparator byte = ','

// maxFormattedNumberLength is the maximum length of a formatted number.
const maxFormattedNumberLength = 26 // Length of the string "-9,223,372,036,854,775,808"

// ======== Private variables ========

// buffer is the place where the formatted number is put.
var buffer [maxFormattedNumberLength]byte

// ======== Private functions ========

// formatUint64WithSeparatorInBuffer formats an uint64 with the specified separator into the buffer.
func formatUint64WithSeparatorInBuffer(aNumber uint64, separator byte) int {
	idx := byte(len(buffer))
	charCount := byte(0)

	last := aNumber

	for {
		if charCount >= 3 {
			idx--
			buffer[idx] = separator
			charCount = 0
		}

		divided := last / 10

		idx--
		buffer[idx] = '0' + byte(last-(divided*10))
		charCount++

		last = divided

		if last == 0 {
			break
		}
	}

	return int(idx) // Strangely, returning int(idx) is faster than returning idx as a byte and converting that in the caller
}

// MinInt64Text returns the correct print string for the minimum Int64 value.
func MinInt64Text(separator byte) string {
	// MinInt64: -9,223,372,036,854,775,808
	result := []byte{
		'-', '9', separator,
		'2', '2', '3', separator,
		'3', '7', '2', separator,
		'0', '3', '6', separator,
		'8', '5', '4', separator,
		'7', '7', '5', separator,
		'8', '0', '8'}

	return string(result)
}

// ======== Public functions ========

// -------- Int functions --------

// FormatInt64WithSeparator formats an int64 with the specified separator.
func FormatInt64WithSeparator(aNumber int64, separator byte) string {
	var positiveNumber uint64
	isNegative := aNumber < 0
	if isNegative {
		// The minimum int64 value can not be made positive, so this is handled on its own.
		if aNumber != math.MinInt64 {
			positiveNumber = uint64(-aNumber)
		} else {
			return MinInt64Text(separator)
		}
	} else {
		positiveNumber = uint64(aNumber)
	}

	idx := formatUint64WithSeparatorInBuffer(positiveNumber, separator)
	if isNegative {
		idx--
		buffer[idx] = '-'
	}

	return string(buffer[idx:])
}

// FormatInt32WithSeparator formats an int32 with the specified separator.
func FormatInt32WithSeparator(aNumber int32, separator byte) string {
	return FormatInt64WithSeparator(int64(aNumber), separator)
}

// FormatIntWithSeparator formats an int with the specified separator.
func FormatIntWithSeparator(aNumber int, separator byte) string {
	return FormatInt64WithSeparator(int64(aNumber), separator)
}

// FormatInt64 formats an int64 with the default separator.
func FormatInt64(aNumber int64) string {
	return FormatInt64WithSeparator(aNumber, defaultSeparator)
}

// FormatInt32 formats an int32 with the default separator.
func FormatInt32(aNumber int32) string {
	return FormatInt32WithSeparator(aNumber, defaultSeparator)
}

// FormatInt formats an int with the default separator.
func FormatInt(aNumber int) string {
	return FormatIntWithSeparator(aNumber, defaultSeparator)
}

// -------- Uint functions --------

// FormatUint64WithSeparator formats an uint64 with the specified separator.
func FormatUint64WithSeparator(aNumber uint64, separator byte) string {
	idx := formatUint64WithSeparatorInBuffer(aNumber, separator)

	return string(buffer[idx:])
}

// FormatUint32WithSeparator formats an uint32 with the specified separator.
func FormatUint32WithSeparator(aNumber uint32, separator byte) string {
	return FormatUint64WithSeparator(uint64(aNumber), separator)
}

// FormatUintWithSeparator formats an uint with the specified separator.
func FormatUintWithSeparator(aNumber uint, separator byte) string {
	return FormatUint64WithSeparator(uint64(aNumber), separator)
}

// FormatUint64 formats an uint64 with the default separator.
func FormatUint64(aNumber uint64) string {
	return FormatUint64WithSeparator(aNumber, defaultSeparator)
}

// FormatUint32 formats an uint32 with the default separator.
func FormatUint32(aNumber uint32) string {
	return FormatUint32WithSeparator(aNumber, defaultSeparator)
}

// FormatUint formats an uint with the default separator.
func FormatUint(aNumber uint) string {
	return FormatUintWithSeparator(aNumber, defaultSeparator)
}
