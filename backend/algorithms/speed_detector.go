package algorithms

import (
	"backend/models"
	"backend/repositories"
	"fmt"
	"time"
)

// SpeedDetector 超速检测算法
type SpeedDetector struct {
	alertRepo *repositories.AlertRepository
	history   map[string][]SpeedRecord // 车辆速度历史
}

// SpeedRecord 速度记录
type SpeedRecord struct {
	Speed     int       `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
	RoadID    uint      `json:"road_id"`
}

// NewSpeedDetector 创建超速检测器
func NewSpeedDetector() *SpeedDetector {
	return &SpeedDetector{
		alertRepo: repositories.NewAlertRepository(),
		history:   make(map[string][]SpeedRecord),
	}
}

// CheckOverspeed 检查超速
func (sd *SpeedDetector) CheckOverspeed(gpsData models.GPSData, roadSegment *models.RoadSegment) bool {
	if roadSegment == nil {
		return false
	}

	// 基本超速检查
	if gpsData.Speed <= roadSegment.MaxSpeed {
		return false
	}

	// 记录速度历史
	sd.recordSpeed(gpsData, roadSegment.ID)

	// 检查是否持续超速
	if sd.isPersistentOverspeed(gpsData.VehicleID) {
		return true
	}

	// 检查超速程度
	overspeedRatio := float64(gpsData.Speed) / float64(roadSegment.MaxSpeed)
	return overspeedRatio > 1.1 // 超过10%才认为是超速
}

// recordSpeed 记录速度
func (sd *SpeedDetector) recordSpeed(gpsData models.GPSData, roadID uint) {
	record := SpeedRecord{
		Speed:     gpsData.Speed,
		Timestamp: gpsData.Timestamp,
		RoadID:    roadID,
	}

	sd.history[gpsData.VehicleID] = append(sd.history[gpsData.VehicleID], record)

	// 只保留最近1小时的数据
	sd.cleanOldRecords(gpsData.VehicleID)
}

// cleanOldRecords 清理旧记录
func (sd *SpeedDetector) cleanOldRecords(vehicleID string) {
	cutoff := time.Now().Add(-time.Hour)
	var validRecords []SpeedRecord

	for _, record := range sd.history[vehicleID] {
		if record.Timestamp.After(cutoff) {
			validRecords = append(validRecords, record)
		}
	}

	sd.history[vehicleID] = validRecords
}

// isPersistentOverspeed 检查是否持续超速
func (sd *SpeedDetector) isPersistentOverspeed(vehicleID string) bool {
	records := sd.history[vehicleID]
	if len(records) < 3 {
		return false
	}

	// 检查最近3次记录是否都超速
	recentRecords := records[len(records)-3:]
	for _, record := range recentRecords {
		// 这里需要根据路段限速判断，简化处理
		if record.Speed <= 60 { // 假设默认限速60
			return false
		}
	}

	return true
}

// GetSpeedHistory 获取速度历史
func (sd *SpeedDetector) GetSpeedHistory(vehicleID string) []SpeedRecord {
	return sd.history[vehicleID]
}

// CalculateAverageSpeed 计算平均速度
func (sd *SpeedDetector) CalculateAverageSpeed(vehicleID string, duration time.Duration) float64 {
	records := sd.history[vehicleID]
	if len(records) == 0 {
		return 0
	}

	cutoff := time.Now().Add(-duration)
	var validRecords []SpeedRecord

	for _, record := range records {
		if record.Timestamp.After(cutoff) {
			validRecords = append(validRecords, record)
		}
	}

	if len(validRecords) == 0 {
		return 0
	}

	totalSpeed := 0
	for _, record := range validRecords {
		totalSpeed += record.Speed
	}

	return float64(totalSpeed) / float64(len(validRecords))
}

// DetectSpeedPattern 检测速度模式
func (sd *SpeedDetector) DetectSpeedPattern(vehicleID string) string {
	records := sd.history[vehicleID]
	if len(records) < 5 {
		return "insufficient_data"
	}

	// 计算速度变化趋势
	var speedChanges []int
	for i := 1; i < len(records); i++ {
		change := records[i].Speed - records[i-1].Speed
		speedChanges = append(speedChanges, change)
	}

	// 分析趋势
	accelerating := 0
	decelerating := 0
	stable := 0

	for _, change := range speedChanges {
		if change > 5 {
			accelerating++
		} else if change < -5 {
			decelerating++
		} else {
			stable++
		}
	}

	total := len(speedChanges)
	if float64(accelerating)/float64(total) > 0.6 {
		return "accelerating"
	} else if float64(decelerating)/float64(total) > 0.6 {
		return "decelerating"
	} else {
		return "stable"
	}
}

// GetOverspeedStatistics 获取超速统计
func (sd *SpeedDetector) GetOverspeedStatistics(roadID uint, duration time.Duration) map[string]interface{} {
	cutoff := time.Now().Add(-duration)
	overspeedCount := 0
	totalCount := 0
	maxSpeed := 0

	for _, records := range sd.history {
		for _, record := range records {
			if record.RoadID == roadID && record.Timestamp.After(cutoff) {
				totalCount++
				if record.Speed > 60 { // 假设默认限速60
					overspeedCount++
					if record.Speed > maxSpeed {
						maxSpeed = record.Speed
					}
				}
			}
		}
	}

	overspeedRate := 0.0
	if totalCount > 0 {
		overspeedRate = float64(overspeedCount) / float64(totalCount)
	}

	return map[string]interface{}{
		"overspeed_count": overspeedCount,
		"total_count":     totalCount,
		"overspeed_rate":  overspeedRate,
		"max_speed":       maxSpeed,
		"duration_hours":  duration.Hours(),
	}
}

// CreateSpeedAlert 创建超速告警
func (sd *SpeedDetector) CreateSpeedAlert(gpsData models.GPSData, roadSegment *models.RoadSegment) error {
	overspeedRatio := float64(gpsData.Speed) / float64(roadSegment.MaxSpeed)

	severity := "medium"
	if overspeedRatio >= 1.5 {
		severity = "critical"
	} else if overspeedRatio >= 1.3 {
		severity = "high"
	} else if overspeedRatio >= 1.1 {
		severity = "medium"
	} else {
		severity = "low"
	}

	alert := models.TrafficAlert{
		AlertType:   "speeding",
		VehicleID:   gpsData.VehicleID,
		RoadSegment: roadSegment,
		AlertValue:  float64(gpsData.Speed),
		Message: fmt.Sprintf("车辆 %s 超速行驶，当前速度 %d km/h，限速 %d km/h",
			gpsData.VehicleID, gpsData.Speed, roadSegment.MaxSpeed),
		Severity:  severity,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	return sd.alertRepo.Create(&alert)
}
