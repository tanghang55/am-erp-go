package usecase

import (
	"fmt"
	"testing"

	"am-erp-go/internal/module/supplier/domain"
)

type stubQuoteRepo struct {
	quotes   map[string]domain.ProductSupplierQuote
	created  *domain.ProductSupplierQuote
	updated  *domain.ProductSupplierQuote
	deleted  [2]uint64
}

func (s *stubQuoteRepo) key(productID, supplierID uint64) string {
	return fmt.Sprintf("%d:%d", productID, supplierID)
}

func (s *stubQuoteRepo) ListByProductIDs(_ []uint64) (map[uint64][]domain.ProductSupplierQuote, error) {
	return map[uint64][]domain.ProductSupplierQuote{}, nil
}

func (s *stubQuoteRepo) ListProductsWithQuotes(_ *domain.QuoteListParams) ([]domain.ProductQuoteRow, int64, error) {
	return nil, 0, nil
}

func (s *stubQuoteRepo) GetByProductSupplier(productID, supplierID uint64) (*domain.ProductSupplierQuote, error) {
	if s.quotes == nil {
		return nil, domain.ErrQuoteNotFound
	}
	quote, ok := s.quotes[s.key(productID, supplierID)]
	if !ok {
		return nil, domain.ErrQuoteNotFound
	}
	copy := quote
	return &copy, nil
}

func (s *stubQuoteRepo) Create(quote *domain.ProductSupplierQuote) error {
	if s.quotes == nil {
		s.quotes = map[string]domain.ProductSupplierQuote{}
	}
	copy := *quote
	s.quotes[s.key(copy.ProductID, copy.SupplierID)] = copy
	s.created = &copy
	return nil
}

func (s *stubQuoteRepo) Update(quote *domain.ProductSupplierQuote) error {
	if s.quotes == nil {
		s.quotes = map[string]domain.ProductSupplierQuote{}
	}
	copy := *quote
	s.quotes[s.key(copy.ProductID, copy.SupplierID)] = copy
	s.updated = &copy
	return nil
}

func (s *stubQuoteRepo) Delete(productID, supplierID uint64) error {
	s.deleted = [2]uint64{productID, supplierID}
	return nil
}

type stubQuoteProductRepo struct {
	defaultSupplierIDs map[uint64]uint64
	updatedSupplierID  struct {
		productID  uint64
		supplierID uint64
	}
	updatedCosts map[uint64]float64
}

func (s *stubQuoteProductRepo) GetDefaultSupplierID(productID uint64) (uint64, error) {
	return s.defaultSupplierIDs[productID], nil
}

func (s *stubQuoteProductRepo) UpdateDefaultSupplierID(productID, supplierID uint64) error {
	s.updatedSupplierID.productID = productID
	s.updatedSupplierID.supplierID = supplierID
	if s.defaultSupplierIDs == nil {
		s.defaultSupplierIDs = map[uint64]uint64{}
	}
	s.defaultSupplierIDs[productID] = supplierID
	return nil
}

func (s *stubQuoteProductRepo) UpdateUnitCost(productID uint64, unitCost float64) error {
	if s.updatedCosts == nil {
		s.updatedCosts = map[uint64]float64{}
	}
	s.updatedCosts[productID] = unitCost
	return nil
}

func TestUpdateQuoteSyncsDefaultSupplierUnitCost(t *testing.T) {
	quoteRepo := &stubQuoteRepo{
		quotes: map[string]domain.ProductSupplierQuote{
			"10:5": {
				ID:         1,
				ProductID:  10,
				SupplierID: 5,
				Price:      12.5,
				Currency:   "USD",
				Status:     domain.ProductSupplierQuoteStatusActive,
				QtyMOQ:     1,
			},
		},
	}
	productRepo := &stubQuoteProductRepo{
		defaultSupplierIDs: map[uint64]uint64{10: 5},
	}
	uc := NewQuoteUsecase(quoteRepo, productRepo, nil)

	updated, err := uc.UpdateQuote(nil, &domain.ProductSupplierQuote{
		ProductID:    10,
		SupplierID:   5,
		Price:        18.2,
		Currency:     "USD",
		QtyMOQ:       1,
		LeadTimeDays: 7,
		Status:       domain.ProductSupplierQuoteStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated == nil || updated.Price != 18.2 {
		t.Fatalf("expected updated quote price 18.2, got %+v", updated)
	}
	if productRepo.updatedCosts[10] != 18.2 {
		t.Fatalf("expected product cost synced to 18.2, got %+v", productRepo.updatedCosts)
	}
}

func TestSetDefaultSupplierSyncsProductUnitCostFromQuote(t *testing.T) {
	quoteRepo := &stubQuoteRepo{
		quotes: map[string]domain.ProductSupplierQuote{
			"10:6": {
				ID:         2,
				ProductID:  10,
				SupplierID: 6,
				Price:      25.4,
				Currency:   "USD",
				Status:     domain.ProductSupplierQuoteStatusActive,
				QtyMOQ:     1,
			},
		},
	}
	productRepo := &stubQuoteProductRepo{
		defaultSupplierIDs: map[uint64]uint64{10: 5},
	}
	uc := NewQuoteUsecase(quoteRepo, productRepo, nil)

	if err := uc.SetDefaultSupplier(nil, 10, 6); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if productRepo.updatedSupplierID.productID != 10 || productRepo.updatedSupplierID.supplierID != 6 {
		t.Fatalf("expected default supplier updated to 6, got %+v", productRepo.updatedSupplierID)
	}
	if productRepo.updatedCosts[10] != 25.4 {
		t.Fatalf("expected product cost synced to 25.4, got %+v", productRepo.updatedCosts)
	}
}

func TestCreateQuoteSyncsDefaultSupplierUnitCost(t *testing.T) {
	quoteRepo := &stubQuoteRepo{}
	productRepo := &stubQuoteProductRepo{
		defaultSupplierIDs: map[uint64]uint64{10: 5},
	}
	uc := NewQuoteUsecase(quoteRepo, productRepo, nil)

	created, err := uc.CreateQuote(nil, &domain.ProductSupplierQuote{
		ProductID:    10,
		SupplierID:   5,
		Price:        13.8,
		Currency:     "USD",
		QtyMOQ:       1,
		LeadTimeDays: 3,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created == nil || created.Price != 13.8 {
		t.Fatalf("expected created quote price 13.8, got %+v", created)
	}
	if productRepo.updatedCosts[10] != 13.8 {
		t.Fatalf("expected product cost synced to 13.8, got %+v", productRepo.updatedCosts)
	}
}
