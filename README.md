# Lumino 数据管理平台

一个从零开始的 Go 数据管理平台 MVP，支持数据上传、数据列表、CSV 预览统计和文件下载。

当前后端使用 SQLite 保存数据集元数据，原始文件继续保存在本地文件系统。服务启动时会创建数据库表，并将已有 `data/metadata.json` 迁移到 SQLite；迁移后不会删除原 JSON 文件。

## 运行

```bash
go run ./cmd/lumino
```

默认监听 `:8080`，浏览器访问 `http://localhost:8080`。

## 环境变量

- `LUMINO_ADDR`: 服务监听地址，默认 `:8080`
- `LUMINO_DATA_DIR`: 数据存储目录，默认 `data`
- `LUMINO_DB_PATH`: SQLite 数据库路径，默认 `data/lumino.db`

## API

- `GET /api/health`: 健康检查
- `GET /api/datasets`: 数据列表
- `POST /api/datasets`: 上传数据，multipart 字段为 `file`、`name`、`description`
- `GET /api/datasets/{id}`: 数据详情和 CSV 预览
- `GET /api/datasets/{id}/download`: 下载原始文件
- `DELETE /api/datasets/{id}`: 删除数据集元数据和本地文件
