package roaring

import (
	"github.com/liwnn/bitset"
)

const (
	ArrayDefaultMaxSize = 4096
	BitmapMaxCapacity   = 65536
)

type sortArray []uint16

func (a sortArray) find(x uint16) (int, bool) {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}

	if i > 0 && a[i-1] == x {
		return i - 1, true
	}
	return i, false
}

func (a *sortArray) remove(x uint16) {
	i, j := 0, len(*a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if (*a)[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}

	if i > 0 && (*a)[i-1] == x {
		*a = append((*a)[:i-1], (*a)[i:]...)
	}
}

func (a *sortArray) insertAt(i int, x uint16) {
	(*a) = append((*a), 0)
	copy((*a)[i+1:], (*a)[i:])
	(*a)[i] = x
}

func (a sortArray) len() int {
	return len(a)
}

type Container interface {
	add(x uint16)
	contains(x uint16) bool
	remove(x uint16)
}

type RoaringArray struct {
	keys   sortArray
	values []Container
}

func (ra *RoaringArray) setContainer(i int, c Container) {
	ra.values[i] = c
}

func (ra *RoaringArray) getIndex(xa uint16) (int, bool) {
	return ra.keys.find(xa)
}

type ArrayContainer struct {
	content sortArray
}

func newArrayContainer() *ArrayContainer {
	return &ArrayContainer{
		content: make(sortArray, 4),
	}
}

func (ac *ArrayContainer) add(x uint16) {
	i, found := ac.content.find(x)
	if found {
		return
	}
	ac.content.insertAt(i, x)
}

func (ac *ArrayContainer) contains(x uint16) bool {
	_, found := ac.content.find(x)
	return found
}

func (ac *ArrayContainer) remove(x uint16) {
	ac.content.remove(x)
}

type BitmapContainer struct {
	bitmap *bitset.BitSet
}

func newBitmapContainer() *BitmapContainer {
	return &BitmapContainer{
		bitmap: bitset.NewSize(BitmapMaxCapacity),
	}
}

func (ac BitmapContainer) add(x uint16) {
	ac.bitmap.Set(uint64(x))
}

func (ac BitmapContainer) contains(x uint16) bool {
	return ac.bitmap.Get(uint64(x))
}

func (ac BitmapContainer) remove(x uint16) {
	ac.bitmap.Clear(uint64(x))
}

type RunContainer struct {
	valueslength []uint64
}

type RoaringBitmap struct {
	highLowConatiner RoaringArray
}

func New() *RoaringBitmap {
	return &RoaringBitmap{}
}

func (rb *RoaringBitmap) Add(x uint32) {
	hb := uint16(x >> 16)
	lb := uint16(x)
	ra := &rb.highLowConatiner
	i, found := ra.getIndex(hb)
	if found {
		switch c := ra.values[i].(type) {
		case *ArrayContainer:
			if c.content.len() < ArrayDefaultMaxSize {
				c.add(lb)
			} else {
				if c.contains(lb) {
					return
				}
				bc := newBitmapContainer()
				for _, v := range c.content {
					bc.add(v)
				}
				bc.add(lb)
				ra.setContainer(i, bc)
			}
		}
	} else {
		ra.keys = append(ra.keys, 0)
		ra.values = append(ra.values, nil)
		copy(ra.keys[i+1:], ra.keys[i:])
		copy(ra.values[i+1:], ra.values[i:])
		newac := newArrayContainer()
		newac.add(lb)
		ra.keys[i] = hb
		ra.setContainer(i, newac)
	}
}

func (rb *RoaringBitmap) Contains(x uint32) bool {
	hb := uint16(x >> 16)
	lb := uint16(x)
	ra := &rb.highLowConatiner
	i, found := ra.getIndex(hb)
	if found {
		return ra.values[i].contains(lb)
	}
	return false
}

func (rb *RoaringBitmap) Remove(x uint32) {
	hb := uint16(x >> 16)
	ra := &rb.highLowConatiner
	i, found := ra.getIndex(hb)
	if found {
		lb := uint16(x)
		ra.values[i].remove(lb)
	}
}
