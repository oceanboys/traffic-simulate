// 基础响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// GPS数据接口
export interface GPSData {
  id?: number
  vehicleId: string
  longitude: number
  latitude: number
  speed: number
  direction: number
  timestamp: string
  roadSegmentId?: number
  vehicleType: 'car' | 'bus' | 'truck'
  createdAt?: string
}

// 路段接口
export interface RoadSegment {
  id?: number
  name: string
  startLng: number
  startLat: number
  endLng: number
  endLat: number
  maxSpeed: number
  capacity: number
  length?: number
  roadType: 'highway' | 'urban' | 'rural'
  createdAt?: string
  updatedAt?: string
}

// 交通告警接口
export interface TrafficAlert {
  id?: number
  alertType: 'overspeed' | 'congestion' | 'accident'
  vehicleId?: string
  roadSegmentId: number
  alertValue?: number
  message: string
  severity: 'low' | 'medium' | 'high'
  resolved: boolean
  timestamp: string
  createdAt?: string
}

// 实时统计接口
export interface RealTimeStats {
  totalVehicles: number
  averageSpeed: number
  congestionLevel: number
  activeAlerts: number
}

// 模拟配置接口
export interface SimulationConfig {
  id?: number
  name: string
  areaLngMin: number
  areaLngMax: number
  areaLatMin: number
  areaLatMax: number
  vehicleCount: number
  speedMin: number
  speedMax: number
  intervalMS: number
  isActive: boolean
}
