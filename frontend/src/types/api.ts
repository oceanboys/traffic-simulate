import type { ApiResponse, GPSData, RoadSegment, TrafficAlert, RealTimeStats, SimulationConfig } from './index'

// API请求参数类型
export interface GPSDataParams extends Omit<GPSData, 'id' | 'createdAt'> {}
export interface RoadSegmentParams extends Omit<RoadSegment, 'id' | 'createdAt' | 'updatedAt'> {}

// API响应类型
export type HealthResponse = ApiResponse<{
  status: string
  service: string
  timestamp: string
}>

export type RoadsResponse = ApiResponse<RoadSegment[]>
export type RoadResponse = ApiResponse<RoadSegment>
export type GPSDataResponse = ApiResponse<GPSData>
export type RecentGPSResponse = ApiResponse<GPSData[]>
export type StatsResponse = ApiResponse<RealTimeStats>
