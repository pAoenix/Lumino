# 计划：数据管理平台前端 UI/UX 优化

## 目标

把当前数据管理平台从“能用的基础界面”优化成更像真实数据管理产品的工作台：

- 信息层级更清晰
- 上传、列表、详情三个核心流程更顺
- 视觉风格更克制、专业、适合数据管理场景
- 移动端和窄屏下不拥挤、不重叠

## 设计系统

已安装并使用 `ui-ux-pro-max` 生成本项目设计系统：

- Skill 安装位置：`.codex/skills/ui-ux-pro-max`
- 设计系统文件：`design-system/lumino-data-platform/MASTER.md`

本次采用 `ui-ux-pro-max` 推荐方向：

- Pattern：Enterprise Gateway 中适合后台产品的可信、清晰结构，但不做营销页
- Style：Data-Dense Dashboard
- 主色：`#1E40AF`
- 辅色：`#3B82F6`
- 强调色：`#F59E0B`
- 背景：`#F8FAFC`
- 文本：`#1E3A8A`
- 字体方向：Fira Sans / Fira Code 气质，但本次不强制加载外部 Google Fonts，优先使用系统字体，避免网络依赖

设计约束：

- 数据表格在移动端使用横向滚动容器，不撑破页面
- 点击元素必须有 hover 和 focus 状态
- 表格行需要 hover 高亮
- 动效控制在 150-300ms，并尊重 `prefers-reduced-motion`
- 亮色模式文本对比度至少满足常规可读性
- 交付前检查 375px、768px、1024px、1440px 主要断点

## 范围

本次只优化前端静态资源：

- `internal/app/web/index.html`
- `internal/app/web/styles.css`
- `internal/app/web/app.js`

后端 API 暂不改动。

## 不做什么

- 不引入前端框架
- 不引入构建工具
- 不接入图表库
- 不做用户登录、权限、删除、编辑等新业务功能
- 不调整后端数据结构和接口协议

## 实现步骤

1. 重构页面信息架构
   - 顶部增加更像产品工作台的状态区
   - 保留上传、数据列表、数据详情三个主要区域
   - 优化空状态、加载状态和错误状态

2. 优化上传体验
   - 文件选择区域更明显
   - 上传按钮状态更清楚
   - 显示已选择文件名
   - 保留名称、描述字段

3. 优化数据列表
   - 数据项改成更紧凑的列表行
   - 强化名称、行列数、大小、时间等关键元信息
   - 当前选中态更明确
   - 空列表时提供明确但不啰嗦的提示

4. 优化详情与可视化
   - 详情顶部展示数据集概览指标
   - 数值字段统计卡片更适合扫描
   - CSV 预览表格增强表头、滚动和密度
   - 下载入口保持清晰

5. 优化视觉系统
   - 使用更稳的中性色背景
   - 减少大面积单色调
   - 控制卡片圆角在 8px 内
   - 避免装饰性渐变和无意义视觉元素
   - 保证按钮、输入框、列表项尺寸稳定

6. 响应式适配
   - 桌面端使用左侧上传、右侧列表、下方详情
   - 窄屏端单列布局
   - 表格保持横向滚动，不挤压文本

## 数据与 API 变化

无。

继续使用现有接口：

- `GET /api/datasets`
- `POST /api/datasets`
- `GET /api/datasets/{id}`
- `GET /api/datasets/{id}/download`

## 验证方式

1. 运行格式/编译检查：

```bash
GOCACHE=/Users/ekko/resp/Lumino/.cache/go-build go test ./...
```

2. 启动本地服务：

```bash
LUMINO_ADDR=127.0.0.1:18080 GOCACHE=/Users/ekko/resp/Lumino/.cache/go-build go run ./cmd/lumino
```

3. 浏览器检查：
   - 首页可以正常打开
   - 数据列表空状态正常
   - CSV 上传成功
   - 上传后列表自动刷新
   - 点击数据项后详情展示正常
   - 下载按钮可用
   - 窄屏下页面不重叠、不溢出

## 风险与待确认

- 当前没有 `ui-ux-pro-max` skill 可安装；如果你有这个 skill 的 GitHub 地址或本地路径，可以再提供，我再安装到 `$CODEX_HOME/skills`。
- 本次按现有项目规则和 Codex 前端设计规范执行，不额外引入 UI 库。
- 如果你希望视觉风格偏“企业后台”、“现代 SaaS”、“深色数据大屏”或“轻量文件管理器”，需要先确认方向。
