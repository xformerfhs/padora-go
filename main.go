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
	"padora/slicehelper"
)

// ******** Private constants ********

// aesBlockSize is the block size of the AES cipher in bytes.
const aesBlockSize = 16

// ******** Main function ********

// main is the main program.
func main() {
	secretMessage := makeSecretMessage(3, aesBlockSize)
	fmt.Printf("Length of secret message is %d\n", len(secretMessage))

	encryptedMessage := PadAndEncrypt(secretMessage, aesBlockSize)

	// Padded length is encrypted length minus initialization vector length.
	paddedLength := len(encryptedMessage) - aesBlockSize
	fmt.Printf("Length of padded encrypted message is %d\n", paddedLength)

	recoveredMessage, count := Crack(encryptedMessage, aesBlockSize)

	fmt.Println()
	if bytes.Compare(secretMessage, recoveredMessage) == 0 {
		fmt.Println(`>>>> Secret message successfully retrieved! <<<<`)
	} else {
		fmt.Println(`!!!! Unable to retrieve secret message!!!!`)
	}
	fmt.Println()
	fmt.Printf("Needed %d decryption calls. This means %d calls per byte.\n", count, int(math.Round(float64(count)/float64(paddedLength))))

	// Clear secret data from memory.
	slicehelper.ClearNumber(secretMessage)
	slicehelper.ClearNumber(recoveredMessage)
}

// ******** private functions ********

// makeSecretMessage builds a random secret message.
func makeSecretMessage(numBlocks int, blockSize int) []byte {
	result := make([]byte, numBlocks*blockSize-rand.Intn(blockSize))
	_, _ = crand.Read(result)
	return result
}
