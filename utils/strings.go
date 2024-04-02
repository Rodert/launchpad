package utils

import (
	"math/rand"
	"sync"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var SBuilder StringBuilder

type (
	StringBuilder struct {
		mutex sync.Mutex
	}
)

func (c *StringBuilder) GetString(length int) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b := make([]byte, length)
	// 生成指定长度的随机字节数组
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	// flush seed
	return string(b)
}

