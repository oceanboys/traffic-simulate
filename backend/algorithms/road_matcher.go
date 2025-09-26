package algorithms

import (
	"backend/models"
	"backend/repositories"
	"math"
)

// RoadMatcher 路段匹配算法
type RoadMatcher struct {
	roadRepo *repositories.RoadRepository
	roads    []models.RoadSegment
}

// NewRoadMatcher 创建路段匹配器
func NewRoadMatcher() *RoadMatcher {
	return &RoadMatcher{
		roadRepo: repositories.NewRoadRepository(),
	}
}

// LoadRoads 加载路段数据
func (rm *RoadMatcher) LoadRoads() error {
	roads, err := rm.roadRepo.GetAll()
	if err != nil {
		return err
	}
	rm.roads = roads
	return nil
}

// FindNearestRoad 查找最近路段
func (rm *RoadMatcher) FindNearestRoad(lng, lat float64) (*models.RoadSegment, float64) {
	if len(rm.roads) == 0 {
		return nil, math.Inf(1)
	}

	var nearestRoad *models.RoadSegment
	minDistance := math.Inf(1)

	for i := range rm.roads {
		distance := rm.calculateDistanceToRoad(lng, lat, &rm.roads[i])
		if distance < minDistance {
			minDistance = distance
			nearestRoad = &rm.roads[i]
		}
	}

	return nearestRoad, minDistance
}

// calculateDistanceToRoad 计算到路段的距离
func (rm *RoadMatcher) calculateDistanceToRoad(lng, lat float64, road *models.RoadSegment) float64 {
	// 计算点到线段的距离
	return rm.pointToLineDistance(lng, lat,
		road.StartLng, road.StartLat,
		road.EndLng, road.EndLat)
}

// pointToLineDistance 计算点到线段的距离
func (rm *RoadMatcher) pointToLineDistance(px, py, x1, y1, x2, y2 float64) float64 {
	// 使用点到直线距离公式
	A := y2 - y1
	B := x1 - x2
	C := x2*y1 - x1*y2

	distance := math.Abs(A*px+B*py+C) / math.Sqrt(A*A+B*B)

	// 检查点是否在线段范围内
	if rm.isPointInSegment(px, py, x1, y1, x2, y2) {
		return distance
	}

	// 如果不在线段范围内，计算到两个端点的距离
	dist1 := rm.haversineDistance(px, py, x1, y1)
	dist2 := rm.haversineDistance(px, py, x2, y2)

	if dist1 < dist2 {
		return dist1
	}
	return dist2
}

// isPointInSegment 检查点是否在线段范围内
func (rm *RoadMatcher) isPointInSegment(px, py, x1, y1, x2, y2 float64) bool {
	// 检查点是否在线段的矩形范围内
	minX, maxX := math.Min(x1, x2), math.Max(x1, x2)
	minY, maxY := math.Min(y1, y2), math.Max(y1, y2)

	return px >= minX && px <= maxX && py >= minY && py <= maxY
}

// haversineDistance 使用Haversine公式计算两点间距离
func (rm *RoadMatcher) haversineDistance(lng1, lat1, lng2, lat2 float64) float64 {
	const R = 6371 // 地球半径（公里）

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// FindRoadsInRadius 查找指定半径内的路段
func (rm *RoadMatcher) FindRoadsInRadius(lng, lat, radiusKm float64) []*models.RoadSegment {
	var roadsInRadius []*models.RoadSegment

	for i := range rm.roads {
		distance := rm.calculateDistanceToRoad(lng, lat, &rm.roads[i])
		if distance <= radiusKm {
			roadsInRadius = append(roadsInRadius, &rm.roads[i])
		}
	}

	return roadsInRadius
}

// GetRoadDirection 获取路段方向
func (rm *RoadMatcher) GetRoadDirection(road *models.RoadSegment) float64 {
	deltaLng := road.EndLng - road.StartLng
	deltaLat := road.EndLat - road.StartLat

	// 计算方向角（度数）
	direction := math.Atan2(deltaLng, deltaLat) * 180 / math.Pi
	if direction < 0 {
		direction += 360
	}

	return direction
}

// IsVehicleOnRoad 判断车辆是否在路段上
func (rm *RoadMatcher) IsVehicleOnRoad(lng, lat float64, road *models.RoadSegment, tolerance float64) bool {
	distance := rm.calculateDistanceToRoad(lng, lat, road)
	return distance <= tolerance
}
