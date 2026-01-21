package router

import (
	"net/http"
	"os"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/upload"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	systemHttp "am-erp-go/internal/module/system/delivery/http"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"

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
) *Router {
	return &Router{
		engine:               gin.Default(),
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

			// Inventory routes
			inventory := protected.Group("/inventory")
			{
				inventory.GET("/warehouses/active", r.warehouseHandler.GetActiveWarehouses)
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
