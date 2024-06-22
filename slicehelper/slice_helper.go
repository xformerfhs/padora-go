//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
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
// Version: 1.0.0
//
// Change history:
//    2024-06-22: V1.0.0: Created.
//

// Package slicehelper implements helper functions for slices.
package slicehelper

// ******** Private constants ********

// powerFillThresholdLen is the slice length where PowerFill is more efficient than SimpleFill.
const powerFillThresholdLen = 74

// ******** Public functions ********

// Fill fills a slice with a value in an efficient way up to its length.
func Fill[S ~[]T, T any](s S, v T) {
	sLen := len(s)

	if sLen > 0 {
		doFill(s, v, sLen)
	}
}

// Concat returns a new slice concatenating the passed in slices.
// This is a streamlined version of the slices.Concat function of Go V1.22.
func Concat[S ~[]T, T any](slices ...S) S {
	// 1. Calculate total size.
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	// 2. Make new slice with the total size as the capacity and 0 length.
	result := make(S, 0, size)

	// 3. Append all source slices.
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}

// CutHead cuts the first [l] elements from the supplied slice.
func CutHead[S ~[]T, T any](s S, l int) (S, S) {
	return s[:l], s[l:]
}

// ******** Private functions ********

// doFill fills a slice in an optimal way.
func doFill[S ~[]T, T any](s S, v T, l int) {
	if l >= powerFillThresholdLen {
		doPowerFill(s, v, l)
	} else {
		doSimpleFill(s, v, l)
	}
}

// doSimpleFill fills a slice in a simple way.
func doSimpleFill[S ~[]T, T any](s S, v T, l int) {
	for i := 0; i < l; i++ {
		s[i] = v
	}
}

// doPowerFill fills a slice in an efficient way.
func doPowerFill[S ~[]T, T any](s S, v T, l int) {
	// Put the value into the first slice element
	s[0] = v

	// Incrementally duplicate the value into the rest of the slice
	for j := 1; j < l; j <<= 1 {
		copy(s[j:], s[:j])
	}
}
