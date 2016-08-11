package sortedmap

import "testing"

func TestLLRBTreeAllButDeleteSimple(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestAllButDeleteSimple(t, llrb)
}

func TestLLRBTreeDeleteByIndexSimple(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestDeleteByIndexSimple(t, llrb)
}

func TestLLRBTreeDeleteByKeySimple(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestDeleteByKeySimple(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByIndexTrivial(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByIndexTrivial(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByIndexSmall(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByIndexSmall(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByIndexLarge(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByIndexLarge(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByIndexHuge(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByIndexHuge(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByKeyTrivial(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByKeyTrivial(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByKeySmall(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByKeySmall(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByKeyLarge(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByKeyLarge(t, llrb)
}

func TestLLRBTreeInsertGetDeleteByKeyHuge(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestInsertGetDeleteByKeyHuge(t, llrb)
}

func TestLLRBTreeBisect(t *testing.T) {
	llrb := NewLLRBTree(CompareInt)
	metaTestBisect(t, llrb)
}
