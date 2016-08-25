package sortedmap

import (
	"encoding/binary"
	"fmt"
	"testing"
)

const specificBPlusTreeTestNumKeysMaxSmall = uint64(2)

type logSegmentChunkStruct struct {
	startingOffset uint64
	chunkByteSlice []byte
}

type specificContextStruct struct {
	t                             *testing.T
	lastLogSegmentNumberGenerated uint64
	lastLogOffsetGenerated        uint64
	logSegmentChunkMap            map[uint64]*logSegmentChunkStruct // Key == logSegmentNumber (only 1 chunk stored per LogSegment)
}

type valueStruct struct {
	u32 uint32
	s8  [8]byte
}

func (context *specificContextStruct) GetNode(logSegmentNumber uint64, logOffset uint64, logLength uint64) (nodeByteSlice []byte, err error) {
	logSegmentChunk, ok := context.logSegmentChunkMap[logSegmentNumber]

	if !ok {
		err = fmt.Errorf("logSegmentNumber not found")
		return
	}

	if logSegmentChunk.startingOffset != logOffset {
		err = fmt.Errorf("logOffset not found")
		return
	}

	if uint64(len(logSegmentChunk.chunkByteSlice)) != logLength {
		err = fmt.Errorf("logLength not found")
		return
	}

	nodeByteSlice = logSegmentChunk.chunkByteSlice
	err = nil

	return
}

func (context *specificContextStruct) PutNode(nodeByteSlice []byte) (logSegmentNumber uint64, logOffset uint64, err error) {
	context.lastLogSegmentNumberGenerated++
	logSegmentNumber = context.lastLogSegmentNumberGenerated

	context.lastLogOffsetGenerated += logSegmentNumber + uint64(len(nodeByteSlice))
	logOffset = context.lastLogOffsetGenerated

	logSegmentChunk := &logSegmentChunkStruct{
		startingOffset: logOffset,
		chunkByteSlice: nodeByteSlice,
	}

	context.logSegmentChunkMap[logSegmentNumber] = logSegmentChunk

	err = nil

	return
}

func (context *specificContextStruct) PackKey(key Key) (packedKey []byte, err error) {
	keyAsUint32, ok := key.(uint32)
	if !ok {
		context.t.Fatalf("PackKey() argument not a uint32")
	}
	packedKey = make([]byte, 4)
	binary.LittleEndian.PutUint32(packedKey, keyAsUint32)
	err = nil
	return
}

func (context *specificContextStruct) UnpackKey(packedKey []byte) (key Key, bytesConsumed uint64, err error) {
	if 4 > len(packedKey) {
		context.t.Fatalf("UnpackKey() called with insufficient packedKey size")
	}
	keyAsUint32 := binary.LittleEndian.Uint32(packedKey[:4])
	key = keyAsUint32
	bytesConsumed = 4
	err = nil
	return
}

func (context *specificContextStruct) PackValue(value Value) (packedValue []byte, err error) {
	valueAsValueStructPtr, ok := value.(valueStruct)
	if !ok {
		context.t.Fatalf("PackValue() argument not a valueStruct")
	}
	u32Packed := make([]byte, 4)
	binary.LittleEndian.PutUint32(u32Packed, valueAsValueStructPtr.u32)
	packedValue = make([]byte, 0, 12)
	packedValue = append(packedValue, u32Packed...)
	packedValue = append(packedValue, valueAsValueStructPtr.s8[:]...)
	err = nil
	return
}

func (context *specificContextStruct) UnpackValue(packedValue []byte) (value Value, bytesConsumed uint64, err error) {
	if 12 > len(packedValue) {
		context.t.Fatalf("UnpackValue() called with insufficient packedValue size")
	}
	valueAsUint32 := binary.LittleEndian.Uint32(packedValue[:4])
	var s8AsArray [8]byte
	copy(s8AsArray[:], packedValue[4:12])
	value = valueStruct{u32: valueAsUint32, s8: s8AsArray}
	bytesConsumed = 12
	err = nil
	return
}

func uint32To8ReplicaByteArray(u32 uint32) (b8 [8]byte) {
	// Assumes u32 < 0x100

	for i := 0; i < 8; i++ {
		b8[i] = byte(u32)
	}

	return
}

