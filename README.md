# 积分管理应用

管理两个小朋友的积分，支持预定义加减分行为和转盘抽奖。

## 技术栈

- **前端**：Vue 3 + TypeScript + Vite + TailwindCSS v3
- **后端**：Go + Gin + SQLite
- **架构**：前后端分离，REST API 通信

## 快速开始

### 后端

```bash
cd backend
go run main.go
```

后端服务运行在 http://localhost:8080

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端服务运行在 http://localhost:5173

## 数据模型

- **Kid**：小朋友（id, name, avatar, points）
- **Behavior**：预定义行为（id, name, type(add/sub), points）
- **Prize**：奖项（id, name, icon, probability）
- **Record**：积分记录（id, kid_id, behavior_id, behavior_name, points, created_at, reversed, reversed_reason, reversed_at）

## API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/kids` | GET | 获取小朋友列表 |
| `/api/kids` | POST | 创建小朋友 |
| `/api/kids/:id` | PUT | 更新小朋友 |
| `/api/kids/:id` | DELETE | 删除小朋友 |
| `/api/kids/:id/points` | GET | 获取积分 |
| `/api/kids/points/add` | POST | 加分 |
| `/api/kids/points/sub` | POST | 减分 |
| `/api/behaviors` | GET | 获取预定义行为 |
| `/api/behaviors` | POST | 创建行为 |
| `/api/behaviors/:id` | PUT | 更新行为 |
| `/api/behaviors/:id` | DELETE | 删除行为 |
| `/api/prizes` | GET | 获取奖项 |
| `/api/prizes` | POST | 创建奖项 |
| `/api/prizes/:id` | PUT | 更新奖项 |
| `/api/prizes/:id` | DELETE | 删除奖项 |
| `/api/draw` | POST | 抽奖 |
| `/api/records` | GET | 行为记录历史 |
| `/api/records/:id/reverse` | POST | 撤销记录 |

## 功能

- 小朋友管理（头像、昵称）
- 积分加减分（通过预定义行为）
- 行为记录历史
- 转盘抽奖（Canvas 动画 + 概率算法）
- 记录撤销与积分回滚
