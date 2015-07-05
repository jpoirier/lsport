
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
//  git clone http://github.com/kezl/lsport
// into the src directory of your GOPATH. You'll then need to alter the #cgo CFLAGS 
// & #cgo LDFLAGS in lsport.go so they are absolute (automation of absolute paths
// is coming to go1.5 as ${SRCDIR}) but for now edit /your/actual/gopath/ to whatever
// it is:
//  #cgo CFLAGS: -I/your/actual/gopath/src/github.com/kezl/lsport/libserialport/
//  #cgo LDFLAGS: /your/actual/gopath/src/github.com/kezl/lsport/libserialport/.libs/libserialport.a
// and then build libserialport with something like:
//  cd libserialport
//  ./autogen.sh
//  ./configure
//  make
// before installing lsport with:
//  go install github.com/kezl/lsport
package lsport

/*
// Use absolute paths here, relative paths are coming to go1.5 as ${SRCDIR} ?
#cgo CFLAGS: -I/home/kerry/Coding/gocode/src/github.com/kezl/lsport/libserialport/
#cgo LDFLAGS: /home/kerry/Coding/gocode/src/github.com/kezl/lsport/libserialport/.libs/libserialport.a

// define _GNU_SOURCE for asprintf
#define _GNU_SOURCE
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <libserialport.h>

static int append(char **str, const char *buf, int size) {
  char *nstr;
  if (*str == NULL) {
    nstr = malloc(size + 1);
    memcpy(nstr, buf, size);
    nstr[size] = '\0';
  }
  else {
    if (asprintf(&nstr, "%s%.*s", *str, size, buf) == -1) return -1;
    free(*str);
  }
  *str = nstr;
  return 0;
}

char *listPorts() {
	char *str = NULL;
	struct sp_port **ports = NULL;
    struct sp_port *port = NULL;
	int result = sp_list_ports(&ports);

    int i = 0;
    while (ports[i] != NULL) {
        port = ports[i];
		char* name = sp_get_port_name(port);
		append(&str, name, strlen(name));
		append(&str, "#", 1);
        i++;
    }
	return str;
}


*/
import "C"

// NB no blank lines between */ and import "C" ;-)

import (
	"errors"
	"strings"
	"unsafe"
)

// Return values for sp_port enums for errors and configurations
const (
	/** Operation completed successfully. */
	SP_OK = 0
	/** Invalid arguments were passed to the function. */
	SP_ERR_ARG = -1
	/** A system error occured while executing the operation. */
	SP_ERR_FAIL = -2
	/** A memory allocation failed while executing the operation. */
	SP_ERR_MEM = -3
	/** The requested operation is not supported by this system or device. */
	SP_ERR_SUPP = -4
	/** Special value to indicate setting should be left alone. */
	SP_PARITY_INVALID = -1
	/** No parity. */
	SP_PARITY_NONE = 0
	/** Odd parity. */
	SP_PARITY_ODD = 1
	/** Even parity. */
	SP_PARITY_EVEN = 2
	/** Mark parity. */
	SP_PARITY_MARK = 3
	/** Space parity. */
	SP_PARITY_SPACE = 4
	/** Open port for read access. */
	SP_MODE_READ = 1
	/** Open port for write access. */
	SP_MODE_WRITE = 2
	/** Open port for read and write access. */
	SP_MODE_READ_WRITE = 3
	/** Input buffer. */
	SP_BUF_INPUT = 1
	/** Output buffer. */
	SP_BUF_OUTPUT = 2
	/** Both buffers. */
	SP_BUF_BOTH = 3
)

// Port a pointer to an sp_port structure containing the file handle and
// port descriptors.
type Port *C.struct_sp_port

// Conf a struct containing a *C.struct_sp_port Port, the old port configuration 
// and the active port configuration (both of type *C.struct_sp_port_config).
type Conf struct {
	port      *C.struct_sp_port
	oldConfig *C.struct_sp_port_config
	newConfig *C.struct_sp_port_config
}

// checkResult error messages for sp_port return values.
func checkResult(result int32) error {
	if result == SP_OK {
		// Operation completed successfully.
		return nil
	} else if result == SP_ERR_ARG {
		return errors.New("Invalid arguments were passed to the function.\n")
	} else if result == SP_ERR_FAIL {
		return errors.New("A system error occured while executing the operation.\n")
	} else if result == SP_ERR_MEM {
		return errors.New("A memory allocation failed while executing the operation.\n")
	} else if result == SP_ERR_SUPP {
		return errors.New("The requested operation is not supported by this system or device.\n")
	}
	return nil
}

