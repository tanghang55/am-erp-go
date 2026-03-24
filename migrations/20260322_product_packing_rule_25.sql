ALTER TABLE product
    ADD COLUMN is_packing_required TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否需要打包：1需要，0免打包直通' AFTER is_inspection_required;

UPDATE product
SET is_packing_required = 1
WHERE is_packing_required IS NULL;

ALTER TABLE inventory_movement
    MODIFY COLUMN movement_type ENUM(
        'PURCHASE_RECEIPT',
        'SALES_SHIPMENT',
        'SALES_ALLOCATE',
        'SALES_RELEASE',
        'SALES_SHIP',
        'STOCK_TAKE_ADJUSTMENT',
        'MANUAL_ADJUSTMENT',
        'DAMAGE_WRITE_OFF',
        'RETURN_RECEIPT',
        'TRANSFER_OUT',
        'TRANSFER_IN',
        'PURCHASE_SHIP',
        'WAREHOUSE_RECEIVE',
        'INSPECTION_PASS',
        'INSPECTION_FAIL',
        'INSPECTION_LOSS',
        'ASSEMBLY_CONSUME',
        'ASSEMBLY_COMPLETE',
        'PACKING_SKIP_COMPLETE',
        'SHIPMENT_ALLOCATE',
        'SHIPMENT_RELEASE',
        'LOGISTICS_SHIP',
        'PLATFORM_RECEIVE',
        'RETURN_INSPECT',
        'SHIPMENT_SHIP'
    ) NOT NULL COMMENT '库存流水类型';
