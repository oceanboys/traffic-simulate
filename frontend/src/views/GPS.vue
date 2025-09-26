<template>
  <div class="gps-page">
    <el-card class="main-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="24"><Location /></el-icon>
            <span>车辆GPS数据管理</span>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="refreshGPSData" class="action-btn">
              <el-icon><Search /></el-icon>
              刷新
            </el-button>
            <el-button type="success" @click="exportGPSData" class="action-btn">
              <el-icon><Van /></el-icon>
              导出数据
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="gps-content">
        <div class="search-bar">
          <div class="search-row">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索车辆ID或位置"
              class="search-input"
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select v-model="selectedVehicle" placeholder="选择车辆" class="filter-select">
              <el-option label="全部车辆" value="" />
              <el-option 
                v-for="vehicle in vehicles" 
                :key="vehicle.id" 
                :label="vehicle.vehicleId" 
                :value="vehicle.vehicleId" 
              />
            </el-select>
          </div>
          <div class="search-row">
            <el-date-picker
              v-model="dateRange"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              class="date-picker"
            />
            <div class="search-actions">
              <el-button type="primary" @click="refreshGPSData">
                <el-icon><Search /></el-icon>
                刷新
              </el-button>
              <el-button type="success" @click="exportGPSData">
                <el-icon><Van /></el-icon>
                导出
              </el-button>
            </div>
          </div>
        </div>
        
        <div class="gps-list">
          <div v-if="filteredGPSData.length === 0" class="no-data">
            <el-icon size="80"><Location /></el-icon>
            <h3>暂无GPS数据</h3>
            <p>当前没有找到符合条件的GPS数据</p>
          </div>
          
          <el-table v-else :data="filteredGPSData" style="width: 100%" stripe>
            <el-table-column prop="vehicleId" label="车辆ID" min-width="120" align="center">
              <template #default="scope">
                <div class="vehicle-id-cell">
                  <el-icon><Van /></el-icon>
                  <span>{{ scope.row.vehicleId }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="location" label="位置信息" min-width="300" show-overflow-tooltip>
              <template #default="scope">
                <div class="location-cell">
                  <div class="coordinates">{{ scope.row.longitude }}, {{ scope.row.latitude }}</div>
                  <div class="location-name">{{ scope.row.location || '未知位置' }}</div>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="speed" label="速度" width="120" align="center">
              <template #default="scope">
                <span>{{ scope.row.speed }} km/h</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="direction" label="方向" width="120" align="center">
              <template #default="scope">
                <span>{{ scope.row.direction }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="accuracy" label="精度" width="120" align="center">
              <template #default="scope">
                <span>{{ scope.row.accuracy }}m</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="timestamp" label="时间" min-width="180" align="center">
              <template #default="scope">
                <span>{{ formatTime(scope.row.timestamp) }}</span>
              </template>
            </el-table-column>
            
            <!-- 填充列，消除右侧空白 -->
            <el-table-column label="" min-width="50">
              <template #default>
                <span></span>
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="180" align="center" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="viewOnMap(scope.row)">地图查看</el-button>
                <el-button size="small" type="danger" @click="deleteGPSData(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
    
    <el-card class="stats-card" shadow="hover" style="margin-top: 20px;">
      <template #header>
        <span>GPS数据统计</span>
      </template>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon total">
            <el-icon><Van /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ gpsData.length }}</div>
            <div class="stat-label">总数据量</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon active">
            <el-icon><DataLine /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ vehicles.length }}</div>
            <div class="stat-label">活跃车辆</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon resolved">
            <el-icon><Aim /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ averageSpeed.toFixed(1) }}</div>
            <div class="stat-label">平均速度</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon high">
            <el-icon><Location /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ averageAccuracy.toFixed(1) }}</div>
            <div class="stat-label">平均精度</div>
          </div>
        </div>
      </div>
    </el-card>
    
    <!-- 地图弹窗 -->
    <el-dialog v-model="showMap" title="地图查看" width="80%" :before-close="closeMap">
      <div class="map-dialog-content">
        <div class="map-dialog-header">
          <div class="vehicle-info">
            <h3>车辆信息</h3>
            <div class="info-list">
              <div class="info-item">
                <span class="label">车辆ID:</span>
                <span class="value">{{ selectedGPS.vehicleId }}</span>
              </div>
              <div class="info-item">
                <span class="label">位置:</span>
                <span class="value">{{ selectedGPS.location || '未知位置' }}</span>
              </div>
              <div class="info-item">
                <span class="label">坐标:</span>
                <span class="value">{{ selectedGPS.longitude }}, {{ selectedGPS.latitude }}</span>
              </div>
              <div class="info-item">
                <span class="label">速度:</span>
                <span class="value">{{ selectedGPS.speed }} km/h</span>
              </div>
              <div class="info-item">
                <span class="label">方向:</span>
                <span class="value">{{ selectedGPS.direction }}</span>
              </div>
              <div class="info-item">
                <span class="label">精度:</span>
                <span class="value">{{ selectedGPS.accuracy }}m</span>
              </div>
              <div class="info-item">
                <span class="label">时间:</span>
                <span class="value">{{ formatTime(selectedGPS.timestamp) }}</span>
              </div>
            </div>
            <div class="map-actions">
              <el-button type="primary" @click="openExternalMap" icon="Location">
                在外部地图中查看
              </el-button>
            </div>
          </div>
        </div>
        
        <div class="map-display">
          <el-card shadow="hover" style="height: 400px;">
            <template #header>
              <span>地图显示</span>
            </template>
            <div class="map-placeholder">
              <MapView 
                v-if="selectedGPS.longitude && selectedGPS.latitude"
                :longitude="selectedGPS.longitude"
                :latitude="selectedGPS.latitude"
                :zoom="15"
                :markers="[{
                  longitude: selectedGPS.longitude,
                  latitude: selectedGPS.latitude,
                  title: selectedGPS.vehicleId,
                  description: selectedGPS.location || '未知位置'
                }]"
              />
            </div>
          </el-card>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Location, Van, DataLine, Aim, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { gpsAPI, trafficAPI } from '../utils/api'
