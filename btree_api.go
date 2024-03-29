// Copyright (c) 2015-2021, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package sortedmap

import (
	"fmt"

	"github.com/NVIDIA/cstruct"
)

// OnDiskByteOrder specifies the endian-ness expected to be used to persist B+Tree data structures
var OnDiskByteOrder = cstruct.LittleEndian

// LayoutReport is a map where key is an objectNumber and value is objectBytes used in that objectNumber
type LayoutReport map[uint64]uint64

type DimensionsReport struct {
	MinKeysPerNode uint64 // only applies to non-Root nodes
	MaxKeysPerNode uint64
	Items          uint64
	Height         uint64
}

// BPlusTree interface declares the available methods available for a B+Tree
type BPlusTree interface {
	SortedMap
	FetchLocation() (rootObjectNumber uint64, rootObjectOffset uint64, rootObjectLength uint64)
	FetchLayoutReport() (layoutReport LayoutReport, err error)
	FetchDimensionsReport() (dimensionsReport DimensionsReport, err error)
	Flush(andPurge bool) (rootObjectNumber uint64, rootObjectOffset uint64, rootObjectLength uint64, err error)
	Purge(full bool) (err error)
	Touch() (err error)
	TouchItem(thisItemIndexToTouch uint64) (nextItemIndexToTouch uint64, err error)
	Prune() (err error)
	Discard() (err error)
}

// BPlusTreeCallbacks specifies the interface to a set of callbacks provided by the client
type BPlusTreeCallbacks interface {
	DumpCallbacks
	GetNode(objectNumber uint64, objectOffset uint64, objectLength uint64) (nodeByteSlice []byte, err error)
	PutNode(nodeByteSlice []byte) (objectNumber uint64, objectOffset uint64, err error)
	DiscardNode(objectNumber uint64, objectOffset uint64, objectLength uint64) (err error)
	PackKey(key Key) (packedKey []byte, err error)
	UnpackKey(payloadData []byte) (key Key, bytesConsumed uint64, err error)
	PackValue(value Value) (packedValue []byte, err error)
	UnpackValue(payloadData []byte) (value Value, bytesConsumed uint64, err error)
}

type BPlusTreeCacheStats struct {
	EvictLowLimit  uint64
	EvictHighLimit uint64
	CleanLRUItems  uint64
	DirtyLRUItems  uint64
	CacheHits      uint64
	CacheMisses    uint64
}

type BPlusTreeCache interface {
	Stats() (bPlusTreeCacheStats *BPlusTreeCacheStats)
	UpdateLimits(evictLowLimit uint64, evictHighLimit uint64)
}

func NewBPlusTreeCache(evictLowLimit uint64, evictHighLimit uint64) (bPlusTreeCache BPlusTreeCache) {
	bPlusTreeCache = &btreeNodeCacheStruct{
		evictLowLimit:  evictLowLimit,
		evictHighLimit: evictHighLimit,
		cleanLRUHead:   nil,
		cleanLRUTail:   nil,
		cleanLRUItems:  0,
		dirtyLRUHead:   nil,
		dirtyLRUTail:   nil,
		dirtyLRUItems:  0,
		drainerActive:  false,
		cacheHits:      0,
		cacheMisses:    0,
	}
	return
}

