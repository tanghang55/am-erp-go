package usecase

import (
	"strings"

	"am-erp-go/internal/module/system/domain"
)

type FieldLabelUseCase struct {
	repo domain.FieldLabelRepository
}

func NewFieldLabelUseCase(repo domain.FieldLabelRepository) *FieldLabelUseCase {
	return &FieldLabelUseCase{repo: repo}
}

func (uc *FieldLabelUseCase) GetLabels(locale string) (map[string]string, error) {
	labels, err := uc.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, label := range labels {
		if label == nil {
			continue
		}
		key := strings.ToLower(label.LabelKey)
		if key == "" {
			continue
		}
		if value, ok := label.Labels[locale]; ok {
			result[key] = value
		}
	}

	return result, nil
}

func (uc *FieldLabelUseCase) List(page, pageSize int, keyword string) ([]*domain.FieldLabel, int64, error) {
	return uc.repo.List(page, pageSize, keyword)
}

func (uc *FieldLabelUseCase) Create(label *domain.FieldLabel) error {
	if label == nil {
		return nil
	}
	label.NormalizeKey()
	return uc.repo.Create(label)
}

func (uc *FieldLabelUseCase) Update(label *domain.FieldLabel) error {
	if label == nil {
		return nil
	}
	label.NormalizeKey()
	return uc.repo.Update(label)
}

func (uc *FieldLabelUseCase) Delete(id uint64) error {
	return uc.repo.Delete(id)
}

func (uc *FieldLabelUseCase) GetByID(id uint64) (*domain.FieldLabel, error) {
	return uc.repo.GetByID(id)
}
