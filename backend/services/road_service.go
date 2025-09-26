package services

import (
	"backend/models"
	"backend/repositories"
)

// RoadService 路段服务
type RoadService struct {
	roadRepo *repositories.RoadRepository
}

func NewRoadService() *RoadService {
	return &RoadService{
		roadRepo: repositories.NewRoadRepository(),
	}
}

func (s *RoadService) GetAllRoads() ([]models.RoadSegment, error) {
	return s.roadRepo.GetAll()
}

func (s *RoadService) GetRoadByID(id uint) (*models.RoadSegment, error) {
	return s.roadRepo.GetByID(id)
}

func (s *RoadService) CreateRoad(road *models.RoadSegment) error {
	return s.roadRepo.Create(road)
}

func (s *RoadService) UpdateRoad(road *models.RoadSegment) error {
	return s.roadRepo.Update(road)
}

func (s *RoadService) DeleteRoad(id uint) error {
	return s.roadRepo.Delete(id)
}
