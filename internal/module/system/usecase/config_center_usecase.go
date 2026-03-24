package usecase

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type ConfigCenterAuditLogger interface {
	RecordFromContext(c *gin.Context, payload AuditLogPayload) error
}

type ConfigCenterUsecase struct {
	repo        systemdomain.ConfigCenterRepository
	auditLogger ConfigCenterAuditLogger
}

var supportedCurrencies = []string{"USD", "CNY", "JPY", "AUD", "EUR", "GBP", "CAD"}

func NewConfigCenterUsecase(repo systemdomain.ConfigCenterRepository, auditLogger ConfigCenterAuditLogger) *ConfigCenterUsecase {
	return &ConfigCenterUsecase{repo: repo, auditLogger: auditLogger}
}

func (uc *ConfigCenterUsecase) ListModules() ([]systemdomain.ConfigCenterModuleSummary, error) {
	definitions, err := uc.repo.ListDefinitions("")
	if err != nil {
		return nil, err
	}
	moduleMap := map[string]systemdomain.ConfigCenterModuleSummary{}
	for _, item := range definitions {
		if item == nil || item.IsActive == 0 {
			continue
		}
		current, ok := moduleMap[item.ModuleCode]
		if !ok || item.Sort < current.Sort {
			moduleMap[item.ModuleCode] = systemdomain.ConfigCenterModuleSummary{
				ModuleCode: item.ModuleCode,
				ModuleName: item.ModuleName,
				Sort:       item.Sort,
			}
		}
	}
	modules := make([]systemdomain.ConfigCenterModuleSummary, 0, len(moduleMap))
	for _, item := range moduleMap {
		modules = append(modules, item)
	}
	sort.Slice(modules, func(i, j int) bool {
		if modules[i].Sort == modules[j].Sort {
			return modules[i].ModuleCode < modules[j].ModuleCode
		}
		return modules[i].Sort < modules[j].Sort
	})
	return modules, nil
}

func (uc *ConfigCenterUsecase) GetModule(moduleCode string) (*systemdomain.ConfigCenterModule, error) {
	moduleCode = strings.TrimSpace(moduleCode)
	if moduleCode == "" {
		return nil, systemdomain.ErrConfigCenterInvalid
	}
	definitions, err := uc.repo.ListDefinitions(moduleCode)
	if err != nil {
		return nil, err
	}
	if len(definitions) == 0 {
		return nil, fmt.Errorf("%w: module not found", systemdomain.ErrConfigCenterInvalid)
	}
	keys := make([]string, 0, len(definitions))
	for _, item := range definitions {
		if item != nil {
			keys = append(keys, item.ConfigKey)
		}
	}
	values, err := uc.repo.ListValues(systemdomain.ConfigScopeGlobal, 0, keys)
	if err != nil {
		return nil, err
	}
	valueMap := map[string]string{}
	for _, item := range values {
		if item != nil {
			valueMap[item.ConfigKey] = strings.TrimSpace(item.ConfigValue)
		}
	}
	return uc.buildModule(definitions, valueMap)
}

