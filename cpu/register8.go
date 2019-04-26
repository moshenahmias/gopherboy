package cpu

// Register8 is a 8bit register
type Register8 struct {
	data byte
}

// rotateLeft rotates the register to the left
func (r *Register8) rotateLeft() {
	r.data = (r.data << 1) | (r.data >> 7)
}

// rotateRight rotates the register to the right
func (r *Register8) rotateRight() {
	r.data = (r.data >> 1) | (r.data << 7)
}

// set the register's data
func (r *Register8) set(data byte) {
	r.data = data
}

// get the register's data
func (r *Register8) get() byte {
	return r.data
}

// increment by 1
func (r *Register8) increment() {
	r.data++
}

// decrement by 1
func (r *Register8) decrement() {
	r.data--
}
