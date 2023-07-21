package pointer

// #include <stdlib.h>
import "C"
import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	mutex sync.RWMutex
	store = map[unsafe.Pointer]interface{}{}
	id    = int32(time.Now().UnixNano())
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
	atomic.AddInt32(&id, 1)
	ptr := unsafe.Pointer(uintptr(id))
	mutex.Lock()
	if _, ok := store[ptr]; ok {
		return nil, errors.New("cgo-pointer has allocated")
	}
	store[ptr] = v
	mutex.Unlock()

	return ptr, nil
}

func Load(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.RLock()
	v = store[ptr]
	mutex.RUnlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(store, ptr)
	mutex.Unlock()
}
