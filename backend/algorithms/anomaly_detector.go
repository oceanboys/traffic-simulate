package algorithms

import (
	"backend/models"
	"backend/repositories"
	"fmt"
	"math"
	"sync"
	"time"
)

// AnomalyDetector 异常检测器
type AnomalyDetector struct {
	alertRepo        *repositories.AlertRepository
	anomalyRules     []AnomalyRule
	detectionHistory map[string][]DetectionRecord
	mutex            sync.RWMutex
}

// AnomalyRule 异常规则
type AnomalyRule struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`        // speed, location, pattern
	Condition   string  `json:"condition"`   // 条件表达式
	Threshold   float64 `json:"threshold"`   // 阈值
	Severity    string  `json:"severity"`    // low, medium, high, critical
	Enabled     bool    `json:"enabled"`     // 是否启用
	Description string  `json:"description"` // 描述
}

// DetectionRecord 检测记录
type DetectionRecord struct {
	VehicleID   string    `json:"vehicle_id"`
	AnomalyType string    `json:"anomaly_type"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Timestamp   time.Time `json:"timestamp"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
}

// NewAnomalyDetector 创建异常检测器
func NewAnomalyDetector() *AnomalyDetector {
	detector := &AnomalyDetector{
		alertRepo:        repositories.NewAlertRepository(),
		detectionHistory: make(map[string][]DetectionRecord),
	}

	// 初始化默认规则
	detector.initDefaultRules()

	return detector
}

// initDefaultRules 初始化默认规则
func (ad *AnomalyDetector) initDefaultRules() {
	ad.anomalyRules = []AnomalyRule{
		{
			ID:          "speed_anomaly_1",
			Name:        "极速异常",
			Type:        "speed",
			Condition:   "speed > threshold",
			Threshold:   150.0,
			Severity:    "critical",
			Enabled:     true,
			Description: "检测到车辆速度超过150km/h",
		},
		{
			ID:          "speed_anomaly_2",
			Name:        "低速异常",
			Type:        "speed",
			Condition:   "speed < threshold",
			Threshold:   5.0,
			Severity:    "medium",
			Enabled:     true,
			Description: "检测到车辆速度低于5km/h（可能停车）",
		},
		{
			ID:          "location_anomaly_1",
			Name:        "位置异常",
			Type:        "location",
			Condition:   "distance > threshold",
			Threshold:   1.0,
			Severity:    "high",
			Enabled:     true,
			Description: "检测到车辆位置异常跳跃",
		},
		{
			ID:          "pattern_anomaly_1",
			Name:        "行为模式异常",
			Type:        "pattern",
			Condition:   "pattern_deviation > threshold",
			Threshold:   0.8,
			Severity:    "medium",
			Enabled:     true,
			Description: "检测到车辆行为模式异常",
		},
	}
}

// DetectAnomalies 检测异常
func (ad *AnomalyDetector) DetectAnomalies(gpsData models.GPSData, roadSegment *models.RoadSegment) []DetectionRecord {
	var anomalies []DetectionRecord

	// 速度异常检测
	speedAnomalies := ad.detectSpeedAnomalies(gpsData)
	anomalies = append(anomalies, speedAnomalies...)

	// 位置异常检测
	locationAnomalies := ad.detectLocationAnomalies(gpsData)
	anomalies = append(anomalies, locationAnomalies...)

	// 行为模式异常检测
	patternAnomalies := ad.detectPatternAnomalies(gpsData)
	anomalies = append(anomalies, patternAnomalies...)

	// 记录检测结果
	for _, anomaly := range anomalies {
		ad.recordDetection(anomaly)
	}

	return anomalies
}

// detectSpeedAnomalies 检测速度异常
func (ad *AnomalyDetector) detectSpeedAnomalies(gpsData models.GPSData) []DetectionRecord {
	var anomalies []DetectionRecord

	for _, rule := range ad.anomalyRules {
		if rule.Type != "speed" || !rule.Enabled {
			continue
		}

		var detected bool
		var value float64

		switch rule.Condition {
		case "speed > threshold":
			detected = float64(gpsData.Speed) > rule.Threshold
			value = float64(gpsData.Speed)
		case "speed < threshold":
			detected = float64(gpsData.Speed) < rule.Threshold
			value = float64(gpsData.Speed)
		}

		if detected {
			anomaly := DetectionRecord{
				VehicleID:   gpsData.VehicleID,
				AnomalyType: rule.Name,
				Value:       value,
				Threshold:   rule.Threshold,
				Timestamp:   time.Now(),
				Severity:    rule.Severity,
				Message:     fmt.Sprintf("%s: 当前值 %.2f, 阈值 %.2f", rule.Description, value, rule.Threshold),
			}
			anomalies = append(anomalies, anomaly)
		}
	}

	return anomalies
}

// detectLocationAnomalies 检测位置异常
func (ad *AnomalyDetector) detectLocationAnomalies(gpsData models.GPSData) []DetectionRecord {
	var anomalies []DetectionRecord

	// 检查位置跳跃
	lastPosition := ad.getLastPosition(gpsData.VehicleID)
	if lastPosition != nil {
		distance := ad.calculateDistance(
			lastPosition.Longitude, lastPosition.Latitude,
			gpsData.Longitude, gpsData.Latitude,
		)

		timeDiff := gpsData.Timestamp.Sub(lastPosition.Timestamp).Seconds()
		if timeDiff > 0 {
			speed := distance / timeDiff * 3.6 // 转换为km/h

			// 如果速度超过300km/h，认为是位置跳跃
			if speed > 300 {
				anomaly := DetectionRecord{
					VehicleID:   gpsData.VehicleID,
					AnomalyType: "位置跳跃",
					Value:       speed,
					Threshold:   300.0,
					Timestamp:   time.Now(),
					Severity:    "high",
					Message:     fmt.Sprintf("检测到位置跳跃，计算速度 %.2f km/h", speed),
				}
				anomalies = append(anomalies, anomaly)
			}
		}
	}

	// 更新位置记录
	ad.updatePosition(gpsData)

	return anomalies
}

// detectPatternAnomalies 检测行为模式异常
func (ad *AnomalyDetector) detectPatternAnomalies(gpsData models.GPSData) []DetectionRecord {
	var anomalies []DetectionRecord

	// 获取车辆历史行为
	history := ad.getVehicleHistory(gpsData.VehicleID)
	if len(history) < 10 {
		return anomalies // 数据不足，无法检测模式
	}

	// 分析速度模式
	speedPattern := ad.analyzeSpeedPattern(history)
	if speedPattern.Deviation > 0.8 {
		anomaly := DetectionRecord{
			VehicleID:   gpsData.VehicleID,
			AnomalyType: "行为模式异常",
			Value:       speedPattern.Deviation,
			Threshold:   0.8,
			Timestamp:   time.Now(),
			Severity:    "medium",
			Message:     fmt.Sprintf("检测到行为模式异常，偏差度 %.2f", speedPattern.Deviation),
		}
		anomalies = append(anomalies, anomaly)
	}

	return anomalies
}

// getLastPosition 获取最后位置
func (ad *AnomalyDetector) getLastPosition(vehicleID string) *models.GPSData {
	ad.mutex.RLock()
	defer ad.mutex.RUnlock()

	history := ad.detectionHistory[vehicleID]
	if len(history) == 0 {
		return nil
	}

	// 简化实现，实际应该从数据库获取
	return nil
}

// updatePosition 更新位置记录
func (ad *AnomalyDetector) updatePosition(_ models.GPSData) {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	// 简化实现，实际应该存储到数据库
}

// calculateDistance 计算距离
func (ad *AnomalyDetector) calculateDistance(lng1, lat1, lng2, lat2 float64) float64 {
	const R = 6371 // 地球半径（公里）

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// recordDetection 记录检测结果
func (ad *AnomalyDetector) recordDetection(record DetectionRecord) {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	ad.detectionHistory[record.VehicleID] = append(
		ad.detectionHistory[record.VehicleID], record)

	// 只保留最近100条记录
	if len(ad.detectionHistory[record.VehicleID]) > 100 {
		ad.detectionHistory[record.VehicleID] = ad.detectionHistory[record.VehicleID][1:]
	}
}

// getVehicleHistory 获取车辆历史
func (ad *AnomalyDetector) getVehicleHistory(_ string) []models.GPSData {
	// 简化实现，实际应该从数据库获取
	return []models.GPSData{}
}

// analyzeSpeedPattern 分析速度模式
func (ad *AnomalyDetector) analyzeSpeedPattern(_ []models.GPSData) SpeedPattern {
	// 简化实现
	return SpeedPattern{
		Average:   50.0,
		Deviation: 0.3,
		Trend:     "stable",
	}
}

// SpeedPattern 速度模式
type SpeedPattern struct {
	Average   float64 `json:"average"`
	Deviation float64 `json:"deviation"`
	Trend     string  `json:"trend"`
}

// CreateAnomalyAlert 创建异常告警
func (ad *AnomalyDetector) CreateAnomalyAlert(record DetectionRecord, roadSegment *models.RoadSegment) error {
	alert := models.TrafficAlert{
		AlertType:   "anomaly",
		VehicleID:   record.VehicleID,
		RoadSegment: roadSegment,
		AlertValue:  record.Value,
		Message:     record.Message,
		Severity:    record.Severity,
		Timestamp:   record.Timestamp,
		CreatedAt:   time.Now(),
	}

	return ad.alertRepo.Create(&alert)
}

// GetAnomalyStatistics 获取异常统计
func (ad *AnomalyDetector) GetAnomalyStatistics(duration time.Duration) map[string]interface{} {
	ad.mutex.RLock()
	defer ad.mutex.RUnlock()

	cutoff := time.Now().Add(-duration)
	totalAnomalies := 0
	anomalyByType := make(map[string]int)
	anomalyBySeverity := make(map[string]int)

	for _, records := range ad.detectionHistory {
		for _, record := range records {
			if record.Timestamp.After(cutoff) {
				totalAnomalies++
				anomalyByType[record.AnomalyType]++
				anomalyBySeverity[record.Severity]++
			}
		}
	}

	return map[string]interface{}{
		"total_anomalies":     totalAnomalies,
		"anomaly_by_type":     anomalyByType,
		"anomaly_by_severity": anomalyBySeverity,
		"duration_hours":      duration.Hours(),
	}
}

// AddAnomalyRule 添加异常规则
func (ad *AnomalyDetector) AddAnomalyRule(rule AnomalyRule) {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	ad.anomalyRules = append(ad.anomalyRules, rule)
}

// GetAnomalyRules 获取异常规则
func (ad *AnomalyDetector) GetAnomalyRules() []AnomalyRule {
	ad.mutex.RLock()
	defer ad.mutex.RUnlock()

	return ad.anomalyRules
}

// UpdateAnomalyRule 更新异常规则
func (ad *AnomalyDetector) UpdateAnomalyRule(ruleID string, rule AnomalyRule) bool {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	for i, r := range ad.anomalyRules {
		if r.ID == ruleID {
			ad.anomalyRules[i] = rule
			return true
		}
	}

	return false
}

// DeleteAnomalyRule 删除异常规则
func (ad *AnomalyDetector) DeleteAnomalyRule(ruleID string) bool {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	for i, r := range ad.anomalyRules {
		if r.ID == ruleID {
			ad.anomalyRules = append(ad.anomalyRules[:i], ad.anomalyRules[i+1:]...)
			return true
		}
	}

	return false
}