func (uc *ConfigCenterUsecase) UpdateModule(c *gin.Context, moduleCode string, values map[string]string, operatorID *uint64) (*systemdomain.ConfigCenterModule, error) {
	moduleCode = strings.TrimSpace(moduleCode)
	if moduleCode == "" {
		return nil, systemdomain.ErrConfigCenterInvalid
	}
	definitions, err := uc.repo.ListDefinitions(moduleCode)
	if err != nil {
		return nil, err
	}
	if len(definitions) == 0 {
		return nil, fmt.Errorf("%w: module not found", systemdomain.ErrConfigCenterInvalid)
	}
	before, err := uc.GetModule(moduleCode)
	if err != nil {
		return nil, err
	}

	normalized, err := uc.normalizeModuleValues(moduleCode, definitions, values)
	if err != nil {
		return nil, err
	}
	if moduleCode == systemdomain.ConfigModuleFinance {
		currentBaseCurrency := strings.ToUpper(strings.TrimSpace(findModuleValue(before, systemdomain.ConfigKeyFinanceDefaultCurrency)))
		nextBaseCurrency := strings.ToUpper(strings.TrimSpace(normalized[systemdomain.ConfigKeyFinanceDefaultCurrency]))
		if currentBaseCurrency != "" && nextBaseCurrency != "" && currentBaseCurrency != nextBaseCurrency {
			return nil, fmt.Errorf("%w: finance.default_currency is immutable after initialization", systemdomain.ErrConfigCenterInvalid)
		}
	}
	items := make([]*systemdomain.ConfigValue, 0, len(normalized))
	for key, value := range normalized {
		items = append(items, &systemdomain.ConfigValue{
			ConfigKey:   key,
			ScopeType:   systemdomain.ConfigScopeGlobal,
			ScopeRefID:  0,
			ConfigValue: value,
			UpdatedBy:   operatorID,
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ConfigKey < items[j].ConfigKey })
	if err := uc.repo.UpsertValues(items); err != nil {
		return nil, err
	}
	after, err := uc.GetModule(moduleCode)
	if err != nil {
		return nil, err
	}
	if uc.auditLogger != nil && c != nil {
		_ = uc.auditLogger.RecordFromContext(c, AuditLogPayload{
			Module:     "System",
			Action:     "CONFIG_CENTER_UPDATE",
			EntityType: "ConfigCenterModule",
			EntityID:   strings.ToUpper(moduleCode),
			Before:     before,
			After:      after,
		})
	}
	return after, nil
}

func (uc *ConfigCenterUsecase) GetDefaultBaseCurrency() string {
	return uc.GetFinanceConfig().DefaultCurrency
}

func (uc *ConfigCenterUsecase) GetExchangeRateScale() uint32 {
	return uc.GetFinanceConfig().ExchangeRateScale
}

func (uc *ConfigCenterUsecase) GetFinanceConfig() systemdomain.ConfigCenterFinance {
	module, err := uc.GetModule(systemdomain.ConfigModuleFinance)
	if err != nil || module == nil {
		return systemdomain.ConfigCenterFinance{
			DefaultCurrency:   "USD",
			ExchangeRateScale: 4,
		}
	}
	config := systemdomain.ConfigCenterFinance{
		DefaultCurrency:   "USD",
		ExchangeRateScale: 4,
	}
	if value := findModuleValue(module, systemdomain.ConfigKeyFinanceDefaultCurrency); value != "" {
		config.DefaultCurrency = value
	}
	config.ExchangeRateScale = clampU32(
		parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyFinanceExchangeRateScale), 4),
		0,
		8,
		4,
	)
	return config
}

func (uc *ConfigCenterUsecase) GetSalesImportDefaults() systemdomain.ConfigCenterSalesImport {
	module, err := uc.GetModule(systemdomain.ConfigModuleSalesImport)
	if err != nil || module == nil {
		return systemdomain.ConfigCenterSalesImport{
			DefaultChannel:     "MANUAL",
			DefaultMarketplace: "US",
		}
	}
	return systemdomain.ConfigCenterSalesImport{
		DefaultChannel:     valueOr(findModuleValue(module, systemdomain.ConfigKeySalesImportDefaultChannel), "MANUAL"),
		DefaultMarketplace: valueOr(findModuleValue(module, systemdomain.ConfigKeySalesImportDefaultMarketplace), "US"),
	}
}

func (uc *ConfigCenterUsecase) GetProcurementDefaults() systemdomain.ConfigCenterProcurement {
	module, err := uc.GetModule(systemdomain.ConfigModuleProcurement)
	if err != nil || module == nil {
		return systemdomain.ConfigCenterProcurement{
			DemandWindowDays:     30,
			DefaultLeadTimeDays:  15,
			DefaultSafetyDays:    7,
			DefaultMOQ:           1,
			DefaultOrderMultiple: 1,
		}
	}
	return systemdomain.ConfigCenterProcurement{
		DemandWindowDays:     parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyProcurementDemandWindowDays), 30),
		DefaultLeadTimeDays:  parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyProcurementDefaultLeadTimeDays), 15),
		DefaultSafetyDays:    parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyProcurementDefaultSafetyDays), 7),
		DefaultMOQ:           parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyProcurementDefaultMOQ), 1),
		DefaultOrderMultiple: parseUint32OrDefault(findModuleValue(module, systemdomain.ConfigKeyProcurementDefaultOrderMultiple), 1),
	}
}

