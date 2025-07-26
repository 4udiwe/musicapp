package postgres

import "time"

type Option func(*Postgres)

func ConnAttempts(a int) Option {
	return func(p *Postgres) {
		p.connAttempts = a
	}
}

func TimeOut(t time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = t
	}
}
