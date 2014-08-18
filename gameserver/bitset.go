package main

type BitSet struct {
    data []uint64
    size uint
}

func NewBitSet(sz uint) BitSet {
    length := (sz / 64) + 1
    return BitSet{
        data: make([]uint64, length),
        size: sz,
    }
}

func (bs BitSet) IsFlag(ith uint) bool {
    idx := ith / 64
    off := ith % 64

    return (bs.data[idx] & (1 << off)) != 0
}

func (bs BitSet) IsDead(ith uint) bool {
    idx := ith / 64
    off := ith % 64

    return (bs.data[idx] & (1 << off)) == 0
}
