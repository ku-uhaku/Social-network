package middleware

import "kuu/internal/service"

type Middleware struct {
	Service *service.Service
}

func New(svc *service.Service) *Middleware {
	return &Middleware{
		Service: svc,
	}
}