// Init attempts to initialises the port represented by the port name passed in.
func Init(s *Conf, name string) (int32, error) {
	var result int32 = SP_OK
	C.sp_new_config(&s.oldConfig)
	C.sp_new_config(&s.newConfig)
	cp := C.CString(name)
	defer C.free(unsafe.Pointer(cp))
	result = C.sp_get_port_by_name(cp, &s.port)
	checkResult(result)
	// Open before setting params
	result = C.sp_open(s.port, SP_MODE_READ|SP_MODE_WRITE)
	return result, checkResult(result)
}

// SetParams sets the common port options: baudrate, bits and stopbits.
func SetParams(s *Conf, baud int, bits int, stopbits int) (int32, error) {
	var result int32 = SP_OK
	C.sp_set_config_baudrate(s.newConfig, C.int(baud))
	C.sp_set_config_bits(s.newConfig, C.int(bits))
	C.sp_set_config_parity(s.newConfig, SP_PARITY_NONE)
	C.sp_set_config_stopbits(s.newConfig, C.int(stopbits))
	result = C.sp_set_config(s.port, s.newConfig)
	return result, checkResult(result)
}

// Close restores the port to how it was before initialisation, closes the port and
// frees resources.
func Close(s *Conf) {
	C.sp_flush(s.port, SP_BUF_BOTH)
	C.sp_free_config(s.newConfig)
	C.sp_set_config(s.port, s.oldConfig)
	C.sp_free_config(s.oldConfig)

	C.sp_close(s.port)
}

// PortSlice returns a string slice of serial ports found on the system.
func PortsSlice() ([]string, error) {
	portsStr := C.GoString(C.listPorts())
	ports := strings.Split(strings.Trim(portsStr, "#"), "#")
	if len(ports) == 0 {
		return ports, errors.New("No serial ports available on the system.\n")
	}
	
	return ports, nil
}

// minInt32 returns the lesser of two int32s.
func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Waiting returns the number of bytes waiting on success or an error code.
func Waiting(port Port) (int32, error) {
	var result int32 = SP_OK
	result = C.sp_input_waiting(port)
	return result, checkResult(result)
}

// BlockingRead attempts to read the lesser of the number of bytes waiting or the
// capacity of the rx buffer, blocks while reading, timeout is the number of mS
// to wait or set 0 to wait indefinitely.
// rxBuf size should not exceed the capacity of int32
func BlockingRead(port Port, rxBuf []byte, timeout uint) (int32, error) {
	var result int32 = SP_OK

	waiting, err := Waiting(port)
	if err == nil {
		// Passing to C so don't exceed the buffer length
		length := minInt32(waiting, int32(len(rxBuf)))

		if waiting > 0 {
			result = C.sp_blocking_read(port, (unsafe.Pointer(&rxBuf[0])), C.size_t(length), C.uint(timeout))
		}
		return result, checkResult(result)
	}
	return waiting, checkResult(waiting)
}

// BlockingWrite attempts to write the string or buffer supplied to the port
// blocks while writing, timeout is the number of mS to wait or set 0 to wait 
// indefinitely.
func BlockingWrite(port Port, txBuf []byte, timeout uint16) (int32, error) {
	var result int32 = SP_OK
	result = C.sp_blocking_write(port, (unsafe.Pointer(&txBuf[0])), C.size_t(len(txBuf)), C.uint(timeout))
	return result, checkResult(result)
}

// Read as BlockingRead but non-blocking and no timeout .
func Read(port Port, rxBuf []byte) (int32, error) {
	var result int32 = SP_OK
	waiting, err := Waiting(port)
	if err == nil {
		// Passing to C so don't exceed the buffer length
		length := minInt32(waiting, int32(len(rxBuf)))

		if waiting > 0 {
			result = C.sp_nonblocking_read(port, (unsafe.Pointer(&rxBuf[0])), C.size_t(length))
		}
		return result, checkResult(result)
	}
	return waiting, checkResult(waiting)

}

// Write as BlockingWrite but non-blocking and no timeout.
func Write(port Port, txBuf []byte) (int32, error) {
	var result int32 = SP_OK
	result = C.sp_nonblocking_write(port, (unsafe.Pointer(&txBuf[0])), C.size_t(len(txBuf)))
	return result, checkResult(result)
}
