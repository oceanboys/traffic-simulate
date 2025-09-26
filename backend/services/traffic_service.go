package services

import (
	"backend/models"
	"math/rand"
	"sync"
	"time"
)

// TrafficService 交通服务
type TrafficService struct {
	vehicles   []models.Vehicle
	alerts     []models.TrafficAlert
	simulating bool
	mu         sync.RWMutex
	stopChan   chan bool
}

// NewTrafficService 创建交通服务
func NewTrafficService() *TrafficService {
	service := &TrafficService{
		vehicles:   make([]models.Vehicle, 0),
		alerts:     make([]models.TrafficAlert, 0),
		simulating: false,
		stopChan:   make(chan bool),
	}

	// 初始化一些默认车辆
	service.initializeDefaultVehicles()

	return service
}

// 初始化默认车辆
func (s *TrafficService) initializeDefaultVehicles() {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicles := []models.Vehicle{
		{
			ID:          1,
			VehicleID:   "V001",
			X:           20.0,
			Y:           20.0,
			Speed:       45.0,
			Direction:   0.0,
			VehicleType: "car",
			Status:      "normal",
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			VehicleID:   "V002",
			X:           40.0,
			Y:           40.0,
			Speed:       85.0,
			Direction:   90.0,
			VehicleType: "truck",
			Status:      "overspeed",
			CreatedAt:   time.Now(),
		},
		{
			ID:          3,
			VehicleID:   "V003",
			X:           60.0,
			Y:           60.0,
			Speed:       35.0,
			Direction:   180.0,
			VehicleType: "bus",
			Status:      "normal",
			CreatedAt:   time.Now(),
		},
	}

	s.vehicles = vehicles
}

// TrafficSummary 交通摘要
type TrafficSummary struct {
	TotalVehicles   int     `json:"total_vehicles"`
	AverageSpeed    float64 `json:"average_speed"`
	CongestionLevel int     `json:"congestion_level"`
	ActiveAlerts    int     `json:"active_alerts"`
	LastUpdate      string  `json:"last_update"`
}

// TrafficStats 交通统计
type TrafficStats struct {
	TotalRoads      int     `json:"total_roads"`
	TotalVehicles   int     `json:"total_vehicles"`
	AverageSpeed    float64 `json:"average_speed"`
	CongestionLevel int     `json:"congestion_level"`
	ActiveAlerts    int     `json:"active_alerts"`
	OverspeedCount  int     `json:"overspeed_count"`
	AccidentCount   int     `json:"accident_count"`
}

// CongestionData 拥堵数据
type CongestionData struct {
	RoadID          uint    `json:"road_id"`
	RoadName        string  `json:"road_name"`
	CongestionLevel float64 `json:"congestion_level"`
	VehicleCount    int     `json:"vehicle_count"`
	AverageSpeed    float64 `json:"average_speed"`
	Status          string  `json:"status"`
}

// VehicleFlow 车流数据
type VehicleFlow struct {
	TotalVehicles    int     `json:"total_vehicles"`
	IncomingVehicles int     `json:"incoming_vehicles"`
	OutgoingVehicles int     `json:"outgoing_vehicles"`
	FlowRate         float64 `json:"flow_rate"`
	PeakHour         string  `json:"peak_hour"`
}

// GetRealTimeSummary 获取实时摘要
func (s *TrafficService) GetRealTimeSummary() TrafficSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalVehicles := len(s.vehicles)
	var totalSpeed float64
	for _, v := range s.vehicles {
		totalSpeed += v.Speed
	}

	activeAlerts := 0
	for _, a := range s.alerts {
		if !a.Resolved {
			activeAlerts++
		}
	}

	return TrafficSummary{
		TotalVehicles:   totalVehicles,
		AverageSpeed:    totalSpeed / float64(totalVehicles),
		CongestionLevel: 1,
		ActiveAlerts:    activeAlerts,
		LastUpdate:      time.Now().Format("2006-01-02 15:04:05"),
	}
}

// GetTrafficStats 获取交通统计
func (s *TrafficService) GetTrafficStats() TrafficStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalVehicles := len(s.vehicles)
	var totalSpeed float64
	overspeedCount := 0

	for _, v := range s.vehicles {
		totalSpeed += v.Speed
		if v.Speed > 80 {
			overspeedCount++
		}
	}

	activeAlerts := 0
	for _, a := range s.alerts {
		if !a.Resolved {
			activeAlerts++
		}
	}

	return TrafficStats{
		TotalRoads:      25,
		TotalVehicles:   totalVehicles,
		AverageSpeed:    totalSpeed / float64(totalVehicles),
		CongestionLevel: 1,
		ActiveAlerts:    activeAlerts,
		OverspeedCount:  overspeedCount,
		AccidentCount:   1,
	}
}

// GetRecentAlerts 获取最近告警
func (s *TrafficService) GetRecentAlerts(limit int) []models.TrafficAlert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.alerts) > limit {
		return s.alerts[:limit]
	}
	return s.alerts
}

