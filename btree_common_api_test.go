package sortedmap

import (
	"encoding/binary"
	"fmt"
	"testing"
)

const (
	commonBPlusTreeTestNumKeysMaxSmall   = uint64(2)
	commonBPlusTreeTestNumKeysMaxModest  = uint64(10)
	commonBPlusTreeTestNumKeysMaxTypical = uint64(100)
	commonBPlusTreeTestNumKeysMaxLarge   = uint64(1000)
)

type commonContextStruct struct {
	t    *testing.T
	tree BPlusTree
}

func (context *commonContextStruct) GetNode(logSegmentNumber uint64, logOffset uint64, logLength uint64) (nodeByteSlice []byte, err error) {
	err = fmt.Errorf("GetNode() not implemented")
	return
}

func (context *commonContextStruct) PutNode(nodeByteSlice []byte) (logSegmentNumber uint64, logOffset uint64, err error) {
	err = fmt.Errorf("PutNode() not implemented")
	return
}

func (context *commonContextStruct) PackKey(key Key) (packedKey []byte, err error) {
	keyAsInt, ok := key.(int)
	if !ok {
		context.t.Fatalf("PackKey() argument not an int")
	}
	keyAsUint64 := uint64(keyAsInt)
	packedKey = make([]byte, 8)
	binary.LittleEndian.PutUint64(packedKey, keyAsUint64)
	err = nil
	return
}

func (context *commonContextStruct) UnpackKey(payloadData []byte) (key Key, bytesConsumed uint64, err error) {
	context.t.Fatalf("UnpackKey() not implemented")
	return
}

func (context *commonContextStruct) PackValue(value Value) (packedValue []byte, err error) {
	valueAsString, ok := value.(string)
	if !ok {
		context.t.Fatalf("PackValue() argument not a string")
	}
	packedValue = []byte(valueAsString)
	err = nil
	return
}

func (context *commonContextStruct) UnpackValue(payloadData []byte) (value Value, bytesConsumed uint64, err error) {
	context.t.Fatalf("UnpackValue() not implemented")
	return
}

func TestBPlusTreeAllButDeleteSimple(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestAllButDeleteSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestAllButDeleteSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestAllButDeleteSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestAllButDeleteSimple(t, context.tree)
}

func TestBPlusTreeDeleteByIndexSimple(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestDeleteByIndexSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestDeleteByIndexSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestDeleteByIndexSimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestDeleteByIndexSimple(t, context.tree)
}

func TestBPlusTreeDeleteByKeySimple(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestDeleteByKeySimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestDeleteByKeySimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestDeleteByKeySimple(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestDeleteByKeySimple(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByIndexTrivial(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByIndexTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByIndexTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByIndexTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByIndexTrivial(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByIndexSmall(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByIndexSmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByIndexSmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByIndexSmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByIndexSmall(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByIndexLarge(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByIndexLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByIndexLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByIndexLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByIndexLarge(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByIndexHuge(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByIndexHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByIndexHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByIndexHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByIndexHuge(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByKeyTrivial(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByKeyTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByKeyTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByKeyTrivial(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByKeyTrivial(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByKeySmall(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByKeySmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByKeySmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByKeySmall(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByKeySmall(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByKeyLarge(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByKeyLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByKeyLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByKeyLarge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByKeyLarge(t, context.tree)
}

func TestBPlusTreeInsertGetDeleteByKeyHuge(t *testing.T) {
	context := &commonContextStruct{t: t}
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxSmall, CompareInt, context)
	metaTestInsertGetDeleteByKeyHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxModest, CompareInt, context)
	metaTestInsertGetDeleteByKeyHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxTypical, CompareInt, context)
	metaTestInsertGetDeleteByKeyHuge(t, context.tree)
	context.tree = NewBPlusTree(commonBPlusTreeTestNumKeysMaxLarge, CompareInt, context)
	metaTestInsertGetDeleteByKeyHuge(t, context.tree)
}
