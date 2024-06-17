package utils

import (
	"sync"
	"time"
)

type TokenBlackList struct {
	tokens map[string]time.Time
	mu     sync.Mutex
}

var (
	once sync.Once
	List *TokenBlackList
)

func InitTokenBlackList() *TokenBlackList {
	once.Do(func() {
		List = &TokenBlackList{
			tokens: make(map[string]time.Time),
		}
	})

	return List
}

func (tl *TokenBlackList) Block(token string) {
	tl.mu.Lock()
	tl.tokens[token] = time.Now()
	tl.mu.Unlock()
}

func (tl *TokenBlackList) Blocked(token string) bool {
	tl.mu.Lock()
	_, exists := tl.tokens[token]
	tl.mu.Unlock()

	return exists
}
