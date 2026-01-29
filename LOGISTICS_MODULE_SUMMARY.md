# 物流模块实现总结

## 已完成

### 1. 数据库迁移（migrations/）
- `007_create_logistics_provider.sql` - 物流供应商表
- `008_create_shipping_rate.sql` - 运费报价表
- `009_alter_shipment_add_logistics_fields.sql` - 发货单表增强

### 2. 物流供应商模块 (logistics/provider)
- ✅ Domain模型 (`domain/provider.go`)
- ✅ Repository接口和实现 (`domain/repository.go`, `repository/provider_repo.go`)
- ✅ Usecase业务逻辑 (`usecase/provider_usecase.go`)
- ✅ HTTP Handler (`delivery/http/provider_handler.go`)
- ✅ 路由注册 (`delivery/http/routes.go`)

**支持的操作：**
- 创建/更新/删除物流供应商
- 列表查询（支持类型、状态、关键词筛选）
- 详情查看

### 3. 运费报价模块 (logistics/shipping_rate)
- ✅ Domain模型 (`domain/shipping_rate.go`)
- ✅ Repository接口和实现 (`repository/shipping_rate_repo.go`)
- ✅ Usecase业务逻辑 (`usecase/shipping_rate_usecase.go`)
- ✅ HTTP Handler (`delivery/http/shipping_rate_handler.go`)
- ✅ 路由注册 (`delivery/http/routes.go`)

**支持的操作：**
- 创建/更新/删除运费报价
- 列表查询（支持供应商、仓库、运输方式等筛选）
- 查询最新有效报价（自动匹配重量/体积区间）

**报价支持的特性：**
- 多种计费方式（按重量/体积/件/固定）
- 燃油附加费
- 偏远地区附加费
- 最低收费
- 重量/体积区间匹配
- 有效期管理

### 4. 发货单模块增强 (shipment)
- ✅ 添加目的地仓库关联 (`destination_warehouse_id`)
- ✅ 添加物流供应商关联 (`logistics_provider_id`)
- ✅ 添加运费报价关联 (`shipping_rate_id`)
- ✅ 添加运输方式 (`transport_mode`)
- ✅ 添加详细时间节点（`confirmed_at`, `shipped_at`, `delivered_at`）
- ✅ 添加操作人记录（`confirmed_by`, `shipped_by`, `delivered_by`）
- ✅ 更新Confirm/MarkShipped/MarkDelivered方法记录时间和操作人

## 待完成

### 1. 后端集成
- [x] 在 `router.go` 中添加 logisticsProviderHandler 和 shippingRateHandler
- [x] 在 `bootstrap/app.go` 中初始化 logistics 模块
- [x] 添加 import 语句
- [x] 编译测试通过

### 2. 发货单创建优化
- [x] 创建发货单时支持选择目的地仓库（自动填充地址信息）
- [x] 创建发货单时支持选择物流供应商
- [x] 根据选择自动查询匹配的运费报价
- [x] 显示运费报价详情（费率、燃油附加费、时效等）

### 3. 前端界面

**物流供应商管理** (`/logistics/providers`)
- [x] 列表页（查询、筛选、新增、编辑、删除）
- [x] 新增/编辑表单
- [x] 路由配置

**运费报价管理** (`/logistics/shipping-rates`)
- [x] 列表页（查询、筛选、新增、编辑、删除）
- [x] 新增/编辑表单（选择供应商、起点仓库、目的地仓库）
- [x] 路由配置

**发货单创建优化**
- [x] 目的地仓库选择器（选择仓库自动填充收货方信息）
- [x] 物流供应商选择器
- [x] 运输方式选择器
- [x] 自动查询报价并显示预估费用
- [x] 显示预计时效
- [x] 更新类型定义（添加物流相关字段）

### 4. 操作日志集成
- [ ] 在发货单操作时调用 `system/audit-logs` 接口记录日志
- [ ] 记录内容：操作人、操作动作、状态变更、时间

