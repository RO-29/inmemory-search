package inmemorycache

// DIContainer is the DI container for inmemory cache.
type DIContainer struct {
	Cache func() Cache
}

func newDIContainer() *DIContainer {
	dic := &DIContainer{}
	dic.Cache = NewDIProvider()
	return dic
}

// NewDIContainer returns a new DIContainer.
func NewDIContainer() *DIContainer {
	return newDIContainer()
}
