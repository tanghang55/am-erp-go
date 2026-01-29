package domain

type LogisticsProviderRepository interface {
	Create(provider *LogisticsProvider) error
	Update(provider *LogisticsProvider) error
	Delete(id uint64) error
	GetByID(id uint64) (*LogisticsProvider, error)
	GetByCode(code string) (*LogisticsProvider, error)
	List(params *LogisticsProviderListParams) ([]*LogisticsProvider, int64, error)
}

type ShippingRateRepository interface {
	Create(rate *ShippingRate) error
	Update(rate *ShippingRate) error
	Delete(id uint64) error
	GetByID(id uint64) (*ShippingRate, error)
	List(params *ShippingRateListParams) ([]*ShippingRate, int64, error)
	QueryLatestRate(params *QueryLatestRateParams) (*ShippingRate, error)
}
