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
// Version: 1.1.0
//
// Change history:
//    2024-06-21: V1.0.0: Created.
//    2024-08-28: V1.0.1: Rename variable to better reflect its meaning.
//    2024-08-29: V1.1.0: Show progress information.
//

// This file contains the cracker functions that perform a padding oracle attack
// on data encrypted with AES-CBC and padded with PKCS#7.
//
// It implements a very simple version of a padding oracle attack,
// just to show how such an attack works in principle.

package main

import (
	"fmt"
	"padora/numberformat"
	"slices"
)

// ======== Private constants ========

// progressStep is the size of step for reporting progress.
const progressStep = 100_000

// ======== Public function ========

// Crack cracks an encrypted message with a CBC/PKCS#7 padding oracle.
func Crack(encryptedMessage []byte, blockSize int) ([]byte, int) {
	result := make([]byte, len(encryptedMessage)-blockSize)
	count := 0
	nextCountShow := progressStep

	fmt.Println()

	// Clone the encrypted message into a buffer that can be manipulated.
	modifiedMessage := slices.Clone(encryptedMessage)

	// Loop through the message block by block, beginning at the last one.
	// The first block (start: 0) is not cracked for two reasons:
	// 1. It is the first block and as such it does not have a previous block,
	//    that can be manipulated,
	// 2. It is the initialization vector and not encrypted data.
	var previousOriginalBlock []byte
	var previousModifiedBlock []byte
	crackedBlock := make([]byte, blockSize)
	isLastBlock := true
	for start := len(encryptedMessage) - blockSize; start >= blockSize; start -= blockSize {
		// Prepare two slices that each point to the block before the current block, as this
		// is the one that is manipulated in this attack.
		previousOriginalBlock = encryptedMessage[start-blockSize : start]
		previousModifiedBlock = modifiedMessage[start-blockSize : start]
		count += crackBlock(modifiedMessage,
			previousOriginalBlock,
			previousModifiedBlock,
			crackedBlock,
			blockSize,
			start,
			isLastBlock)

		copy(result[start-blockSize:], crackedBlock)

		if isLastBlock {
			isLastBlock = false
		}

		if count >= nextCountShow {
			fmt.Printf("Guess count: %9s\n", numberformat.FormatInt(count))
			nextCountShow += progressStep
		}
	}

	result, _ = Unpad(result, blockSize)

	return result, count
}

// crackBlock cracks one block.
func crackBlock(
	modifiedMessage []byte,
	previousOriginalBlock []byte,
	previousModifiedBlock []byte,
	crackedBlock []byte,
	blockSize int,
	start int,
	isLastBlock bool) int {
	// Shorten the modified message so that the block we want to crack is the last block.
	modifiedMessage = modifiedMessage[:start+blockSize]

	count := 0
	for pos := blockSize - 1; pos >= 0; pos-- {
		// This is the padding value that is forced upon the end of the modified message.
		wantedPaddingLength := byte(blockSize - pos)

		// 1. Set all encrypted bytes following the current byte so that they are
		//    decrypted to [wantedPaddingLength].
		prepareKnownPadding(
			previousOriginalBlock,
			previousModifiedBlock,
			crackedBlock,
			pos,
			blockSize,
			wantedPaddingLength)

		// 2. Guess the current byte.
		count += guessValue(
			modifiedMessage,
			previousOriginalBlock,
			previousModifiedBlock,
			crackedBlock,
			pos,
			blockSize,
			wantedPaddingLength,
			isLastBlock)
	}

	// Restore previous modified block to contain the original data again.
	// It is the next block to be attacked, so the original content is needed.
	copy(previousModifiedBlock, previousOriginalBlock)

	return count
}

// prepareKnownPadding sets the bytes following the current byte
// so that they are decrypted to the wanted valid padding bytes.
func prepareKnownPadding(
	previousOriginalBlock []byte,
	previousModifiedBlock []byte,
	crackedBlock []byte,
	pos int,
	blockSize int,
	wantedPaddingLength byte) {
	for preparePos := pos + 1; preparePos < blockSize; preparePos++ {
		previousModifiedBlock[preparePos] = previousOriginalBlock[preparePos] ^
			crackedBlock[preparePos] ^
			wantedPaddingLength
	}
}

// guessValue finds the correct byte by guessing it and asking the padding oracle,
// if the guess is correct.
func guessValue(
	modifiedMessage []byte,
	previousOriginalBlock []byte,
	previousModifiedBlock []byte,
	crackedBlock []byte,
	pos int,
	blockSize int,
	wantedPaddingLength byte,
	isLastBlock bool) int {
	count := 0
	foundValue := false
	for guess := 0; guess < 256; guess++ {
		guessByte := byte(guess)
		// The following does not work if this is the last padded block and
		// guessByte == wantedPaddingLength, so skip the guess in this case.
		if isLastBlock && (guessByte == wantedPaddingLength) {
			continue
		}

		// This is the modification of the previous block that modifies the decryption
		// of the current block.
		previousModifiedBlock[pos] = previousOriginalBlock[pos] ^
			guessByte ^
			wantedPaddingLength

		count++

		// Now ask the oracle: Did we construct a valid padding?
		_, err := DecryptAndUnpad(modifiedMessage, blockSize)
		if err == nil {
			// There was no padding error, so this is a candidate.
			// However, sometimes this is a match that is caused by the byte before the current one.
			// E.g., if we try to force a 0x01 in the last byte and the second-to-last byte
			// is a 0x02 and our guess produces a 0x02, this is a valid padding, but not the intended one.
			// Disturb the byte before the current one to check, if the match still holds.
			if pos > 0 {
				previousModifiedBlock[pos-1] ^= 0xff
				count++
				_, err = DecryptAndUnpad(modifiedMessage, blockSize)
				if err != nil {
					// Disturbing the byte before this one gave a padding error.
					// So this was an accidental match caused by the previous byte.
					continue
				}
			}

			// It was a real match!
			foundValue = true
			crackedBlock[pos] = guessByte
			break
		}
	}

	// If the loop did not find a value, the correct value is wantedPaddingLength.
	if !foundValue {
		crackedBlock[pos] = wantedPaddingLength
	}

	return count
}
