package buzhash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

// Test the rolling hash property of the buzhash
func TestRollingHash64(t *testing.T) {
	phrase1 := "Aenean massa. Cum sociis natoque"
	phrase2 := "Phasellus leo dolor, tempus non, auctor et, hendrerit quis, nisi"

	hasher1 := NewBuzHash64(32)
	fmt.Fprint(hasher1, phrase1)
	p1sum := hasher1.Sum64()

	hasher1.Reset()
	found := false
	for idx, b := range []byte(loremipsum1) {
		ssum := hasher1.HashByte(b)
		if (ssum == p1sum) && (idx-32 == 91) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find '%s' by its checksum %08x.", phrase1, p1sum)
	}

	hasher2 := NewBuzHash64(64)
	fmt.Fprint(hasher2, phrase2)
	p2sum := hasher2.Sum64()

	hasher2.Reset()
	found = false
	for idx, b := range []byte(loremipsum2) {
		ssum := hasher2.HashByte(b)
		if (ssum == p2sum) && (idx-64 == 592) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find '%s' by its checksum %08x.", phrase2, p2sum)
	}
}

func TestSum64(t *testing.T) {
	h := NewBuzHash64(16)
	fmt.Fprint(h, loremipsum1)

	sum32 := h.Sum64()
	sumBytes := h.Sum(nil)

	if l := len(sumBytes); l != 8 {
		t.Fatalf("h.Sum() returned slice of len %d, expected 4", l)
	}

	var sumBytesAsNum uint64
	if err := binary.Read(bytes.NewBuffer(sumBytes), binary.LittleEndian, &sumBytesAsNum); err != nil {
		t.Fatalf("Could not read binary number? %s", err)
	}

	if sum32 != sumBytesAsNum {
		t.Errorf("Sum32 (%08x) and Sum (%08x) returned different sums!", sum32, sumBytesAsNum)
	}
}

func TestWrite64(t *testing.T) {
	h1 := NewBuzHash64(16)
	h2 := NewBuzHash64(16)

	data := []byte(loremipsum1)
	for _, b := range data {
		h1.HashByte(b)
	}

	h2.Write(data)

	if h2.Sum64() != h1.Sum64() {
		t.Errorf(" got %x want %x", h2.Sum64(), h1.Sum64())
	}
}

func BenchmarkWrite(b *testing.B) {
	h := NewBuzHash64(16)

	data := make([]byte, b.N)

	b.ResetTimer()

	h.Write(data)

}
