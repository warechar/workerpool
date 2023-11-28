package deque

import "fmt"

type Seq[T any] interface {
	Get() int32
	Compare(T) bool
}

type DequeSeq[T Seq[any]] struct {
	buf   []T
	head  int
	tail  int
	count int
	cap   int
}

func NewSeq[T Seq[any]]() *DequeSeq[T] {
	return &DequeSeq[T]{
		buf: make([]T, minCapacity),
		cap: minCapacity,
	}
}

func (d *DequeSeq[T]) Len() int {
	return d.count
}

func (d *DequeSeq[T]) Cap() int {
	return len(d.buf)
}

func (d *DequeSeq[T]) Push(elem T, leftIdx, rightIdx int) {

	d.ifFull()
	fmt.Println(cap(d.buf))

	if d.count == 0 {
		d.buf[d.tail] = elem
		d.tail = 1
		d.count++
		return
	}

	// 如果head 小于elem
	if !d.buf[leftIdx].Compare(elem) {
		// 插入最左边
		if leftIdx == 0 {
			fmt.Println(cap(d.buf))
			newBuf := make([]T, cap(d.buf))
			newBuf[0] = elem
			copy(newBuf[d.head+1:], d.buf[:d.tail])
			d.buf = newBuf
			d.tail++
		}

		d.count++
		return
	}

	// 如果右边大于
	if d.buf[rightIdx].Compare(elem) {
		d.buf[d.tail] = elem
		d.tail++
		d.count++
		return
	}

	if (rightIdx-leftIdx) == 1 && d.buf[leftIdx].Compare(elem) && !d.buf[rightIdx].Compare(elem) {
		fmt.Println(leftIdx, rightIdx)
		newb := make([]T, cap(d.buf))
		copy(newb, d.buf[:leftIdx+1])
		newb[leftIdx+1] = elem
		copy(newb[leftIdx+2:], d.buf[rightIdx:])
		d.buf = newb
		d.tail++
		d.count++
		return
	}

	length := leftIdx + (rightIdx-leftIdx)/2
	middleVal := d.buf[length]

	if middleVal.Compare(elem) {
		d.Push(elem, length, rightIdx)
		return
	} else {

		d.Push(elem, leftIdx, length)
		return
	}
}

func (d *DequeSeq[T]) Pop() T {
	if d.count <= 0 {
		panic("deque: Pop() called on empty queue")
	}
	//var null T
	ret := d.buf[d.head]
	copy(d.buf[d.head:], d.buf[d.head+1:])
	d.tail--
	d.count--

	d.ifCompress()
	return ret
}

func (d *DequeSeq[T]) ifFull() {
	fmt.Println(d.count, cap(d.buf))
	if d.count != cap(d.buf) {
		return
	}

	if cap(d.buf) == 0 {
		if d.cap == 0 {
			d.cap = minCapacity
		}
		d.buf = make([]T, d.cap)
		return
	}

	d.resize()
}

// pop后判断容量是否合适，如果过大则进行压缩
func (d *DequeSeq[T]) ifCompress() {
	if d.Cap() > d.cap && d.count<<2 == d.Cap() {
		d.resize()
	}
}

// 只有full 或者 count = 1/4的时候触发
func (d *DequeSeq[T]) resize() {

	newBuf := make([]T, d.count<<1)

	if d.tail > d.head {
		copy(newBuf, d.buf[d.head:d.tail])
	} else {
		n := copy(newBuf, d.buf[d.head:])
		copy(newBuf[n:], d.buf[:d.tail])
	}

	d.head = 0
	d.tail = d.count
	d.buf = newBuf
}

func (d *DequeSeq[T]) Front() T {
	return d.buf[d.head]
}

func (d *DequeSeq[T]) GetHead() int {
	return d.head
}

func (d *DequeSeq[T]) GetTail() int {
	return d.tail
}

func (d *DequeSeq[T]) G() int {
	return cap(d.buf)
}
