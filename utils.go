package main

import (
	"bytes"
	"log"
	"unsafe"
)

// String join from byte[]
// https://segmentfault.com/q/1010000006058923
// https://www.flysnow.org/2018/11/05/golang-concat-strings-performance-analysis.html
// https://www.vzhima.com/2020/01/07/golang-convert-a-byte-to-string.html
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// handle it if err nut nil
func logErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// BytesCombine join the bytes
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
