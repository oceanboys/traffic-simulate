# 智慧交通实时数据平台

这是一个基于Vue 3 + Go + MySQL的智慧交通实时数据平台，提供实时交通监控、路段管理、GPS数据管理和告警管理功能。

## 技术栈

### 后端
- Go 1.21+
- Beego 2.x 框架
- MySQL 数据库
- ORM 数据访问层

### 前端
- Vue 3 + TypeScript
- Element Plus UI组件库
- Vue Router 路由管理
- Axios HTTP客户端

## 功能特性

### 实时监控
- 在线车辆统计
- 超速告警统计
- 拥堵路段统计
- 平均速度统计
- 实时地图显示（待集成）
- 实时告警列表
- 路段状态监控

### 路段管理
- 路段信息的增删改查
- 路段类型分类（高速公路、城市道路、乡村道路）
- 路段坐标信息管理
- 限速和容量设置

### GPS数据管理
- GPS数据的录入和查询
- 车辆轨迹追踪
- 速度监控
- 车型分类管理
- 时间范围筛选

### 告警管理
- 超速告警
- 拥堵告警
- 事故告警
- 告警严重程度分级
- 告警处理状态管理

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 数据库配置
1. 创建MySQL数据库：
`sql
CREATE DATABASE traffic_insights CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
`

2. 修改后端配置文件 ackend/conf/app.conf：
`ini
# mysql数据库信息
db.driver = mysql
db.host = localhost
db.port = 3306
db.user = root
db.password = "your_password"
db.name = traffic_insights
db.prefix =
`

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

`
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
`

## 开发说明

### 后端开发
- 使用Beego框架的MVC架构
- 数据库操作使用ORM
- 支持自动建表（开发环境）
- 支持CORS跨域请求

### 前端开发
- 使用Vue 3 Composition API
- TypeScript类型安全
- Element Plus组件库
- 响应式设计

## 部署说明

### 生产环境配置
1. 修改 ackend/conf/app.conf 中的数据库配置
2. 设置 unmode = prod
3. 编译后端：go build -o backend
4. 构建前端：
pm run build
5. 部署到服务器

## 许可证

MIT License
