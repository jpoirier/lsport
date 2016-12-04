// Package lsport is a go wrapper for the cross-platform C library libserialport.
//
// Libserialport is a fairly thorough library and lsport only implements a subset
// of the full functionality available, it should however be fairly straightforward
// to add the remainder should you wish to do so.
// Documentation for libserialport can be found at http://sigrok.org/wiki/Libserialport
//
// Installation
//
// Using go get is (I think) impossible at the moment, at least in a cross platform
// way, as libserialport must be configured on the target system. So instead you
// should proceed as follows
//  git clone http://github.com/jpoirier/lsport
// into the src directory of your GOPATH. You'll then need to alter the #cgo CFLAGS
// & #cgo LDFLAGS in lsport.go so they are absolute (automation of absolute paths
// is coming to go1.5 as ${SRCDIR}) but for now edit /your/actual/gopath/ to whatever
// it is:
//  -I/your/actual/gopath/src/github.com/jpoirier/lsport/libserialport/
//  /your/actual/gopath/src/github.com/jpoirier/lsport/libserialport/.libs/libserialport.a
// and then build libserialport with something like:
//  cd libserialport
//  ./autogen.sh
//  ./configure
//  make
// before installing lsport with:
//  go install github.com/jpoirier/lsport
package lsport

/*
// Use absolute paths here, relative paths are coming to go1.5 as ${SRCDIR} ?
#cgo CFLAGS: -I/your/actual/gopath/src/github.com/jpoirier/lsport/libserialport/
#cgo LDFLAGS: /your/actual/gopath/src/github.com/jpoirier/lsport/libserialport/.libs/libserialport.a

// define _GNU_SOURCE for asprintf
#define _GNU_SOURCE
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <libserialport.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Term represents an asynchronous communications port.
type Term struct {
	port *C.struct_sp_port
}

// Return values for sp_port enums for errors and configurations
const (
	/** Special value to indicate setting should be left alone. */
	ParityIvalid = C.SP_PARITY_INVALID
	/** No parity. */
	ParityNone = C.SP_PARITY_NONE
	/** Odd parity. */
	ParityOdd = C.SP_PARITY_ODD
	/** Even parity. */
	ParityEven = C.SP_PARITY_EVEN
	/** Mark parity. */
	ParityMarkK = C.SP_PARITY_MARK
	/** Space parity. */
	ParitySpace = C.SP_PARITY_SPACE
	/** Open port for read access. */
	ModeRead = C.SP_MODE_READ
	/** Open port for write access. */
	ModeWrite = C.SP_MODE_WRITE
	/** Open port for read and write access. */
	ModeReadWrite = C.SP_MODE_READ_WRITE
	/** Flush input buffer. */
	BufInput = C.SP_BUF_INPUT
	/** Flush output buffer. */
	BufOutPut = C.SP_BUF_OUTPUT
	/** Flush both buffers. */
	BufBoth = C.SP_BUF_BOTH
)

func getError(err C.enum_sp_return) error {
	switch err {
	case C.SP_OK:
		return nil
	case C.SP_ERR_ARG:
		return errors.New("Invalid arguments were passed to the function")
	case C.SP_ERR_FAIL:
		return errors.New("A system error occured while executing the operation")
	case C.SP_ERR_MEM:
		return errors.New("A memory allocation failed while executing the operation")
	case C.SP_ERR_SUPP:
		return errors.New("The requested operation is not supported by this system or device")
	default:
		if err > 0 {
			return nil
		}
	}
	return errors.New("Error unknown")
}

// Open opens the specified serial port.
func Open(port string) (*Term, error) {
	var p *C.struct_sp_port
	cp := C.CString(port)
	defer C.free(unsafe.Pointer(cp))
	if err := getError(C.sp_get_port_by_name(cp, &p)); err != nil {
		return nil, err
	}
	if err := getError(C.sp_open(p, C.SP_MODE_READ_WRITE)); err != nil {
		return nil, err
	}
	if err := portConfig(p, 115200, 8, 1); err != nil {
		return nil, err
	}

	return &Term{port: p}, nil
}

// SetParams sets the common port options: baudrate, bits and stopbits.
func portConfig(p *C.struct_sp_port, baud int, bits int, stopbits int) error {
	if err := getError(C.sp_set_baudrate(p, C.int(baud))); err != nil {
		return err
	}
	if err := getError(C.sp_set_bits(p, C.int(bits))); err != nil {
		return err
	}
	if err := getError(C.sp_set_parity(p, C.SP_PARITY_NONE)); err != nil {
		return err
	}
	if err := getError(C.sp_set_stopbits(p, C.int(stopbits))); err != nil {
		return err
	}
	if err := getError(C.sp_set_flowcontrol(p, C.SP_FLOWCONTROL_NONE)); err != nil {
		return err
	}

	return nil
}

// Close closes the port.
func (t *Term) Close() error {
	err := getError(C.sp_close(t.port))
	C.sp_free_port(t.port)
	return err
}

// InputWaiting returns the number of bytes waiting in the input buffer.
func (t *Term) InputWaiting() (int32, error) {
	count := C.sp_input_waiting(t.port)
	return count, getError(count)
}

// OutputWaiting returns the number of bytes waiting in the output buffer.
func (t *Term) OutputWaiting() (int32, error) {
	count := C.sp_input_waiting(t.port)
	return count, getError(count)
}

// BlockingRead reads bytes from the serial port, blocking until complete.
// timeout in milliseconds, or zero to wait indefinitely.
func (t *Term) BlockingRead(rxBuf []byte, timeout uint) (int, error) {
	r := C.sp_blocking_read(t.port, (unsafe.Pointer(&rxBuf[0])), C.size_t(len(rxBuf)), C.uint(timeout))
	return int(r), getError(r)
}

// BlockingWrite writes bytes to the serial port, blocking until complete.
// timeout in milliseconds, or zero to wait indefinitely.
func (t *Term) BlockingWrite(txBuf []byte, timeout uint) (int, error) {
	r := C.sp_blocking_write(t.port, (unsafe.Pointer(&txBuf[0])), C.size_t(len(txBuf)), C.uint(timeout))
	return int(r), getError(r)
}

// Flush flushes serial port buffers,
// BufInput, BufOutput, BufBoth
func (t *Term) Flush(buffer C.enum_sp_buffer) error {
	return getError(C.sp_flush(t.port, buffer))
}

// Drain waits for buffered data to be transmitted.
func (t *Term) Drain() error {
	return getError(C.sp_drain(t.port))
}

// Read reads bytes from the serial port, without blocking.
func (t *Term) Read(rxBuf []byte) (int, error) {
	r := C.sp_nonblocking_read(t.port, (unsafe.Pointer(&rxBuf[0])), C.size_t(len(rxBuf)))
	return int(r), getError(r)
}

// Write writes bytes to the serial port, without blocking.
func (t *Term) Write(txBuf []byte) (int, error) {
	r := C.sp_nonblocking_write(t.port, (unsafe.Pointer(&txBuf[0])), C.size_t(len(txBuf)))
	return int(r), getError(r)
}
