package sortedmap

import (
	"fmt"
	"strings"
)

func (tree *llrbTreeStruct) Dump() (err error) {
	tree.Lock()
	defer tree.Unlock()

	err = nil

	dumpInFlatForm(tree.root)
	dumpInTreeForm(tree)

	return
}

func dumpInFlatForm(node *llrbNodeStruct) {
	if nil == node {
		return
	}

	nodeLeftKey := "nil"
	if nil != node.left {
		nodeLeftKey = fmt.Sprintf("%v", node.left.Key)
	}

	nodeRightKey := "nil"
	if nil != node.right {
		nodeRightKey = fmt.Sprintf("%v", node.right.Key)
	}

	var colorString string
	if RED == node.color {
		colorString = "RED"
	} else { // BLACK == node.color
		colorString = "BLACK"
	}

	fmt.Printf("%v Node Key == %v Node.left.Key == %v Node.right.Key == %v len == %v\n", colorString, node.Key, nodeLeftKey, nodeRightKey, node.len)

	dumpInFlatForm(node.left)
	dumpInFlatForm(node.right)
}

func dumpInTreeForm(tree *llrbTreeStruct) {
	if nil == tree.root {
		return
	}

	if nil != tree.root.right {
		dumpInTreeFormNode(tree.root.right, true, "")
	}
	fmt.Println(tree.root.Key)
	if nil != tree.root.left {
		dumpInTreeFormNode(tree.root.left, false, "")
	}
}

func dumpInTreeFormNode(node *llrbNodeStruct, isRight bool, indent string) {
	var indentAppendage string
	var nextIndent string

	if nil != node.right {
		if isRight {
			indentAppendage = "        "
		} else {
			indentAppendage = " |      "
		}
		nextIndent = strings.Join([]string{indent, indentAppendage}, "")
		dumpInTreeFormNode(node.right, true, nextIndent)
	}
	fmt.Printf("%v", indent)
	if isRight {
		fmt.Printf(" /")
	} else {
		fmt.Printf(" \\")
	}
	fmt.Println("-----", node.Key)
	if nil != node.left {
		if isRight {
			indentAppendage = " |      "
		} else {
			indentAppendage = "        "
		}
		nextIndent = strings.Join([]string{indent, indentAppendage}, "")
		dumpInTreeFormNode(node.left, false, nextIndent)
	}
}
