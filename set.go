package gocom

type Empty struct{}

var empty Empty

type Set struct {
	M map[interface{}]Empty
}

func NewSet() *Set{
	return &Set{
		M: map[interface{}]Empty{},
	}
}

//添加元素
func (s *Set) Add(key interface{}) {
	s.M[key] = empty
}

//删除元素
func (s *Set) Remove(key interface{}) {
	delete(s.M, key)
}

//检查是否存在
func (s *Set) Exist(key interface{}) bool {
	_, ok := s.M[key]
	return ok
}

//获取长度
func (s *Set) Len() int {
	return len(s.M)
}

//清空set
func (s *Set) Clear() {
	s.M = make(map[interface{}]Empty)
}
