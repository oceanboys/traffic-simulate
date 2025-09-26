# 交通业务算法文档

## 概述
交通业务算法模块为智慧交通实时数据平台提供核心的交通分析和处理功能，包括路段匹配、超速检测、拥堵计算和异常检测等算法。

## 核心算法

### 1. 路段匹配算法 (RoadMatcher)

**功能描述：**
- 根据GPS坐标找到最近的路段
- 计算车辆到路段的距离
- 判断车辆是否在指定路段上

**核心方法：**
```go
// 查找最近路段
func (rm *RoadMatcher) FindNearestRoad(lng, lat float64) (*models.RoadSegment, float64)

// 计算到路段的距离
func (rm *RoadMatcher) calculateDistanceToRoad(lng, lat float64, road *models.RoadSegment) float64

// 判断车辆是否在路段上
func (rm *RoadMatcher) IsVehicleOnRoad(lng, lat float64, road *models.RoadSegment, tolerance float64) bool
```

**算法特点：**
- 使用Haversine公式计算地理距离
- 支持点到线段距离计算
- 支持指定半径内的路段查找

### 2. 超速检测算法 (SpeedDetector)

**功能描述：**
- 实时检测车辆超速行为
- 记录车辆速度历史
- 分析速度模式和趋势
- 生成超速告警

**核心方法：**
```go
// 检查超速
func (sd *SpeedDetector) CheckOverspeed(gpsData models.GPSData, roadSegment *models.RoadSegment) bool

// 计算平均速度
func (sd *SpeedDetector) CalculateAverageSpeed(vehicleID string, duration time.Duration) float64

// 检测速度模式
func (sd *SpeedDetector) DetectSpeedPattern(vehicleID string) string
```

**检测规则：**
- 基础超速检测：速度 > 路段限速
- 持续超速检测：连续多次超速
- 极速异常检测：速度 > 150 km/h
- 低速异常检测：速度 < 5 km/h

### 3. 拥堵计算算法 (CongestionCalculator)

**功能描述：**
- 计算路段拥堵指数
- 分析交通流量
- 预测拥堵趋势
- 提供拥堵等级分类

**核心方法：**
```go
// 计算拥堵指数
func (cc *CongestionCalculator) CalculateCongestion(roadID uint) float64

// 获取拥堵等级
func (cc *CongestionCalculator) GetCongestionLevel(roadID uint) CongestionLevel

// 预测拥堵
func (cc *CongestionCalculator) PredictCongestion(roadID uint, futureMinutes int) float64
```

**拥堵等级：**
- 畅通 (0-0.2): 交通流畅
- 轻微拥堵 (0.2-0.4): 车流量适中
- 中度拥堵 (0.4-0.6): 车流量较大
- 重度拥堵 (0.6-0.8): 车流量很大
- 严重拥堵 (0.8-1.0): 交通瘫痪

**计算公式：**
```
拥堵指数 = (1-速度比)  0.7 + 密度比  0.3
其中：
- 速度比 = 平均速度 / 路段限速
- 密度比 = 车辆数 / 路段容量
```

### 4. 异常检测算法 (AnomalyDetector)

**功能描述：**
- 检测各种交通异常行为
- 支持自定义异常规则
- 生成异常告警
- 统计分析异常数据

**检测类型：**
- 速度异常：极速、低速、急加速/减速
- 位置异常：位置跳跃、异常轨迹
- 行为模式异常：偏离正常驾驶模式

**核心方法：**
```go
// 检测异常
func (ad *AnomalyDetector) DetectAnomalies(gpsData models.GPSData, roadSegment *models.RoadSegment) []DetectionRecord

// 添加异常规则
func (ad *AnomalyDetector) AddAnomalyRule(rule AnomalyRule)

// 获取异常统计
func (ad *AnomalyDetector) GetAnomalyStatistics(duration time.Duration) map[string]interface{}
```

## 交通分析服务

### TrafficAnalysisService

**功能描述：**
整合所有算法，提供统一的交通分析服务。

**核心功能：**
- 实时数据处理
- 路段分析
- 车辆分析
- 交通预测
- 系统概览

