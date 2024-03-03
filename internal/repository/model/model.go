package model

import "time"

type RedisRequestModel struct {
	Key    string
	Value  interface{}
	Expire time.Duration
}
