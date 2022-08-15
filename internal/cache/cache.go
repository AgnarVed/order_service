package cache

import (
	"encoding/json"
	"errors"
	lru "github.com/hashicorp/golang-lru"
	"github.com/mitchellh/mapstructure"
	"log"
	"tests2/internal/models"
	"tests2/internal/service"
)

type Cache struct {
	cache   *lru.Cache
	size    int
	logger  *log.Logger
	service *service.Service
}

func NewCache(size int) (*Cache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	cache := &Cache{
		cache: c,
		size:  size,
	}
	return cache, nil
}

func (c *Cache) Add(key interface{}, insert interface{}) bool {
	input, err := json.Marshal(insert)
	if err != nil {
		return false
	}
	ok := c.cache.Add(key, input)
	return ok
}

func (c *Cache) Get(key interface{}) ([]byte, bool, error) {
	val, ok := c.cache.Get(key)
	order, err := Convert(val)
	if err != nil {
		return nil, ok, err
	}
	if order == nil {
		return nil, ok, errors.New("order is nil")
	}
	return order, ok, nil
}

func Convert(input interface{}) ([]byte, error) {
	var order []byte
	err := mapstructure.Decode(input, &order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (c *Cache) UploadCache(orders []*models.OrderDB) error {
	for _, order := range orders {
		c.cache.Add(order.OrderUID, order.Data)
	}
	return nil
}
