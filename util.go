package bloom

import "unsafe"

type stringStruct struct {
	str unsafe.Pointer
	len int
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

func memhashSliceKey(data []byte, k uint64) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&data))
	return uint64(memhash(ss.str, uintptr(k), uintptr(ss.len)))
}

func memhashStringKey(str string, k uint64) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(memhash(ss.str, uintptr(k), uintptr(ss.len)))
}
