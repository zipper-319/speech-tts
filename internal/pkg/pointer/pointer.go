package pointer

// #include <stdlib.h>
import "C"
import (
	"github.com/pkg/errors"
	"sync"
	"unsafe"
)

var (
	mutex sync.RWMutex
	store = map[int32]interface{}{}
	id    = int32(0)
)

func Save(v interface{}) (unsafe.Pointer, error) {
	if v == nil {
		return nil, errors.New("tts object is null")
	}

	// Generate real fake C pointer.
	// This pointer will not store any data, but will bi used for indexing purposes.
	// Since Go doest allow to cast dangling pointer to unsafe.Pointer, we do rally allocate one byte.
	// Why we need indexing, because Go doest allow C code to store pointers to Go data.
	//var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	//if ptr == nil {
	//	panic("can't allocate 'cgo-pointer hack index pointer': ptr == nil")
	//}
	mutex.Lock()
	defer mutex.Unlock()
	id += 1
	ptr := unsafe.Pointer(uintptr(id))

	if _, ok := store[id]; ok {
		return nil, errors.New("cgo-pointer has allocated")
	}
	store[id] = v

	return ptr, nil
}

func Load(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.RLock()
	v = store[int32(uintptr(ptr))]
	mutex.RUnlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(store, int32(uintptr(ptr)))
	mutex.Unlock()
}
