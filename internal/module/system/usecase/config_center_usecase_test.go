package usecase

import (
	"testing"

	"am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type fakeConfigCenterRepo struct {
	definitions []*domain.ConfigDefinition
	values      []*domain.ConfigValue
	upserts     []*domain.ConfigValue
}

func (f *fakeConfigCenterRepo) ListDefinitions(moduleCode string) ([]*domain.ConfigDefinition, error) {
	if moduleCode == "" {
		return f.definitions, nil
	}
	result := make([]*domain.ConfigDefinition, 0)
	for _, item := range f.definitions {
		if item != nil && item.ModuleCode == moduleCode {
			result = append(result, item)
		}
	}
	return result, nil
}

func (f *fakeConfigCenterRepo) ListValues(scopeType string, scopeRefID uint64, keys []string) ([]*domain.ConfigValue, error) {
	keySet := map[string]struct{}{}
	for _, key := range keys {
		keySet[key] = struct{}{}
	}
	result := make([]*domain.ConfigValue, 0)
	for _, item := range f.values {
		if item == nil || item.ScopeType != scopeType || item.ScopeRefID != scopeRefID {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[item.ConfigKey]; !ok {
				continue
			}
		}
		result = append(result, item)
	}
	return result, nil
}

func (f *fakeConfigCenterRepo) UpsertValues(items []*domain.ConfigValue) error {
	f.upserts = append([]*domain.ConfigValue{}, items...)
	next := make([]*domain.ConfigValue, 0, len(items))
	for _, item := range items {
		copied := *item
		next = append(next, &copied)
	}
	f.values = next
	return nil
}

type fakeConfigCenterAuditLogger struct {
	payloads []AuditLogPayload
}

func (f *fakeConfigCenterAuditLogger) RecordFromContext(_ *gin.Context, payload AuditLogPayload) error {
	f.payloads = append(f.payloads, payload)
	return nil
}

func TestConfigCenterUsecaseGetModuleUsesFallbackDefaults(t *testing.T) {
	options := `["USD","EUR"]`
	repo := &fakeConfigCenterRepo{
		definitions: []*domain.ConfigDefinition{
			{
				ConfigKey:    domain.ConfigKeyFinanceDefaultCurrency,
				ModuleCode:   domain.ConfigModuleFinance,
				ModuleName:   "财务配置",
				GroupCode:    "base",
				GroupName:    "财务基础",
				Label:        "本位币",
				ValueType:    domain.ConfigValueTypeEnum,
				DefaultValue: "USD",
				OptionsJSON:  &options,
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         10,
			},
			{
				ConfigKey:    domain.ConfigKeyFinanceExchangeRateScale,
				ModuleCode:   domain.ConfigModuleFinance,
				ModuleName:   "财务配置",
				GroupCode:    "base",
				GroupName:    "财务基础",
				Label:        "汇率小数位",
				ValueType:    domain.ConfigValueTypeInt,
				DefaultValue: "4",
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         20,
			},
		},
	}
	uc := NewConfigCenterUsecase(repo, nil)

	module, err := uc.GetModule(domain.ConfigModuleFinance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if module.ModuleCode != domain.ConfigModuleFinance || len(module.Groups) != 1 || len(module.Groups[0].Items) != 2 {
		t.Fatalf("unexpected module payload: %+v", module)
	}
	if module.Groups[0].Items[0].Value != "USD" {
		t.Fatalf("expected fallback value USD, got %s", module.Groups[0].Items[0].Value)
	}
	financeConfig := uc.GetFinanceConfig()
	if financeConfig.ExchangeRateScale != 4 {
		t.Fatalf("expected fallback exchange rate scale 4, got %d", financeConfig.ExchangeRateScale)
	}
}

func TestConfigCenterUsecaseUpdateModuleWritesAudit(t *testing.T) {
	repo := &fakeConfigCenterRepo{
		definitions: []*domain.ConfigDefinition{
			{
				ConfigKey:    domain.ConfigKeySalesImportDefaultChannel,
				ModuleCode:   domain.ConfigModuleSalesImport,
				ModuleName:   "订单导入配置",
				GroupCode:    "defaults",
				GroupName:    "默认值",
				Label:        "默认渠道",
				ValueType:    domain.ConfigValueTypeString,
				DefaultValue: "MANUAL",
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         10,
			},
			{
				ConfigKey:    domain.ConfigKeySalesImportDefaultMarketplace,
				ModuleCode:   domain.ConfigModuleSalesImport,
				ModuleName:   "订单导入配置",
				GroupCode:    "defaults",
				GroupName:    "默认值",
				Label:        "默认站点",
				ValueType:    domain.ConfigValueTypeString,
				DefaultValue: "US",
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         20,
			},
		},
	}
	auditLogger := &fakeConfigCenterAuditLogger{}
	uc := NewConfigCenterUsecase(repo, auditLogger)

	module, err := uc.UpdateModule(nil, domain.ConfigModuleSalesImport, map[string]string{
		domain.ConfigKeySalesImportDefaultChannel:     "MANUAL-IMPORT",
		domain.ConfigKeySalesImportDefaultMarketplace: "DE",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.upserts) != 2 {
		t.Fatalf("unexpected upserts: %+v", repo.upserts)
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected no audit without gin context, got %+v", auditLogger.payloads)
	}
	if module.Groups[0].Items[0].Value != "MANUAL-IMPORT" {
		t.Fatalf("expected module value updated, got %+v", module.Groups[0].Items[0])
	}
}

func TestConfigCenterUsecaseRejectsFinanceBaseCurrencyChange(t *testing.T) {
	options := `["USD","EUR"]`
	repo := &fakeConfigCenterRepo{
		definitions: []*domain.ConfigDefinition{
			{
				ConfigKey:    domain.ConfigKeyFinanceDefaultCurrency,
				ModuleCode:   domain.ConfigModuleFinance,
				ModuleName:   "财务配置",
				GroupCode:    "base",
				GroupName:    "财务基础",
				Label:        "本位币",
				ValueType:    domain.ConfigValueTypeEnum,
				DefaultValue: "USD",
				OptionsJSON:  &options,
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         10,
			},
			{
				ConfigKey:    domain.ConfigKeyFinanceExchangeRateScale,
				ModuleCode:   domain.ConfigModuleFinance,
				ModuleName:   "财务配置",
				GroupCode:    "base",
				GroupName:    "财务基础",
				Label:        "汇率小数位",
				ValueType:    domain.ConfigValueTypeInt,
				DefaultValue: "4",
				ScopeType:    domain.ConfigScopeGlobal,
				IsActive:     1,
				Sort:         20,
			},
		},
		values: []*domain.ConfigValue{
			{ConfigKey: domain.ConfigKeyFinanceDefaultCurrency, ScopeType: domain.ConfigScopeGlobal, ScopeRefID: 0, ConfigValue: "USD"},
			{ConfigKey: domain.ConfigKeyFinanceExchangeRateScale, ScopeType: domain.ConfigScopeGlobal, ScopeRefID: 0, ConfigValue: "4"},
		},
	}
	uc := NewConfigCenterUsecase(repo, nil)

	_, err := uc.UpdateModule(nil, domain.ConfigModuleFinance, map[string]string{
		domain.ConfigKeyFinanceDefaultCurrency:   "EUR",
		domain.ConfigKeyFinanceExchangeRateScale: "4",
	}, nil)
	if err == nil {
		t.Fatalf("expected finance base currency update rejected")
	}
	if len(repo.upserts) != 0 {
		t.Fatalf("expected no upserts, got %+v", repo.upserts)
	}
}
