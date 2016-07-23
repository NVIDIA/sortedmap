package sortedmap

type LLRBTree interface {
	SortedMap
}

func NewLLRBTree(compare Compare) (tree LLRBTree) {
	tree = &llrbTreeStruct{Compare: compare, root: nil}
	return
}
