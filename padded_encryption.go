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

// This file contains the functions that process encryption and padding.

package main

// PadAndEncrypt pads and encrypts a clear message.
func PadAndEncrypt(clearMessage []byte, blockSize int) []byte {
	return Encrypt(Pad(clearMessage, blockSize))
}

// DecryptAndUnpad decrypts and unpads a concatenation of an initialization vector and an encrypted message.
func DecryptAndUnpad(compoundEncryptedMessage []byte, blockSize int) ([]byte, error) {
	decryptedMessage := Decrypt(compoundEncryptedMessage)
	return Unpad(decryptedMessage, blockSize)
}
