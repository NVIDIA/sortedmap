// Copyright (c) 2015-2025, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package sortedmap

import "fmt"

func (tree *btreeTreeStruct) Dump() (err error) {
	tree.Lock()
	defer tree.Unlock()

	if tree.nodeCache != nil {
		fmt.Printf("B+Tree @ %p has Node Cache @ %p\n", tree, tree.nodeCache)
		fmt.Printf("  .evictLowLimit  = %v\n", tree.nodeCache.evictLowLimit)
		fmt.Printf("  .evictHighLimit = %v\n", tree.nodeCache.evictHighLimit)
		if tree.nodeCache.cleanLRUHead == nil {
			fmt.Printf("  .cleanLRUHead   = nil\n")
		} else {
			fmt.Printf("  .cleanLRUHead   = %p\n", tree.nodeCache.cleanLRUHead)
		}
		if tree.nodeCache.cleanLRUTail == nil {
			fmt.Printf("  .cleanLRUTail   = nil\n")
		} else {
			fmt.Printf("  .cleanLRUTail   = %p\n", tree.nodeCache.cleanLRUTail)
		}
		fmt.Printf("  .cleanLRUItems  = %v\n", tree.nodeCache.cleanLRUItems)
		if tree.nodeCache.dirtyLRUHead == nil {
			fmt.Printf("  .dirtyLRUHead   = nil\n")
		} else {
			fmt.Printf("  .dirtyLRUHead   = %p\n", tree.nodeCache.dirtyLRUHead)
		}
		if tree.nodeCache.dirtyLRUTail == nil {
			fmt.Printf("  .dirtyLRUTail   = nil\n")
		} else {
			fmt.Printf("  .dirtyLRUTail   = %p\n", tree.nodeCache.dirtyLRUTail)
		}
		fmt.Printf("  .dirtyLRUItems  = %v\n", tree.nodeCache.dirtyLRUItems)
		fmt.Printf("  .drainerActive  = %v\n", tree.nodeCache.drainerActive)
	}

	fmt.Printf("B+Tree @ %p has Root Node @ %p\n", tree, tree.root)
	fmt.Printf("  .minKeysPerNode      = %v\n", tree.minKeysPerNode)
	fmt.Printf("  .maxKeysPerNode      = %v\n", tree.maxKeysPerNode)

	err = tree.dumpNode(tree.root, "")

	return
}

func (tree *btreeTreeStruct) dumpNode(node *btreeNodeStruct, indent string) (err error) {
	if !node.loaded {
		err = node.tree.loadNode(node)
		if err != nil {
			return
		}
	}

	if tree.nodeCache != nil {
		switch node.btreeNodeCacheElement.btreeNodeCacheTag {
		case noLRU:
			fmt.Printf("%v  .btreeNodeCacheTag   = noLRU (%v)\n", indent, noLRU)
		case cleanLRU:
			fmt.Printf("%v  .btreeNodeCacheTag   = cleanLRU (%v)\n", indent, cleanLRU)
		case dirtyLRU:
			fmt.Printf("%v  .btreeNodeCacheTag   = dirtyLRU (%v)\n", indent, dirtyLRU)
		default:
			fmt.Printf("%v  .btreeNodeCacheTag   = <unknown> (%v)\n", indent, node.btreeNodeCacheElement.btreeNodeCacheTag)
		}

		if noLRU != node.btreeNodeCacheElement.btreeNodeCacheTag {
			if node.btreeNodeCacheElement.nextBTreeNode == nil {
				fmt.Printf("%v  .nextBTreeNode       = nil\n", indent)
			} else {
				fmt.Printf("%v  .nextBTreeNode       = %p\n", indent, node.btreeNodeCacheElement.nextBTreeNode)
			}
			if node.btreeNodeCacheElement.prevBTreeNode == nil {
				fmt.Printf("%v  .prevBTreeNode       = nil\n", indent)
			} else {
				fmt.Printf("%v  .prevBTreeNode       = %p\n", indent, node.btreeNodeCacheElement.prevBTreeNode)
			}
		}
	}

	fmt.Printf("%v  .objectNumber        = 0x%016x\n", indent, node.objectNumber)
	fmt.Printf("%v  .objectOffset        = 0x%016x\n", indent, node.objectOffset)
	fmt.Printf("%v  .objectLength        = 0x%016x\n", indent, node.objectLength)
	fmt.Printf("%v  .items               = %v\n", indent, node.items)
	fmt.Printf("%v  .loaded              = %v\n", indent, node.loaded)
	fmt.Printf("%v  .dirty               = %v\n", indent, node.dirty)
	fmt.Printf("%v  .root                = %v\n", indent, node.root)
	fmt.Printf("%v  .leaf                = %v\n", indent, node.leaf)

	if node.parentNode == nil {
		fmt.Printf("%v  .parentNode          = nil\n", indent)
	} else {
		fmt.Printf("%v  .parentNode          = %p\n", indent, node.parentNode)
	}

	if !node.leaf {
		if node.rootPrefixSumChild == nil {
			fmt.Printf("%v  .rootPrefixSumChild  = nil\n", indent)
		} else {
			fmt.Printf("%v  .rootPrefixSumChild  = %p\n", indent, node.rootPrefixSumChild)
		}
	}

	if !node.root {
		fmt.Printf("%v  .prefixSumItems      = %v\n", indent, node.prefixSumItems)
		fmt.Printf("%v  .prefixSumKVIndex    = %v\n", indent, node.prefixSumKVIndex)
		if node.prefixSumParent == nil {
			fmt.Printf("%v  .prefixSumParent     = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumParent     = %p\n", indent, node.prefixSumParent)
		}
		if node.prefixSumLeftChild == nil {
			fmt.Printf("%v  .prefixSumLeftChild  = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumLeftChild  = %p\n", indent, node.prefixSumLeftChild)
		}
		if node.prefixSumRightChild == nil {
			fmt.Printf("%v  .prefixSumRightChild = nil\n", indent)
		} else {
			fmt.Printf("%v  .prefixSumRightChild = %p\n", indent, node.prefixSumRightChild)
		}
	}

	if !node.leaf {
		if node.nonLeafLeftChild != nil {
			fmt.Printf("%v  .nonLeafLeftChild    = %p\n", indent, node.nonLeafLeftChild)
			tree.dumpNode(node.nonLeafLeftChild, "    "+indent)
		}
	}

	numKVentries, lenErr := node.kvLLRB.Len()
	if lenErr != nil {
		err = lenErr
		return
	}
	for i := range numKVentries {
		key, value, _, getByIndexErr := node.kvLLRB.GetByIndex(i)
		if getByIndexErr != nil {
			err = getByIndexErr
			return
		}
		keyAsString, nonShadowingErr := tree.DumpKey(key)
		if nonShadowingErr != nil {
			err = nonShadowingErr
			return
		}
		fmt.Printf("%v  .kvLLRB[%v].Key       = %v\n", indent, i, keyAsString)
		if node.leaf {
			valueAsString, nonShadowingErr := tree.DumpValue(value)
			if nonShadowingErr != nil {
				err = nonShadowingErr
				return
			}
			fmt.Printf("%v  .kvLLRB[%v].Value     = %v\n", indent, i, valueAsString)
		} else {
			childNode := value.(*btreeNodeStruct)
			fmt.Printf("%v  .kvLLRB[%v].Value     = %p\n", indent, i, childNode)
			tree.dumpNode(childNode, "    "+indent)
		}
	}

	err = nil

	return
}
