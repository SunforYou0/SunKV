package dict

type Consumer func(key string, val interface{}) bool

type Dict interface {
	Get(key string) (val interface{}, exsist bool)
	Len() (length int)
	Put(key string, val interface{}) (result int) //res:存进了几个
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExist(key string, val interface{}) (result int)
	Remove(key string) (result int)
	Removes(keys ...string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
