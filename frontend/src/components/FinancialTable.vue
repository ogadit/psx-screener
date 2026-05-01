<script setup>
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
</script>

<template>
	<div class="fin-table-wrap">
		<div class="fin-table-title">{{ title }}</div>
		<table>
			<thead>
				<tr>
					<th></th>
					<th v-for="h in headers" :key="h">{{ h }}</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="item in items" :key="item.label">
					<td>{{ item.label }}</td>
					<td v-for="h in headers" :key="h" :class="{ na: fmt(item, h) === 'N/A' }">
						{{ fmt(item, h) }}
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
