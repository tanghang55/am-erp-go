package bootstrap

import (
	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/router"
	"am-erp-go/internal/infrastructure/upload"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	identityRepo "am-erp-go/internal/module/identity/repository"
	identityUsecase "am-erp-go/internal/module/identity/usecase"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	menuRepo "am-erp-go/internal/module/menu/repository"
	menuUsecase "am-erp-go/internal/module/menu/usecase"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	procurementRepo "am-erp-go/internal/module/procurement/repository"
	procurementUsecase "am-erp-go/internal/module/procurement/usecase"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	productRepo "am-erp-go/internal/module/product/repository"
	productUsecase "am-erp-go/internal/module/product/usecase"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	supplierRepo "am-erp-go/internal/module/supplier/repository"
	supplierUsecase "am-erp-go/internal/module/supplier/usecase"
	systemHttp "am-erp-go/internal/module/system/delivery/http"
	systemRepo "am-erp-go/internal/module/system/repository"
	systemUsecase "am-erp-go/internal/module/system/usecase"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
	Config *config.Config
}

func Build() (*App, error) {
	cfg := config.Load()

	database, err := db.NewMySQL(&cfg.Database)
	if err != nil {
		return nil, err
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	userRepo := identityRepo.NewUserRepository(database)
	authUsecase := identityUsecase.NewAuthUsecase(userRepo, jwtManager)
	authHandler := identityHttp.NewAuthHandler(authUsecase)

	menuRepository := menuRepo.NewMenuRepository(database)
	menuSvc := menuUsecase.NewMenuUsecase(menuRepository, userRepo)
	menuHandler := menuHttp.NewMenuHandler(menuSvc)

	fieldLabelRepository := systemRepo.NewFieldLabelRepository(database)
	fieldLabelSvc := systemUsecase.NewFieldLabelUseCase(fieldLabelRepository)
	fieldLabelHandler := systemHttp.NewFieldLabelHandler(fieldLabelSvc)

	auditLogRepository := systemRepo.NewAuditLogRepository(database)
	auditLogUsecase := systemUsecase.NewAuditLogUsecase(auditLogRepository)
	auditLogHandler := systemHttp.NewAuditLogHandler(auditLogUsecase)

	productRepository := productRepo.NewProductRepository(database)
	productParentRepository := productRepo.NewProductParentRepository(database)
	productSvc := productUsecase.NewProductUsecase(
		productRepository,
		productParentRepository,
	)
	imageRepository := productRepo.NewProductImageRepository(database)
	imageUsecase := productUsecase.NewProductImageUsecase(imageRepository, productRepository)
	productHandler := productHttp.NewProductHandler(productSvc, imageUsecase)
	comboRepository := productRepo.NewProductComboRepository(database)
	comboUsecase := productUsecase.NewProductComboUsecase(comboRepository, productRepository, productRepository)
	comboHandler := productHttp.NewComboHandler(comboUsecase)

	uploadHandler := upload.NewUploadHandler()

	supplierRepository := supplierRepo.NewSupplierRepository(database)
	supplierTypeRepository := supplierRepo.NewSupplierTypeRepository(database)
	supplierContactRepository := supplierRepo.NewSupplierContactRepository(database)
	supplierAccountRepository := supplierRepo.NewSupplierAccountRepository(database)
	supplierTagRepository := supplierRepo.NewSupplierTagRepository(database)
	supplierSvc := supplierUsecase.NewSupplierUsecase(
		supplierRepository,
		supplierTypeRepository,
		supplierContactRepository,
		supplierAccountRepository,
		supplierTagRepository,
	)
	supplierHandler := supplierHttp.NewSupplierHandler(supplierSvc)
	quoteRepository := supplierRepo.NewQuoteRepository(database)
	quoteUsecase := supplierUsecase.NewQuoteUsecase(quoteRepository, productRepository, auditLogUsecase)
	quoteHandler := supplierHttp.NewQuoteHandler(quoteUsecase)

	purchaseOrderRepository := procurementRepo.NewPurchaseOrderRepository(database)
	movementRepository := procurementRepo.NewMovementRepository(database)
	purchaseOrderUsecase := procurementUsecase.NewPurchaseOrderUsecase(
		purchaseOrderRepository,
		productRepository,
		comboRepository,
		movementRepository,
		auditLogUsecase,
	)
	purchaseOrderHandler := procurementHttp.NewPurchaseOrderHandler(purchaseOrderUsecase)

	warehouseRepository := inventoryRepo.NewWarehouseRepository(database)
	warehouseUsecase := inventoryUsecase.NewWarehouseUsecase(warehouseRepository)
	warehouseHandler := inventoryHttp.NewWarehouseHandler(warehouseUsecase)

	r := router.NewRouter(
		jwtManager,
		authHandler,
		menuHandler,
		productHandler,
		comboHandler,
		supplierHandler,
		quoteHandler,
		purchaseOrderHandler,
		fieldLabelHandler,
		auditLogHandler,
		uploadHandler,
		warehouseHandler,
	)
	engine := r.Setup()

	return &App{
		Engine: engine,
		Config: cfg,
	}, nil
}
