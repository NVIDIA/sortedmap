// Copyright (c) 2015-2025, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package sortedmap

import "fmt"

func (tree *llrbTreeStruct) Validate() (err error) {
	tree.Lock()
	defer tree.Unlock()

	err = tree.root.validate()

	return
}

func (node *llrbNodeStruct) validate() (err error) {
	if node == nil {
		return nil
	}

	computedHeight := 1

	if node.left != nil {
		computedHeight += node.left.len

		err = node.left.validate()

		if err != nil {
			return
		}
	}

	if node.right != nil {
		computedHeight += node.right.len

		err = node.right.validate()

		if err != nil {
			return
		}
	}

	if computedHeight != node.len {
		err = fmt.Errorf("For node.Key == %v, computedHeight(%v) != node.len(%v)", node.Key, computedHeight, node.len)
		return
	}

	return nil
}
