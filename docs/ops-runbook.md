# ERP 运维最小手册

## 数据库迁移

在 `am-erp-go` 目录执行：

```powershell
go run .\cmd\migrate\main.go
```

规则：

- 只执行 `migrations/` 下文件名以数字开头的 `.sql`
- 执行记录写入 `schema_migration`
- 已执行过的版本不会重复执行
- `baseline/legacy/am-erp-legacy.sql` 这类全量导出文件不会被当成迁移执行

注意：

- 当前 `baseline/legacy/am-erp-legacy.sql` 是 legacy 全量导出，不应继续作为新环境初始化基线
- 新环境初始化请使用“baseline bundle”流程，不要假设历史增量迁移文件完整存在

如果数据库是历史库，已经有业务表但还没有 `schema_migration` 记录，先执行基线化：

```powershell
go run .\cmd\migrate\main.go -baseline
```

说明：

- `-baseline` 只记录历史迁移版本，不执行 SQL
- 基线化完成后，后续再执行普通迁移命令

## 数据库备份

在 `am-erp-go` 目录执行：

```powershell
.\scripts\backup_db.ps1
```

说明：

- 默认读取当前目录 `.env`
- 默认输出到 `backups\db`
- 文件名格式：`数据库名-yyyyMMdd-HHmmss.sql`

可选参数：

```powershell
.\scripts\backup_db.ps1 -EnvFile .env -OutputDir backups\db
```

## 数据库恢复

恢复是破坏性操作，必须显式带 `-Force`：

```powershell
.\scripts\restore_db.ps1 -BackupFile .\backups\db\am-erp-20260308-140408.sql -Force
```

说明：

- 默认读取当前目录 `.env`
- `BackupFile` 可以传相对路径或绝对路径
- 不带 `-Force` 会直接拒绝执行

建议流程：

1. 先执行一次备份
2. 确认目标库正确
3. 再执行恢复

## Baseline Bundle（推荐的新环境初始化方式）

导出当前可恢复 baseline：

```powershell
.\scripts\export_baseline_bundle.ps1
```

输出内容：

- `structure.sql`
- `minimal_seed.sql`
- `schema_migration_versions.txt`

恢复 baseline bundle 到目标库：

```powershell
.\scripts\restore_baseline_bundle.ps1 -BundleDir .\backups\baseline\20260308-145500 -Force
```

说明：

- 先恢复 `structure.sql`
- 再用 `schema_migration_versions.txt` 重建 `schema_migration`
- 再执行 `minimal_seed.sql`，并动态创建随机密码的默认管理员
- 最后自动执行当前 `migrations/` 目录下仍存在的增量迁移

恢复完成后会在 bundle 目录下生成：

- `admin_credentials.txt`

这个文件只用于首次登录，拿到密码后应立即改密。

## 最小种子初始化

如果你已经有一套纯结构库，只想补系统启动必需数据：

```powershell
.\scripts\init_minimal_seed.ps1
```

如果菜单、权限或角色有调整，先重新导出最小种子：

```powershell
.\scripts\export_minimal_seed.ps1
```

说明：

- 默认读取 `baseline\minimal\minimal_seed.sql`
- 动态创建 `admin` 用户
- 随机生成密码并输出到 `admin_credentials.txt`
- 不会导入 `field_label`
- 不会导入任何业务数据或运行日志

## 发布前最低检查

```powershell
.\scripts\preflight_check.ps1
```

## 任务日志保留策略

系统会自动清理以下运行记录：

- `job_run`
- `system_log`

默认策略：

- `JOB_RUN_RETENTION_DAYS=30`
- `SYSTEM_LOG_RETENTION_DAYS=30`
- `LOG_RETENTION_CLEANUP_INTERVAL_MINUTES=1440`
- `LOG_RETENTION_ENABLED=true`

说明：

- 这是运维级配置，不走前端系统设置页
- 当前“系统设置”前端页还没有完整后端配置中心支撑，不适合承载日志保留策略
- 修改这些值后，需要重启后端进程

建议：

- 开发环境可缩短到 `7`
- 试运行/生产环境建议至少保留 `30`
