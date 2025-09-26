# 智慧交通实时数据平台

这是一个基于Vue 3 + Go + MySQL的智慧交通实时数据平台，提供实时交通监控、路段管理、GPS数据管理和告警管理功能。

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 数据库配置
1. 创建MySQL数据库：

2. 修改后端配置文件 backend/conf/app.conf：

### 启动后端服务
`ash
cd backend
go mod tidy
go build -o backend.exe
./backend.exe
`

后端服务将在 http://localhost:8080 启动

### 启动前端服务
`ash
cd frontend
npm install
npm run dev
`

前端服务将在 http://localhost:5173 启动

## API接口

### 健康检查
- GET /api/health - 服务健康状态

### 路段管理
- GET /api/roads - 获取所有路段
- POST /api/roads - 创建路段
- GET /api/roads/:id - 获取指定路段
- PUT /api/roads/:id - 更新路段
- DELETE /api/roads/:id - 删除路段

### GPS数据
- POST /api/gps - 创建GPS数据
- GET /api/gps/road/:roadId - 获取指定路段的GPS数据
- GET /api/gps/vehicle/:vehicleId - 获取指定车辆的GPS数据

## 项目结构


traffic-insights/
 backend/                 # 后端Go代码
    algorithms/         # 算法模块
    controllers/        # 控制器
    models/            # 数据模型
    repositories/      # 数据访问层
    services/          # 业务逻辑层
    routers/           # 路由配置
    utils/             # 工具函数
    conf/              # 配置文件
 frontend/               # 前端Vue代码
    src/
       api/           # API接口
       components/    # 组件
       views/         # 页面
       router/        # 路由
       types/         # 类型定义
    public/            # 静态资源
 README.md


## 开发说明

### 后端开发
- 使用Beego框架架构
- 数据库操作使用ORM
- 支持自动建表（开发环境）
- 支持CORS跨域请求

### 前端开发
- 使用Vue 3 Composition API
- TypeScript类型安全
- Element Plus组件库
- 响应式设计



