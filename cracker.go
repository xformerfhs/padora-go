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

// This file contains the cracker functions that perform a padding oracle attack
// on data encrypted with AES-CBC and padded with PKCS#7.
//
// It implements a very simple version of a padding oracle attack,
// just to show how such an attack works in principle.

package main

import (
	"padora/slicehelper"
)

// Crack cracks an encrypted message with a CBC/PKCS#7 padding oracle.
func Crack(encryptedMessage []byte, blockSize int) ([]byte, int) {
	result := make([]byte, 0)
	count := 0

	// Copy the encrypted message in a buffer that can be manipulated.
	modifiedMessage := make([]byte, len(encryptedMessage))
	copy(modifiedMessage, encryptedMessage)

	// Loop through the message block by block, beginning at the last one.
	for start := len(encryptedMessage) - blockSize; start >= blockSize; start -= blockSize {
		var crackedBlock []byte
		crackedBlock, count = crackBlock(encryptedMessage, modifiedMessage, blockSize, start, count)
		result = slicehelper.Concat(crackedBlock, result)
		slicehelper.ClearNumber(crackedBlock)
	}

	result, _ = Unpad(result, blockSize)
	return result, count
}

// crackBlock cracks one block.
func crackBlock(encryptedMessage []byte, modifiedMessage []byte, blockSize int, start int, count int) ([]byte, int) {
	result := make([]byte, blockSize)

	// Prepare two slices that each point to the block before the current block, as this
	// is the one that is manipulated in this attack.
	previousMessageBlock := encryptedMessage[start-blockSize : start]
	previousModifiedBlock := modifiedMessage[start-blockSize : start]

	// Shorten the message so that the block we want to crack is the last block.
	modifiedMessage = modifiedMessage[:start+blockSize]

	for pos := blockSize - 1; pos >= 0; pos-- {
		// This is the padding value that is forced upon the end of the modified message.
		wantedPaddingLength := byte(blockSize - pos)

		// 1. Set all encrypted bytes following the current byte so that they are
		//    decrypted to wantedPaddingLength.
		prepareKnownPadding(
			previousMessageBlock,
			previousModifiedBlock,
			result,
			pos,
			blockSize,
			wantedPaddingLength)

		// 2. Guess the current byte.
		count += guessValue(
			modifiedMessage,
			previousMessageBlock,
			previousModifiedBlock,
			result,
			pos,
			blockSize,
			wantedPaddingLength)
	}

	// Restore previous modified block to contain the original data again.
	// It is the next block to be attacked, so the original content is needed.
	copy(previousModifiedBlock, previousMessageBlock)

	return result, count
}

// prepareKnownPadding sets the bytes following the current byte
// so that they are decrypted to the wanted valid padding bytes.
func prepareKnownPadding(
	previousMessageBlock []byte,
	previousModifiedBlock []byte,
	result []byte,
	pos int,
	blockSize int,
	wantedPaddingLength byte) {
	for preparePos := pos + 1; preparePos < blockSize; preparePos++ {
		previousModifiedBlock[preparePos] = previousMessageBlock[preparePos] ^
			result[preparePos] ^
			wantedPaddingLength
	}
}

// guessValue finds the correct byte by guessing it and asking the padding oracle,
// if the guess is correct.
func guessValue(
	modifiedMessage []byte,
	previousMessageBlock []byte,
	previousModifiedBlock []byte,
	result []byte,
	pos int,
	blockSize int,
	wantedPaddingLength byte) int {
	count := 0
	foundValue := false
	for guess := 0; guess < 256; guess++ {
		guessByte := byte(guess)
		// The following does not work if guessByte == wantedPaddingLength,
		// so skip the guess in this case.
		if guessByte == wantedPaddingLength {
			continue
		}

		// This is the modification of the previous block that modifies the decryption
		// of the current block.
		previousModifiedBlock[pos] = previousMessageBlock[pos] ^
			guessByte ^
			wantedPaddingLength

		// Now ask the oracle: Did we construct a valid padding, or not?
		count++
		_, err := DecryptAndUnpad(modifiedMessage, blockSize)
		if err == nil {
			// There was no padding error. Strike!
			foundValue = true
			result[pos] = guessByte
			break
		}
	}

	// If the loop did not find a value, the correct value is wantedPaddingLength.
	if !foundValue {
		result[pos] = wantedPaddingLength
	}

	return count
}
