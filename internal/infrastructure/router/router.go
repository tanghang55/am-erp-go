package router

import (
	"net/http"
	"os"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/middleware"
	"am-erp-go/internal/infrastructure/upload"
	financeHttp "am-erp-go/internal/module/finance/delivery/http"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	integrationHttp "am-erp-go/internal/module/integration/delivery/http"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"
	logisticsHttp "am-erp-go/internal/module/logistics/delivery/http"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	packagingHttp "am-erp-go/internal/module/packaging/delivery/http"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	salesHttp "am-erp-go/internal/module/sales/delivery/http"
	shipmentHttp "am-erp-go/internal/module/shipment/delivery/http"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	systemHttp "am-erp-go/internal/module/system/delivery/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine                      *gin.Engine
	jwtManager                  *auth.JWTManager
	permissionRepo              auth.PermissionRepository
	authHandler                 *identityHttp.AuthHandler
	userHandler                 *identityHttp.UserHandler
	menuHandler                 *menuHttp.MenuHandler
	productHandler              *productHttp.ProductHandler
	productComboHandler         *productHttp.ComboHandler
	supplierHandler             *supplierHttp.SupplierHandler
	supplierQuoteHandler        *supplierHttp.QuoteHandler
	purchaseOrderHandler        *procurementHttp.PurchaseOrderHandler
	replenishmentHandler        *procurementHttp.ReplenishmentHandler
	configCenterHandler         *systemHttp.ConfigCenterHandler
	fieldLabelHandler           *systemHttp.FieldLabelHandler
	auditLogHandler             *systemHttp.AuditLogHandler
	monitorHandler              *systemHttp.MonitorHandler
	uploadHandler               *upload.UploadHandler
	warehouseHandler            *inventoryHttp.WarehouseHandler
	inventoryHandler            *inventoryHttp.InventoryHandler
	shipmentHandler             *shipmentHttp.ShipmentHandler
	packageSpecHandler          *shipmentHttp.PackageSpecHandler
	logisticsProviderHandler    *logisticsHttp.LogisticsProviderHandler
	shippingRateHandler         *logisticsHttp.ShippingRateHandler
	logisticsServiceHandler     *logisticsHttp.LogisticsServiceHandler
	packagingHandler            *packagingHttp.PackagingHandler
	packagingProcurementHandler *packagingHttp.PackagingProcurementHandler
	salesOrderHandler           *salesHttp.SalesOrderHandler
	financeHandler              *financeHttp.FinanceHandler
	orderSyncHandler            *integrationHttp.OrderSyncHandler
	refundSyncHandler           *integrationHttp.RefundSyncHandler
	authorizationHandler        *integrationHttp.AuthorizationHandler
	skuMappingHandler           *integrationHttp.SKUMappingHandler
}

