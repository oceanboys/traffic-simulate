package services

import (
	"backend/models"
	"backend/repositories"
	"time"
)

// GPSService GPS服务
type GPSService struct {
	gpsRepo *repositories.GPSRepository
}

func NewGPSService() *GPSService {
	return &GPSService{
		gpsRepo: repositories.NewGPSRepository(),
	}
}

func (s *GPSService) CreateGPSData(gpsData *models.GPSData) error {
	return s.gpsRepo.Create(gpsData)
}

func (s *GPSService) GetRecentGPSData(limit, minutes int) ([]models.GPSData, error) {
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	return s.gpsRepo.FindRecent(limit, since)
}

func (s *GPSService) GetGPSDataByVehicle(vehicleId string, limit int) ([]models.GPSData, error) {
	return s.gpsRepo.FindByVehicle(vehicleId, limit)
}

func (s *GPSService) GetGPSDataByRoad(roadId uint, minutes int) ([]models.GPSData, error) {
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	return s.gpsRepo.FindByRoad(roadId, since)
}

func (s *GPSService) GetRealTimeStats() (map[string]interface{}, error) {
	// 实现实时统计逻辑
	stats := map[string]interface{}{
		"totalVehicles":   0,
		"averageSpeed":    0,
		"congestionLevel": 0,
		"activeAlerts":    0,
	}

	// 这里可以添加实际的统计计算逻辑
	return stats, nil
}