// NewBPlusTree is used to construct an in-memory B+Tree supporting the Tree interface
//
// The supplied value of maxKeysPerNode, being precisely twice the computed value of
// (uint64) minKeysPerNode, must be even. In addition, minKeysPerNode must be computed
// to be greater than one to ensure that during node merge operations, there is always
// at least one sibling node with which to merge. Hence, maxKeysPerNode must also be
// at least four.
//
// Note that if the B+Tree will reside only in memory, callback argument may be nil.
// That said, there is no advantage to using a B+Tree for an in-memory collection over
// the llrb-provided collection implementing the same APIs.
func NewBPlusTree(maxKeysPerNode uint64, compare Compare, callbacks BPlusTreeCallbacks, bPlusTreeCache BPlusTreeCache) (tree BPlusTree) {
	minKeysPerNode := maxKeysPerNode >> 1
	if (4 > maxKeysPerNode) || ((2 * minKeysPerNode) != maxKeysPerNode) {
		err := fmt.Errorf("maxKeysPerNode (%v) invalid - must be an even positive number greater than 3", maxKeysPerNode)
		panic(err)
	}

	rootNode := &btreeNodeStruct{
		objectNumber:        0, //                            To be filled in once root node is posted
		objectOffset:        0, //                            To be filled in once root node is posted
		objectLength:        0, //                            To be filled in once root node is posted
		items:               0,
		loaded:              true, //                         Special case in that objectNumber == 0 means it has no onDisk copy
		dirty:               true, //                         To be set just below
		root:                true,
		leaf:                true,
		tree:                nil, //                          To be set just below
		parentNode:          nil,
		kvLLRB:              NewLLRBTree(compare, callbacks),
		nonLeafLeftChild:    nil,
		rootPrefixSumChild:  nil,
		prefixSumItems:      0,   //                          Not applicable to root node
		prefixSumParent:     nil, //                          Not applicable to root node
		prefixSumLeftChild:  nil, //                          Not applicable to root node
		prefixSumRightChild: nil, //                          Not applicable to root node
	}

	rootNode.btreeNodeCacheElement.btreeNodeCacheTag = noLRU

	treePtr := &btreeTreeStruct{
		minKeysPerNode:            minKeysPerNode,
		maxKeysPerNode:            maxKeysPerNode,
		Compare:                   compare,
		BPlusTreeCallbacks:        callbacks,
		root:                      rootNode,
		staleOnDiskReferencesList: nil,
	}

	if nil == bPlusTreeCache {
		treePtr.nodeCache = nil
	} else {
		treePtr.nodeCache = bPlusTreeCache.(*btreeNodeCacheStruct)
	}

	rootNode.tree = treePtr

	treePtr.initNodeAsEvicted(rootNode)
	treePtr.markNodeDirty(rootNode)

	tree = treePtr

	return
}

// OldBPlusTree is used to re-construct a B+Tree previously persisted
func OldBPlusTree(rootObjectNumber uint64, rootObjectOffset uint64, rootObjectLength uint64, compare Compare, callbacks BPlusTreeCallbacks, bPlusTreeCache BPlusTreeCache) (tree BPlusTree, err error) {
	rootNode := &btreeNodeStruct{
		objectNumber:        rootObjectNumber,
		objectOffset:        rootObjectOffset,
		objectLength:        rootObjectLength,
		items:               0, //             To be filled in once root node is loaded
		loaded:              false,
		dirty:               false,
		root:                true,
		leaf:                true, //          To be updated once root node is loaded
		tree:                nil,  //          To be set just below
		parentNode:          nil,
		kvLLRB:              nil, //           To be filled in once root node is loaded
		nonLeafLeftChild:    nil, //           To be filled in once root node is loaded
		rootPrefixSumChild:  nil, //           To be filled in once root node is loaded (nil if root is also leaf)
		prefixSumItems:      0,   //           Not applicable to root node
		prefixSumParent:     nil, //           Not applicable to root node
		prefixSumLeftChild:  nil, //           Not applicable to root node
		prefixSumRightChild: nil, //           Not applicable to root node
	}

	rootNode.btreeNodeCacheElement.btreeNodeCacheTag = noLRU

	treePtr := &btreeTreeStruct{
		minKeysPerNode:            0, //              To be filled in once root node is loaded
		maxKeysPerNode:            0, //              To be filled in once root node is loaded
		Compare:                   compare,
		BPlusTreeCallbacks:        callbacks,
		root:                      rootNode,
		staleOnDiskReferencesList: nil,
	}

	if nil == bPlusTreeCache {
		treePtr.nodeCache = nil
	} else {
		treePtr.nodeCache = bPlusTreeCache.(*btreeNodeCacheStruct)
	}

	rootNode.tree = treePtr

	treePtr.initNodeAsEvicted(rootNode)

	tree = treePtr

	err = treePtr.loadNode(rootNode) // Return from loadNode() sufficient for return from this func

	return
}