func (uc *ConfigCenterUsecase) UpdateProcurementDefaults(c *gin.Context, defaults systemdomain.ConfigCenterProcurement, operatorID *uint64) error {
	_, err := uc.UpdateModule(c, systemdomain.ConfigModuleProcurement, map[string]string{
		systemdomain.ConfigKeyProcurementDemandWindowDays:     strconv.FormatUint(uint64(defaults.DemandWindowDays), 10),
		systemdomain.ConfigKeyProcurementDefaultLeadTimeDays:  strconv.FormatUint(uint64(defaults.DefaultLeadTimeDays), 10),
		systemdomain.ConfigKeyProcurementDefaultSafetyDays:    strconv.FormatUint(uint64(defaults.DefaultSafetyDays), 10),
		systemdomain.ConfigKeyProcurementDefaultMOQ:           strconv.FormatUint(uint64(defaults.DefaultMOQ), 10),
		systemdomain.ConfigKeyProcurementDefaultOrderMultiple: strconv.FormatUint(uint64(defaults.DefaultOrderMultiple), 10),
	}, operatorID)
	return err
}

func (uc *ConfigCenterUsecase) buildModule(definitions []*systemdomain.ConfigDefinition, currentValues map[string]string) (*systemdomain.ConfigCenterModule, error) {
	module := &systemdomain.ConfigCenterModule{
		ModuleCode: definitions[0].ModuleCode,
		ModuleName: definitions[0].ModuleName,
	}
	groupOrder := make([]string, 0)
	groupMap := map[string]*systemdomain.ConfigCenterGroup{}
	for _, item := range definitions {
		if item == nil || item.IsActive == 0 {
			continue
		}
		group, ok := groupMap[item.GroupCode]
		if !ok {
			group = &systemdomain.ConfigCenterGroup{
				GroupCode: item.GroupCode,
				GroupName: item.GroupName,
				Items:     make([]systemdomain.ConfigCenterItem, 0),
			}
			groupMap[item.GroupCode] = group
			groupOrder = append(groupOrder, item.GroupCode)
		}
		currentValue := strings.TrimSpace(currentValues[item.ConfigKey])
		if currentValue == "" {
			currentValue = item.DefaultValue
		}
		group.Items = append(group.Items, systemdomain.ConfigCenterItem{
			ConfigKey:    item.ConfigKey,
			Label:        item.Label,
			Description:  safeString(item.Description),
			ValueType:    item.ValueType,
			ScopeType:    item.ScopeType,
			DefaultValue: item.DefaultValue,
			Value:        currentValue,
			Options:      parseStringOptions(item.OptionsJSON),
			Sort:         item.Sort,
		})
	}
	for _, code := range groupOrder {
		group := groupMap[code]
		sort.Slice(group.Items, func(i, j int) bool {
			if group.Items[i].Sort == group.Items[j].Sort {
				return group.Items[i].ConfigKey < group.Items[j].ConfigKey
			}
			return group.Items[i].Sort < group.Items[j].Sort
		})
		module.Groups = append(module.Groups, *group)
	}
	return module, nil
}

