package sortedmap

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type testBPlusTreeStruct struct {
	sync.Mutex
	nextObjectNumber uint64
	objectMap        map[uint64][]byte
}

func (tree *testBPlusTreeStruct) DumpKey(key Key) (keyAsString string, err error) {
	var (
		ok          bool
		keyAsUint16 uint16
	)

	keyAsUint16, ok = key.(uint16)
	if !ok {
		err = fmt.Errorf("DumpKey() expected key of type uint16... instead it was of type %v", reflect.TypeOf(key))
		return
	}

	keyAsString = fmt.Sprintf("0x%04X", keyAsUint16)

	err = nil
	return
}

func (tree *testBPlusTreeStruct) DumpValue(value Value) (valueAsString string, err error) {
	var (
		ok            bool
		valueAsUint32 uint32
	)

	valueAsUint32, ok = value.(uint32)
	if !ok {
		err = fmt.Errorf("DumpValue() expected value of type uint32... instead it was of type %v", reflect.TypeOf(value))
		return
	}

	valueAsString = fmt.Sprintf("0x%08X", valueAsUint32)

	err = nil
	return
}

func (tree *testBPlusTreeStruct) GetNode(objectNumber uint64, objectOffset uint64, objectLength uint64) (nodeByteSlice []byte, err error) {
	var (
		ok bool
	)

	tree.Lock()
	defer tree.Unlock()

	nodeByteSlice, ok = tree.objectMap[objectNumber]
	if !ok {
		err = fmt.Errorf("GetNode() called for non-existent objectNumber 0x%016X", objectNumber)
		return
	}
	if uint64(0) != objectOffset {
		err = fmt.Errorf("GetNode() called for non-zero objectOffset 0x%016X", objectOffset)
		return
	}
	if uint64(len(nodeByteSlice)) != objectLength {
		err = fmt.Errorf("GetNode() called for objectLength 0x%016X... for node of length 0x%016X", objectLength, uint64(len(nodeByteSlice)))
		return
	}

	err = nil
	return
}

func (tree *testBPlusTreeStruct) PutNode(nodeByteSlice []byte) (objectNumber uint64, objectOffset uint64, err error) {
	tree.Lock()
	defer tree.Unlock()

	objectNumber = tree.nextObjectNumber
	tree.nextObjectNumber++
	tree.objectMap[objectNumber] = nodeByteSlice

	objectOffset = 0

	err = nil
	return
}

func (tree *testBPlusTreeStruct) DiscardNode(objectNumber uint64, objectOffset uint64, objectLength uint64) (err error) {
	var (
		ok            bool
		nodeByteSlice []byte
	)

	tree.Lock()
	defer tree.Unlock()

	nodeByteSlice, ok = tree.objectMap[objectNumber]
	if !ok {
		err = fmt.Errorf("DiscardNode() called for non-existent objectNumber 0x%016X", objectNumber)
		return
	}
	if uint64(0) != objectOffset {
		err = fmt.Errorf("DiscardNode() called for non-zero objectOffset 0x%016X", objectOffset)
		return
	}
	if uint64(len(nodeByteSlice)) != objectLength {
		err = fmt.Errorf("DiscardNode() called for objectLength 0x%016X... for node of length 0x%016X", objectLength, uint64(len(nodeByteSlice)))
		return
	}

	err = nil
	return
}

func (tree *testBPlusTreeStruct) PackKey(key Key) (packedKey []byte, err error) {
	var (
		ok          bool
		keyAsUint16 uint16
	)

	keyAsUint16, ok = key.(uint16)
	if !ok {
		err = fmt.Errorf("PackKey() expected key of type uint16... instead it was of type %v", reflect.TypeOf(key))
		return
	}

	packedKey = make([]byte, 2)
	packedKey[0] = uint8(keyAsUint16 & uint16(0x00FF) >> 0)
	packedKey[1] = uint8(keyAsUint16 & uint16(0xFF00) >> 8)

	err = nil
	return
}

func (tree *testBPlusTreeStruct) UnpackKey(payloadData []byte) (key Key, bytesConsumed uint64, err error) {
	if len(payloadData) != 2 {
		err = fmt.Errorf("UnpackKey() called for length %v... expected length of 2", len(payloadData))
		return
	}

	key = uint32(payloadData[0]) | (uint32(payloadData[1]) << 8)

	bytesConsumed = 2

	err = nil
	return
}

func (tree *testBPlusTreeStruct) PackValue(value Value) (packedValue []byte, err error) {
	var (
		ok            bool
		valueAsUint32 uint32
	)

	valueAsUint32, ok = value.(uint32)
	if !ok {
		err = fmt.Errorf("PackValue() expected value of type uint32... instead it was of type %v", reflect.TypeOf(value))
		return
	}

	packedValue = make([]byte, 4)
	packedValue[0] = uint8(valueAsUint32 & uint32(0x000000FF) >> 0)
	packedValue[1] = uint8(valueAsUint32 & uint32(0x0000FF00) >> 8)
	packedValue[2] = uint8(valueAsUint32 & uint32(0x00FF0000) >> 16)
	packedValue[3] = uint8(valueAsUint32 & uint32(0xFF000000) >> 24)

	err = nil
	return
}

func (tree *testBPlusTreeStruct) UnpackValue(payloadData []byte) (value Value, bytesConsumed uint64, err error) {
	if len(payloadData) != 4 {
		err = fmt.Errorf("UnpackValue() called for length %v... expected length of 4", len(payloadData))
		return
	}

	value = uint32(payloadData[0]) | (uint32(payloadData[1]) << 8) | (uint32(payloadData[2]) << 16) | (uint32(payloadData[3]) << 24)

	bytesConsumed = 4

	err = nil
	return
}

func TestBPlusTreeCache(t *testing.T) {
	var (
		tree       BPlusTree // map[uint16]uint32
		treeCache  BPlusTreeCache
		treeStruct *testBPlusTreeStruct
	)

	treeStruct = &testBPlusTreeStruct{
		nextObjectNumber: uint64(0),
		objectMap:        make(map[uint64][]byte),
	}

	treeCache = NewBPlusTreeCache(0, 1000)

	tree = NewBPlusTree(4, CompareUint16, treeStruct, treeCache)

	_, _ = tree.Put(uint16(0x0000), uint32(0x00000000))
	_, _ = tree.Put(uint16(0x0001), uint32(0x00000001))
	_, _ = tree.Put(uint16(0x0002), uint32(0x00000002))
	_, _ = tree.Put(uint16(0x0003), uint32(0x00000003))
	_, _ = tree.Put(uint16(0x0004), uint32(0x00000004))
	_, _ = tree.Put(uint16(0x0005), uint32(0x00000005))
	_, _ = tree.Put(uint16(0x0006), uint32(0x00000006))
	_, _ = tree.Put(uint16(0x0007), uint32(0x00000007))
	_, _ = tree.Put(uint16(0x0008), uint32(0x00000008))
	_, _ = tree.Put(uint16(0x0009), uint32(0x00000009))
	_, _ = tree.Put(uint16(0x000A), uint32(0x0000000A))
	_, _ = tree.Put(uint16(0x000B), uint32(0x0000000B))
	_, _ = tree.Put(uint16(0x000C), uint32(0x0000000C))
	_, _ = tree.Put(uint16(0x000D), uint32(0x0000000D))
	_, _ = tree.Put(uint16(0x000E), uint32(0x0000000E))
	_, _ = tree.Put(uint16(0x000F), uint32(0x0000000F))

	// TODO: Need to finish TestBPlusTreeCache()
}
