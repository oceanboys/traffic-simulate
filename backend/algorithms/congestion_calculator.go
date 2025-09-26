package algorithms

import (
	"backend/models"
	"backend/repositories"
	"sync"
	"time"
)

// CongestionCalculator 拥堵计算器
type CongestionCalculator struct {
	gpsRepo   *repositories.GPSRepository
	roadRepo  *repositories.RoadRepository
	roadStats map[uint]*RoadStatistics
	mutex     sync.RWMutex
}

// RoadStatistics 路段统计
type RoadStatistics struct {
	RoadID          uint      `json:"road_id"`
	VehicleCount    int       `json:"vehicle_count"`
	AverageSpeed    float64   `json:"average_speed"`
	MaxSpeed        int       `json:"max_speed"`
	CongestionLevel float64   `json:"congestion_level"`
	LastUpdate      time.Time `json:"last_update"`
}

// CongestionLevel 拥堵等级
type CongestionLevel struct {
	Level       string  `json:"level"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

// NewCongestionCalculator 创建新的拥堵计算器
func NewCongestionCalculator() *CongestionCalculator {
	return &CongestionCalculator{
		gpsRepo:   repositories.NewGPSRepository(),
		roadRepo:  repositories.NewRoadRepository(),
		roadStats: make(map[uint]*RoadStatistics),
	}
}

// CalculateCongestion 计算拥堵情况
func (cc *CongestionCalculator) CalculateCongestion(roadID uint) float64 {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// 获取路段信息
	road, err := cc.roadRepo.GetByID(roadID)
	if err != nil {
		return 0
	}

	// 获取最近5分钟的数据
	since := time.Now().Add(-5 * time.Minute)
	gpsData, err := cc.gpsRepo.FindByRoad(roadID, since)
	if err != nil {
		return 0
	}

	if len(gpsData) == 0 {
		return 0
	}

	// 计算路段统计
	stats := cc.calculateRoadStatistics(roadID, gpsData, road)

	// 计算拥堵等级
	congestionLevel := cc.calculateCongestionLevel(stats, road)

	// 更新路段统计
	cc.roadStats[roadID] = stats
	stats.CongestionLevel = congestionLevel

	return congestionLevel
}

// calculateRoadStatistics 计算路段统计信息
func (cc *CongestionCalculator) calculateRoadStatistics(roadID uint, gpsData []models.GPSData, road *models.RoadSegment) *RoadStatistics {
	stats := &RoadStatistics{
		RoadID:       roadID,
		VehicleCount: len(gpsData),
		MaxSpeed:     road.MaxSpeed,
		LastUpdate:   time.Now(),
	}

	if len(gpsData) == 0 {
		stats.AverageSpeed = 0
		return stats
	}

	// 计算总速度
	totalSpeed := 0
	for _, data := range gpsData {
		totalSpeed += data.Speed
	}
	stats.AverageSpeed = float64(totalSpeed) / float64(len(gpsData))

	return stats
}

// calculateCongestionLevel 计算拥堵等级
func (cc *CongestionCalculator) calculateCongestionLevel(stats *RoadStatistics, road *models.RoadSegment) float64 {
	if stats.AverageSpeed == 0 {
		return 1.0 // 完全拥堵
	}

	// 速度比率
	speedRatio := stats.AverageSpeed / float64(road.MaxSpeed)

	// 密度比率
	densityRatio := float64(stats.VehicleCount) / float64(road.Capacity)

	// 拥堵评分 (0-1之间)
	congestionScore := (1-speedRatio)*0.7 + densityRatio*0.3

	// 限制在0-1之间
	if congestionScore < 0 {
		congestionScore = 0
	}
	if congestionScore > 1 {
		congestionScore = 1
	}

	return congestionScore
}

// GetCongestionLevel 获取拥堵等级
func (cc *CongestionCalculator) GetCongestionLevel(roadID uint) CongestionLevel {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	stats, exists := cc.roadStats[roadID]
	if !exists {
		return CongestionLevel{
			Level:       "unknown",
			Score:       0,
			Description: "未知",
		}
	}

	score := stats.CongestionLevel

	if score < 0.2 {
		return CongestionLevel{
			Level:       "free",
			Score:       score,
			Description: "畅通",
		}
	} else if score < 0.4 {
		return CongestionLevel{
			Level:       "light",
			Score:       score,
			Description: "轻微拥堵",
		}
	} else if score < 0.6 {
		return CongestionLevel{
			Level:       "moderate",
			Score:       score,
			Description: "中度拥堵",
		}
	} else if score < 0.8 {
		return CongestionLevel{
			Level:       "heavy",
			Score:       score,
			Description: "严重拥堵",
		}
	} else {
		return CongestionLevel{
			Level:       "severe",
			Score:       score,
			Description: "极度拥堵",
		}
	}
}

// GetRoadStatistics 获取路段统计
func (cc *CongestionCalculator) GetRoadStatistics(roadID uint) *RoadStatistics {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	stats, exists := cc.roadStats[roadID]
	if !exists {
		return nil
	}

	return &RoadStatistics{
		RoadID:          stats.RoadID,
		VehicleCount:    stats.VehicleCount,
		AverageSpeed:    stats.AverageSpeed,
		MaxSpeed:        stats.MaxSpeed,
		CongestionLevel: stats.CongestionLevel,
		LastUpdate:      stats.LastUpdate,
	}
}

// GetAllRoadStatistics 获取所有路段统计
func (cc *CongestionCalculator) GetAllRoadStatistics() map[uint]*RoadStatistics {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	result := make(map[uint]*RoadStatistics)
	for roadID, stats := range cc.roadStats {
		result[roadID] = &RoadStatistics{
			RoadID:          stats.RoadID,
			VehicleCount:    stats.VehicleCount,
			AverageSpeed:    stats.AverageSpeed,
			MaxSpeed:        stats.MaxSpeed,
			CongestionLevel: stats.CongestionLevel,
			LastUpdate:      stats.LastUpdate,
		}
	}

	return result
}

// CalculateTrafficFlow 计算交通流量
func (cc *CongestionCalculator) CalculateTrafficFlow(roadID uint, duration time.Duration) float64 {
	since := time.Now().Add(-duration)
	gpsData, err := cc.gpsRepo.FindByRoad(roadID, since)
	if err != nil {
		return 0
	}

	// 计算流量 (车辆数/小时)
	flowRate := float64(len(gpsData)) / duration.Hours()
	return flowRate
}

// PredictCongestion 预测拥堵情况
func (cc *CongestionCalculator) PredictCongestion(roadID uint, futureMinutes int) float64 {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	stats, exists := cc.roadStats[roadID]
	if !exists {
		return 0
	}

	// 基于当前拥堵情况预测
	currentCongestion := stats.CongestionLevel

	// 时间因子 (简化预测模型)
	timeFactor := 1.0 + float64(futureMinutes)*0.1

	// 预测拥堵情况
	predictedCongestion := currentCongestion * timeFactor

	// 限制在0-1之间
	if predictedCongestion > 1 {
		predictedCongestion = 1
	}

	return predictedCongestion
}

// GetCongestionTrend 获取拥堵趋势
func (cc *CongestionCalculator) GetCongestionTrend(roadID uint, hours int) []float64 {
	var trend []float64

	// 获取历史数据趋势 (简化实现)
	for i := 0; i < hours; i++ {
		// 使用简化的趋势计算
		congestion := 0.3 + float64(i)*0.1
		if congestion > 1 {
			congestion = 1
		}
		trend = append(trend, congestion)
	}

	return trend
}

// UpdateRoadStatistics 更新路段统计
func (cc *CongestionCalculator) UpdateRoadStatistics(roadID uint) {
	cc.CalculateCongestion(roadID)
}
