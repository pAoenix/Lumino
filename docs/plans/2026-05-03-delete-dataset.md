# 计划：数据集删除功能

## 目标

为数据管理平台增加删除数据集能力：

- 用户可以从前端删除指定数据集
- 后端同时删除元数据和本地上传文件
- 删除后数据列表、概览统计和详情区域自动刷新

## 范围

涉及文件：

- `internal/app/store.go`
- `internal/app/server.go`
- `internal/app/web/index.html`
- `internal/app/web/styles.css`
- `internal/app/web/app.js`

## 不做什么

- 不做批量删除
- 不做回收站
- 不做权限校验
- 不做软删除
- 不改变现有上传、列表、详情、下载接口行为

## 实现步骤

1. 后端存储层增加删除方法
   - 根据 dataset id 查找记录
   - 删除 `data/uploads/` 下对应存储文件
   - 从内存 map 删除元数据
   - 持久化更新 `data/metadata.json`

2. 后端 API 增加删除接口
   - 新增 `DELETE /api/datasets/{id}`
   - 数据不存在返回 404
   - 删除成功返回 204
   - 文件已不存在时允许继续删除元数据，避免坏数据卡死

3. 前端列表增加删除入口
   - 每个数据项增加删除按钮
   - 删除前使用浏览器确认弹窗
   - 删除时禁用按钮，避免重复提交
   - 删除成功后刷新列表和概览

4. 详情区域状态处理
   - 如果删除的是当前选中的数据集，自动切换到最新一条数据
   - 如果删除后没有数据集，隐藏详情区域

5. 样式优化
   - 删除按钮使用低干扰危险态
   - hover/focus 状态清晰
   - 移动端不挤压数据项内容

## 数据与 API 变化

新增 API：

```http
DELETE /api/datasets/{id}
```

成功：

```http
204 No Content
```

失败：

```json
{"error":"dataset not found"}
```

## 验证方式

1. 编译检查：

```bash
GOCACHE=/Users/ekko/resp/Lumino/.cache/go-build go test ./...
```

2. 接口冒烟：

- 上传测试 CSV
- 调用 `DELETE /api/datasets/{id}`
- 确认 `GET /api/datasets` 不再返回该记录
- 确认下载接口对该 id 返回 404

3. 前端检查：

- 删除按钮可见
- 点击删除会确认
- 删除后列表和统计刷新
- 删除当前详情项后详情区域状态正确

## 风险与待确认

- 删除是不可恢复操作，本计划先使用浏览器确认弹窗作为保护。
- 如果你希望更稳妥，可以改成二次确认输入数据集名称，或先实现软删除/回收站。
