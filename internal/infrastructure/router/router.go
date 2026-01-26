package router

import (
	"net/http"
	"os"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/middleware"
	"am-erp-go/internal/infrastructure/upload"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	systemHttp "am-erp-go/internal/module/system/delivery/http"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"
	shipmentHttp "am-erp-go/internal/module/shipment/delivery/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine               *gin.Engine
	jwtManager           *auth.JWTManager
	authHandler          *identityHttp.AuthHandler
	menuHandler          *menuHttp.MenuHandler
	productHandler       *productHttp.ProductHandler
	productComboHandler  *productHttp.ComboHandler
	supplierHandler      *supplierHttp.SupplierHandler
	supplierQuoteHandler *supplierHttp.QuoteHandler
	purchaseOrderHandler *procurementHttp.PurchaseOrderHandler
	fieldLabelHandler    *systemHttp.FieldLabelHandler
	auditLogHandler      *systemHttp.AuditLogHandler
	uploadHandler        *upload.UploadHandler
	warehouseHandler     *inventoryHttp.WarehouseHandler
	inventoryHandler     *inventoryHttp.InventoryHandler
	shipmentHandler      *shipmentHttp.ShipmentHandler
	packageSpecHandler   *shipmentHttp.PackageSpecHandler
}

func NewRouter(
	jwtManager *auth.JWTManager,
	authHandler *identityHttp.AuthHandler,
	menuHandler *menuHttp.MenuHandler,
	productHandler *productHttp.ProductHandler,
	productComboHandler *productHttp.ComboHandler,
	supplierHandler *supplierHttp.SupplierHandler,
	supplierQuoteHandler *supplierHttp.QuoteHandler,
	purchaseOrderHandler *procurementHttp.PurchaseOrderHandler,
	fieldLabelHandler *systemHttp.FieldLabelHandler,
	auditLogHandler *systemHttp.AuditLogHandler,
	uploadHandler *upload.UploadHandler,
	warehouseHandler *inventoryHttp.WarehouseHandler,
	inventoryHandler *inventoryHttp.InventoryHandler,
	shipmentHandler *shipmentHttp.ShipmentHandler,
	packageSpecHandler *shipmentHttp.PackageSpecHandler,
) *Router {
	// 使用 gin.New() 而不是 gin.Default()，手动配置中间件
	engine := gin.New()
	engine.Use(gin.Logger())           // 日志中间件
	engine.Use(middleware.Recovery())  // 自定义恢复中间件（统一错误响应）
	engine.Use(middleware.ErrorHandler()) // 错误处理中间件

	return &Router{
		engine:               engine,
		jwtManager:           jwtManager,
		authHandler:          authHandler,
		menuHandler:          menuHandler,
		productHandler:       productHandler,
		productComboHandler:  productComboHandler,
		supplierHandler:      supplierHandler,
		supplierQuoteHandler: supplierQuoteHandler,
		purchaseOrderHandler: purchaseOrderHandler,
		fieldLabelHandler:    fieldLabelHandler,
		auditLogHandler:      auditLogHandler,
		uploadHandler:        uploadHandler,
		warehouseHandler:     warehouseHandler,
		inventoryHandler:     inventoryHandler,
		shipmentHandler:      shipmentHandler,
		packageSpecHandler:   packageSpecHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(corsMiddleware())

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	r.engine.Static("/uploads", uploadDir)

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
			menusGroup.GET("", r.menuHandler.ListMenus)
			menusGroup.POST("", r.menuHandler.CreateMenu)
			menusGroup.PUT("/:id", r.menuHandler.UpdateMenu)
			menusGroup.PATCH("/:id/status", r.menuHandler.UpdateMenuStatus)
			menusGroup.DELETE("/:id", r.menuHandler.DeleteMenu)
		}

		systemPublic := api.Group("/system")
		{
			systemPublic.GET("/field-labels", r.fieldLabelHandler.GetLabels)
		}

		systemProtected := api.Group("/system")
		systemProtected.Use(auth.AuthMiddleware(r.jwtManager))
		{
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

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(r.jwtManager))
		{
			// Product (SKU) routes
			products := protected.Group("/products")
			{
				products.GET("", r.productHandler.ListProducts)
				products.GET("/:id", r.productHandler.GetProduct)
				products.POST("", r.productHandler.CreateProduct)
				products.PUT("/:id", r.productHandler.UpdateProduct)
				products.DELETE("/:id", r.productHandler.DeleteProduct)
				products.GET("/:id/images", r.productHandler.ListProductImages)
				products.PUT("/:id/images/reorder", r.productHandler.SaveProductImages)
			}

			// Product Parent routes
			productParents := protected.Group("/product-parents")
			{
				productParents.GET("", r.productHandler.ListProductParents)
				productParents.GET("/:id", r.productHandler.GetProductParent)
				productParents.POST("", r.productHandler.CreateProductParent)
				productParents.PUT("/:id", r.productHandler.UpdateProductParent)
				productParents.DELETE("/:id", r.productHandler.DeleteProductParent)
			}

			// Product Combo routes
			productCombos := protected.Group("/product-combos")
			{
				productCombos.GET("", r.productComboHandler.ListCombos)
				productCombos.GET("/:id", r.productComboHandler.GetCombo)
				productCombos.POST("", r.productComboHandler.CreateCombo)
				productCombos.PUT("/:id", r.productComboHandler.UpdateCombo)
				productCombos.DELETE("/:id", r.productComboHandler.DeleteCombo)
			}

			// Supplier routes
			suppliers := protected.Group("/suppliers")
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
			{
				quotes.GET("", r.supplierQuoteHandler.ListProductQuotes)
				quotes.POST("", r.supplierQuoteHandler.CreateQuote)
				quotes.PUT("", r.supplierQuoteHandler.UpdateQuote)
				quotes.DELETE("", r.supplierQuoteHandler.DeleteQuote)
				quotes.POST("/default", r.supplierQuoteHandler.SetDefaultSupplier)
			}

			registerProcurementRoutes(protected, r.purchaseOrderHandler)

			// Shipment routes
			shipments := protected.Group("/shipments")
			{
				shipments.GET("", r.shipmentHandler.ListShipments)
				shipments.GET("/:id", r.shipmentHandler.GetShipment)
				shipments.POST("", r.shipmentHandler.CreateShipment)
				shipments.POST("/:id/confirm", r.shipmentHandler.ConfirmShipment)     // DRAFT → CONFIRMED
				shipments.POST("/:id/pack", r.shipmentHandler.PackShipment)           // CONFIRMED → PACKED
				shipments.POST("/:id/ship", r.shipmentHandler.MarkShipmentShipped)    // PACKED → SHIPPED
				shipments.POST("/:id/delivered", r.shipmentHandler.MarkShipmentDelivered) // SHIPPED → DELIVERED
				shipments.POST("/:id/cancel", r.shipmentHandler.CancelShipment)       // Cancel with rollback
				shipments.DELETE("/:id", r.shipmentHandler.DeleteShipment)            // Delete DRAFT or CANCELLED
			}

			// Package Spec routes (装箱规格)
			r.packageSpecHandler.RegisterRoutes(protected)

			// Inventory routes
			inventory := protected.Group("/inventory")
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
				inventory.GET("/balances/sku/:sku_id/warehouse/:warehouse_id", r.inventoryHandler.GetSkuBalance)

				// Movements
				inventory.GET("/movements", r.inventoryHandler.ListMovements)
				inventory.GET("/movements/:id", r.inventoryHandler.GetMovement)
				inventory.POST("/movements/purchase-receipt", r.inventoryHandler.CreatePurchaseReceipt)
				inventory.POST("/movements/sales-shipment", r.inventoryHandler.CreateSalesShipment)
				inventory.POST("/movements/stock-take", r.inventoryHandler.CreateStockTakeAdjustment)
				inventory.POST("/movements/manual-adjustment", r.inventoryHandler.CreateManualAdjustment)
				inventory.POST("/movements/damage-write-off", r.inventoryHandler.CreateDamageWriteOff)
				inventory.POST("/movements/return-receipt", r.inventoryHandler.CreateReturnReceipt)
				inventory.POST("/movements/transfer", r.inventoryHandler.CreateTransfer)

				// 库存状态流转
				inventory.POST("/movements/purchase-ship", r.inventoryHandler.CreatePurchaseShip)           // 供应商发货
				inventory.POST("/movements/warehouse-receive", r.inventoryHandler.CreateWarehouseReceive)   // 到仓收货
				inventory.POST("/movements/inspection-pass", r.inventoryHandler.CreateInspectionPass)       // 质检通过
				inventory.POST("/movements/inspection-fail", r.inventoryHandler.CreateInspectionFail)       // 质检不合格
				inventory.POST("/movements/assembly-complete", r.inventoryHandler.CreateAssemblyComplete)   // 组装完成
				inventory.POST("/movements/logistics-ship", r.inventoryHandler.CreateLogisticsShip)         // 物流发货
				inventory.POST("/movements/platform-receive", r.inventoryHandler.CreatePlatformReceive)     // 平台上架
				inventory.POST("/movements/return-receive", r.inventoryHandler.CreateReturnReceive)         // 退货入库
				inventory.POST("/movements/return-inspect", r.inventoryHandler.CreateReturnInspect)         // 退货质检
			}
		}
	}

	legacy := r.engine.Group("/api")
	legacy.Use(auth.AuthMiddleware(r.jwtManager))
	{
		registerProcurementRoutes(legacy, r.purchaseOrderHandler)
	}

	return r.engine
}

func registerProcurementRoutes(group *gin.RouterGroup, handler *procurementHttp.PurchaseOrderHandler) {
	procurement := group.Group("/procurement/purchase-orders")
	{
		procurement.GET("", handler.ListPurchaseOrders)
		procurement.GET("/:id", handler.GetPurchaseOrder)
		procurement.POST("", handler.CreatePurchaseOrder)
		procurement.PUT("/:id", handler.UpdatePurchaseOrder)
		procurement.DELETE("/:id", handler.DeletePurchaseOrder)
		procurement.POST("/:id/submit", handler.SubmitPurchaseOrder)
		procurement.POST("/:id/ship", handler.MarkPurchaseOrderShipped)
		procurement.POST("/:id/receive", handler.ReceivePurchaseOrder)
		procurement.POST("/:id/close", handler.ClosePurchaseOrder)
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
