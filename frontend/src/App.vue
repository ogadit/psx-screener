<script setup>
import { ref, computed, onMounted } from 'vue'
import { GetData } from '../wailsjs/go/main/App'
import FinancialTable from './components/FinancialTable.vue'

const tickerInput = ref('')
const results = ref([])   // [{ ticker, data, error }]
const selectedTicker = ref('')
const loading = ref(false)
const inputRef = ref(null)

const statementHeaders = ['TTM', 'FY 2025', 'FY 2024', 'FY 2023', 'FY 2022', 'FY 2021']
const ratioHeaders = ['Current', 'FY 2025', 'FY 2024', 'FY 2023', 'FY 2022', 'FY 2021']

const current = computed(() => results.value.find(r => r.ticker === selectedTicker.value) ?? null)

onMounted(() => inputRef.value?.focus())

async function fetchAll() {
  const tickers = [...new Set(
    tickerInput.value.split(',').map(t => t.trim().toUpperCase()).filter(Boolean)
  )]
  if (!tickers.length) return

  loading.value = true
  results.value = []
  selectedTicker.value = ''

  const settled = await Promise.allSettled(
    tickers.map(t => GetData(t).then(data => ({ ticker: t, data, error: null })))
  )

  results.value = settled.map((r, i) =>
    r.status === 'fulfilled'
      ? r.value
      : { ticker: tickers[i], data: null, error: r.reason instanceof Error ? r.reason.message : String(r.reason) }
  )

  const first = results.value.find(r => r.data)
  selectedTicker.value = (first ?? results.value[0]).ticker

  loading.value = false
}
</script>

<template>
  <div>
    <div class="header">
      <span class="header-logo">PSX Screener</span>
      <div class="search-bar">
        <input ref="inputRef" v-model="tickerInput" placeholder="e.g. FFC, ENGRO, LUCK" @keyup.enter="fetchAll" />
        <button @click="fetchAll" :disabled="loading">
          {{ loading ? 'Fetching...' : 'Fetch' }}
        </button>
      </div>
    </div>

    <p v-if="loading" class="status-loading">● Fetching data...</p>

    <template v-if="!loading && results.length">
      <div v-if="results.length > 1" class="stock-selector">
        <select v-model="selectedTicker">
          <option v-for="r in results" :key="r.ticker" :value="r.ticker">
            {{ r.data ? r.data.companyName : r.ticker }}{{ r.error ? ' — failed' : '' }}
          </option>
        </select>
      </div>

      <p v-if="current?.error" class="status-error">{{ current.error }}</p>

      <div v-if="current?.data">
        <div class="company-header">
          <span class="company-name">{{ current.data.companyName }}</span>
          <span class="company-currency">PKR</span>
          <span class="company-price">{{ current.data.currentPrice.toLocaleString() }}</span>
        </div>
        <div class="tables-grid">
          <FinancialTable title="Income Statement" :items="current.data.incomeStatement" :headers="statementHeaders" />
          <FinancialTable title="Balance Sheet" :items="current.data.balanceSheet" :headers="statementHeaders" />
          <FinancialTable title="Ratios" :items="current.data.ratios" :headers="ratioHeaders" />
        </div>
      </div>
    </template>
  </div>
</template>
