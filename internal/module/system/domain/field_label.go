package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LabelMap map[string]string

func (m LabelMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(raw), nil
}

func (m *LabelMap) Scan(value interface{}) error {
	if value == nil {
		*m = LabelMap{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported label map type: %T", value)
	}

	if len(data) == 0 {
		*m = LabelMap{}
		return nil
	}

	return json.Unmarshal(data, m)
}

type FieldLabel struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	LabelKey  string    `json:"label_key" gorm:"column:label_key;size:100;not null"`
	Module    string    `json:"module,omitempty" gorm:"column:module;size:50"`
	Scene     string    `json:"scene,omitempty" gorm:"column:scene;size:50"`
	Status    string    `json:"status,omitempty" gorm:"column:status;size:20"`
	Remark    string    `json:"remark,omitempty" gorm:"column:remark;type:text"`
	Labels    LabelMap  `json:"labels" gorm:"column:labels;type:json"`
	CreatedAt time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (FieldLabel) TableName() string {
	return "field_label"
}

func (f *FieldLabel) NormalizeKey() {
	f.LabelKey = strings.ToLower(strings.TrimSpace(f.LabelKey))
}
