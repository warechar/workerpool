package deque

// 2^k  x % n == x & (n - 1)
const minCapacity = 8

var next = func(index int, len int) int {
	return (index + 1) & (len - 1)
}

var prev = func(index int, len int) int {
	return (index - 1) & (len - 1)
}

type Deque[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
	cap   int
}

func New[T any]() *Deque[T] {
	return &Deque[T]{
		buf: make([]T, minCapacity),
		cap: minCapacity,
	}
}

func (d *Deque[T]) Len() int {
	return d.count
}

func (d *Deque[T]) Cap() int {
	return len(d.buf)
}

func (d *Deque[T]) Push(elem T) {
	d.ifFull()

	d.buf[d.tail] = elem
	d.tail = next(d.tail, d.Cap())
	d.count++
}

func (d *Deque[T]) Pop() T {
	if d.count <= 0 {
		panic("deque: Pop() called on empty queue")
	}
	var null T
	ret := d.buf[d.head]
	d.buf[d.head] = null
	d.head = next(d.head, d.Cap())
	d.count--

	d.ifCompress()
	return ret
}

func (d *Deque[T]) ifFull() {
	if d.count != len(d.buf) {
		return
	}

	if len(d.buf) == 0 {
		if d.cap == 0 {
			d.cap = minCapacity
		}
		d.buf = make([]T, d.cap)
		return
	}

	d.resize()
}

// pop后判断容量是否合适，如果过大则进行压缩
func (d *Deque[T]) ifCompress() {
	if d.Cap() > d.cap && d.count<<2 == d.Cap() {
		d.resize()
	}
}

// 只有full 或者 count = 1/4的时候触发
func (d *Deque[T]) resize() {

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

func (d *Deque[T]) Front() T {
	return d.buf[d.head]
}
