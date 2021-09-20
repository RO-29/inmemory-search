package inmemorycache

import "testing"

func TestDI(t *testing.T) {
	dic := NewDIContainer()
	_ = dic.Cache()
}
