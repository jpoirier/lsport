Package lsport is a go wrapper for the cross-platform C library libserialport.

Libserialport is a fairly thorough library and lsport currently only implements
a subset of the full functionality available: about enough for a basic connection
without modifying DTR and CTS. It should however be fairly straightforward to 
add the remainder should you wish to do so.

Documentation for libserialport can be found at http://sigrok.org/wiki/Libserialport

Installation

Using go get is, I think, impossible at the moment, at least in a cross platform
way, as libserialport must be configured on the target system. So instead you 
should proceed as follows

	git clone http://github.com/kezl/lsport

into the src directory of your GOPATH. 

Building libserialport
######################
Build libserialport with something like:

	cd libserialport
	./autogen.sh
	./configure
	make

The desired outcome is a static libserialport.a that can be linked into your
binary by cgo/go, so there is no need for a make install step unless you intend to 
have libserialport available throughout your system.

Windows builds require the MinGW-w64 toolchain which is probably best used in 
conjunction with msys2.

http://mingw-w64.org/doku.php
http://msys2.github.io/

libserial port is under the GNU LESSER GENERAL PUBLIC LICENSE
see COPYING in the libserialport directory.

Installing lsport
#################
You'll then need to alter the #cgo CFLAGS & #cgo LDFLAGS in lsport.go so they are 
absolute, automation of absolute paths is coming to go1.5 as ${SRCDIR}, but for 
now edit /your/actual/gopath/ below to whatever it really is:

#cgo CFLAGS: -I/your/actual/gopath/src/github.com/kezl/lsport/libserialport/
#cgo LDFLAGS: /your/actual/gopath/src/github.com/kezl/lsport/libserialport/.libs/libserialport.a

 
Finally install lsport with:

	go install github.com/kezl/lsport
