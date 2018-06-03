package store

type Event struct {
	Key string
}

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Watch(key string) (chan Event, chan error)
	Close() error
}

type common struct {
	root string
}

func (c *common) SetRoot(root string) {
	c.root = root
}

func (c *common) realKey(key string) string {
	return c.root + key
}

func (c *common) orignalKey(key string) string {
	return key[len(c.root):]
}

func (c *common) Get(key string) (string, error) {
	panic("abstract method")
	return "", nil
}

func (c *common) Set(key, value string) error {
	panic("abstract method")
	return nil
}

func (c *common) Watch(key string) (chan Event, chan error) {
	panic("abstract method")
	return nil, nil
}

func (c *common) Close() error {
	panic("abstract method")
	return nil
}
