# Minimal Seed

这个目录只放“系统能启动并登录”的最小种子，不放业务数据。

包含：

- `permission`
- `role`
- `role_permission`
- `menu`

不包含：

- `field_label`
- `user`
- `user_role`
- 所有业务表
- 所有运行日志表

说明：

- 默认管理员用户由 `scripts/init_minimal_seed.ps1` / `cmd/initseed/main.go` 在初始化时动态创建。
- 密码不是固定值，会随机生成并输出到凭据文件。
- `minimal_seed.sql` 应通过 `scripts/export_minimal_seed.ps1` 重新生成，不要手工改表后忘记同步种子。
- 这个最小种子应与 `structure.sql` 搭配使用，不应用来覆盖正在运行的生产库。
