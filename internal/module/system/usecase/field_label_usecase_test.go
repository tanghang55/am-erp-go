package usecase_test

import (
	"testing"

	"am-erp-go/internal/module/system/domain"
	"am-erp-go/internal/module/system/usecase"
)

type fakeRepo struct {
	labels []*domain.FieldLabel
}

func (f *fakeRepo) GetAll() ([]*domain.FieldLabel, error) { return f.labels, nil }
func (f *fakeRepo) List(page, pageSize int, keyword string) ([]*domain.FieldLabel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRepo) GetByID(id uint64) (*domain.FieldLabel, error) {
	return nil, domain.ErrFieldLabelNotFound
}
func (f *fakeRepo) GetByKey(key string) (*domain.FieldLabel, error) {
	return nil, domain.ErrFieldLabelNotFound
}
func (f *fakeRepo) Create(label *domain.FieldLabel) error {
	f.labels = append(f.labels, label)
	return nil
}
func (f *fakeRepo) Update(label *domain.FieldLabel) error { return nil }
func (f *fakeRepo) Delete(id uint64) error                { return nil }

func TestGetLabels_NormalizesKey(t *testing.T) {
	repo := &fakeRepo{labels: []*domain.FieldLabel{{
		LabelKey: "Product.List.Title",
		Labels:   domain.LabelMap{"zh-CN": "Product List"},
	}}}
	uc := usecase.NewFieldLabelUseCase(repo)

	labels, err := uc.GetLabels("zh-CN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if labels["product.list.title"] != "Product List" {
		t.Fatalf("expected normalized key")
	}
}
