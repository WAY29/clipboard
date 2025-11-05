// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build ios

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework MobileCoreServices

#import <stdlib.h>
void clipboard_write_string(char *s);
char *clipboard_read_string();
unsigned int clipboard_read_image(void **out);
unsigned int clipboard_read_image_jpeg(void **out);
int clipboard_write_image(const void *bytes, unsigned long n);
int clipboard_write_image_jpeg(const void *bytes, unsigned long n);
*/
import "C"
import (
	"bytes"
	"context"
	"time"
	"unsafe"
)

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	var (
		data unsafe.Pointer
		n    C.uint
	)
	switch t {
	case FmtText:
		return []byte(C.GoString(C.clipboard_read_string())), nil
	case FmtImage:
		n = C.clipboard_read_image(&data)
	case FmtImageJPEG:
		n = C.clipboard_read_image_jpeg(&data)
	default:
		return nil, errUnsupported
	}
	if t == FmtImage || t == FmtImageJPEG {
		if data == nil {
			return nil, errUnavailable
		}
		defer C.free(unsafe.Pointer(data))
		if n == 0 {
			return nil, nil
		}
		return C.GoBytes(data, C.int(n)), nil
	}
	return nil, errUnsupported
}

// SetContent sets the clipboard content for iOS
func write(t Format, buf []byte) (<-chan struct{}, error) {
	done := make(chan struct{}, 1)
	var ok C.int
	switch t {
	case FmtText:
		cs := C.CString(string(buf))
		defer C.free(unsafe.Pointer(cs))

		C.clipboard_write_string(cs)
		return done, nil
	case FmtImage:
		if len(buf) == 0 {
			ok = C.clipboard_write_image(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_image(unsafe.Pointer(&buf[0]),
				C.ulong(len(buf)))
		}
	case FmtImageJPEG:
		if len(buf) == 0 {
			ok = C.clipboard_write_image_jpeg(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_image_jpeg(unsafe.Pointer(&buf[0]),
				C.ulong(len(buf)))
		}
	default:
		return nil, errUnsupported
	}
	if ok != 0 {
		return nil, errUnavailable
	}
	return done, nil
}

func watch(ctx context.Context, t Format) <-chan []byte {
	recv := make(chan []byte, 1)
	ti := time.NewTicker(time.Second)
	last := Read(t)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				b := Read(t)
				if b == nil {
					continue
				}
				if bytes.Compare(last, b) != 0 {
					recv <- b
					last = b
				}
			}
		}
	}()
	return recv
}
