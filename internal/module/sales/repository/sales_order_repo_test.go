package repository

import (
	"reflect"
	"strings"
	"testing"

	"am-erp-go/internal/module/sales/domain"
)

func TestSalesOrderModelTableNames(t *testing.T) {
	var order domain.SalesOrder
	if order.TableName() != "sales_order" {
		t.Fatalf("table name mismatch: got %s", order.TableName())
	}

	var item domain.SalesOrderItem
	if item.TableName() != "sales_order_item" {
		t.Fatalf("table name mismatch: got %s", item.TableName())
	}
}

func TestDefaultStockPoolBySourceType(t *testing.T) {
	if defaultStockPoolBySourceType("AMAZON_API") != domain.StockPoolSellable {
		t.Fatalf("expected AMAZON_API to use SELLABLE pool")
	}
	if defaultStockPoolBySourceType("AMAZON_IMPORT") != domain.StockPoolSellable {
		t.Fatalf("expected AMAZON_IMPORT to use SELLABLE pool")
	}
	if defaultStockPoolBySourceType("MANUAL_IMPORT") != domain.StockPoolAvailable {
		t.Fatalf("expected MANUAL_IMPORT to use AVAILABLE pool")
	}
}

func TestSalesOrderItemDisplayFieldsAreReadOnly(t *testing.T) {
	itemType := reflect.TypeOf(domain.SalesOrderItem{})
	for _, fieldName := range []string{"SellerSKU", "ProductTitle", "ProductImageURL"} {
		field, ok := itemType.FieldByName(fieldName)
		if !ok {
			t.Fatalf("field %s not found", fieldName)
		}
		tag := field.Tag.Get("gorm")
		if !strings.Contains(tag, "->") {
			t.Fatalf("field %s must be read-only, got gorm tag %q", fieldName, tag)
		}
	}
}
