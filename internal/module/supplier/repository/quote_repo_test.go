package repository

import (
	"testing"

	"am-erp-go/internal/module/supplier/domain"
)

func TestQuoteRepositoryImplementsInterface(t *testing.T) {
	var _ domain.QuoteRepository = (*quoteRepository)(nil)
}
