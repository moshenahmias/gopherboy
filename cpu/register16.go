/*H**********************************************************************
* FILENAME :        register16.go
*
* PACKAGE :			cpu
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package cpu

// Register16 is a 16bit register constructed from two 8bit registers
type Register16 struct {
	h Register8
	l Register8
}

// rotateLeft rotates the register to the left
func (r *Register16) rotateLeft() {
	r.rotateHighLeft()
	r.rotateLowLeft()
	if r.h.get()&0x01 != r.l.get()&0x01 {
		r.h.set(r.h.get() ^ 0x01)
		r.l.set(r.l.get() ^ 0x01)
	}
}

// rotateRight rotates the register to the right
func (r *Register16) rotateRight() {
	r.rotateHighRight()
	r.rotateLowRight()
	if r.h.get()&0x80 != r.l.get()&0x80 {
		r.h.set(r.h.get() ^ 0x80)
		r.l.set(r.l.get() ^ 0x80)
	}
}

// rotateHighLeft rotates the high part to the left
func (r *Register16) rotateHighLeft() {
	r.h.rotateLeft()
}

// rotateHighRight rotates the high part to the right
func (r *Register16) rotateHighRight() {
	r.h.rotateRight()
}

// rotateLowLeft rotates the low part to the left
func (r *Register16) rotateLowLeft() {
	r.l.rotateLeft()
}

// rotateLowRight rotates the low part to the right
func (r *Register16) rotateLowRight() {
	r.l.rotateRight()
}

// increment by 1
func (r *Register16) increment() {
	r.set(r.get() + 1)
}

// decrement by 1
func (r *Register16) decrement() {
	r.set(r.get() - 1)
}

// set the register's data
func (r *Register16) set(data uint16) {
	r.l.set(byte(data))
	r.h.set(byte(data >> 8))
}

// get the register's data
func (r *Register16) get() uint16 {
	retval := uint16(r.h.get())
	retval = (retval << 8) | uint16(r.l.get())
	return retval
}

// setHigh sets the register's high part
func (r *Register16) setHigh(data byte) {
	r.h.set(data)
}

// incrementHigh by 1
func (r *Register16) incrementhigh() {
	r.h.increment()
}

// decrementHigh by 1
func (r *Register16) decrementhigh() {
	r.h.decrement()
}

// high gets the register's high part
func (r *Register16) high() *Register8 {
	return &r.h
}

// highByte gets the register's high byte
func (r *Register16) highByte() byte {
	return r.h.get()
}

// setLow sets the register's low part
func (r *Register16) setLow(data byte) {
	r.l.set(data)
}

// incrementlow by 1
func (r *Register16) incrementlow() {
	r.l.increment()
}

// decrementlow by 1
func (r *Register16) decrementlow() {
	r.l.decrement()
}

// low gets the register's low part
func (r *Register16) low() *Register8 {
	return &r.l
}

// lowByte gets the register's high byte
func (r *Register16) lowByte() byte {
	return r.l.get()
}