func TestBPlusTreeSpecific(t *testing.T) {
	var (
		btreeLen                   int
		err                        error
		layoutReportExpected       LayoutReport
		layoutReportReturned       LayoutReport
		logSegmentBytesExpected    uint64
		logSegmentBytesReturned    uint64
		logSegmentChunk            *logSegmentChunkStruct
		logSegmentNumber           uint64
		ok                         bool
		rootLogSegmentNumber       uint64
		rootLogOffset              uint64
		rootLogLength              uint64
		valueAsValueStructExpected valueStruct
		valueAsValueStructReturned valueStruct
		valueAsValueStructToInsert valueStruct
		valueAsValueReturned       Value
	)

	persistentContext := &specificContextStruct{t: t, lastLogSegmentNumberGenerated: 0, lastLogOffsetGenerated: 0, logSegmentChunkMap: make(map[uint64]*logSegmentChunkStruct)}

	btreeNew := NewBPlusTree(specificBPlusTreeTestNumKeysMaxSmall, CompareUint32, persistentContext)

	valueAsValueStructToInsert = valueStruct{u32: 5, s8: uint32To8ReplicaByteArray(5)}
	ok, err = btreeNew.Put(uint32(5), valueAsValueStructToInsert)
	if nil != err {
		t.Fatalf("btreeNew.Put(uint32(5) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.Put(uint32(5), valueAsValueStructToInsert).ok should have been true")
	}

	valueAsValueStructToInsert = valueStruct{u32: 3, s8: uint32To8ReplicaByteArray(3)}
	ok, err = btreeNew.Put(uint32(3), valueAsValueStructToInsert)
	if nil != err {
		t.Fatalf("btreeNew.Put(uint32(3) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.Put(uint32(3), valueAsValueStructToInsert).ok should have been true")
	}

	valueAsValueStructToInsert = valueStruct{u32: 7, s8: uint32To8ReplicaByteArray(7)}
	ok, err = btreeNew.Put(uint32(7), valueAsValueStructToInsert)
	if nil != err {
		t.Fatalf("btreeNew.Put(uint32(7) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.Put(uint32(7), valueAsValueStructToInsert)).ok should have been true")
	}

	rootLogSegmentNumber, rootLogOffset, rootLogLength, err = btreeNew.Flush(false)
	if nil != err {
		t.Fatalf("btreeNew.Flush(false) should not have failed")
	}

	valueAsValueReturned, ok, err = btreeNew.GetByKey(uint32(5))
	if nil != err {
		t.Fatalf("btreeNew.GetByKey(uint32(5)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.GetByKey(uint32(5)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 5, s8: uint32To8ReplicaByteArray(5)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeNew.GetByKey(uint32(5)).value should have been valueAsValueStructExpected")
	}

	rootLogSegmentNumber, rootLogOffset, rootLogLength, err = btreeNew.Flush(true)
	if nil != err {
		t.Fatalf("btreeNew.Flush(true) should not have failed")
	}

	valueAsValueReturned, ok, err = btreeNew.GetByKey(uint32(3))
	if nil != err {
		t.Fatalf("btreeNew.GetByKey(uint32(3)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.GetByKey(uint32(3)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 3, s8: uint32To8ReplicaByteArray(3)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeNew.GetByKey(uint32(3)).value should have been valueAsValueStructExpected")
	}

	layoutReportExpected = make(map[uint64]uint64)
	for logSegmentNumber, logSegmentChunk = range persistentContext.logSegmentChunkMap {
		logSegmentBytesExpected = uint64(len(logSegmentChunk.chunkByteSlice))
		layoutReportExpected[logSegmentNumber] = logSegmentBytesExpected // Note: assumes no chunks are stale
	}
	layoutReportReturned, err = btreeNew.FetchLayoutReport()
	if nil != err {
		t.Fatalf("btreeNew.FetchLayoutReport() should not have failed")
	}
	if len(layoutReportExpected) != len(layoutReportReturned) {
		t.Fatalf("btreeNew.FetchLayoutReport() returned unexpected LayoutReport")
	}
	for logSegmentNumber, logSegmentBytesReturned = range layoutReportReturned {
		logSegmentBytesExpected, ok = layoutReportExpected[logSegmentNumber]
		if (!ok) || (logSegmentBytesExpected != logSegmentBytesReturned) {
			t.Fatalf("btreeNew.FetchLayoutReport() returned unexpected LayoutReport")
		}
	}

	btreeNew.Purge()

	valueAsValueReturned, ok, err = btreeNew.GetByKey(uint32(7))
	if nil != err {
		t.Fatalf("btreeNew.GetByKey(uint32(7)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeNew.GetByKey(uint32(7)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 7, s8: uint32To8ReplicaByteArray(7)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeNew.GetByKey(uint32(3)).value should have been valueAsValueStructExpected")
	}

	btreeOld := OldBPlusTree(rootLogSegmentNumber, rootLogOffset, rootLogLength, CompareUint32, persistentContext)

	btreeLen, err = btreeOld.Len()
	if nil != err {
		t.Fatalf("btreeOld.Len() should not have failed")
	}
	if 3 != btreeLen {
		t.Fatalf("btreeOld.Len() should have been 3")
	}

	valueAsValueReturned, ok, err = btreeOld.GetByKey(uint32(5))
	if nil != err {
		t.Fatalf("btreeOld.GetByKey(uint32(5)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeOld.GetByKey(uint32(5)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 5, s8: uint32To8ReplicaByteArray(5)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeNew.GetByKey(uint32(5)).value should have been valueAsValueStructExpected")
	}

	valueAsValueReturned, ok, err = btreeOld.GetByKey(uint32(3))
	if nil != err {
		t.Fatalf("btreeOld.GetByKey(uint32(3)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeOld.GetByKey(uint32(3)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 3, s8: uint32To8ReplicaByteArray(3)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeOld.GetByKey(uint32(3)).value should have been valueAsValueStructExpected")
	}

	valueAsValueReturned, ok, err = btreeOld.GetByKey(uint32(7))
	if nil != err {
		t.Fatalf("btreeOld.GetByKey(uint32(7)) should not have failed")
	}
	if !ok {
		t.Fatalf("btreeOld.GetByKey(uint32(7)).ok should have been true")
	}
	valueAsValueStructReturned = valueAsValueReturned.(valueStruct)
	valueAsValueStructExpected = valueStruct{u32: 7, s8: uint32To8ReplicaByteArray(7)}
	if valueAsValueStructReturned != valueAsValueStructExpected {
		t.Fatalf("btreeOld.GetByKey(uint32(3)).value should have been valueAsValueStructExpected")
	}
}
