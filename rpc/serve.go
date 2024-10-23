package rpc

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/server"
)

type Args struct {
	A int
	B int
}

type Cache struct {
	cache map[string]string
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]string),
	}
}

func (c *Cache) Get(ctx context.Context, k Key, v *Value) error {
	v.Value = c.cache[k.Key]
	return nil
}

func (c *Cache) Post(ctx context.Context, p Pair, f *Flag) error {
	c.cache[p.Key] = p.Value
	return nil
}

func (c *Cache) Delete(ctx context.Context, k Key, f *Flag) error {
	if _, ok := c.cache[k.Key]; !ok {
		f.Flag = false
		return nil
	}

	delete(c.cache, k.Key)
	f.Flag = true
	return nil
}

func (c *Cache) Query(ctx context.Context, key Key, f *Flag) error {
	if _, ok := c.cache[key.Key]; !ok {
		f.Flag = false
	} else {
		f.Flag = true
	}
	return nil
}

func StartXServer(port int) {
	addr := fmt.Sprintf("localhost:%d", port)
	s := server.NewServer()
	s.Register(NewCache(), "")
	go func() {
		if err := s.Serve("tcp", addr); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	fmt.Printf("XServer is running on http://localhost:%d\n", port)
}
