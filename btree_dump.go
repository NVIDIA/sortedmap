package sortedmap

import "fmt"

func (tree *btreeTreeStruct) Dump() (err error) {
	fmt.Printf("Root Node @ %p\n", tree.root)
	err = tree.dumpNode(tree.root, "")
	return
}

func (tree *btreeTreeStruct) dumpNode(node *btreeNodeStruct, indent string) (err error) {
	if !node.loaded {
		node.tree.loadNode(node)
	}

	fmt.Printf("%v  .items               = %v\n", indent, node.items)
	fmt.Printf("%v  .dirty               = %v\n", indent, node.dirty)
	fmt.Printf("%v  .root                = %v\n", indent, node.root)
	fmt.Printf("%v  .leaf                = %v\n", indent, node.leaf)

	if nil == node.parentNode {
		fmt.Printf("%v  .parentNode          = nil\n", indent)
	} else {
		fmt.Printf("%v  .parentNode          = %p\n", indent, node.parentNode)
	}

	if !node.leaf {
		if nil == node.rootPrefixSumChild {
			fmt.Printf("%v  .rootPrefixSumChild  = nil\n", indent)
		} else {
			fmt.Printf("%v  .rootPrefixSumChild  = %p\n", indent, node.rootPrefixSumChild)
		}
	}

	if !node.root {
		fmt.Printf("%v  .prefixSumItems      = %v\n", indent, node.prefixSumItems)
		fmt.Printf("%v  .prefixSumKVIndex    = %v\n", indent, node.prefixSumKVIndex)
		if nil == node.prefixSumParent {
			fmt.Printf("%v  .prefixSumParent     = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumParent     = %p\n", indent, node.prefixSumParent)
		}
		if nil == node.prefixSumLeftChild {
			fmt.Printf("%v  .prefixSumLeftChild  = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumLeftChild  = %p\n", indent, node.prefixSumLeftChild)
		}
		if nil == node.prefixSumRightChild {
			fmt.Printf("%v  .prefixSumRightChild = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumRightChild = %p\n", indent, node.prefixSumRightChild)
		}
	}

	if !node.leaf {
		if nil != node.nonLeafLeftChild {
			fmt.Printf("%v  .nonLeafLeftChild    = %p\n", indent, node.nonLeafLeftChild)
			tree.dumpNode(node.nonLeafLeftChild, "    "+indent)
		}
	}

	numKVentries, lenErr := node.kvLLRB.Len()
	if nil != lenErr {
		err = lenErr
		return
	}
	for i := 0; i < numKVentries; i++ {
		key, value, _, getByIndexErr := node.kvLLRB.GetByIndex(i)
		if nil != getByIndexErr {
			err = getByIndexErr
			return
		}
		packedKey, packedKeyErr := tree.PackKey(key)
		if nil != packedKeyErr {
			err = packedKeyErr
			return
		}
		fmt.Printf("%v  .kvLLRB[%v].Key       = %v\n", indent, i, packedKey)
		if node.leaf {
			packedValue, packedValueErr := tree.PackValue(value)
			if nil != packedValueErr {
				err = packedValueErr
				return
			}
			fmt.Printf("%v  .kvLLRB[%v].Value     = %v\n", indent, i, packedValue)
		} else {
			childNode := value.(*btreeNodeStruct)
			fmt.Printf("%v  .kvLLRB[%v].Value     = %p\n", indent, i, childNode)
			tree.dumpNode(childNode, "    "+indent)
		}
	}

	err = nil

	return
}