func NewRouter(
	jwtManager *auth.JWTManager,
	permissionRepo auth.PermissionRepository,
	authHandler *identityHttp.AuthHandler,
	userHandler *identityHttp.UserHandler,
	menuHandler *menuHttp.MenuHandler,
	productHandler *productHttp.ProductHandler,
	productComboHandler *productHttp.ComboHandler,
	supplierHandler *supplierHttp.SupplierHandler,
	supplierQuoteHandler *supplierHttp.QuoteHandler,
	purchaseOrderHandler *procurementHttp.PurchaseOrderHandler,
	replenishmentHandler *procurementHttp.ReplenishmentHandler,
	configCenterHandler *systemHttp.ConfigCenterHandler,
	fieldLabelHandler *systemHttp.FieldLabelHandler,
	auditLogHandler *systemHttp.AuditLogHandler,
	monitorHandler *systemHttp.MonitorHandler,
	uploadHandler *upload.UploadHandler,
	warehouseHandler *inventoryHttp.WarehouseHandler,
	inventoryHandler *inventoryHttp.InventoryHandler,
	shipmentHandler *shipmentHttp.ShipmentHandler,
	packageSpecHandler *shipmentHttp.PackageSpecHandler,
	logisticsProviderHandler *logisticsHttp.LogisticsProviderHandler,
	shippingRateHandler *logisticsHttp.ShippingRateHandler,
	logisticsServiceHandler *logisticsHttp.LogisticsServiceHandler,
	packagingHandler *packagingHttp.PackagingHandler,
	packagingProcurementHandler *packagingHttp.PackagingProcurementHandler,
	salesOrderHandler *salesHttp.SalesOrderHandler,
	financeHandler *financeHttp.FinanceHandler,
	orderSyncHandler *integrationHttp.OrderSyncHandler,
	refundSyncHandler *integrationHttp.RefundSyncHandler,
	authorizationHandler *integrationHttp.AuthorizationHandler,
	skuMappingHandler *integrationHttp.SKUMappingHandler,
) *Router {
	// 使用 gin.New() 而不是 gin.Default()，手动配置中间件
	engine := gin.New()
	engine.Use(gin.Logger())              // 日志中间件
	engine.Use(middleware.Recovery())     // 自定义恢复中间件（统一错误响应）
	engine.Use(middleware.ErrorHandler()) // 错误处理中间件

	return &Router{
		engine:                      engine,
		jwtManager:                  jwtManager,
		permissionRepo:              permissionRepo,
		authHandler:                 authHandler,
		userHandler:                 userHandler,
		menuHandler:                 menuHandler,
		productHandler:              productHandler,
		productComboHandler:         productComboHandler,
		supplierHandler:             supplierHandler,
		supplierQuoteHandler:        supplierQuoteHandler,
		purchaseOrderHandler:        purchaseOrderHandler,
		replenishmentHandler:        replenishmentHandler,
		configCenterHandler:         configCenterHandler,
		fieldLabelHandler:           fieldLabelHandler,
		auditLogHandler:             auditLogHandler,
		monitorHandler:              monitorHandler,
		uploadHandler:               uploadHandler,
		warehouseHandler:            warehouseHandler,
		inventoryHandler:            inventoryHandler,
		shipmentHandler:             shipmentHandler,
		packageSpecHandler:          packageSpecHandler,
		logisticsProviderHandler:    logisticsProviderHandler,
		shippingRateHandler:         shippingRateHandler,
		logisticsServiceHandler:     logisticsServiceHandler,
		packagingHandler:            packagingHandler,
		packagingProcurementHandler: packagingProcurementHandler,
		salesOrderHandler:           salesOrderHandler,
		financeHandler:              financeHandler,
		orderSyncHandler:            orderSyncHandler,
		refundSyncHandler:           refundSyncHandler,
		authorizationHandler:        authorizationHandler,
		skuMappingHandler:           skuMappingHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(corsMiddleware())

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	r.engine.Static(upload.ResolveURLBase(), uploadDir)

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	api := r.engine.Group("/api/v1")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", r.authHandler.Login)

			protected := authGroup.Group("")
			protected.Use(auth.AuthMiddleware(r.jwtManager))
			{
				protected.GET("/me", r.authHandler.GetCurrentUser)
			}
		}

		// Menu routes
		menusGroup := api.Group("/menus")
		menusGroup.Use(auth.AuthMiddleware(r.jwtManager))
		{
			menusGroup.GET("/tree", r.menuHandler.GetMenuTree)
			menuManage := menusGroup.Group("")
			menuManage.Use(auth.RequirePermission(r.permissionRepo, "system.menu.manage"))
			{
				menuManage.GET("", r.menuHandler.ListMenus)
				menuManage.POST("", r.menuHandler.CreateMenu)
				menuManage.PUT("/:id", r.menuHandler.UpdateMenu)
				menuManage.PATCH("/:id/status", r.menuHandler.UpdateMenuStatus)
				menuManage.DELETE("/:id", r.menuHandler.DeleteMenu)
			}
		}

		systemPublic := api.Group("/system")
		{
			systemPublic.GET("/field-labels", r.fieldLabelHandler.GetLabels)
		}

		systemProtected := api.Group("/system")
		systemProtected.Use(auth.AuthMiddleware(r.jwtManager))
		systemProtected.Use(auth.RequirePermission(r.permissionRepo, "system.manage"))
		{
			systemProtected.GET("/monitor/overview", r.monitorHandler.GetOverview)
			systemProtected.GET("/monitor/jobs", r.monitorHandler.ListRecentJobs)
			systemProtected.GET("/monitor/logs", r.monitorHandler.ListRecentLogs)
			systemProtected.GET("/config-center/modules", r.configCenterHandler.ListModules)
			systemProtected.GET("/config-center/modules/:module", r.configCenterHandler.GetModule)
			systemProtected.PUT("/config-center/modules/:module", r.configCenterHandler.UpdateModule)
			systemProtected.GET("/logs", r.auditLogHandler.List)
			systemProtected.GET("/field-labels/manage", r.fieldLabelHandler.List)
			systemProtected.POST("/field-labels", r.fieldLabelHandler.Create)
			systemProtected.PUT("/field-labels/:id", r.fieldLabelHandler.Update)
			systemProtected.DELETE("/field-labels/:id", r.fieldLabelHandler.Delete)
		}

		uploadGroup := api.Group("/upload")
		uploadGroup.Use(auth.AuthMiddleware(r.jwtManager))
		{
			uploadGroup.POST("/image", r.uploadHandler.UploadImage)
		}

		identityGroup := api.Group("/identity")
		identityGroup.Use(auth.AuthMiddleware(r.jwtManager))
		identityGroup.Use(auth.RequirePermission(r.permissionRepo, "identity.manage"))
		{
			identityGroup.GET("/users", r.userHandler.ListUsers)
			identityGroup.GET("/users/:id", r.userHandler.GetUser)
			identityGroup.POST("/users", r.userHandler.CreateUser)
			identityGroup.PUT("/users/:id", r.userHandler.UpdateUser)
			identityGroup.DELETE("/users/:id", r.userHandler.DeleteUser)
			identityGroup.POST("/users/:id/roles", r.userHandler.AssignUserRoles)
			identityGroup.GET("/roles", r.userHandler.ListRoles)
			identityGroup.GET("/permissions", r.userHandler.ListPermissions)
		}

		if r.authorizationHandler != nil {
			integrationProtected := api.Group("/integrations/authorizations")
			integrationProtected.Use(auth.AuthMiddleware(r.jwtManager))
			integrationProtected.Use(auth.RequirePermission(r.permissionRepo, "integration.manage"))
			{
				integrationProtected.GET("/providers", r.authorizationHandler.ListProviders)
				integrationProtected.GET("", r.authorizationHandler.ListAuthorizations)
				integrationProtected.POST("/start", r.authorizationHandler.StartAuthorization)
				integrationProtected.POST("/:id/refresh", r.authorizationHandler.ManualRefresh)
				if r.skuMappingHandler != nil {
					integrationProtected.GET("/sku-mappings", r.skuMappingHandler.List)
					integrationProtected.POST("/sku-mappings", r.skuMappingHandler.Create)
					integrationProtected.PUT("/sku-mappings/:id", r.skuMappingHandler.Update)
				}
			}

			integrationPublic := api.Group("/integrations/providers/:provider/oauth")
			{
				integrationPublic.GET("/callback", r.authorizationHandler.HandleOAuthCallback)
			}
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(r.jwtManager))
		{
			// Product routes
			products := protected.Group("/products")
			products.Use(auth.RequirePermission(r.permissionRepo, "product.manage"))
			{
				products.GET("", r.productHandler.ListProducts)
				products.GET("/:id", r.productHandler.GetProduct)
				products.POST("", r.productHandler.CreateProduct)
				products.PUT("/:id", r.productHandler.UpdateProduct)
				products.DELETE("/:id", r.productHandler.DeleteProduct)
				products.GET("/:id/images", r.productHandler.ListProductImages)
				products.PUT("/:id/images/reorder", r.productHandler.SaveProductImages)
				products.GET("/:id/packaging-items", r.productHandler.GetProductPackagingItems)
				products.PUT("/:id/packaging-items", r.productHandler.SaveProductPackagingItems)
			}

			productConfigs := protected.Group("/product-configs")
			productConfigs.Use(auth.RequirePermission(r.permissionRepo, "product.manage"))
			{
				productConfigs.GET("", r.productHandler.ListProductConfigs)
				productConfigs.POST("", r.productHandler.CreateProductConfig)
				productConfigs.PUT("/:id", r.productHandler.UpdateProductConfig)
				productConfigs.DELETE("/:id", r.productHandler.DeleteProductConfig)
			}

			productCategories := protected.Group("/product-categories")
			productCategories.Use(auth.RequirePermission(r.permissionRepo, "product.manage"))
			{
				productCategories.GET("", r.productHandler.ListProductCategories)
				productCategories.POST("", r.productHandler.CreateProductCategory)
				productCategories.PUT("/:id", r.productHandler.UpdateProductCategory)
				productCategories.DELETE("/:id", r.productHandler.DeleteProductCategory)
			}

			// Product Parent routes
			productParents := protected.Group("/product-parents")
			productParents.Use(auth.RequirePermission(r.permissionRepo, "product.manage"))
			{
				productParents.GET("", r.productHandler.ListProductParents)
				productParents.GET("/:id", r.productHandler.GetProductParent)
				productParents.POST("", r.productHandler.CreateProductParent)
				productParents.PUT("/:id", r.productHandler.UpdateProductParent)
				productParents.DELETE("/:id", r.productHandler.DeleteProductParent)
				productParents.POST("/:id/children", r.productHandler.AttachProductParentChildren)
				productParents.DELETE("/:id/children/:childId", r.productHandler.DetachProductParentChild)
			}

			// Product Combo routes
			productCombos := protected.Group("/product-combos")
			productCombos.Use(auth.RequirePermission(r.permissionRepo, "product.manage"))
			{
				productCombos.GET("", r.productComboHandler.ListCombos)
				productCombos.GET("/:id", r.productComboHandler.GetCombo)
				productCombos.POST("", r.productComboHandler.CreateCombo)
				productCombos.PUT("/:id", r.productComboHandler.UpdateCombo)
				productCombos.DELETE("/:id", r.productComboHandler.DeleteCombo)
			}

			// Supplier routes
			suppliers := protected.Group("/suppliers")
			suppliers.Use(auth.RequirePermission(r.permissionRepo, "supplier.manage"))
			{
				suppliers.GET("", r.supplierHandler.ListSuppliers)
				suppliers.GET("/:id", r.supplierHandler.GetSupplier)
				suppliers.POST("", r.supplierHandler.CreateSupplier)
				suppliers.PUT("/:id", r.supplierHandler.UpdateSupplier)
				suppliers.DELETE("/:id", r.supplierHandler.DeleteSupplier)
				suppliers.POST("/:id/contacts", r.supplierHandler.CreateSupplierContact)
				suppliers.PUT("/:id/contacts", r.supplierHandler.UpdateSupplierContact)
				suppliers.DELETE("/:id/contacts", r.supplierHandler.DeleteSupplierContact)
				suppliers.POST("/:id/accounts", r.supplierHandler.CreateSupplierAccount)
				suppliers.PUT("/:id/accounts", r.supplierHandler.UpdateSupplierAccount)
				suppliers.DELETE("/:id/accounts", r.supplierHandler.DeleteSupplierAccount)
				suppliers.POST("/:id/tags", r.supplierHandler.CreateSupplierTag)
				suppliers.PUT("/:id/tags", r.supplierHandler.UpdateSupplierTag)
				suppliers.DELETE("/:id/tags", r.supplierHandler.DeleteSupplierTag)
			}

			// Supplier product quotes
			quotes := protected.Group("/suppliers/product-quotes")
			quotes.Use(auth.RequirePermission(r.permissionRepo, "supplier.manage"))
			{
				quotes.GET("", r.supplierQuoteHandler.ListProductQuotes)
				quotes.GET("/detail", r.supplierQuoteHandler.GetProductQuote)
				quotes.POST("", r.supplierQuoteHandler.CreateQuote)
				quotes.PUT("", r.supplierQuoteHandler.UpdateQuote)
				quotes.DELETE("", r.supplierQuoteHandler.DeleteQuote)
				quotes.POST("/default", r.supplierQuoteHandler.SetDefaultSupplier)
			}

			registerProcurementRoutes(protected, r.permissionRepo, r.purchaseOrderHandler, r.replenishmentHandler)

			// Shipment routes
			shipments := protected.Group("/shipments")
			shipments.Use(auth.RequirePermission(r.permissionRepo, "shipping.manage"))
			{
				shipments.GET("", r.shipmentHandler.ListShipments)
				shipments.GET("/:id", r.shipmentHandler.GetShipment)
				shipments.POST("", r.shipmentHandler.CreateShipment)
				shipments.PUT("/:id", r.shipmentHandler.UpdateShipment)
				shipments.POST("/:id/confirm", r.shipmentHandler.ConfirmShipment)         // DRAFT → CONFIRMED
				shipments.POST("/:id/ship", r.shipmentHandler.MarkShipmentShipped)        // CONFIRMED → SHIPPED
				shipments.POST("/:id/delivered", r.shipmentHandler.MarkShipmentDelivered) // SHIPPED → DELIVERED
				shipments.POST("/:id/cancel", r.shipmentHandler.CancelShipment)           // Cancel with rollback
				shipments.DELETE("/:id", r.shipmentHandler.DeleteShipment)                // Delete DRAFT or CANCELLED
			}

			// Package Spec routes (装箱规格)
			packageSpecGroup := protected.Group("")
			packageSpecGroup.Use(auth.RequirePermission(r.permissionRepo, "shipping.manage"))
			r.packageSpecHandler.RegisterRoutes(packageSpecGroup)

			// Logistics routes
			logisticsGroup := protected.Group("")
			logisticsGroup.Use(auth.RequirePermission(r.permissionRepo, "logistics.manage"))
			r.logisticsProviderHandler.RegisterProviderRoutes(logisticsGroup)
			r.shippingRateHandler.RegisterRateRoutes(logisticsGroup)
			r.logisticsServiceHandler.RegisterServiceRoutes(logisticsGroup)

			// Packaging routes
			packaging := protected.Group("/packaging")
			packaging.Use(auth.RequirePermission(r.permissionRepo, "packaging.manage"))
			{
				packaging.GET("/items", r.packagingHandler.ListItems)
				packaging.GET("/items/low-stock", r.packagingHandler.GetLowStockItems)
				packaging.GET("/items/:id", r.packagingHandler.GetItem)
				packaging.POST("/items", r.packagingHandler.CreateItem)
				packaging.PUT("/items/:id", r.packagingHandler.UpdateItem)
				packaging.DELETE("/items/:id", r.packagingHandler.DeleteItem)

				packaging.GET("/ledger", r.packagingHandler.ListLedgers)
				packaging.GET("/ledger/:id", r.packagingHandler.GetLedger)
				packaging.POST("/ledger/inbound", r.packagingHandler.CreateInboundLedger)
				packaging.POST("/ledger/outbound", r.packagingHandler.CreateOutboundLedger)
				packaging.POST("/ledger/adjustment", r.packagingHandler.CreateAdjustmentLedger)
				packaging.GET("/ledger/usage-summary", r.packagingHandler.GetUsageSummary)

				packaging.GET("/procurement/plans", r.packagingProcurementHandler.ListPlans)
				packaging.GET("/procurement/runs", r.packagingProcurementHandler.ListRuns)
				packaging.POST("/procurement/plans/generate", r.packagingProcurementHandler.GeneratePlans)
				packaging.POST("/procurement/plans/convert", r.packagingProcurementHandler.ConvertPlans)
				packaging.GET("/procurement/orders", r.packagingProcurementHandler.ListPurchaseOrders)
				packaging.GET("/procurement/orders/:id", r.packagingProcurementHandler.GetPurchaseOrder)
				packaging.POST("/procurement/orders/:id/submit", r.packagingProcurementHandler.SubmitPurchaseOrder)
				packaging.POST("/procurement/orders/:id/receive", r.packagingProcurementHandler.ReceivePurchaseOrder)
			}

			// Sales order routes
			sales := protected.Group("/sales/orders")
			sales.Use(auth.RequirePermission(r.permissionRepo, "sales.manage"))
			{
				sales.GET("", r.salesOrderHandler.ListSalesOrders)
				sales.GET("/:id", r.salesOrderHandler.GetSalesOrder)
				sales.POST("", r.salesOrderHandler.CreateSalesOrder)
				sales.POST("/import", r.salesOrderHandler.ImportSalesOrders)
				sales.GET("/imports", r.salesOrderHandler.ListImportBatches)
				sales.GET("/imports/:id", r.salesOrderHandler.GetImportBatch)
				sales.GET("/imports/:id/errors", r.salesOrderHandler.ListImportBatchErrors)
				sales.PUT("/:id", r.salesOrderHandler.UpdateSalesOrder)
				sales.POST("/:id/confirm", r.salesOrderHandler.ConfirmSalesOrder)
				sales.POST("/:id/allocate", r.salesOrderHandler.AllocateSalesOrder)
				sales.POST("/:id/ship", r.salesOrderHandler.ShipSalesOrder)
				sales.POST("/:id/deliver", r.salesOrderHandler.DeliverSalesOrder)
				sales.POST("/:id/cancel", r.salesOrderHandler.CancelSalesOrder)
				sales.POST("/:id/return", r.salesOrderHandler.ReturnSalesOrder)
			}

			registerFinanceRoutes(protected, r.permissionRepo, r.financeHandler)
			registerIntegrationRoutes(protected, r.permissionRepo, r.orderSyncHandler, r.refundSyncHandler)

			// Inventory routes
			inventory := protected.Group("/inventory")
			inventory.Use(auth.RequirePermission(r.permissionRepo, "inventory.manage"))
			{
				// Warehouses - 完整的CRUD接口
				inventory.GET("/warehouses", r.warehouseHandler.ListWarehouses)
				inventory.GET("/warehouses/active", r.warehouseHandler.GetActiveWarehouses)
				inventory.GET("/warehouses/:id", r.warehouseHandler.GetWarehouse)
				inventory.POST("/warehouses", r.warehouseHandler.CreateWarehouse)
				inventory.PUT("/warehouses/:id", r.warehouseHandler.UpdateWarehouse)
				inventory.DELETE("/warehouses/:id", r.warehouseHandler.DeleteWarehouse)

				// Balances
				inventory.GET("/balances", r.inventoryHandler.ListBalances)
				inventory.GET("/lots", r.inventoryHandler.ListLots)
				inventory.GET("/balances/product/:product_id/warehouse/:warehouse_id", r.inventoryHandler.GetProductBalance)

				// Movements
				inventory.GET("/movements", r.inventoryHandler.ListMovements)
				inventory.GET("/movements/:id", r.inventoryHandler.GetMovement)
				inventory.POST("/movements/stock-take", r.inventoryHandler.CreateStockTakeAdjustment)
				inventory.POST("/movements/manual-adjustment", r.inventoryHandler.CreateManualAdjustment)
				inventory.POST("/movements/damage-write-off", r.inventoryHandler.CreateDamageWriteOff)
				inventory.POST("/movements/transfer", r.inventoryHandler.CreateTransfer)

				// 仅保留专用业务页仍在使用的入口
				inventory.POST("/movements/assembly-complete", r.inventoryHandler.CreateAssemblyComplete) // 组装完成
				inventory.POST("/movements/platform-receive", r.inventoryHandler.CreatePlatformReceive)   // 平台上架
			}
		}
	}

	legacy := r.engine.Group("/api")
	legacy.Use(auth.AuthMiddleware(r.jwtManager))
	{
		registerProcurementRoutes(legacy, r.permissionRepo, r.purchaseOrderHandler, r.replenishmentHandler)
		registerFinanceRoutes(legacy, r.permissionRepo, r.financeHandler)
		registerIntegrationRoutes(legacy, r.permissionRepo, r.orderSyncHandler, r.refundSyncHandler)
	}

	return r.engine
}