// GetCongestionData 获取拥堵数据
func (s *TrafficService) GetCongestionData() []CongestionData {
	return []CongestionData{
		{
			RoadID:          1,
			RoadName:        "主干道A",
			CongestionLevel: 0.8,
			VehicleCount:    45,
			AverageSpeed:    25.5,
			Status:          "拥堵",
		},
		{
			RoadID:          2,
			RoadName:        "主干道B",
			CongestionLevel: 0.3,
			VehicleCount:    20,
			AverageSpeed:    55.2,
			Status:          "畅通",
		},
	}
}

// GetVehicleFlow 获取车流数据
func (s *TrafficService) GetVehicleFlow() VehicleFlow {
	return VehicleFlow{
		TotalVehicles:    len(s.vehicles),
		IncomingVehicles: 45,
		OutgoingVehicles: 38,
		FlowRate:         0.85,
		PeakHour:         "17:00-18:00",
	}
}

// GetVehicles 获取车辆列表
func (s *TrafficService) GetVehicles() []models.Vehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vehicles
}

// AddVehicle 添加车辆
func (s *TrafficService) AddVehicle(vehicle models.Vehicle) models.Vehicle {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle.ID = uint(len(s.vehicles) + 1)
	vehicle.CreatedAt = time.Now()
	s.vehicles = append(s.vehicles, vehicle)

	return vehicle
}

// RemoveVehicle 移除车辆
func (s *TrafficService) RemoveVehicle(vehicleID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.vehicles {
		if v.VehicleID == vehicleID {
			s.vehicles = append(s.vehicles[:i], s.vehicles[i+1:]...)
			return true
		}
	}
	return false
}

// StartSimulation 开始模拟
func (s *TrafficService) StartSimulation() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.simulating {
		return
	}

	s.simulating = true
	go s.runSimulation()
}

// StopSimulation 停止模拟
func (s *TrafficService) StopSimulation() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.simulating {
		return
	}

	s.simulating = false
	s.stopChan <- true
}

// GetSimulationStatus 获取模拟状态
func (s *TrafficService) GetSimulationStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"simulating":    s.simulating,
		"vehicle_count": len(s.vehicles),
		"alert_count":   len(s.alerts),
		"last_update":   time.Now().Format("2006-01-02 15:04:05"),
	}
}

// 运行模拟
func (s *TrafficService) runSimulation() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.updateVehicles()
			s.generateAlerts()
		case <-s.stopChan:
			return
		}
	}
}

// 更新车辆状态
func (s *TrafficService) updateVehicles() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.vehicles {
		// 随机移动车辆
		speedFactor := s.vehicles[i].Speed / 100.0
		moveDistance := speedFactor * 2.0

		//radians := s.vehicles[i].Direction * 3.14159 / 180.0
		s.vehicles[i].X += float64(moveDistance) * float64(rand.Float64()-0.5)
		s.vehicles[i].Y += float64(moveDistance) * float64(rand.Float64()-0.5)

		// 边界检查
		if s.vehicles[i].X < 0 {
			s.vehicles[i].X = 100
		}
		if s.vehicles[i].X > 100 {
			s.vehicles[i].X = 0
		}
		if s.vehicles[i].Y < 0 {
			s.vehicles[i].Y = 100
		}
		if s.vehicles[i].Y > 100 {
			s.vehicles[i].Y = 0
		}

		// 随机改变方向
		if rand.Float64() < 0.1 {
			s.vehicles[i].Direction = rand.Float64() * 360
		}

		// 随机改变速度
		if rand.Float64() < 0.2 {
			speedChange := (rand.Float64() - 0.5) * 20
			s.vehicles[i].Speed = s.vehicles[i].Speed + speedChange
			if s.vehicles[i].Speed < 10 {
				s.vehicles[i].Speed = 10
			}
			if s.vehicles[i].Speed > 120 {
				s.vehicles[i].Speed = 120
			}
		}

		// 更新状态
		if s.vehicles[i].Speed > 80 {
			s.vehicles[i].Status = "overspeed"
		} else if s.vehicles[i].Speed < 20 {
			s.vehicles[i].Status = "slow"
		} else {
			s.vehicles[i].Status = "normal"
		}
	}
}

// 生成告警
func (s *TrafficService) generateAlerts() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查超速
	for _, vehicle := range s.vehicles {
		if vehicle.Speed > 80 {
			// 检查是否已有未解决的超速告警
			hasAlert := false
			for _, alert := range s.alerts {
				if alert.VehicleID == vehicle.VehicleID && alert.AlertType == "overspeed" && !alert.Resolved {
					hasAlert = true
					break
				}
			}

			if !hasAlert {
				alert := models.TrafficAlert{
					ID:         uint(len(s.alerts) + 1),
					AlertType:  "overspeed",
					VehicleID:  vehicle.VehicleID,
					AlertValue: vehicle.Speed,
					Message:    "车辆" + vehicle.VehicleID + "超速行驶",
					Severity:   "high",
					Resolved:   false,
					Timestamp:  time.Now(),
				}
				s.alerts = append(s.alerts, alert)
			}
		}
	}

	// 限制告警数量
	if len(s.alerts) > 20 {
		s.alerts = s.alerts[len(s.alerts)-20:]
	}
}
