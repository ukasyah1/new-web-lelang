package faq

import "context"

type Repository interface {
	GetAll(ctx context.Context, lang string) ([]Category, error)
}
