package faq

import (
	"context"
	"errors"
	"strings"
)

var ErrUnsupportedLanguage = errors.New("lang harus id atau en")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context, lang string) ([]Category, error) {
	normalizedLang, err := normalizeLang(lang)
	if err != nil {
		return nil, err
	}
	return s.repository.GetAll(ctx, normalizedLang)
}

func normalizeLang(lang string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(lang))
	if value == "" {
		return "id", nil
	}
	if value != "id" && value != "en" {
		return "", ErrUnsupportedLanguage
	}
	return value, nil
}
