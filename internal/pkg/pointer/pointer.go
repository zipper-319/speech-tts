package pointer

// #include <stdlib.h>
import "C"
import (
	"github.com/pkg/errors"
	"sync"
)

var (
	mutex sync.RWMutex
	store = map[int32]interface{}{}
	id    = int32(0)
)

func Save(v interface{}) (int32, error) {
	if v == nil {
		return 0, errors.New("tts object is null")
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

	if _, ok := store[id]; ok {
		return 0, errors.New("cgo-pointer has allocated")
	}
	store[id] = v

	return id, nil
}

func Load(id int32) (v interface{}) {

	mutex.RLock()
	v = store[id]
	mutex.RUnlock()
	return
}

func Unref(id int32) {

	mutex.Lock()
	delete(store, id)
	mutex.Unlock()
}
