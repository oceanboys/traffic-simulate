# 智慧交通实时数据平台 - 数据库设计文档

## 概述
本文档描述了智慧交通实时数据平台的数据库设计，包括表结构、索引、外键关系等。

## 数据库表设计

### 1. road_segments (路段表)
存储道路路段信息，包括起点、终点坐标、限速、容量等。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | INT UNSIGNED | 主键 | AUTO_INCREMENT |
| name | VARCHAR(100) | 路段名称 | NOT NULL |
| start_lng | DECIMAL(10,6) | 起点经度 | NOT NULL |
| start_lat | DECIMAL(10,6) | 起点纬度 | NOT NULL |
| end_lng | DECIMAL(10,6) | 终点经度 | NOT NULL |
| end_lat | DECIMAL(10,6) | 终点纬度 | NOT NULL |
| max_speed | INT | 限速值(km/h) | DEFAULT 60 |
| capacity | INT | 道路容量 | DEFAULT 1000 |
| length | DECIMAL(8,2) | 路段长度(km) | NULL |
| road_type | VARCHAR(10) | 道路类型 | DEFAULT urban |
| created_at | DATETIME | 创建时间 | DEFAULT CURRENT_TIMESTAMP |
| updated_at | DATETIME | 更新时间 | ON UPDATE CURRENT_TIMESTAMP |

**索引:**
- PRIMARY KEY (id)
- INDEX idx_name (name)
- INDEX idx_road_type (road_type)

### 2. gps_data (GPS数据表)
存储车辆GPS定位数据，包括位置、速度、方向等信息。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | INT UNSIGNED | 主键 | AUTO_INCREMENT |
| vehicle_id | VARCHAR(20) | 车辆ID | NOT NULL |
| longitude | DECIMAL(10,6) | 经度 | NOT NULL |
| latitude | DECIMAL(10,6) | 纬度 | NOT NULL |
| speed | INT | 速度(km/h) | DEFAULT 0 |
| direction | INT | 方向(0-360度) | NULL |
| timestamp | DATETIME | GPS时间戳 | NOT NULL |
| road_segment_id | INT UNSIGNED | 关联路段ID | NULL |
| vehicle_type | VARCHAR(10) | 车辆类型 | DEFAULT car |
| created_at | DATETIME | 创建时间 | DEFAULT CURRENT_TIMESTAMP |

**索引:**
- PRIMARY KEY (id)
- INDEX idx_vehicle_id (vehicle_id)
- INDEX idx_timestamp (timestamp)
- INDEX idx_road_segment (road_segment_id)
- INDEX idx_location (longitude, latitude)
- FOREIGN KEY (road_segment_id) REFERENCES road_segments(id) ON DELETE SET NULL

### 3. traffic_alerts (交通告警表)
存储交通告警信息，包括超速、拥堵、事故等告警。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | INT UNSIGNED | 主键 | AUTO_INCREMENT |
| alert_type | VARCHAR(20) | 告警类型 | NOT NULL |
| vehicle_id | VARCHAR(20) | 车辆ID | NULL |
| road_segment_id | INT UNSIGNED | 路段ID | NOT NULL |
| alert_value | DECIMAL(10,2) | 告警数值 | NULL |
| message | VARCHAR(200) | 告警消息 | NOT NULL |
| severity | VARCHAR(10) | 严重程度 | DEFAULT medium |
| resolved | BOOLEAN | 是否已解决 | DEFAULT FALSE |
| timestamp | DATETIME | 告警时间 | NOT NULL |
| created_at | DATETIME | 创建时间 | DEFAULT CURRENT_TIMESTAMP |

**索引:**
- PRIMARY KEY (id)
- INDEX idx_alert_type (alert_type)
- INDEX idx_road_segment (road_segment_id)
- INDEX idx_timestamp (timestamp)
- INDEX idx_resolved (resolved)
- FOREIGN KEY (road_segment_id) REFERENCES road_segments(id) ON DELETE CASCADE

## 模型方法

### GPSData 模型方法
- `GetLocation() (float64, float64)` - 获取位置信息
- `IsSpeeding(maxSpeed int) bool` - 判断是否超速

### RoadSegment 模型方法
- `GetCenterPoint() (float64, float64)` - 获取路段中心点
- `GetLength() float64` - 计算路段长度
- `IsVehicleInSegment(lng, lat float64) bool` - 判断车辆是否在路段内

### TrafficAlert 模型方法
- `IsHighSeverity() bool` - 判断是否为高严重程度告警
- `IsRecent() bool` - 判断是否为最近告警（1小时内）
- `GetSeverityLevel() int` - 获取严重程度等级
- `Resolve()` - 解决告警

## 数据库迁移

### 执行迁移
```bash
# 进入后端目录
cd backend

# 执行数据库迁移
go run cmd/migrate.go
```

### 迁移内容
1. 创建表结构
2. 建立索引和外键约束
3. 插入初始数据（示例路段）

## 初始数据

系统会自动插入以下示例路段数据：
- 中山路-东段 (城市道路，限速60km/h)
- 解放路-南段 (城市道路，限速50km/h)
- 环城高速-西段 (高速公路，限速80km/h)

## 性能优化建议

1. **GPS数据表**：建议定期清理历史数据，保留最近30天的数据
2. **索引优化**：根据查询模式调整索引策略
3. **分区表**：对于大量GPS数据，可考虑按时间分区
4. **缓存策略**：热点路段数据可考虑Redis缓存
