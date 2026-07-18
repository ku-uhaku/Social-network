package middleware

import "kuu/internal/repository"

type Middleware struct {
	Repo *repository.Repository
}

func New(repo *repository.Repository) *Middleware {
	return &Middleware{
		Repo: repo,
	}
}
