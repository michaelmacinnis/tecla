# Released under an MIT-style license. See the LICENSE file for details.

include $(GOROOT)/src/Make.inc

TARG=tecla
CGOFILES=tecla.go
#CGO_OFILES=handler_c.o
CGO_CFLAGS=-I/usr/local/lib
CGO_LDFLAGS=-ltecla -lcurses

include $(GOROOT)/src/Make.pkg
