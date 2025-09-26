<template>
  <div id="app">
    <div class="app-header">
      <div class="header-content">
        <div class="logo">
          <el-icon size="32"><Van /></el-icon>
          <h1>交通数据模拟平台</h1>
        </div>
      </div>
    </div>

    <div class="app-container">
      <nav class="sidebar">
        <div class="nav-menu">
          <router-link to="/dashboard" class="nav-item">
            <el-icon><DataBoard /></el-icon>
            <span>实时监控</span>
          </router-link>
          <router-link to="/roads" class="nav-item">
            <el-icon><Van /></el-icon>
            <span>路段管理</span>
          </router-link>
          <router-link to="/gps" class="nav-item">
            <el-icon><Location /></el-icon>
            <span>GPS数据</span>
          </router-link>
          <router-link to="/alerts" class="nav-item">
            <el-icon><Warning /></el-icon>
            <span>告警管理</span>
          </router-link>
        </div>
      </nav>

      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Van, DataBoard, Location, Warning } from '@element-plus/icons-vue'
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.app-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  padding: 15px 30px;
  width: 100%;
  margin: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 15px;
}

.logo h1 {
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.app-container {
  display: flex;
  width: 100%;
  margin: 0;
  min-height: calc(100vh - 80px);
}

.sidebar {
  width: 250px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 0;
  margin: 0;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 30px 0;
}

.nav-menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 20px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px 20px;
  text-decoration: none;
  color: #606266;
  border-radius: 12px;
  transition: all 0.3s ease;
  font-weight: 500;
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(64, 158, 255, 0.1), transparent);
  transition: left 0.5s;
}

.nav-item:hover::before {
  left: 100%;
}

.nav-item:hover {
  background: linear-gradient(135deg, #409eff, #67c23a);
  color: white;
  transform: translateX(5px);
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.main-content {
  flex: 1;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 0;
  margin: 0;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-container {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    margin: 0;
    border-radius: 0;
  }

  .main-content {
    margin: 0;
    border-radius: 0;
  }

  .nav-menu {
    flex-direction: row;
    overflow-x: auto;
    padding: 10px 20px;
  }

  .nav-item {
    white-space: nowrap;
    min-width: 120px;
  }
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #409eff, #67c23a);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #67c23a, #409eff);
}
</style>

