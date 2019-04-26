package memory

// Echoer echos the read and write requests to the destination address
type Echoer struct {
	mmu  *MMU
	dest uint16
}

// NewEchoer creates Echoer instance
func NewEchoer(mmu *MMU, dest uint16) *Echoer {
	e := Echoer{mmu: mmu, dest: dest}
	return &e
}

// Read echos the request to the destination address
func (e *Echoer) Read(addr uint16) (byte, error) {
	return e.mmu.Read(e.dest)
}

// Write echos the request to the destination address
func (e *Echoer) Write(addr uint16, data byte) error {
	return e.mmu.Write(e.dest, data)
}
