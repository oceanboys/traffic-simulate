<template>
  <div class="Vans-page">
    <div class="page-header">
      <h2>路段管理系统</h2>
      <p>管理城市道路网络，监控交通流量和路况信息</p>
    </div>
    
    <el-card class="main-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="24"><Van /></el-icon>
            <span>路段数据管理</span>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="refreshRoads" class="action-btn">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="success" @click="showAddDialog = true" class="action-btn">
              <el-icon><Plus /></el-icon>
              添加路段
            </el-button>
            <el-button type="info" @click="exportData" class="action-btn">
              <el-icon><Download /></el-icon>
              导出数据
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="Vans-content">
        <div class="search-bar">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索路段名称、编号或位置"
            class="search-input"
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="selectedType" placeholder="路段类型" class="filter-select">
            <el-option label="全部类型" value="" />
            <el-option label="高速公路" value="highway" />
            <el-option label="城市道路" value="urban" />
            <el-option label="乡村道路" value="rural" />
          </el-select>
          <el-select v-model="selectedStatus" placeholder="路段状态" class="filter-select">
            <el-option label="全部状态" value="" />
            <el-option label="正常" value="normal" />
            <el-option label="维护中" value="maintenance" />
            <el-option label="封闭" value="closed" />
          </el-select>
        </div>
        
        <div class="Vans-list">
          <div v-if="filteredVans.length === 0" class="no-data">
            <el-icon size="80"><Van /></el-icon>
            <h3>暂无路段数据</h3>
            <p>当前没有找到符合条件的路段信息</p>
          </div>
          
          <el-table v-else :data="filteredVans" style="width: 100%" stripe>
            <el-table-column prop="name" label="路段名称" min-width="150" show-overflow-tooltip>
              <template #default="scope">
                <div class="road-name">
                  <el-icon><Location /></el-icon>
                  <span>{{ scope.row.name }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="code" label="路段编号" width="120" align="center">
              <template #default="scope">
                <span>{{ scope.row.code }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="startPoint" label="起点" min-width="120" show-overflow-tooltip>
              <template #default="scope">
                <span>{{ scope.row.startPoint }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="endPoint" label="终点" min-width="120" show-overflow-tooltip>
              <template #default="scope">
                <span>{{ scope.row.endPoint }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="length" label="长度" width="100" align="center">
              <template #default="scope">
                <span>{{ scope.row.length }} km</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="speedLimit" label="限速" width="100" align="center">
              <template #default="scope">
                <span>{{ scope.row.speedLimit }} km/h</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="type" label="类型" width="120" align="center">
              <template #default="scope">
                <el-tag :type="getTypeTagType(scope.row.type)" size="small">
                  {{ getTypeText(scope.row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status) as any" size="small">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="currentTraffic" label="当前流量" width="120" align="center">
              <template #default="scope">
                <span>{{ scope.row.currentTraffic }} 辆/h</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="congestionLevel" label="拥堵程度" width="120" align="center">
              <template #default="scope">
                <el-tag :type="getCongestionTagType(scope.row.congestionLevel) as any" size="small">
                  {{ getCongestionText(scope.row.congestionLevel) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <!-- 填充列，消除右侧空白 -->
            <el-table-column label="" min-width="50">
              <template #default>
                <span></span>
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="200" align="center" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="editRoad(scope.row)">编辑</el-button>
                <el-button size="small" type="success" @click="viewDetails(scope.row)">详情</el-button>
                <el-button size="small" type="danger" @click="deleteRoad(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
    
    <!-- 添加路段对话框 -->
    <el-dialog v-model="showAddDialog" title="添加路段" width="600px">
      <el-form :model="newRoad" label-width="100px">
        <el-form-item label="路段名称" required>
          <el-input v-model="newRoad.name" placeholder="请输入路段名称" />
        </el-form-item>
        <el-form-item label="路段编号" required>
          <el-input v-model="newRoad.code" placeholder="请输入路段编号" />
        </el-form-item>
        <el-form-item label="起点" required>
          <el-input v-model="newRoad.startPoint" placeholder="请输入起点" />
        </el-form-item>
        <el-form-item label="终点" required>
          <el-input v-model="newRoad.endPoint" placeholder="请输入终点" />
        </el-form-item>
        <el-form-item label="长度(km)" required>
          <el-input-number v-model="newRoad.length" :min="0" :precision="1" />
        </el-form-item>
        <el-form-item label="限速(km/h)" required>
          <el-input-number v-model="newRoad.speedLimit" :min="0" :max="200" />
        </el-form-item>
        <el-form-item label="路段类型" required>
          <el-select v-model="newRoad.type" placeholder="请选择路段类型">
            <el-option label="高速公路" value="highway" />
            <el-option label="城市道路" value="urban" />
            <el-option label="乡村道路" value="rural" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="newRoad.description" type="textarea" placeholder="请输入路段描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addRoad">确定</el-button>
      </template>
    </el-dialog>
    
    <!-- 编辑路段对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑路段" width="600px">
      <el-form :model="editingRoad" label-width="100px">
        <el-form-item label="路段名称" required>
          <el-input v-model="editingRoad.name" placeholder="请输入路段名称" />
        </el-form-item>
        <el-form-item label="路段编号" required>
          <el-input v-model="editingRoad.code" placeholder="请输入路段编号" />
        </el-form-item>
        <el-form-item label="起点" required>
          <el-input v-model="editingRoad.startPoint" placeholder="请输入起点" />
        </el-form-item>
        <el-form-item label="终点" required>
          <el-input v-model="editingRoad.endPoint" placeholder="请输入终点" />
        </el-form-item>
        <el-form-item label="长度(km)" required>
          <el-input-number v-model="editingRoad.length" :min="0" :precision="1" />
        </el-form-item>
        <el-form-item label="限速(km/h)" required>
          <el-input-number v-model="editingRoad.speedLimit" :min="0" :max="200" />
        </el-form-item>
        <el-form-item label="路段类型" required>
          <el-select v-model="editingRoad.type" placeholder="请选择路段类型">
            <el-option label="高速公路" value="highway" />
            <el-option label="城市道路" value="urban" />
            <el-option label="乡村道路" value="rural" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editingRoad.description" type="textarea" placeholder="请输入路段描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="updateRoad">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Location, DataLine, Check, Tools, Warning, Search, Van, Plus, 
  Refresh, Download, Edit, View, Switch, Delete, DataBoard
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { roadsAPI } from '../utils/api'

const Vans = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const selectedType = ref('')
const selectedStatus = ref('')
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const editingRoad = ref({})

// 新增路段表单
const newRoad = ref({
  name: '',
  code: '',
  startPoint: '',
  endPoint: '',
  length: 0,
  speedLimit: 60,
  type: 'urban',
  description: ''
})

// 从API获取路段数据
const fetchRoads = async () => {
  loading.value = true
  try {
    const response = await roadsAPI.getAllRoads()
    Vans.value = response.data || []
  } catch (error) {
    console.error('获取路段数据失败:', error)
    ElMessage.error('获取路段数据失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshRoads = () => {
  fetchRoads()
}

// 添加路段
const addRoad = async () => {
  try {
    await roadsAPI.createRoad(newRoad.value)
    ElMessage.success('路段添加成功')
    showAddDialog.value = false
    newRoad.value = {
      name: '',
      code: '',
      startPoint: '',
      endPoint: '',
      length: 0,
      speedLimit: 60,
      type: 'urban',
      description: ''
    }
    await fetchRoads()
  } catch (error) {
    console.error('添加路段失败:', error)
    ElMessage.error('添加路段失败')
  }
}

// 编辑路段
const editRoad = (road: any) => {
  editingRoad.value = { ...road }
  showEditDialog.value = true
}

// 更新路段
const updateRoad = async () => {
  try {
    await roadsAPI.updateRoad(editingRoad.value.id, editingRoad.value)
    ElMessage.success('路段更新成功')
    showEditDialog.value = false
    await fetchRoads()
  } catch (error) {
    console.error('更新路段失败:', error)
    ElMessage.error('更新路段失败')
  }
}

// 删除路段
const deleteRoad = async (road: any) => {
  try {
    await roadsAPI.deleteRoad(road.id)
    ElMessage.success('路段删除成功')
    await fetchRoads()
  } catch (error) {
    console.error('删除路段失败:', error)
    ElMessage.error('删除路段失败')
  }
}

// 查看详情
const viewDetails = (road: any) => {
  ElMessage.info('查看路段详情: ' + road.name)
}

// 导出数据
const exportData = () => {
  ElMessage.success('数据导出功能开发中')
}

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在computed中处理
}

// 搜索过滤
const filteredVans = computed(() => {
  let result = Vans.value
  
  if (searchKeyword.value) {
    result = result.filter(road => 
      road.name.includes(searchKeyword.value) ||
      road.code.includes(searchKeyword.value) ||
      road.startPoint.includes(searchKeyword.value) ||
      road.endPoint.includes(searchKeyword.value)
    )
  }
  
  if (selectedType.value) {
    result = result.filter(road => road.type === selectedType.value)
  }
  
  if (selectedStatus.value) {
    result = result.filter(road => road.status === selectedStatus.value)
  }
  
  return result
})

// 获取拥堵标签类型
const getCongestionTagType = (level: number) => {
  if (level < 0.3) return 'success'
  if (level < 0.7) return 'warning'
  return 'danger'
}

// 获取拥堵文本
const getCongestionText = (level: number) => {
  if (level < 0.3) return '畅通'
  if (level < 0.7) return '缓慢'
  return '拥堵'
}

// 获取状态标签类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'normal': return 'success'
    case 'maintenance': return 'warning'
    case 'closed': return 'danger'
    default: return 'info'
  }
}

// 获取状态文本
const getStatusText = (status: string) => {
  switch (status) {
    case 'normal': return '正常'
    case 'maintenance': return '维护中'
    case 'closed': return '封闭'
    default: return '未知'
  }
}

// 获取类型标签类型
const getTypeTagType = (type: string) => {
  switch (type) {
    case 'highway': return 'danger'
    case 'urban': return 'primary'
    case 'rural': return 'success'
    default: return 'info'
  }
}

// 获取类型文本
const getTypeText = (type: string) => {
  switch (type) {
    case 'highway': return '高速公路'
    case 'urban': return '城市道路'
    case 'rural': return '乡村道路'
    default: return '未知'
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchRoads()
})
</script>

<style scoped>
.Vans-page {
  padding: 15px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

.page-header {
  text-align: center;
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 15px;
  color: white;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.page-header h2 {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.page-header p {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
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

.Vans-content {
  padding: 20px;
}

.search-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  align-items: center;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.filter-select {
  width: 150px;
}

.Vans-list {
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

.road-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

.el-table tbody tr:hover > td {
  background-color: #f5f7fa !important;
}

@media (max-width: 768px) {
  .Vans-page {
    padding: 10px;
  }
  
  .search-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input, .filter-select {
    width: 100%;
  }
}
</style>
