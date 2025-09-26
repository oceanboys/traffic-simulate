<template>
  <div class="map-container">
    <div ref="mapContainer" class="map"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

interface Props {
  longitude: number
  latitude: number
  zoom?: number
  markers?: Array<{
    longitude: number
    latitude: number
    title?: string
    description?: string
  }>
}

const props = withDefaults(defineProps<Props>(), {
  zoom: 13
})

const mapContainer = ref<HTMLElement>()
let map: any = null

// 初始化地图
const initMap = () => {
  if (!mapContainer.value) return

  // 创建地图
  map = L.map(mapContainer.value).setView([props.latitude, props.longitude], props.zoom)

  // 添加OpenStreetMap图层
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: ' OpenStreetMap contributors',
    maxZoom: 19
  }).addTo(map)

  // 添加标记
  if (props.markers && props.markers.length > 0) {
    props.markers.forEach(marker => {
      const popupContent = '<div><h4>' + (marker.title || '位置') + '</h4><p>' + (marker.description || '') + '</p><p>坐标: ' + marker.latitude + ', ' + marker.longitude + '</p></div>'
      L.marker([marker.latitude, marker.longitude])
        .addTo(map)
        .bindPopup(popupContent)
    })
  } else {
    // 添加默认标记
    const popupContent = '<div><h4>当前位置</h4><p>坐标: ' + props.latitude + ', ' + props.longitude + '</p></div>'
    L.marker([props.latitude, props.longitude])
      .addTo(map)
      .bindPopup(popupContent)
  }
}

// 更新地图中心
const updateMapCenter = () => {
  if (map) {
    map.setView([props.latitude, props.longitude], props.zoom)
  }
}

// 监听属性变化
watch(() => [props.longitude, props.latitude], updateMapCenter)

onMounted(() => {
  // 动态加载Leaflet CSS和JS
  const loadLeaflet = () => {
    return new Promise((resolve) => {
      // 检查是否已加载
      if (window.L) {
        resolve(true)
        return
      }

      // 加载CSS
      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
      document.head.appendChild(link)

      // 加载JS
      const script = document.createElement('script')
      script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
      script.onload = () => resolve(true)
      document.head.appendChild(script)
    })
  }

  loadLeaflet().then(() => {
    initMap()
  })
})

onUnmounted(() => {
  if (map) {
    map.remove()
  }
})
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 100%;
}

.map {
  width: 100%;
  height: 100%;
  border-radius: 8px;
}
</style>