func (uc *ConfigCenterUsecase) normalizeModuleValues(moduleCode string, definitions []*systemdomain.ConfigDefinition, rawValues map[string]string) (map[string]string, error) {
	if len(definitions) == 0 {
		return nil, systemdomain.ErrConfigCenterInvalid
	}
	result := map[string]string{}
	definitionMap := map[string]*systemdomain.ConfigDefinition{}
	for _, item := range definitions {
		if item != nil {
			definitionMap[item.ConfigKey] = item
		}
	}
	for key, definition := range definitionMap {
		value := strings.TrimSpace(rawValues[key])
		if value == "" {
			value = definition.DefaultValue
		}
		normalized, err := normalizeConfigValue(definition, value)
		if err != nil {
			return nil, err
		}
		result[key] = normalized
	}
	switch moduleCode {
	case systemdomain.ConfigModuleFinance:
		if !containsString(supportedCurrencies, result[systemdomain.ConfigKeyFinanceDefaultCurrency]) {
			return nil, fmt.Errorf("%w: unsupported finance.default_currency", systemdomain.ErrConfigCenterInvalid)
		}
		result[systemdomain.ConfigKeyFinanceExchangeRateScale] = strconv.FormatUint(
			uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyFinanceExchangeRateScale], 4), 0, 8, 4)),
			10,
		)
	case systemdomain.ConfigModuleSalesImport:
		if result[systemdomain.ConfigKeySalesImportDefaultChannel] == "" {
			return nil, fmt.Errorf("%w: sales_import.default_channel is required", systemdomain.ErrConfigCenterInvalid)
		}
		if result[systemdomain.ConfigKeySalesImportDefaultMarketplace] == "" {
			return nil, fmt.Errorf("%w: sales_import.default_marketplace is required", systemdomain.ErrConfigCenterInvalid)
		}
	case systemdomain.ConfigModuleProcurement:
		result[systemdomain.ConfigKeyProcurementDemandWindowDays] = strconv.FormatUint(uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyProcurementDemandWindowDays], 30), 1, 365, 30)), 10)
		result[systemdomain.ConfigKeyProcurementDefaultLeadTimeDays] = strconv.FormatUint(uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyProcurementDefaultLeadTimeDays], 15), 0, 365, 15)), 10)
		result[systemdomain.ConfigKeyProcurementDefaultSafetyDays] = strconv.FormatUint(uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyProcurementDefaultSafetyDays], 7), 0, 365, 7)), 10)
		result[systemdomain.ConfigKeyProcurementDefaultMOQ] = strconv.FormatUint(uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyProcurementDefaultMOQ], 1), 1, 1000000, 1)), 10)
		result[systemdomain.ConfigKeyProcurementDefaultOrderMultiple] = strconv.FormatUint(uint64(clampU32(parseUint32OrDefault(result[systemdomain.ConfigKeyProcurementDefaultOrderMultiple], 1), 1, 1000000, 1)), 10)
	default:
		return nil, fmt.Errorf("%w: unsupported module", systemdomain.ErrConfigCenterInvalid)
	}
	return result, nil
}

func normalizeConfigValue(definition *systemdomain.ConfigDefinition, value string) (string, error) {
	switch definition.ValueType {
	case systemdomain.ConfigValueTypeString:
		return value, nil
	case systemdomain.ConfigValueTypeEnum:
		value = strings.ToUpper(value)
		options := parseStringOptions(definition.OptionsJSON)
		if len(options) > 0 && !containsString(options, value) {
			return "", fmt.Errorf("%w: invalid enum %s", systemdomain.ErrConfigCenterInvalid, definition.ConfigKey)
		}
		return value, nil
	case systemdomain.ConfigValueTypeInt:
		if _, err := strconv.ParseUint(value, 10, 32); err != nil {
			return "", fmt.Errorf("%w: invalid integer %s", systemdomain.ErrConfigCenterInvalid, definition.ConfigKey)
		}
		return value, nil
	case systemdomain.ConfigValueTypeBool:
		lowered := strings.ToLower(value)
		if lowered == "1" || lowered == "true" {
			return "1", nil
		}
		if lowered == "0" || lowered == "false" {
			return "0", nil
		}
		return "", fmt.Errorf("%w: invalid boolean %s", systemdomain.ErrConfigCenterInvalid, definition.ConfigKey)
	default:
		return "", fmt.Errorf("%w: unsupported value type %s", systemdomain.ErrConfigCenterInvalid, definition.ValueType)
	}
}

func parseStringOptions(raw *string) []string {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil
	}
	var options []string
	if err := json.Unmarshal([]byte(*raw), &options); err != nil {
		return nil
	}
	return options
}

func findModuleValue(module *systemdomain.ConfigCenterModule, key string) string {
	if module == nil {
		return ""
	}
	for _, group := range module.Groups {
		for _, item := range group.Items {
			if item.ConfigKey == key {
				return strings.TrimSpace(item.Value)
			}
		}
	}
	return ""
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func parseUint32OrDefault(raw string, fallback uint32) uint32 {
	parsed, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 32)
	if err != nil {
		return fallback
	}
	return uint32(parsed)
}

func valueOr(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func clampU32(value uint32, min uint32, max uint32, fallback uint32) uint32 {
	if value == 0 {
		return fallback
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
