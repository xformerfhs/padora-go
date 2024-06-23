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

// This is the main program of the padding oracle demonstration.

package main

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math"
	"math/rand"
	"padora/numberformat"
	"time"
)

// ******** Private constants ********

// aesBlockSize is the block size of the AES cipher in bytes.
const aesBlockSize = 16

// ******** Main function ********

// main is the main program.
func main() {
	numBlocks := GetNumBlocks()

	secretMessage := makeSecretMessage(numBlocks, aesBlockSize)
	fmt.Printf("\nLength of secret message is %s bytes\n", numberformat.FormatInt(len(secretMessage)))

	encryptedMessage := PadAndEncrypt(secretMessage, aesBlockSize)

	// Padded length is encrypted length minus initialization vector length.
	paddedLength := len(encryptedMessage) - aesBlockSize
	fmt.Printf("Length of padded encrypted message is %s bytes\n", numberformat.FormatInt(paddedLength))

	startTime := time.Now()
	recoveredMessage, count := Crack(encryptedMessage, aesBlockSize)
	elapsedTime := time.Since(startTime)

	fmt.Println()
	if bytes.Compare(secretMessage, recoveredMessage) == 0 {
		fmt.Println(`>>>> Secret message successfully retrieved! <<<<`)
	} else {
		fmt.Println(`!!!! Unable to retrieve secret message!!!!`)
		showDiff(secretMessage, recoveredMessage)
	}
	fmt.Println()
	fmt.Printf("%s decryption calls needed %v. This means %d calls per byte.\n",
		numberformat.FormatInt(count),
		elapsedTime,
		int(math.Round(float64(count)/float64(paddedLength))))
}

// ******** Private functions ********

// makeSecretMessage builds a random secret message.
func makeSecretMessage(numBlocks int, blockSize int) []byte {
	result := make([]byte, numBlocks*blockSize-rand.Intn(blockSize))
	_, _ = crand.Read(result)
	return result
}

// showDiff shows the difference between two byte slices.
func showDiff(a []byte, b []byte) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("%d: %02x != %02x\n", i, a[i], b[i])
		}
	}
}