import MapView from '../components/MapView.vue'

const gpsData = ref([])
const vehicles = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const selectedVehicle = ref('')
const dateRange = ref([])
const showMap = ref(false)
const selectedGPS = ref({})

// 从API获取GPS数据
const fetchGPSData = async () => {
  loading.value = true
  try {
    const response = await gpsAPI.getGPSDataByVehicle('all', 100)
    gpsData.value = response.data || []
  } catch (error) {
    console.error('获取GPS数据失败:', error)
    ElMessage.error('获取GPS数据失败')
  } finally {
    loading.value = false
  }
}

// 从API获取车辆列表
const fetchVehicles = async () => {
  try {
    const response = await trafficAPI.getVehicles()
    vehicles.value = response.data || []
  } catch (error) {
    console.error('获取车辆列表失败:', error)
  }
}

// 刷新GPS数据
const refreshGPSData = () => {
  fetchGPSData()
}

// 导出GPS数据
const exportGPSData = () => {
  ElMessage.success('数据导出功能开发中')
}

// 查看地图
const viewOnMap = (gps: any) => {
  selectedGPS.value = gps
  showMap.value = true
}

// 关闭地图
const closeMap = () => {
  showMap.value = false
  selectedGPS.value = {}
}

// 打开外部地图
const openExternalMap = () => {
  if (selectedGPS.value.longitude && selectedGPS.value.latitude) {
    const url = 'https://map.baidu.com/?q=' + selectedGPS.value.latitude + ',' + selectedGPS.value.longitude
    window.open(url, '_blank')
  }
}

// 删除GPS数据
const deleteGPSData = (gps: any) => {
  ElMessage.success('删除功能开发中')
}

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在computed中处理
}

// 过滤GPS数据
const filteredGPSData = computed(() => {
  let filtered = gpsData.value
  
  if (searchKeyword.value) {
    filtered = filtered.filter(gps => 
      gps.vehicleId.includes(searchKeyword.value) ||
      (gps.location && gps.location.includes(searchKeyword.value))
    )
  }
  
  if (selectedVehicle.value) {
    filtered = filtered.filter(gps => gps.vehicleId === selectedVehicle.value)
  }
  
  if (dateRange.value && dateRange.value.length === 2) {
    const [start, end] = dateRange.value
    filtered = filtered.filter(gps => {
      const gpsTime = new Date(gps.timestamp)
      return gpsTime >= start && gpsTime <= end
    })
  }
  
  return filtered
})

// 计算平均速度
const averageSpeed = computed(() => {
  if (gpsData.value.length === 0) return 0
  const total = gpsData.value.reduce((sum, gps) => sum + gps.speed, 0)
  return total / gpsData.value.length
})

// 计算平均精度
const averageAccuracy = computed(() => {
  if (gpsData.value.length === 0) return 0
  const total = gpsData.value.reduce((sum, gps) => sum + gps.accuracy, 0)
  return total / gpsData.value.length
})

// 格式化时间
const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchGPSData()
  fetchVehicles()
})
</script>

<style scoped>
.gps-page {
  padding: 15px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

.main-card, .stats-card {
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-bottom: 15px;
  transition: all 0.3s ease;
}

.main-card:hover, .stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  border-radius: 25px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.gps-content {
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
}

.search-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  align-items: center;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.filter-select {
  width: 150px;
}

.date-picker {
  flex: 1;
  min-width: 300px;
  max-width: 400px;
}

.search-actions {
  display: flex;
  gap: 10px;
}

.gps-list {
  margin-top: 20px;
}

.no-data {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #909399;
  text-align: center;
  background: linear-gradient(135deg, #f8f9fa, #e9ecef);
  border-radius: 15px;
  border: 2px dashed #d1d5db;
}

.no-data h3 {
  margin: 20px 0 10px;
  color: #909399;
  font-size: 24px;
}

.no-data p {
  margin-bottom: 20px;
  font-size: 16px;
}

.vehicle-id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
}

.location-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.coordinates {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
}

.location-name {
  font-weight: 500;
  color: #606266;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

.el-table tbody tr:hover > td {
  background-color: #f5f7fa !important;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  padding: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa, #e9ecef);
  border-radius: 15px;
  transition: all 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  color: white;
  font-size: 20px;
}

.stat-icon.total {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.stat-icon.active {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-icon.resolved {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.stat-icon.high {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.map-dialog-content {
  display: flex;
  gap: 20px;
  height: 500px;
}

.map-dialog-header {
  flex: 1;
  min-width: 300px;
}

.vehicle-info {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 10px;
  height: 100%;
}

.vehicle-info h3 {
  margin: 0 0 20px 0;
  color: #303133;
  font-size: 18px;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #e4e7ed;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  font-weight: 500;
  color: #606266;
}

.value {
  color: #303133;
  font-weight: 600;
}

.map-actions {
  margin-top: 20px;
}

.map-display {
  flex: 2;
  min-width: 400px;
}

.map-placeholder {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

@media (max-width: 768px) {
  .map-dialog-content {
    flex-direction: column;
    height: auto;
  }
  
  .map-display {
    min-width: auto;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
