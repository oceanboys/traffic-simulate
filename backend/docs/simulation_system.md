# 数据模拟生成系统文档

## 概述
数据模拟生成系统为智慧交通实时数据平台提供车辆GPS数据模拟功能，支持可配置的模拟参数和实时数据流处理。

## 核心组件

### 1. SimulationService - 基础模拟服务
提供基本的车辆GPS数据模拟功能。

**主要功能：**
- 多车辆并发模拟
- 可配置的模拟参数
- 实时数据流处理
- 自动告警检测

**配置参数：**
```go
type SimulationConfig struct {
    VehicleCount    int           // 车辆数量
    Interval        time.Duration // 数据生成间隔
    SpeedRange      SpeedRange    // 速度范围
    RoadSegments    []uint        // 路段ID列表
    VehicleTypes    []string      // 车辆类型
    EnableAlerts    bool          // 是否启用告警
    AlertThreshold  float64       // 告警阈值
}
```

### 2. RealisticSimulationService - 真实感模拟服务
基于时间模式和交通规律的增强模拟服务。

**特色功能：**
- 基于时间的交通模式（早高峰、晚高峰、夜间等）
- 天气条件影响
- 交通事件模拟
- 更真实的车辆行为

**时间模式：**
- 早高峰 (7-9点): 高拥堵，低速度
- 午间 (10-14点): 正常交通
- 晚高峰 (17-19点): 高拥堵，低速度
- 夜间 (20-6点): 低流量，高速度

### 3. DataPipeline - 数据流处理管道
提供可扩展的数据处理管道。

**处理器类型：**
- ValidationProcessor: 数据验证
- StatisticsProcessor: 统计计算
- AlertProcessor: 告警检测

### 4. SimulationConfigManager - 配置管理器
管理模拟配置的创建、验证和存储。

**预设配置：**
- light: 轻量级模拟 (5辆车)
- medium: 中等强度 (20辆车)
- heavy: 高强度 (50辆车)
- test: 测试配置 (2辆车)

## 使用方法

### 1. 基础使用
```go
// 创建模拟服务
simulationService := services.NewSimulationService()

// 配置参数
config := services.SimulationConfig{
    VehicleCount:   10,
    Interval:       5 * time.Second,
    SpeedRange:     services.SpeedRange{Min: 30, Max: 80},
    VehicleTypes:   []string{"car", "truck"},
    EnableAlerts:   true,
}

// 启动模拟
err := simulationService.StartSimulation(config)
```

### 2. 真实感模拟
```go
// 创建真实感模拟服务
realisticService := services.NewRealisticSimulationService()

// 设置天气条件
weather := services.WeatherCondition{
    Type:        "rainy",
    Visibility:  0.7,
    SpeedImpact: 0.8,
}
realisticService.SetWeather(weather)

// 添加交通事件
event := services.TrafficEvent{
    Type:        "accident",
    Location:     services.Location{Lng: 120.1551, Lat: 30.2741, Radius: 0.5},
    StartTime:   time.Now(),
    EndTime:     time.Now().Add(2 * time.Hour),
    Impact:      0.8,
    Description: "交通事故",
}
realisticService.AddTrafficEvent(event)
```

### 3. 数据管道使用
```go
// 创建数据管道
pipeline := services.CreateDefaultPipeline()

// 启动管道
pipeline.Start()

// 输入数据
pipeline.Input(gpsData)

// 获取输出数据
for data := range pipeline.Output() {
    // 处理数据
    fmt.Printf("处理数据: %+v\n", data)
}
```

## API接口

### 模拟控制接口
- `POST /api/simulation/start` - 启动模拟
- `POST /api/simulation/stop` - 停止模拟
- `GET /api/simulation/status` - 获取状态
- `GET /api/simulation/config` - 获取配置
- `PUT /api/simulation/config` - 更新配置
- `GET /api/simulation/statistics` - 获取统计

### 请求示例
```json
// 启动模拟
POST /api/simulation/start
{
    "vehicle_count": 20,
    "interval": "5s",
    "speed_range": {
        "min": 30,
        "max": 80
    },
    "vehicle_types": ["car", "truck", "bus"],
    "enable_alerts": true
}
```

## 性能优化

### 1. 并发控制
- 使用goroutine池控制并发数量
- 避免创建过多协程导致资源耗尽

### 2. 内存管理
- 使用带缓冲的channel避免阻塞
- 定期清理历史数据

### 3. 数据库优化
- 批量插入GPS数据
- 使用连接池管理数据库连接

## 监控和调试

### 1. 状态监控
```go
status := simulationService.GetSimulationStatus()
fmt.Printf("运行状态: %+v\n", status)
```

### 2. 管道统计
```go
stats := pipeline.GetStats()
fmt.Printf("管道统计: %+v\n", stats)
```

### 3. 日志记录
系统会自动记录关键事件和错误信息。

## 扩展开发

### 1. 自定义处理器
```go
type CustomProcessor struct{}

func (p *CustomProcessor) Process(data models.GPSData) (models.GPSData, error) {
    // 自定义处理逻辑
    return data, nil
}

func (p *CustomProcessor) Name() string {
    return "CustomProcessor"
}
```

### 2. 自定义交通模式
```go
// 添加自定义时间模式
pattern := services.TimeBasedPattern{
    Hour: 12,
    Pattern: services.TrafficPattern{
        Name:        "午休时间",
        SpeedFactor: 0.9,
        Volume:      3,
        Congestion:  0.2,
    },
    Probability: 0.8,
}
```

## 故障排除

### 常见问题
1. **模拟启动失败**: 检查数据库连接和配置参数
2. **数据生成缓慢**: 调整车辆数量和间隔时间
3. **内存占用过高**: 减少车辆数量或增加数据清理频率

### 调试技巧
1. 使用测试配置进行调试
2. 监控系统资源使用情况
3. 查看日志文件定位问题

## 最佳实践

1. **开发环境**: 使用轻量级配置
2. **测试环境**: 使用中等强度配置
3. **生产环境**: 根据实际需求调整配置
4. **监控**: 定期检查系统状态和性能指标
