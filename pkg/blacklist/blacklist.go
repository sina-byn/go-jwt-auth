package blacklist

import (
	"sync"
	"time"
)

type TokenBlackList struct {
	tokens map[string]time.Time
	mu     sync.Mutex
}

var (
	BlockedTokens *TokenBlackList
	once          sync.Once
)

func InitTokenBlackList() *TokenBlackList {
	once.Do(func() {
		BlockedTokens = &TokenBlackList{
			tokens: make(map[string]time.Time),
		}
	})

	return BlockedTokens
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
