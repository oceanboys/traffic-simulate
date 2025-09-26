<template>
  <div class="alerts-page">
    <el-card style="margin-top: 20px;">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon total">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.totalAlerts }}</div>
            <div class="stat-label">总告警数</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon active">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.activeAlerts }}</div>
            <div class="stat-label">活跃告警</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon resolved">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.resolvedAlerts }}</div>
            <div class="stat-label">已处理</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon high">
            <el-icon><WarningFilled /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.highSeverityAlerts }}</div>
            <div class="stat-label">高危告警</div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="main-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="24"><Warning /></el-icon>
            <span>交通告警列表</span>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="refreshAlerts" class="action-btn">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="success" @click="markAllResolved" class="action-btn">
              <el-icon><Check /></el-icon>
              全部处理
            </el-button>
            <el-button type="danger" @click="clearAllAlerts" class="action-btn">
              <el-icon><Delete /></el-icon>
              清除全部
            </el-button>
          </div>
        </div>
      </template>

      <div v-if="alerts.length === 0" class="no-data">
        <el-icon size="80"><Warning /></el-icon>
        <h3>暂无告警</h3>
        <p>当前没有交通告警信息</p>
      </div>

      <el-table v-else :data="alerts" style="width: 100%" height="600" stripe>
        <el-table-column prop="alertType" label="告警类型" width="120" align="center">
          <template #default="scope">
            <div class="alert-type-cell">
              <el-icon>
                <component :is="getAlertTypeIcon(scope.row.alertType)" />
              </el-icon>
              <span>{{ getAlertTypeText(scope.row.alertType) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="message" label="告警信息" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.message }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="vehicleId" label="车辆ID" width="120" align="center">
          <template #default="scope">
            <span>{{ scope.row.vehicleId }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="severity" label="严重程度" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getSeverityTagType(scope.row.severity) as any" size="small">
              {{ getSeverityText(scope.row.severity) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="location" label="位置" width="150" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.location || '未知' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="timestamp" label="时间" width="180" align="center">
          <template #default="scope">
            <span>{{ formatTime(scope.row.timestamp) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="resolved" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.resolved ? 'success' : 'danger'" size="small">
              {{ scope.row.resolved ? '已处理' : '未处理' }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 填充列，消除右侧空白 -->
        <el-table-column label="" min-width="50">
          <template #default>
            <span></span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="scope">
            <el-button
              size="small"
              :type="scope.row.resolved ? 'success' : 'primary'"
              @click="toggleResolve(scope.row)"
            >
              {{ scope.row.resolved ? '已处理' : '处理' }}
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="deleteAlert(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Warning, Opportunity, InfoFilled, Check, Clock, WarningFilled, Refresh, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { trafficAPI } from '../utils/api'

const alerts = ref([])
const loading = ref(false)

// 从API获取告警数据
const fetchAlerts = async () => {
  loading.value = true
  try {
    const response = await trafficAPI.getAlerts()
    alerts.value = response.data || []
  } catch (error) {
    console.error('获取告警数据失败:', error)
    ElMessage.error('获取告警数据失败')
  } finally {
    loading.value = false
  }
}

// 统计数据
const stats = computed(() => {
  const totalAlerts = alerts.value.length
  const activeAlerts = alerts.value.filter(a => !a.resolved).length
  const resolvedAlerts = alerts.value.filter(a => a.resolved).length
  const highSeverityAlerts = alerts.value.filter(a => a.severity === 'high' && !a.resolved).length

  return {
    totalAlerts,
    activeAlerts,
    resolvedAlerts,
    highSeverityAlerts
  }
})

// 刷新告警
const refreshAlerts = () => {
  fetchAlerts()
}

// 全部处理
const markAllResolved = () => {
  alerts.value.forEach(alert => {
    alert.resolved = true
  })
  ElMessage.success('所有告警已标记为已处理')
}

// 清除全部
const clearAllAlerts = () => {
  alerts.value = []
  ElMessage.success('所有告警已清除')
}

// 切换处理状态
const toggleResolve = (alert: any) => {
  alert.resolved = !alert.resolved
  ElMessage.success(alert.resolved ? '告警已处理' : '告警已重新激活')
}

// 删除告警
const deleteAlert = (alert: any) => {
  const index = alerts.value.findIndex(a => a.id === alert.id)
  if (index > -1) {
    alerts.value.splice(index, 1)
    ElMessage.success('告警已删除')
  }
}

// 获取告警类型图标
const getAlertTypeIcon = (type: string) => {
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

// 获取严重程度标签类型
const getSeverityTagType = (severity: string) => {
  switch (severity) {
    case 'high': return 'danger'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'info'
  }
}

// 获取严重程度文本
const getSeverityText = (severity: string) => {
  switch (severity) {
    case 'high': return '高危'
    case 'medium': return '中等'
    case 'low': return '低危'
    default: return '未知'
  }
}

// 格式化时间
const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchAlerts()
})
</script>

<style scoped>
.alerts-page {
  padding: 15px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

.main-card {
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-bottom: 15px;
  transition: all 0.3s ease;
}

.main-card:hover {
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

.alert-type-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
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
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.stat-icon.resolved {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-icon.high {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
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

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

.el-table tbody tr:hover > td {
  background-color: #f5f7fa !important;
}
</style>
