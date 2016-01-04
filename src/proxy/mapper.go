package proxy

import (
	"initializers"
	"reflect"
	"strings"
	"sync"
)

var lock sync.RWMutex

func serviceLookup(name, version string) (string, error) {
	lock.RLock()
	target, err := getField(initializers.Registry(), name)[version]
	lock.RUnlock()

	if !err {
		return "", initializers.StatusNotFound
	}
	return target, nil
}

func getField(v *initializers.SoaRegistry, field string) map[string]string {
	f := reflect.Indirect(reflect.ValueOf(v)).FieldByName(strings.Title(field))
	if !f.IsValid() {
		return nil
	}
	h := f.Interface().(map[string]string)
	return h
}
