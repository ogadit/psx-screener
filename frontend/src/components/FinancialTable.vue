<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: String,
  items: Array,
  headers: Array,
})

function fmt(item, header) {
  const val = item.values?.[header]
  if (val === undefined || val === null) return 'N/A'
  return item.isPercent ? val.toFixed(2) + '%' : val.toLocaleString()
}

function rawRow(item) {
  return [item.label, ...props.headers.map(h => fmt(item, h))].join('\t')
}

const menu = ref({ visible: false, x: 0, y: 0, value: '', item: null })

function onContextMenu(e, item, header) {
  e.preventDefault()
  menu.value = { visible: true, x: e.clientX, y: e.clientY, value: fmt(item, header), item }
}

function copyValue() {
  navigator.clipboard.writeText(menu.value.value)
  menu.value.visible = false
}

function copyRow() {
  navigator.clipboard.writeText(rawRow(menu.value.item))
  menu.value.visible = false
}

function hideMenu() {
  menu.value.visible = false
}
</script>

<template>
  <div class="fin-table-wrap">
    <div class="fin-table-title">{{ title }}</div>
    <div class="table-scroll">
      <table>
        <thead>
          <tr>
            <th></th>
            <th v-for="h in headers" :key="h">{{ h }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.label">
            <td class="row-label">{{ item.label }}</td>
            <td
              v-for="h in headers"
              :key="h"
              :class="{ na: fmt(item, h) === 'N/A' }"
              @contextmenu="onContextMenu($event, item, h)"
            >{{ fmt(item, h) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="menu.visible" class="ctx-backdrop" @click="hideMenu" @contextmenu.prevent="hideMenu" />
    <div v-if="menu.visible" class="ctx-menu" :style="{ top: menu.y + 'px', left: menu.x + 'px' }">
      <div class="ctx-value">{{ menu.value }}</div>
      <button class="ctx-item" @click="copyValue">Copy value</button>
      <button class="ctx-item" @click="copyRow">Copy row</button>
    </div>
  </Teleport>
</template>
