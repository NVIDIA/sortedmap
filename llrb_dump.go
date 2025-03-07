// Copyright (c) 2015-2025, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package sortedmap

import (
	"fmt"
)

func (tree *llrbTreeStruct) Dump() (err error) {
	tree.Lock()
	defer tree.Unlock()

	err = nil

	err = tree.dumpInFlatForm(tree.root)
	if err != nil {
		err = fmt.Errorf("dumpInFlatForm() failed: %v", err)
		fmt.Printf("\n%v\n", err)
		return
	}

	err = tree.dumpInTreeForm()
	if err != nil {
		err = fmt.Errorf("dumpInTreeForm() failed: %v", err)
		fmt.Printf("\n%v\n", err)
		return
	}

	err = nil
	return
}

func (tree *llrbTreeStruct) dumpInFlatForm(node *llrbNodeStruct) (err error) {
	if node == nil {
		err = nil
		return
	}

	nodeLeftKey := "nil"
	if node.left != nil {
		nodeLeftKey, err = tree.DumpKey(node.left.Key)
		if err != nil {
			return
		}
	}

	nodeRightKey := "nil"
	if node.right != nil {
		nodeRightKey, err = tree.DumpKey(node.right.Key)
		if err != nil {
			return
		}
	}

	var colorString string
	// if RED == node.color {
	if node.color {
		colorString = "RED"
	} else { // BLACK == node.color
		colorString = "BLACK"
	}

	nodeThisKey, err := tree.DumpKey(node.Key)
	if err != nil {
		return
	}

	fmt.Printf("%v Node Key == %v Node.left.Key == %v Node.right.Key == %v len == %v\n", colorString, nodeThisKey, nodeLeftKey, nodeRightKey, node.len)

	err = tree.dumpInFlatForm(node.left)
	if err != nil {
		return
	}

	err = tree.dumpInFlatForm(node.right)
	if err != nil {
		return
	}

	err = nil
	return
}

func (tree *llrbTreeStruct) dumpInTreeForm() (err error) {
	if tree.root == nil {
		err = nil
		return
	}

	if tree.root.right != nil {
		err = tree.dumpInTreeFormNode(tree.root.right, true, "")
		if err != nil {
			return
		}
	}

	rootKey, err := tree.DumpKey(tree.root.Key)
	if err != nil {
		return
	}

	fmt.Printf("%v\n", rootKey)

	if tree.root.left != nil {
		err = tree.dumpInTreeFormNode(tree.root.left, false, "")
		if err != nil {
			return
		}
	}

	err = nil
	return
}

func (tree *llrbTreeStruct) dumpInTreeFormNode(node *llrbNodeStruct, isRight bool, indent string) (err error) {
	var indentAppendage string
	var nextIndent string

	if node.right != nil {
		if isRight {
			indentAppendage = "        "
		} else {
			indentAppendage = " |      "
		}
		nextIndent = indent + indentAppendage
		err = tree.dumpInTreeFormNode(node.right, true, nextIndent)
		if err != nil {
			return
		}
	}

	fmt.Printf("%v", indent)
	if isRight {
		fmt.Printf(" /")
	} else {
		fmt.Printf(" \\")
	}

	nodeKey, err := tree.DumpKey(node.Key)
	if err != nil {
		return
	}

	fmt.Printf("----- %v\n", nodeKey)

	if node.left != nil {
		if isRight {
			indentAppendage = " |      "
		} else {
			indentAppendage = "        "
		}
		nextIndent = indent + indentAppendage
		err = tree.dumpInTreeFormNode(node.left, false, nextIndent)
		if err != nil {
			return
		}
	}

	err = nil
	return
}
