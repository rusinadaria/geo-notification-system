package repository

import "context"

type DBPinger interface {
	PingContext(ctx context.Context) error
}

type RedisPinger interface {
	Ping(ctx context.Context) error
}
