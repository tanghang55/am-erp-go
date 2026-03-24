package domain

import "time"

const (
	ConfigScopeGlobal = "GLOBAL"

	ConfigModuleFinance     = "finance"
	ConfigModuleSalesImport = "sales_import"
	ConfigModuleProcurement = "procurement"

	ConfigValueTypeString = "STRING"
	ConfigValueTypeInt    = "INT"
	ConfigValueTypeBool   = "BOOL"
	ConfigValueTypeEnum   = "ENUM"

	ConfigKeyFinanceDefaultCurrency          = "finance.default_currency"
	ConfigKeyFinanceExchangeRateScale        = "finance.exchange_rate_scale"
	ConfigKeySalesImportDefaultChannel       = "sales_import.default_channel"
	ConfigKeySalesImportDefaultMarketplace   = "sales_import.default_marketplace"
	ConfigKeyProcurementDemandWindowDays     = "procurement.demand_window_days"
	ConfigKeyProcurementDefaultLeadTimeDays  = "procurement.default_lead_time_days"
	ConfigKeyProcurementDefaultSafetyDays    = "procurement.default_safety_days"
	ConfigKeyProcurementDefaultMOQ           = "procurement.default_moq"
	ConfigKeyProcurementDefaultOrderMultiple = "procurement.default_order_multiple"
)

type ConfigDefinition struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	ConfigKey    string    `json:"config_key" gorm:"column:config_key;size:120;not null;uniqueIndex"`
	ModuleCode   string    `json:"module_code" gorm:"column:module_code;size:50;not null;index"`
	ModuleName   string    `json:"module_name" gorm:"column:module_name;size:100;not null"`
	GroupCode    string    `json:"group_code" gorm:"column:group_code;size:50;not null"`
	GroupName    string    `json:"group_name" gorm:"column:group_name;size:100;not null"`
	Label        string    `json:"label" gorm:"column:label;size:120;not null"`
	Description  *string   `json:"description" gorm:"column:description;size:255"`
	ValueType    string    `json:"value_type" gorm:"column:value_type;size:20;not null"`
	ScopeType    string    `json:"scope_type" gorm:"column:scope_type;size:20;not null;default:'GLOBAL'"`
	DefaultValue string    `json:"default_value" gorm:"column:default_value;size:255;not null"`
	OptionsJSON  *string   `json:"options_json" gorm:"column:options_json;type:json"`
	Sort         uint32    `json:"sort" gorm:"column:sort;not null;default:100"`
	IsActive     uint8     `json:"is_active" gorm:"column:is_active;not null;default:1"`
	GmtCreate    time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ConfigDefinition) TableName() string {
	return "config_definition"
}

type ConfigValue struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	ConfigKey   string    `json:"config_key" gorm:"column:config_key;size:120;not null;uniqueIndex:uk_config_scope"`
	ScopeType   string    `json:"scope_type" gorm:"column:scope_type;size:20;not null;default:'GLOBAL';uniqueIndex:uk_config_scope"`
	ScopeRefID  uint64    `json:"scope_ref_id" gorm:"column:scope_ref_id;not null;default:0;uniqueIndex:uk_config_scope"`
	ConfigValue string    `json:"config_value" gorm:"column:config_value;size:255;not null"`
	UpdatedBy   *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate   time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ConfigValue) TableName() string {
	return "config_value"
}

type ConfigCenterModuleSummary struct {
	ModuleCode string `json:"module_code"`
	ModuleName string `json:"module_name"`
	Sort       uint32 `json:"sort"`
}

type ConfigCenterModule struct {
	ModuleCode string              `json:"module_code"`
	ModuleName string              `json:"module_name"`
	Groups     []ConfigCenterGroup `json:"groups"`
}

type ConfigCenterGroup struct {
	GroupCode string             `json:"group_code"`
	GroupName string             `json:"group_name"`
	Items     []ConfigCenterItem `json:"items"`
}

type ConfigCenterItem struct {
	ConfigKey    string   `json:"config_key"`
	Label        string   `json:"label"`
	Description  string   `json:"description"`
	ValueType    string   `json:"value_type"`
	ScopeType    string   `json:"scope_type"`
	DefaultValue string   `json:"default_value"`
	Value        string   `json:"value"`
	Options      []string `json:"options"`
	Sort         uint32   `json:"sort"`
}

type ConfigCenterUpdateInput struct {
	Values map[string]string `json:"values"`
}

type ConfigCenterFinance struct {
	DefaultCurrency   string `json:"default_currency"`
	ExchangeRateScale uint32 `json:"exchange_rate_scale"`
}

type ConfigCenterSalesImport struct {
	DefaultChannel     string `json:"default_channel"`
	DefaultMarketplace string `json:"default_marketplace"`
}

type ConfigCenterProcurement struct {
	DemandWindowDays     uint32 `json:"demand_window_days"`
	DefaultLeadTimeDays  uint32 `json:"default_lead_time_days"`
	DefaultSafetyDays    uint32 `json:"default_safety_days"`
	DefaultMOQ           uint32 `json:"default_moq"`
	DefaultOrderMultiple uint32 `json:"default_order_multiple"`
}
