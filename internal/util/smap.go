package util

import "sync"

type SMap[K comparable, V any] sync.Map

func (m *SMap[K, V]) Delete(key K) {
	(*sync.Map)(m).Delete(key)
}

func (m *SMap[K, V]) Load(key K) (v V, ok bool) {
	i, ok := (*sync.Map)(m).Load(key)
	if !ok {
		return
	}
	return i.(V), true
}

func (m *SMap[K, V]) LoadAndDelete(key K) (v V, ok bool) {
	i, ok := (*sync.Map)(m).LoadAndDelete(key)
	if !ok {
		return
	}
	return i.(V), true
}

func (m *SMap[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	i, loaded := (*sync.Map)(m).LoadOrStore(key, val)
	if !loaded {
		return
	}
	return i.(V), true
}

func (m *SMap[K, V]) Range(f func(K, V) bool) {
	(*sync.Map)(m).Range(func(key, val any) bool {
		return f(key.(K), val.(V))
	})
}

func (m *SMap[K, V]) Store(key K, val V) {
	(*sync.Map)(m).Store(key, val)
}