func registerIntegrationRoutes(group *gin.RouterGroup, permissionRepo auth.PermissionRepository, handler *integrationHttp.OrderSyncHandler, refundHandler *integrationHttp.RefundSyncHandler) {
	if handler == nil && refundHandler == nil {
		return
	}

	integration := group.Group("/integrations/providers/:provider")
	integration.Use(auth.RequirePermission(permissionRepo, "integration.manage"))
	{
		if handler != nil {
			integration.GET("/orders/sync/state", handler.GetState)
			integration.GET("/orders/sync/runs", handler.ListRuns)
			integration.POST("/orders/sync", handler.SyncOrders)
		}
		if refundHandler != nil {
			integration.GET("/refunds/sync/state", refundHandler.GetState)
			integration.GET("/refunds/sync/runs", refundHandler.ListRuns)
			integration.POST("/refunds/sync", refundHandler.SyncRefunds)
		}
	}
}

func registerFinanceRoutes(group *gin.RouterGroup, permissionRepo auth.PermissionRepository, handler *financeHttp.FinanceHandler) {
	if handler == nil {
		return
	}

	finance := group.Group("/finance")
	finance.Use(auth.RequirePermission(permissionRepo, "finance.manage"))
	{
		finance.GET("/cash-ledger", handler.ListCashLedger)
		finance.GET("/cash-ledger/summary", handler.GetCashLedgerSummary)
		finance.GET("/cash-ledger/summary-by-category", handler.GetCashLedgerSummaryByCategory)
		finance.GET("/cash-ledger/:id", handler.GetCashLedger)
		finance.POST("/cash-ledger", handler.CreateCashLedger)
		finance.POST("/cash-ledger/:id/reverse", handler.ReverseCashLedger)
		finance.PUT("/cash-ledger/:id", handler.UpdateCashLedger)
		finance.DELETE("/cash-ledger/:id", handler.DeleteCashLedger)
		finance.GET("/exchange-rates", handler.ListExchangeRates)
		finance.POST("/exchange-rates", handler.CreateExchangeRate)
		finance.PATCH("/exchange-rates/:id/status", handler.UpdateExchangeRateStatus)
		finance.GET("/profit/dashboard", handler.GetProfitDashboard)
		finance.GET("/profit/orders", handler.ListOrderProfits)
		finance.GET("/profit/orders/:id", handler.GetOrderProfitDetail)
		finance.POST("/profit/rebuild", handler.RebuildDailyProfit)
		finance.GET("/product-cost/ledger", handler.ListProductCostLedger)
		finance.GET("/product-cost/summary", handler.GetProductCostSummary)

		finance.GET("/costing/snapshots", handler.ListCostingSnapshots)
		finance.GET("/costing/snapshots/:id", handler.GetCostingSnapshot)
		finance.POST("/costing/snapshots", handler.CreateCostingSnapshot)
		finance.PUT("/costing/snapshots/:id", handler.UpdateCostingSnapshot)
		finance.DELETE("/costing/snapshots/:id", handler.DeleteCostingSnapshot)
		finance.GET("/costing/current/:product_id", handler.GetCurrentCost)
		finance.GET("/costing/current/:product_id/all", handler.GetAllCurrentCosts)
	}
}

