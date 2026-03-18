---
name: remove-sqlite-use-memory-storage
overview: 移除 SQLite 依赖，改用内存 map 存储数据源，并增加前端列表和详情查看功能
design:
  architecture:
    framework: html
  styleKeywords:
    - Minimalism
    - Glassmorphism
    - Gradient
    - Card Layout
  fontSystem:
    fontFamily: Inter
    heading:
      size: 24px
      weight: 600
    subheading:
      size: 16px
      weight: 500
    body:
      size: 14px
      weight: 400
  colorSystem:
    primary:
      - "#667EEA"
      - "#764BA2"
    background:
      - "#F9FAFB"
      - "#FFFFFF"
    text:
      - "#1F2937"
      - "#6B7280"
    functional:
      - "#10B981"
      - "#EF4444"
      - "#3B82F6"
todos:
  - id: modify-models
    content: 修改 models/datasource.go，移除 SQLite，改用内存 Map 存储
    status: completed
  - id: modify-services
    content: 修改 services/datasource_service.go，移除 SQLite 驱动导入
    status: completed
  - id: modify-handlers
    content: 修改 handlers/datasource_handler.go，新增列表和详情接口
    status: completed
    dependencies:
      - modify-models
  - id: modify-main
    content: 修改 main.go，移除 InitDB 调用，新增 API 路由
    status: completed
    dependencies:
      - modify-handlers
  - id: modify-frontend
    content: 修改 static/index.html，新增数据源列表和详情查看功能
    status: completed
    dependencies:
      - modify-handlers
  - id: build-verify
    content: 清理依赖并重新编译验证
    status: completed
    dependencies:
      - modify-main
      - modify-frontend
---

## 产品概述

将现有的数据源管理系统从 SQLite 存储改为内存 Map 存储，并增加前端数据源列表展示和详情查看功能。

## 核心功能

1. **移除 SQLite 依赖**：删除所有 SQLite 相关代码和驱动导入
2. **内存存储**：使用并发安全的 Map 存储数据源，以 ID 作为 key
3. **新增 API 接口**：

- `GET /api/datasource/list` - 获取数据源列表
- `GET /api/datasource/:id` - 获取单个数据源详情

4. **前端列表展示**：新增数据源列表区域，显示已保存的数据源
5. **详情查看**：点击列表项可查看数据源详细信息

## 技术栈

- **后端框架**: Gin
- **存储方式**: 内存 Map（使用 `sync.RWMutex` 保证并发安全）
- **前端**: 纯 HTML + JavaScript + Tailwind CSS

## 实现方案

### 存储层改造

将 `models/datasource.go` 中的 SQLite 存储改为内存 Map：

- 使用 `map[int64]*Datasource` 存储数据
- 使用 `sync.RWMutex` 保证并发读写安全
- 使用原子计数器生成自增 ID
- 移除 `InitDB()` 函数和 SQLite 驱动导入

### 新增 API

在 `handlers/datasource_handler.go` 中新增：

- `GetDatasourceList` - 返回所有数据源列表
- `GetDatasourceByID` - 根据 ID 返回单个数据源详情

### 前端改造

在 `static/index.html` 中新增：

- 数据源列表展示区域（卡片形式）
- 点击卡片弹出详情模态框
- 列表实时刷新（保存成功后自动刷新）

## 目录结构

```
script_server/
├── main.go                      # [MODIFY] 移除 InitDB 调用，新增路由
├── models/
│   └── datasource.go            # [MODIFY] 改为内存 Map 存储
├── services/
│   └── datasource_service.go    # [MODIFY] 移除 SQLite 驱动导入
├── handlers/
│   └── datasource_handler.go    # [MODIFY] 新增列表和详情接口
└── static/
    └── index.html               # [MODIFY] 新增列表展示和详情查看
```

## 设计风格

沿用现有的现代简约风格，紫色渐变主题，玻璃态效果卡片。

## 页面布局

采用左右分栏布局：

- **左侧**：数据源配置表单（现有功能）
- **右侧**：已保存的数据源列表（新增）

## 新增 UI 组件

1. **数据源列表卡片**：显示数据源名称、类型、主机端口，hover 效果
2. **详情模态框**：点击卡片弹出，显示完整数据源信息
3. **类型图标**：不同数据库类型显示不同图标和颜色标识