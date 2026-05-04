# 计划：数据库元数据与后端工程化重构

## 目标

把当前 MVP 后端从单包脚本式结构升级为更正式的工程结构：

- 元数据从 `data/metadata.json` 改为数据库记录
- 数据本体继续保持本地文件存储，不把文件内容写入数据库
- 数据库记录元信息和文件目录信息
- 后端按 handler / service / repository / storage 分层
- 保持现有前端 API 行为尽量不变

## 技术选择

数据库使用 PostgreSQL。

理由：

- 项目已安装 PostgreSQL，不再使用 SQLite 运行时文件
- PostgreSQL 是正式服务型数据库，便于后续支持多用户、权限、查询、迁移和部署
- 数据本体仍在文件系统，PostgreSQL 只保存 metadata 和文件路径
- 本地开发通过 `DATABASE_URL` 连接 PostgreSQL，避免把数据库文件提交到仓库

新增依赖：

- `github.com/jackc/pgx/v5/stdlib`

说明：

- 使用 pgx 的 database/sql driver，保持 repository 实现简单
- 默认连接串建议为 `postgres://postgres@localhost:5432/lumino?sslmode=disable`

## 目标目录结构

计划重构为：

```text
cmd/lumino/main.go
internal/config/config.go
internal/httpserver/server.go
internal/httpserver/handlers/dataset_handler.go
internal/httpserver/middleware.go
internal/domain/dataset.go
internal/service/dataset_service.go
internal/repository/dataset_repository.go
internal/repository/postgres/dataset_repository.go
internal/storage/file_storage.go
internal/csvutil/profile.go
internal/web/
```

说明：

- `domain`: 业务模型和通用错误
- `service`: 上传、删除、预览、下载路径等业务编排
- `repository`: 元数据持久化接口
- `repository/postgres`: PostgreSQL 实现
- `storage`: 文件系统保存、删除、路径解析
- `httpserver`: 路由、handler、middleware、静态资源
- `csvutil`: CSV profile 逻辑
- `web`: 前端静态资源

## 数据库设计

数据库：

```text
PostgreSQL database: lumino
```

数据表：

```sql
CREATE TABLE IF NOT EXISTS datasets (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  original_file_name TEXT NOT NULL,
  stored_file_name TEXT NOT NULL,
  relative_path TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  rows INTEGER NOT NULL DEFAULT 0,
  columns_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  profile_json JSONB NOT NULL DEFAULT '{"numeric":[]}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_datasets_created_at ON datasets(created_at DESC);
```

字段说明：

- `original_file_name`: 用户上传时的原始文件名
- `stored_file_name`: 服务端保存后的文件名，例如 `<id>.csv`
- `relative_path`: 相对 `data/` 的路径，例如 `uploads/<id>.csv`
- `columns_json`: CSV 字段数组
- `profile_json`: CSV 数值统计 JSON
- `rows`: CSV 行数；非 CSV 文件为 0

## 兼容与迁移

当前已有 `data/metadata.json` 的情况，需要迁移。

计划启动时执行：

1. 创建数据库和数据表
2. 如果存在 `data/metadata.json`
3. 读取 JSON 元数据
4. 将记录 upsert 到 PostgreSQL
5. 不删除 `metadata.json`，避免误删用户数据；可以保留作为备份

迁移规则：

- 已有 `storedName` 优先作为 `stored_file_name`
- 没有 `storedName` 时兼容旧记录的 `fileName`
- `relative_path` 统一写成 `uploads/<stored_file_name>`

## API 变化

尽量保持现有 API 不变：

- `GET /api/health`
- `GET /api/datasets`
- `POST /api/datasets`
- `GET /api/datasets/{id}`
- `DELETE /api/datasets/{id}`
- `GET /api/datasets/{id}/download`

响应 JSON 字段保持兼容：

- `id`
- `name`
- `description`
- `fileName`
- `storedName`
- `contentType`
- `size`
- `rows`
- `columns`
- `profile`
- `createdAt`
- `updatedAt`
- `preview`

## 实现步骤

1. 新增 PostgreSQL 依赖
   - 更新 `go.mod`
   - 生成 `go.sum`

2. 建立新目录结构
   - 拆分 domain、config、httpserver、service、repository、storage、csvutil、web
   - 将当前 `internal/app/web` 移到 `internal/web`

3. 实现数据库初始化和迁移
   - 创建 `datasets` 表和索引
   - 从 `data/metadata.json` 迁移已有记录
   - 保留 JSON 文件作为备份

4. 实现 repository 接口
   - `List(ctx)`
   - `Get(ctx, id)`
   - `Create(ctx, dataset)`
   - `Delete(ctx, id)`

5. 实现 file storage
   - 保存上传文件
   - 删除文件
   - 返回绝对路径和相对路径
   - 防止路径穿越

6. 实现 service
   - 上传文件：保存文件、CSV profile、写入数据库
   - 获取详情：读数据库、按需生成 preview
   - 删除数据集：先查元数据，再删除文件，最后删除数据库记录
   - 下载：返回文件路径和下载名

7. 重写 HTTP handler
   - handler 只负责 HTTP 入参、出参和状态码
   - 业务逻辑下沉到 service

8. 清理旧 `internal/app`
   - 确保没有重复逻辑
   - 保持前端静态资源可嵌入

9. 更新 README
   - 增加 SQLite 元数据说明
   - 增加 `LUMINO_DB_PATH` 配置说明
   - 保留 `LUMINO_DATA_DIR`

## 配置变化

新增环境变量：

- `DATABASE_URL`: PostgreSQL 连接串，默认 `postgres://postgres@localhost:5432/lumino?sslmode=disable`

保留：

- `LUMINO_ADDR`: 服务监听地址，默认 `:8080`
- `LUMINO_DATA_DIR`: 数据目录，默认 `data`

## 验证方式

1. 编译测试：

```bash
GOCACHE=/Users/ekko/resp/Lumino/.cache/go-build go test ./...
```

2. 启动服务：

```bash
LUMINO_ADDR=127.0.0.1:18080 GOCACHE=/Users/ekko/resp/Lumino/.cache/go-build go run ./cmd/lumino
```

3. 接口冒烟：

- `GET /api/health`
- 上传 CSV
- `GET /api/datasets` 返回数据库记录
- `GET /api/datasets/{id}` 返回 preview
- `GET /api/datasets/{id}/download` 能下载本体文件
- `DELETE /api/datasets/{id}` 删除数据库记录和本体文件

4. 数据库验证：

```bash
psql "$DATABASE_URL" -c '\d datasets'
psql "$DATABASE_URL" -c 'select id,name,relative_path from datasets;'
```

## 风险与待确认

- 新增 PostgreSQL driver 需要下载依赖，可能需要网络权限。
- 本地需要 PostgreSQL 服务可用，并需要提前创建 `lumino` 数据库或提供有效 `DATABASE_URL`。
- 删除功能目前是硬删除；数据库化后仍按硬删除实现。
- 迁移会读取 `data/metadata.json`，但不会删除它，避免误删历史元数据。
- 这次会是较大后端重构，文件移动较多，但前端 API 兼容。

## 建议确认项

请确认：

1. 数据库采用 PostgreSQL。
2. 是否接受新增 `github.com/jackc/pgx/v5/stdlib` 依赖。
3. 是否保留 `data/metadata.json` 作为迁移备份，不自动删除。