func registerProcurementRoutes(group *gin.RouterGroup, permissionRepo auth.PermissionRepository, handler *procurementHttp.PurchaseOrderHandler, replenishmentHandler *procurementHttp.ReplenishmentHandler) {
	procurement := group.Group("/procurement/purchase-orders")
	procurement.Use(auth.RequirePermission(permissionRepo, "procurement.manage"))
	{
		procurement.GET("", handler.ListPurchaseOrders)
		procurement.GET("/:id", handler.GetPurchaseOrder)
		procurement.POST("", handler.CreatePurchaseOrder)
		procurement.POST("/batch", handler.CreatePurchaseOrderBatch)
		procurement.PUT("/:id", handler.UpdatePurchaseOrder)
		procurement.DELETE("/:id", handler.DeletePurchaseOrder)
		procurement.POST("/:id/submit", handler.SubmitPurchaseOrder)
		procurement.POST("/:id/ship", handler.MarkPurchaseOrderShipped)
		procurement.POST("/:id/receive", handler.ReceivePurchaseOrder)
		procurement.POST("/:id/inspect", handler.InspectPurchaseOrder)
		procurement.POST("/:id/close", handler.ClosePurchaseOrder)
		procurement.POST("/:id/force-complete", handler.ForceCompletePurchaseOrder)
	}

	if replenishmentHandler == nil {
		return
	}
	replenishment := group.Group("/procurement/replenishment")
	replenishment.Use(auth.RequirePermission(permissionRepo, "procurement.manage"))
	{
		replenishment.GET("/config", replenishmentHandler.GetConfig)
		replenishment.PUT("/config", replenishmentHandler.UpdateConfig)
		replenishment.GET("/strategies", replenishmentHandler.ListStrategies)
		replenishment.POST("/strategies", replenishmentHandler.UpsertStrategy)
		replenishment.GET("/plans", replenishmentHandler.ListPlans)
		replenishment.DELETE("/plans/:id", replenishmentHandler.DeletePlan)
		replenishment.GET("/runs", replenishmentHandler.ListRuns)
		replenishment.POST("/plans/generate", replenishmentHandler.GeneratePlans)
		replenishment.POST("/plans/convert", replenishmentHandler.ConvertPlans)
		replenishment.POST("/plans/cleanup", replenishmentHandler.CleanupPlans)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
