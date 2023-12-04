package deque

type TimerInterface[T any] interface {
	Get() any
	Compare(T) bool
}

type DequeTimer[T TimerInterface[T]] struct {
	buf   []T
	head  int
	tail  int
	count int
	cap   int
}

func NewTimer[T TimerInterface[T]]() *DequeTimer[T] {
	return &DequeTimer[T]{
		buf: make([]T, minCapacity),
		cap: minCapacity,
	}
}

func (d *DequeTimer[T]) Len() int {
	return d.count
}

func (d *DequeTimer[T]) Cap() int {
	return len(d.buf)
}

func (d *DequeTimer[T]) Push(elem T, leftIdx, rightIdx int) {

	d.ifFull()

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

func (d *DequeTimer[T]) Pop() T {
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

func (d *DequeTimer[T]) ifFull() {

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
func (d *DequeTimer[T]) ifCompress() {
	if d.Cap() > d.cap && d.count<<2 == d.Cap() {
		d.resize()
	}
}

// 只有full 或者 count = 1/4的时候触发
func (d *DequeTimer[T]) resize() {

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

func (d *DequeTimer[T]) Front() T {
	return d.buf[d.head]
}

func (d *DequeTimer[T]) GetHead() int {
	return d.head
}

func (d *DequeTimer[T]) GetTail() int {
	return d.tail
}
