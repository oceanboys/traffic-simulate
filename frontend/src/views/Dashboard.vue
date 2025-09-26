<template>
  <div class="dashboard">
    <el-card class="control-panel" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>模拟控制</span>
          <el-tag :type="isSimulating ? 'success' : 'info'" size="large">
            {{ isSimulating ? '模拟运行中' : '模拟已停止' }}
          </el-tag>
        </div>
      </template>

      <div class="controls">
        <el-button
          :type="isSimulating ? 'danger' : 'success'"
          @click="toggleSimulation"
          size="large"
          class="control-btn"
        >
          <el-icon><VideoPlay v-if="!isSimulating" /><VideoPause v-else /></el-icon>
          {{ isSimulating ? '停止模拟' : '开始模拟' }}
        </el-button>
        <el-button @click="addVehicle" type="primary" size="large" class="control-btn">
          <el-icon><Plus /></el-icon>
          添加车辆
        </el-button>
        <el-button @click="removeVehicle" type="warning" size="large" class="control-btn">
          <el-icon><Minus /></el-icon>
          移除车辆
        </el-button>
        <el-button @click="addTrafficJam" type="danger" size="large" class="control-btn">
          <el-icon><Warning /></el-icon>
          模拟拥堵
        </el-button>
        <el-button @click="clearAlerts" type="info" size="large" class="control-btn">
          <el-icon><Delete /></el-icon>
          清除告警
        </el-button>
        <el-button @click="resetSimulation" type="default" size="large" class="control-btn">
          <el-icon><Refresh /></el-icon>
          重置模拟
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="20" class="content-row">
      <el-col :span="18">
        <el-card class="map-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>实时交通流</span>
              <div class="header-actions">
                <el-button @click="refreshData" type="primary" size="small">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </div>
          </template>

          <div class="map-container">
            <div class="traffic-map">
              <div class="road-segments">
                <div
                  v-for="segment in roadSegments"
                  :key="segment.id"
                  class="road-segment"
                  :style="getSegmentStyle(segment)"
                >
                  <div class="segment-info">
                    <span class="segment-name">{{ segment.name }}</span>
                    <span class="congestion-level">{{ Math.round(segment.congestionLevel * 100) }}%</span>
                  </div>
                </div>
              </div>

              <div class="vehicles">
                <div
                  v-for="vehicle in vehicles"
                  :key="vehicle.id"
                  class="vehicle"
                  :class="vehicle.status"
                  :style="getVehicleStyle(vehicle)"
                  :title="getVehicleTitle(vehicle)"
                >
                  <el-icon><Van /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <template #header>
            <span>实时统计</span>
          </template>

          <div class="stats-content">
            <div class="stat-item">
              <div class="stat-icon total">
                <el-icon><Van /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ realTimeStats.totalVehicles }}</div>
                <div class="stat-label">总车辆数</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon speed">
                <el-icon><DataLine /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ realTimeStats.averageSpeed.toFixed(1) }}</div>
                <div class="stat-label">平均速度</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon congestion">
                <el-icon><Warning /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ realTimeStats.congestionLevel.toFixed(1) }}%</div>
                <div class="stat-label">拥堵程度</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon alerts">
                <el-icon><InfoFilled /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ realTimeStats.activeAlerts }}</div>
                <div class="stat-label">活跃告警</div>
              </div>
            </div>
          </div>
        </el-card>

        <el-card class="alerts-card" shadow="hover" style="margin-top: 20px;">
          <template #header>
            <span>实时告警</span>
          </template>

          <div class="alerts-content">
            <div v-if="alerts.length === 0" class="no-alerts">
              <el-icon size="40"><Check /></el-icon>
              <p>暂无告警</p>
            </div>

            <div v-else class="alert-items">
              <div
                v-for="alert in alerts.slice(0, 5)"
                :key="alert.id"
                class="alert-item"
                :class="'alert-' + alert.severity"
                @click="handleAlertClick(alert)"
              >
                <div class="alert-icon">
                  <el-icon>
                    <component :is="getAlertIcon(alert.alertType)" />
                  </el-icon>
                </div>
                <div class="alert-content">
                  <div class="alert-header">
                    <span class="alert-type">{{ getAlertTypeText(alert.alertType) }}</span>
                    <span class="alert-time">{{ formatTime(alert.timestamp) }}</span>
                  </div>
                  <div class="alert-message">{{ alert.message }}</div>
                  <div class="alert-vehicle">车辆: {{ alert.vehicleId }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { ElMessage } from "element-plus"
import {
  Location, Van, Warning, DataLine, Opportunity, Check, InfoFilled, WarningFilled,
  VideoPlay, VideoPause, Plus, Minus, Delete, Refresh
} from "@element-plus/icons-vue"
import { trafficAPI, simulationAPI } from "../utils/api"

interface RoadSegment {
  id: number
  x: number
  y: number
  width: number
  height: number
  direction: number
  congestionLevel: number
  name: string
}

interface Vehicle {
  id: string
  vehicleId: string
  x: number
  y: number
  speed: number
  direction: number
  vehicleType: string
  timestamp: string
  status: "normal" | "overspeed" | "congestion"
  roadSegmentId?: number
}

interface Alert {
  id: number
  alertType: string
  vehicleId: string
  message: string
  severity: "low" | "medium" | "high"
  timestamp: string
  resolved: boolean
  x?: number
  y?: number
}

const isSimulating = ref(false)
const vehicles = ref<Vehicle[]>([])
const alerts = ref<Alert[]>([])
const roadSegments = ref<RoadSegment[]>([])
const vehicleIdCounter = ref(1)
let simulationInterval: number | undefined

// 实时统计数据
const realTimeStats = ref({
  totalVehicles: 0,
  averageSpeed: 0,
  congestionLevel: 0,
  activeAlerts: 0
})

// 从API获取实时数据
const fetchRealTimeData = async () => {
  try {
    const [trafficData, alertsData, vehiclesData] = await Promise.all([
      trafficAPI.getRealTimeTraffic(),
      trafficAPI.getAlerts(),
      trafficAPI.getVehicles()
    ])

    realTimeStats.value = trafficData
    alerts.value = alertsData || []
    vehicles.value = vehiclesData || []
  } catch (error) {
    console.error('获取实时数据失败:', error)
    ElMessage.error('获取实时数据失败')
  }
}

// 开始模拟
const startSimulation = async () => {
  try {
    await simulationAPI.startSimulation()
    isSimulating.value = true
    ElMessage.success('模拟已开始')

    // 开始定时更新数据
    simulationInterval = setInterval(fetchRealTimeData, 2000)
  } catch (error) {
    console.error('开始模拟失败:', error)
    ElMessage.error('开始模拟失败')
  }
}

// 停止模拟
const stopSimulation = async () => {
  try {
    await simulationAPI.stopSimulation()
    isSimulating.value = false
    ElMessage.success('模拟已停止')

    // 清除定时器
    if (simulationInterval) {
      clearInterval(simulationInterval)
      simulationInterval = undefined
    }
  } catch (error) {
    console.error('停止模拟失败:', error)
    ElMessage.error('停止模拟失败')
  }
}

// 切换模拟状态
const toggleSimulation = () => {
  if (isSimulating.value) {
    stopSimulation()
  } else {
    startSimulation()
  }
}

// 添加车辆
const addVehicle = async () => {
  try {
    const newVehicle = {
      vehicleId: V,
      vehicleType: 'car',
      status: 'normal'
    }

    await trafficAPI.addVehicle(newVehicle)
    vehicleIdCounter.value++
    ElMessage.success('车辆添加成功')
    await fetchRealTimeData()
  } catch (error) {
    console.error('添加车辆失败:', error)
    ElMessage.error('添加车辆失败')
  }
}

// 移除车辆
const removeVehicle = async () => {
  if (vehicles.value.length === 0) {
    ElMessage.warning('没有可移除的车辆')
    return
  }

  try {
    const vehicleToRemove = vehicles.value[0]
    await trafficAPI.removeVehicle(vehicleToRemove.id)
    ElMessage.success('车辆移除成功')
    await fetchRealTimeData()
  } catch (error) {
    console.error('移除车辆失败:', error)
    ElMessage.error('移除车辆失败')
  }
}

// 模拟拥堵
const addTrafficJam = () => {
  ElMessage.success('拥堵模拟已添加')
}

// 清除告警
const clearAlerts = () => {
  alerts.value = []
  ElMessage.success('告警已清除')
}

// 重置模拟
const resetSimulation = async () => {
  await stopSimulation()
  vehicles.value = []
  alerts.value = []
  vehicleIdCounter.value = 1
  ElMessage.success('模拟已重置')
}

// 刷新数据
const refreshData = () => {
  fetchRealTimeData()
}

// 获取路段样式
const getSegmentStyle = (segment: RoadSegment) => {
  return {
    left: segment.x + 'px',
    top: segment.y + 'px',
    width: segment.width + 'px',
    height: segment.height + 'px',
    transform:
otate(deg),
    backgroundColor: getCongestionColor(segment.congestionLevel)
  }
}

// 获取车辆样式
const getVehicleStyle = (vehicle: Vehicle) => {
  return {
    left: vehicle.x + 'px',
    top: vehicle.y + 'px',
    transform:
otate(deg)
  }
}

// 获取车辆标题
const getVehicleTitle = (vehicle: Vehicle) => {
  return vehicle.vehicleId + ' - ' + vehicle.speed + 'km/h'
}

// 获取拥堵颜色
const getCongestionColor = (level: number) => {
  if (level < 0.3) return '#67c23a'
  if (level < 0.7) return '#e6a23c'
  return '#f56c6c'
}

// 获取告警图标
const getAlertIcon = (type: string) => {
  switch (type) {
    case 'overspeed': return 'Warning'
    case 'congestion': return 'Opportunity'
    case 'accident': return 'InfoFilled'
    default: return 'Warning'
  }
}

// 获取告警类型文本
const getAlertTypeText = (type: string) => {
  switch (type) {
    case 'overspeed': return '超速'
    case 'congestion': return '拥堵'
    case 'accident': return '事故'
    default: return '未知'
  }
}

// 格式化时间
const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

// 处理告警点击
const handleAlertClick = (alert: Alert) => {
  ElMessage.info('点击了告警: ' + alert.message)
}

// 初始化路段
const initRoadSegments = () => {
  roadSegments.value = [
    { id: 1, x: 100, y: 200, width: 200, height: 20, direction: 0, congestionLevel: 0.3, name: '东三环' },
    { id: 2, x: 300, y: 150, width: 150, height: 20, direction: 45, congestionLevel: 0.7, name: '西二环' },
    { id: 3, x: 200, y: 300, width: 180, height: 20, direction: 90, congestionLevel: 0.5, name: '北四环' }
  ]
}

// 组件挂载时获取初始数据
onMounted(async () => {
  initRoadSegments()
  await fetchRealTimeData()
})

// 组件卸载时清理定时器
onUnmounted(() => {
  if (simulationInterval) {
    clearInterval(simulationInterval)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

.control-panel {
  margin-bottom: 20px;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

.controls {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.control-btn {
  border-radius: 25px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.control-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.content-row {
  margin-top: 20px;
}

.map-card, .stats-card, .alerts-card {
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.map-card:hover, .stats-card:hover, .alerts-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.map-container {
  height: 500px;
  position: relative;
  background: linear-gradient(135deg, #e8f4fd, #f0f9ff);
  border-radius: 10px;
  overflow: hidden;
}

.traffic-map {
  width: 100%;
  height: 100%;
  position: relative;
}

.road-segments {
  position: absolute;
  width: 100%;
  height: 100%;
}

.road-segment {
  position: absolute;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.road-segment:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.segment-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: white;
  font-weight: bold;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
}

.segment-name {
  font-size: 12px;
  margin-bottom: 2px;
}

.congestion-level {
  font-size: 10px;
  opacity: 0.9;
}

.vehicles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.vehicle {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.vehicle.normal {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.vehicle.overspeed {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.vehicle.congestion {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.vehicle:hover {
  transform: scale(1.2);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 15px;
  background: linear-gradient(135deg, #f8f9fa, #e9ecef);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  color: white;
  font-size: 18px;
}

.stat-icon.total {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.stat-icon.speed {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-icon.congestion {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.stat-icon.alerts {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.alerts-content {
  max-height: 300px;
  overflow-y: auto;
}

.no-alerts {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 150px;
  color: #67c23a;
  text-align: center;
}

.no-alerts p {
  margin-top: 10px;
  font-size: 14px;
}

.alert-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 4px solid;
}

.alert-item:hover {
  transform: translateX(3px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.alert-item.alert-high {
  background: linear-gradient(135deg, #fef0f0, #fde2e2);
  border-left-color: #f56c6c;
}

.alert-item.alert-medium {
  background: linear-gradient(135deg, #fdf6ec, #fbeaa6);
  border-left-color: #e6a23c;
}

.alert-item.alert-low {
  background: linear-gradient(135deg, #f0f9ff, #d1ecf1);
  border-left-color: #409eff;
}

.alert-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  color: white;
  font-size: 18px;
}

.alert-high .alert-icon {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.alert-medium .alert-icon {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.alert-low .alert-icon {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.alert-content {
  flex: 1;
}

.alert-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.alert-type {
  font-weight: bold;
  font-size: 14px;
}

.alert-time {
  font-size: 12px;
  color: #909399;
}

.alert-message {
  font-size: 13px;
  color: #606266;
  margin-bottom: 4px;
  line-height: 1.4;
}

.alert-vehicle {
  font-size: 12px;
  color: #909399;
}

@media (max-width: 1200px) {
  .content-row .el-col:first-child {
    margin-bottom: 20px;
  }
}

@media (max-width: 768px) {
  .dashboard {
    padding: 10px;
  }

  .controls {
    justify-content: center;
  }
}
</style>