**主要方法：**
```go
// 处理实时数据
func (tas *TrafficAnalysisService) ProcessRealTimeData(gpsData models.GPSData) *TrafficAnalysisResult

// 获取路段分析
func (tas *TrafficAnalysisService) GetRoadAnalysis(roadID uint) map[string]interface{}

// 获取车辆分析
func (tas *TrafficAnalysisService) GetVehicleAnalysis(vehicleID string) map[string]interface{}

// 获取交通预测
func (tas *TrafficAnalysisService) GetTrafficPredictions(roadID uint, futureMinutes int) map[string]interface{}
```

## API接口

### 交通分析接口
- `POST /api/traffic/process` - 处理GPS数据
- `GET /api/traffic/overview` - 获取系统概览
- `GET /api/traffic/road/:road_id` - 获取路段分析
- `GET /api/traffic/vehicle/:vehicle_id` - 获取车辆分析
- `GET /api/traffic/predict/:road_id` - 获取交通预测
- `GET /api/traffic/stats` - 获取实时统计
- `GET /api/traffic/congestion` - 获取拥堵等级
- `GET /api/traffic/anomalies` - 获取异常统计

### 请求示例
```json
// 处理GPS数据
POST /api/traffic/process
{
    "vehicle_id": "VH0001",
    "longitude": 120.1601,
    "latitude": 30.2791,
    "speed": 60,
    "direction": 90,
    "timestamp": "2024-01-01T12:00:00Z"
}
```

### 响应示例
```json
{
    "success": true,
    "data": {
        "road_segment": {
            "id": 1,
            "name": "中山路-东段",
            "max_speed": 60
        },
        "distance": 0.025,
        "is_overspeed": false,
        "congestion_level": {
            "level": "light",
            "score": 0.3,
            "description": "轻微拥堵"
        },
        "anomalies": [],
        "stats": {
            "road_id": 1,
            "vehicle_count": 15,
            "average_speed": 45.5,
            "congestion_level": 0.3
        }
    }
}
```

## 性能优化

### 1. 算法优化
- 使用空间索引加速路段匹配
- 缓存计算结果减少重复计算
- 批量处理提高吞吐量

### 2. 内存管理
- 定期清理历史数据
- 使用对象池减少GC压力
- 合理设置缓存大小

### 3. 并发处理
- 使用读写锁保护共享数据
- 异步处理非关键计算
- 合理设置工作协程数量

## 监控和调试

### 1. 性能监控
```go
// 获取性能指标
metrics := analysisService.GetPerformanceMetrics()
fmt.Printf("性能指标: %+v\n", metrics)
```

### 2. 算法测试
```bash
# 运行算法测试
go run cmd/test_algorithms.go
```

### 3. 日志记录
系统会自动记录关键事件和错误信息，便于问题定位。

## 扩展开发

### 1. 自定义算法
```go
// 实现自定义算法接口
type CustomAlgorithm struct{}

func (ca *CustomAlgorithm) Process(data models.GPSData) (result interface{}, err error) {
    // 自定义处理逻辑
    return result, nil
}
```

### 2. 添加新的异常规则
```go
rule := algorithms.AnomalyRule{
    ID:          "custom_rule_1",
    Name:        "自定义规则",
    Type:        "custom",
    Condition:   "custom_condition",
    Threshold:   100.0,
    Severity:    "medium",
    Enabled:     true,
    Description: "自定义异常检测规则",
}
anomalyDetector.AddAnomalyRule(rule)
```

### 3. 扩展拥堵计算
可以基于更多因素计算拥堵指数，如：
- 天气条件
- 时间段
- 历史数据
- 事件影响

## 最佳实践

1. **算法选择**：根据实际需求选择合适的算法
2. **参数调优**：根据实际数据调整算法参数
3. **性能监控**：定期检查算法性能
4. **数据质量**：确保输入数据的准确性
5. **异常处理**：合理处理算法异常情况

## 故障排除

### 常见问题
1. **路段匹配不准确**：检查路段数据质量
2. **拥堵计算异常**：验证输入数据格式
3. **性能问题**：检查算法复杂度和数据量
4. **内存泄漏**：定期清理历史数据

### 调试技巧
1. 使用测试数据验证算法
2. 监控关键性能指标
3. 查看详细日志信息
4. 使用性能分析工具
