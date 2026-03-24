package domain

import "time"

type OrderProfitListParams struct {
	Page        int
	PageSize    int
	DateFrom    *time.Time
	DateTo      *time.Time
	Marketplace string
	Keyword     string
}

type OrderProfitSummary struct {
	SalesOrderID         uint64    `json:"sales_order_id"`
	OrderNo              string    `json:"order_no"`
	Marketplace          string    `json:"marketplace"`
	BaseCurrency         string    `json:"base_currency"`
	SalesIncomeAmount    float64   `json:"sales_income_amount"`
	COGSAmount           float64   `json:"cogs_amount"`
	GrossProfitAmount    float64   `json:"gross_profit_amount"`
	OrderExpenseAmount   float64   `json:"order_expense_amount"`
	OrderNetProfitAmount float64   `json:"order_net_profit_amount"`
	OccurredAt           time.Time `json:"occurred_at"`
}

type OrderProfitLine struct {
	SalesOrderItemID  uint64  `json:"sales_order_item_id"`
	ProductID         uint64  `json:"product_id"`
	SellerSKU         string  `json:"seller_sku"`
	ProductTitle      string  `json:"product_title"`
	ProductImageURL   string  `json:"product_image_url"`
	QtyShipped        uint64  `json:"qty_shipped"`
	UnitPrice         float64 `json:"unit_price"`
	IncomeAmount      float64 `json:"income_amount"`
	COGSAmount        float64 `json:"cogs_amount"`
	GrossProfitAmount float64 `json:"gross_profit_amount"`
}

type OrderProfitExpense struct {
	ID           uint64    `json:"id"`
	Category     string    `json:"category"`
	BaseCurrency string    `json:"base_currency"`
	BaseAmount   float64   `json:"base_amount"`
	OccurredAt   time.Time `json:"occurred_at"`
	Remark       *string   `json:"remark"`
}

type OrderProfitDetail struct {
	Summary  OrderProfitSummary   `json:"summary"`
	Lines    []OrderProfitLine    `json:"lines"`
	Expenses []OrderProfitExpense `json:"expenses"`
}
