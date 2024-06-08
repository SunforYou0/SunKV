package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exsist bool) {
	val, ok := dict.m.Load(key)
	return val, ok
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0
	}
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, exist := dict.m.Load(key)
	if exist {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExist(key string, val interface{}) (result int) {
	_, exist := dict.m.Load(key)
	if exist {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exist := dict.m.Load(key)

	if exist {
		dict.m.Delete(key)
		return 1
	}
	return 0
}
func (dict *SyncDict) Removes(keys ...string) (result int) {

	for i := 0; i < len(keys); i++ {
		result += dict.Remove(keys[i])
	}
	return
}
func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	//TODO 这个方法可能有很大问题
	keys := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			keys[i] = key.(string)
			return false
		})
	}
	return keys
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	keys := make([]string, limit)
	i := 0
	dict.m.Range(func(key, value any) bool {
		keys[i] = key.(string)
		i++
		return i < limit
	})
	return keys
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
	// dict=MakeSyncDict() go函数参数是值传递，这样只能改变副本dict的值
}
