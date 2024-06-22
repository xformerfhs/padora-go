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

// This file contains the AES-CBC encryption and decryption functions.

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"padora/slicehelper"
)

// ******** Private constants ********

// key is the encryption key.
// This is only a simple demonstration program.
// NEVER hard-code a key in production systems.
var key = []byte{
	0xe4, 0x15, 0x03, 0xed,
	0x35, 0xe0, 0x43, 0x65,
	0xe5, 0xbd, 0xf1, 0x36,
	0x2c, 0x72, 0x54, 0x3f,
}

// ******** Private variables ********

// modAesCipher is the AES cipher.
// It is instanced only once and reused for every call.
var modAesCipher cipher.Block

// ******** Public functions ********

// Encrypt encrypts a clear message and returns a concatenation
// of the initialization vector and the encrypted data.
func Encrypt(clearMessage []byte) []byte {
	aesCipher := getCipher()

	iv := make([]byte, aes.BlockSize)
	_, _ = rand.Read(iv)

	cbcCipher := cipher.NewCBCEncrypter(aesCipher, iv)

	encryptedBytes := make([]byte, len(clearMessage))
	cbcCipher.CryptBlocks(encryptedBytes, clearMessage)

	return slicehelper.Concat(iv, encryptedBytes)
}

// Decrypt decrypts a concatenation of an initialization vector and an encrypted message.
func Decrypt(compoundEncryptedMessage []byte) []byte {
	aesCipher := getCipher()

	iv, encryptedMessage := slicehelper.CutHead(compoundEncryptedMessage, aes.BlockSize)

	cbcCipher := cipher.NewCBCDecrypter(aesCipher, iv)

	decryptedBytes := make([]byte, len(encryptedMessage))
	cbcCipher.CryptBlocks(decryptedBytes, encryptedMessage)

	return decryptedBytes
}

// ******** Public functions ********

// getCipher returns the AES cipher and creates it, if it does not exist, yet.
func getCipher() cipher.Block {
	if modAesCipher == nil {
		modAesCipher, _ = aes.NewCipher(key)
	}

	return modAesCipher
}
