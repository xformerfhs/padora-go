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
//    2024-06-23: V1.0.0: Created.
//    2024-08-29: V1.1.0: Print used number of blocks.
//

// This file contains the functions to process the command line arguments.

package main

import (
	"fmt"
	"os"
	"padora/numberformat"
	"strconv"
)

// ******** Private constants ********

// errMsgInvalidNoOfBlocks is the error message for an invalid number of blocks.
const errMsgInvalidNoOfBlocks = "Invalid number of blocks: %d\n"

// defaultNumBlocks is the default number of blocks for secret message.
const defaultNumBlocks = 3

// minNumBlocks is the minimum allowed number of blocks.
const minNumBlocks = 1

// minNumBlocks is the maximum allowed number of blocks.
const maxNumBlocks = 4_000

// ******** Public functions ********

// GetNumBlocks gets the number of blocks to generate from the command line argument.
func GetNumBlocks() int {
	var err error

	fmt.Println()

	numBlocks := defaultNumBlocks

	if len(os.Args) > 1 {
		numBlocks, err = strconv.Atoi(os.Args[1])

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, errMsgInvalidNoOfBlocks, numBlocks)
			numBlocks = defaultNumBlocks
		}

		if numBlocks < minNumBlocks {
			_, _ = fmt.Fprintf(os.Stderr, errMsgInvalidNoOfBlocks, numBlocks)
			numBlocks = minNumBlocks
		}

		if numBlocks > maxNumBlocks {
			_, _ = fmt.Fprintf(os.Stderr, errMsgInvalidNoOfBlocks, numBlocks)
			numBlocks = maxNumBlocks
		}
	}

	fmt.Printf("Using %s blocks\n", numberformat.FormatInt(numBlocks))

	return numBlocks
}
