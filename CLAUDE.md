# 积分管理应用 - 项目计划

## 概述
管理两个小朋友的积分，支持预定义加减分行为和转盘抽奖。

## 技术栈
- **前端**：Vue 3 + TypeScript + Vite + TailwindCSS v3 + Vue Reactive
- **后端**：Go + Gin + SQLite
- **架构**：前后端分离，REST API 通信

## 数据模型
- **Kid**：id, name, avatar, points
- **Behavior**：id, name, type(add/sub), points
- **Prize**：id, name, icon, probability (0-100%)
- **Record**：id, kid_id, behavior_id, behavior_name, points, created_at, reversed, reversed_reason, reversed_at

## API 设计
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/kids` | GET | 获取小朋友列表 |
| `/api/kids` | POST | 创建小朋友 |
| `/api/kids/:id` | PUT | 更新小朋友 |
| `/api/kids/:id` | DELETE | 删除小朋友 |
| `/api/kids/:id/points` | GET | 获取积分（支持 id 或 name） |
| `/api/kids/points/add` | POST | 加分（body: id 或 name + behavior_id） |
| `/api/kids/points/sub` | POST | 减分（body: id 或 name + behavior_id） |
| `/api/behaviors` | GET | 获取预定义行为 |
| `/api/behaviors` | POST | 创建行为 |
| `/api/behaviors/:id` | PUT | 更新行为 |
| `/api/behaviors/:id` | DELETE | 删除行为 |
| `/api/prizes` | GET | 获取奖项 |
| `/api/prizes` | POST | 创建奖项 |
| `/api/prizes/:id` | PUT | 更新奖项 |
| `/api/prizes/:id` | DELETE | 删除奖项 |
| `/api/draw` | POST | 抽奖（body: kid_id） |
| `/api/records` | GET | 行为记录历史 |
| `/api/records/:id/reverse` | POST | 撤销记录并回滚积分（body: reason） |

## 开发阶段

### Phase 1 - 后端基础（Go + Gin + SQLite）✅
- [x] 1.1 项目骨架 + go mod init + config
- [x] 1.2 数据模型定义（models/）
- [x] 1.3 数据库初始化 + Repository 层
- [x] 1.4 Service 层（含 unit test）
- [x] 1.5 Handler 层 + 路由（含 unit test）
- [x] 1.6 main.go 组装 + 联调
- [x] 1.7 Kids CRUD（POST/PUT/DELETE）+ unit test
- [x] 1.8 CORS 中间件配置

### Phase 2 - 前端初始化 ✅
- [x] 2.1 Vite + Vue 3 + TypeScript 项目创建
- [x] 2.2 TailwindCSS v3 + UI 组件配置
- [x] 2.3 基础布局 + 小朋友切换组件
- [x] 2.4 Vue Reactive 状态管理 + API 服务层

### Phase 3 - 配置管理 ✅
- [x] 3.1 行为增删改查
- [x] 3.2 奖项增删改查

### Phase 4 - 积分加减与展示 ✅
- [x] 4.1 积分展示页面
- [x] 4.2 预定义行为按钮（加分/减分）
- [x] 4.3 积分变更 API 对接

### Phase 5 - 行为记录历史 ✅
- [x] 5.1 记录列表展示
- [x] 5.2 按小朋友筛选
- [x] 5.3 时间倒序排列

### Phase 6 - 转盘界面 + 抽奖 API ✅
- [x] 6.1 Canvas 转盘 UI 设计
- [x] 6.2 抽奖 API 对接
- [x] 6.3 抽奖结果展示

### Phase 7 - 转盘动画 + 中奖判定 ✅
- [x] 7.1 Canvas 转盘动画
- [x] 7.2 概率算法实现
- [x] 7.3 中奖动效

## 后端目录结构
```
backend/
├── main.go
├── go.mod
├── go.sum
├── config/
│   └── config.go
├── models/
│   ├── kid.go
│   ├── behavior.go
│   ├── prize.go
│   └── record.go
├── database/
│   └── init.go
├── repository/
│   ├── kid.go
│   ├── behavior.go
│   ├── prize.go
│   └── record.go
├── service/
│   ├── kid.go
│   ├── kid_test.go
│   ├── draw.go
│   ├── draw_test.go
│   ├── record.go
│   └── record_test.go
├── handlers/
│   ├── kid.go
│   ├── kid_test.go
│   ├── behavior.go
│   ├── prize.go
│   ├── draw.go
│   ├── draw_test.go
│   └── record.go
└── router/
    └── router.go
```

## 前端目录结构
```
frontend/
├── src/
│   ├── api/
│   │   └── index.ts          # REST API 服务层
│   ├── components/
│   │   ├── ui/
│   │   │   └── EmojiPicker.vue  # Emoji 选择器
│   │   ├── KidSwitcher.vue      # 小朋友切换
│   │   ├── PointsDisplay.vue    # 积分展示
│   │   ├── BehaviorButtons.vue  # 加减分按钮
│   │   ├── BehaviorManager.vue  # 行为管理
│   │   ├── PrizeManager.vue     # 奖项管理
│   │   ├── RecordList.vue       # 记录列表
│   │   └── LuckyWheel.vue       # 抽奖转盘（Canvas）
│   ├── stores/
│   │   └── app.ts             # Vue Reactive 状态管理
│   ├── types/
│   │   └── index.ts           # TypeScript 类型
│   ├── App.vue
│   ├── main.ts
│   └── style.css               # TailwindCSS 主题
├── index.html
├── package.json
├── tailwind.config.js
├── postcss.config.js
└── vite.config.ts
```

## 已知问题/注意事项
- 后端已配置 CORS，允许 localhost:5173/5174/5175
- 前后端分离部署时需各自启动（后端 :8080，前端 :5173+）
