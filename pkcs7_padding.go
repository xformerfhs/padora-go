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
//    2024-06-21: V1.0.0: Created.
//

// This file contains the PKCS#7 padding and unpadding functions.

package main

import (
	"errors"
	"padora/slicehelper"
)

// ******** Public constants ********

// ErrInvalidPadding signals an invalid padding.
var ErrInvalidPadding = errors.New(`invalid padding`)

// ******** Public functions ********

// Pad pads an unpadded message.
func Pad(unpaddedMessage []byte, blockSize int) []byte {
	paddingLength := byte(blockSize - len(unpaddedMessage)%blockSize)
	padding := make([]byte, paddingLength)
	slicehelper.Fill(padding, paddingLength)
	return slicehelper.Concat(unpaddedMessage, padding)
}

// Unpad unpads a padded message.
func Unpad(paddedMessage []byte, blockSize int) ([]byte, error) {
	maxIndex := len(paddedMessage) - 1

	// 1. Check padding.

	// Get last byte.
	lastByte := paddedMessage[maxIndex]
	intLastByte := int(lastByte)

	// Has last byte an invalid value?
	if lastByte == 0 || intLastByte > blockSize {
		return nil, ErrInvalidPadding
	}

	// Check if all expected padding bytes are present.
	for i := maxIndex - 1; i > maxIndex-intLastByte; i-- {
		if paddedMessage[i] != lastByte {
			return nil, ErrInvalidPadding
		}
	}

	// 2. Now unpad.
	return paddedMessage[:maxIndex-intLastByte+1], nil
}
