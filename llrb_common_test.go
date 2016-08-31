package sortedmap

import (
	"fmt"
	"testing"
)

type commonLLRBTreeContextStruct struct {
	t    *testing.T
	tree LLRBTree
}

func (context *commonLLRBTreeContextStruct) DumpKey(key Key) (keyAsString string, err error) {
	keyAsInt, ok := key.(int)
	if !ok {
		context.t.Fatalf("DumpKey() argument not an int")
	}
	keyAsString = fmt.Sprintf("%v", keyAsInt)
	err = nil
	return
}

func (context *commonLLRBTreeContextStruct) DumpValue(value Value) (valueAsString string, err error) {
	valueAsString, ok := value.(string)
	if !ok {
		context.t.Fatalf("PackValue() argument not a string")
	}
	err = nil
	return
}

func TestLLRBTreeAllButDeleteSimple(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestAllButDeleteSimple(t, context.tree)
}

func TestLLRBTreeDeleteByIndexSimple(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestDeleteByIndexSimple(t, context.tree)
}

func TestLLRBTreeDeleteByKeySimple(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestDeleteByKeySimple(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByIndexTrivial(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByIndexTrivial(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByIndexSmall(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByIndexSmall(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByIndexLarge(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByIndexLarge(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByIndexHuge(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByIndexHuge(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByKeyTrivial(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByKeyTrivial(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByKeySmall(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByKeySmall(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByKeyLarge(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByKeyLarge(t, context.tree)
}

func TestLLRBTreeInsertGetDeleteByKeyHuge(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestInsertGetDeleteByKeyHuge(t, context.tree)
}

func TestLLRBTreeBisect(t *testing.T) {
	context := &commonLLRBTreeContextStruct{t: t}
	context.tree = NewLLRBTree(CompareInt, context)
	metaTestBisect(t, context.tree)
}