## 数据库字段说明

### logistics_provider 表
- `provider_type`: 供应商类型（FREIGHT_FORWARDER货代/COURIER快递/SHIPPING_LINE船公司/AIRLINE航空）
- `service_types`: 服务类型（EXPRESS,AIR,SEA,RAIL，逗号分隔）
- `account_number`: 客户账号
- `credit_days`: 账期天数

### shipping_rate 表
- `transport_mode`: 运输方式（EXPRESS快递/AIR空运/SEA海运/RAIL铁路/TRUCK卡车）
- `pricing_method`: 计费方式（PER_KG按公斤/PER_CBM按立方/PER_PACKAGE按件/FIXED固定）
- `base_rate`: 基础费率
- `fuel_surcharge_rate`: 燃油附加费率(%)
- `min_weight/max_weight`: 重量区间（用于阶梯报价）
- `min_volume/max_volume`: 体积区间
- `effective_date/expiry_date`: 有效期

### shipment 表新增字段
- `destination_warehouse_id`: 目的地仓库ID
- `logistics_provider_id`: 物流供应商ID
- `shipping_rate_id`: 运费报价ID
- `transport_mode`: 运输方式
- `confirmed_at/shipped_at/delivered_at`: 详细时间节点
- `confirmed_by/shipped_by/delivered_by`: 操作人记录

## API 接口列表

### 物流供应商
- `GET /api/v1/logistics-providers` - 列表
- `GET /api/v1/logistics-providers/:id` - 详情
- `POST /api/v1/logistics-providers` - 创建
- `PUT /api/v1/logistics-providers/:id` - 更新
- `DELETE /api/v1/logistics-providers/:id` - 删除

### 运费报价
- `GET /api/v1/shipping-rates` - 列表
- `GET /api/v1/shipping-rates/:id` - 详情
- `POST /api/v1/shipping-rates` - 创建
- `PUT /api/v1/shipping-rates/:id` - 更新
- `DELETE /api/v1/shipping-rates/:id` - 删除
- `GET /api/v1/shipping-rates/query-latest` - 查询最新报价

## 前端文件清单

### 物流模块 (am-erp-vue/src/modules/logistics/)
- `types/index.ts` - 类型定义（供应商、运费报价、枚举等）
- `api/index.ts` - API 接口（供应商和运费报价的CRUD）
- `views/ProviderList.vue` - 物流供应商管理页面
- `views/ShippingRateList.vue` - 运费报价管理页面

### 发货单模块更新
- `am-erp-vue/src/modules/shipping/types/index.ts` - 添加物流相关字段到 Shipment 和 CreateShipmentParams
- `am-erp-vue/src/modules/shipping/views/ShipmentCreate.vue` - 发货单创建页面优化：
  - 添加目的地仓库选择器
  - 添加物流供应商选择器
  - 添加运输方式选择器
  - 自动查询并显示运费报价信息
  - 选择仓库后自动填充收货方信息

### 路由配置
- `am-erp-vue/src/router/index.ts` - 添加物流模块路由：
  - `/logistics/providers` - 物流供应商管理
  - `/logistics/shipping-rates` - 运费报价管理

## 使用流程

1. **维护基础数据**
   - 添加物流供应商
   - 添加目的地仓库（使用现有warehouse表）
   - 维护运费报价

2. **创建发货单**
   - 选择起点仓库
   - 选择目的地仓库（自动填充地址）
   - 选择物流供应商
   - 选择运输方式
   - 系统自动查询匹配的最新报价
   - 显示预估运费和时效

3. **发货单流程**
   - DRAFT（草稿）→ CONFIRMED（确认，记录确认时间和人）
   - CONFIRMED → SHIPPED（发货，记录发货时间和人）
   - SHIPPED → DELIVERED（送达，记录送达时间和人）

4. **查看日志**
   - 通过 audit_log 表查看所有操作记录
