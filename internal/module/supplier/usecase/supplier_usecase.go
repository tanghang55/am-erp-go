package usecase

import "am-erp-go/internal/module/supplier/domain"

type SupplierUsecase struct {
	supplierRepo        domain.SupplierRepository
	supplierTypeRepo    domain.SupplierTypeRepository
	supplierContactRepo domain.SupplierContactRepository
	supplierAccountRepo domain.SupplierAccountRepository
	supplierTagRepo     domain.SupplierTagRepository
}

func NewSupplierUsecase(
	supplierRepo domain.SupplierRepository,
	supplierTypeRepo domain.SupplierTypeRepository,
	supplierContactRepo domain.SupplierContactRepository,
	supplierAccountRepo domain.SupplierAccountRepository,
	supplierTagRepo domain.SupplierTagRepository,
) *SupplierUsecase {
	return &SupplierUsecase{
		supplierRepo:        supplierRepo,
		supplierTypeRepo:    supplierTypeRepo,
		supplierContactRepo: supplierContactRepo,
		supplierAccountRepo: supplierAccountRepo,
		supplierTagRepo:     supplierTagRepo,
	}
}

func (uc *SupplierUsecase) ListSuppliers(params *domain.SupplierListParams) ([]domain.SupplierListItem, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	suppliers, total, err := uc.supplierRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	typeMap := map[uint64][]string{}
	if uc.supplierTypeRepo != nil {
		typeMap, err = uc.supplierTypeRepo.ListBySupplierIDs(extractSupplierIDs(suppliers))
		if err != nil {
			return nil, 0, err
		}
	}
	items := make([]domain.SupplierListItem, 0, len(suppliers))
	for _, supplier := range suppliers {
		items = append(items, domain.SupplierListItem{
			Supplier: supplier,
			Types:    typeMap[supplier.ID],
		})
	}
	return items, total, nil
}

func (uc *SupplierUsecase) GetSupplier(id uint64) (*domain.SupplierDetail, error) {
	supplier, err := uc.supplierRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	types, err := uc.getSupplierTypes(id)
	if err != nil {
		return nil, err
	}
	contacts, err := uc.getSupplierContacts(id)
	if err != nil {
		return nil, err
	}
	accounts, err := uc.getSupplierAccounts(id)
	if err != nil {
		return nil, err
	}
	tags, err := uc.getSupplierTags(id)
	if err != nil {
		return nil, err
	}
	return &domain.SupplierDetail{
		Supplier: *supplier,
		Types:    types,
		Contacts: contacts,
		Accounts: accounts,
		Tags:     tags,
	}, nil
}

func (uc *SupplierUsecase) CreateSupplier(supplier *domain.Supplier, types []string) (*domain.SupplierDetail, error) {
	if err := uc.supplierRepo.Create(supplier); err != nil {
		return nil, err
	}
	if uc.supplierTypeRepo != nil {
		if err := uc.supplierTypeRepo.ReplaceBySupplierID(supplier.ID, types); err != nil {
			return nil, err
		}
	}
	return &domain.SupplierDetail{
		Supplier: *supplier,
		Types:    types,
		Contacts: []domain.SupplierContact{},
		Accounts: []domain.SupplierAccount{},
		Tags:     []domain.SupplierTag{},
	}, nil
}

func (uc *SupplierUsecase) UpdateSupplier(supplier *domain.Supplier, types []string) (*domain.SupplierDetail, error) {
	if err := uc.supplierRepo.Update(supplier); err != nil {
		return nil, err
	}
	if uc.supplierTypeRepo != nil {
		if err := uc.supplierTypeRepo.ReplaceBySupplierID(supplier.ID, types); err != nil {
			return nil, err
		}
	}
	return uc.GetSupplier(supplier.ID)
}

func (uc *SupplierUsecase) DeleteSupplier(id uint64) error {
	return uc.supplierRepo.Delete(id)
}

func (uc *SupplierUsecase) CreateSupplierContact(supplierID uint64, contact *domain.SupplierContact) (*domain.SupplierContact, error) {
	contact.SupplierID = supplierID
	if err := uc.supplierContactRepo.Create(contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (uc *SupplierUsecase) UpdateSupplierContact(supplierID uint64, contact *domain.SupplierContact) (*domain.SupplierContact, error) {
	contact.SupplierID = supplierID
	if err := uc.supplierContactRepo.Update(contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (uc *SupplierUsecase) DeleteSupplierContact(supplierID uint64, contactID uint64) error {
	return uc.supplierContactRepo.Delete(contactID, supplierID)
}

func (uc *SupplierUsecase) CreateSupplierAccount(supplierID uint64, account *domain.SupplierAccount) (*domain.SupplierAccount, error) {
	account.SupplierID = supplierID
	if err := uc.supplierAccountRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (uc *SupplierUsecase) UpdateSupplierAccount(supplierID uint64, account *domain.SupplierAccount) (*domain.SupplierAccount, error) {
	account.SupplierID = supplierID
	if err := uc.supplierAccountRepo.Update(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (uc *SupplierUsecase) DeleteSupplierAccount(supplierID uint64, accountID uint64) error {
	return uc.supplierAccountRepo.Delete(accountID, supplierID)
}

func (uc *SupplierUsecase) CreateSupplierTag(supplierID uint64, tag *domain.SupplierTag) (*domain.SupplierTag, error) {
	tag.SupplierID = supplierID
	if err := uc.supplierTagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (uc *SupplierUsecase) UpdateSupplierTag(supplierID uint64, tag *domain.SupplierTag) (*domain.SupplierTag, error) {
	tag.SupplierID = supplierID
	if err := uc.supplierTagRepo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (uc *SupplierUsecase) DeleteSupplierTag(supplierID uint64, tagID uint64) error {
	return uc.supplierTagRepo.Delete(tagID, supplierID)
}

func extractSupplierIDs(items []domain.Supplier) []uint64 {
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func (uc *SupplierUsecase) getSupplierTypes(id uint64) ([]string, error) {
	if uc.supplierTypeRepo == nil {
		return []string{}, nil
	}
	return uc.supplierTypeRepo.ListBySupplierID(id)
}

func (uc *SupplierUsecase) getSupplierContacts(id uint64) ([]domain.SupplierContact, error) {
	if uc.supplierContactRepo == nil {
		return []domain.SupplierContact{}, nil
	}
	return uc.supplierContactRepo.ListBySupplierID(id)
}

func (uc *SupplierUsecase) getSupplierAccounts(id uint64) ([]domain.SupplierAccount, error) {
	if uc.supplierAccountRepo == nil {
		return []domain.SupplierAccount{}, nil
	}
	return uc.supplierAccountRepo.ListBySupplierID(id)
}

func (uc *SupplierUsecase) getSupplierTags(id uint64) ([]domain.SupplierTag, error) {
	if uc.supplierTagRepo == nil {
		return []domain.SupplierTag{}, nil
	}
	return uc.supplierTagRepo.ListBySupplierID(id)
}
