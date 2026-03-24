package domain

import (
	"reflect"
	"strings"
	"testing"
)

func TestJobRunAuditTimestampFieldsUseAutoTimeTags(t *testing.T) {
	assertFieldHasTag(t, reflect.TypeOf(JobRun{}), "GmtCreate", "autoCreateTime")
	assertFieldHasTag(t, reflect.TypeOf(JobRun{}), "GmtModified", "autoUpdateTime")
}

func TestSystemLogAuditTimestampFieldsUseAutoTimeTags(t *testing.T) {
	assertFieldHasTag(t, reflect.TypeOf(SystemLog{}), "GmtCreate", "autoCreateTime")
	assertFieldHasTag(t, reflect.TypeOf(SystemLog{}), "GmtModified", "autoUpdateTime")
}

func assertFieldHasTag(t *testing.T, typ reflect.Type, fieldName string, expected string) {
	t.Helper()

	field, ok := typ.FieldByName(fieldName)
	if !ok {
		t.Fatalf("field %s not found on %s", fieldName, typ.Name())
	}
	gormTag := field.Tag.Get("gorm")
	if !strings.Contains(gormTag, expected) {
		t.Fatalf("expected %s.%s gorm tag to contain %q, got %q", typ.Name(), fieldName, expected, gormTag)
	}
}
