package config

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
)

var lock = &sync.Mutex{}

type RedisContext struct {
	Ctx context.Context
}

var singleInstance *RedisContext

func RedisContextGetInstance() *RedisContext {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Info("creating redis context singleton instance")
			singleInstance = &RedisContext{Ctx: context.Background()}
		}
	}

	return singleInstance
}
