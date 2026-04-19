<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <header class="relative z-20 border-b border-gray-100 bg-white/80 backdrop-blur-xl dark:border-dark-800 dark:bg-dark-950/80">
      <nav class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-8 w-8 overflow-hidden rounded-lg">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-bold text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-2">
          <router-link to="/key-usage"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('home.keyUsage') }}
          </router-link>
          <router-link to="/monitoring"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('admin.monitoring.title') }}
          </router-link>
          <LocaleSwitcher />
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('home.docs') }}
          </a>
          <button @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-50 dark:text-dark-400 dark:hover:bg-dark-800">
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link v-if="isAuthenticated" :to="dashboardPath"
            class="ml-1 inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100">
            {{ t('home.dashboard') }}
            <Icon name="arrowRight" size="xs" :stroke-width="2" />
          </router-link>
          <router-link v-else to="/login"
            class="ml-1 inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100">
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="mx-auto w-full max-w-7xl flex-1 px-6 py-8">
      <div class="space-y-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('admin.monitoring.title') }}</h1>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.description') }}</p>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('admin.monitoring.last24h') }}</span>
            <button @click="refresh" :disabled="loading"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:opacity-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700">
              <svg class="h-4 w-4" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
              </svg>
              {{ t('common.refresh') }}
            </button>
          </div>
        </div>

        <!-- Summary Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-5">
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.todayRequests') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data?.total_requests_today?.toLocaleString() ?? '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-green-600 dark:text-green-400">Success</p>
            <p class="mt-1 text-2xl font-bold text-green-600 dark:text-green-400">{{ data?.success_count_today?.toLocaleString() ?? '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-red-500 dark:text-red-400">Errors</p>
            <p class="mt-1 text-2xl font-bold text-red-600 dark:text-red-400">{{ data?.error_count_today?.toLocaleString() ?? '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.avgLatency') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data ? Math.round(data.avg_latency_ms_today) + 'ms' : '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.totalGroups') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data?.groups?.length ?? 0 }}</p>
          </div>
        </div>

        <!-- Group × Model Matrix -->
        <div v-for="group in groupedModels" :key="group.groupId" class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900 overflow-hidden">
          <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ group.groupName }}</span>
            <span class="rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ group.models.length }} {{ group.models.length === 1 ? 'model' : 'models' }}
            </span>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">Model</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">Requests</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-green-600 dark:text-green-400">Success</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-red-500 dark:text-red-400">Errors</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">Success Rate</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">AVG</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-amber-500 dark:text-amber-400">P95</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">TTFT</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in group.models" :key="m.model" class="border-b border-gray-50 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white font-mono">{{ m.model }}</td>
                  <td class="px-6 py-3 text-center text-sm text-gray-700 dark:text-dark-200">{{ m.request_count.toLocaleString() }}</td>
                  <td class="px-6 py-3 text-center">
                    <span class="text-sm font-semibold" :class="m.success_count > 0 ? 'text-green-600 dark:text-green-400' : 'text-gray-400'">{{ m.success_count }}</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <span :class="m.error_count > 0 ? 'text-sm font-semibold text-red-600 dark:text-red-400' : 'text-sm text-gray-400'">{{ m.error_count }}</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <div class="flex items-center justify-center gap-2">
                      <div class="h-1.5 w-16 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-700">
                        <div class="h-full rounded-full" :class="successRateColor(m.success_count, m.request_count)" :style="{ width: successRateWidth(m.success_count, m.request_count) }"></div>
                      </div>
                      <span class="text-sm" :class="successRateTextColor(m.success_count, m.request_count)">{{ successRate(m.success_count, m.request_count) }}</span>
                    </div>
                  </td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm" :class="m.p95_latency_ms > 10000 ? 'text-amber-600 dark:text-amber-400 font-semibold' : 'text-gray-700 dark:text-dark-200'">{{ formatMs(m.p95_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_ttft) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="data && (!data.group_models || data.group_models.length === 0)" class="rounded-xl border border-gray-200 bg-white py-16 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.noData') }}</p>
        </div>

        <!-- Model Latency Overview -->
        <div v-if="data?.model_latencies?.length" class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900 overflow-hidden">
          <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.monitoring.modelLatency') }}</h2>
            <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('admin.monitoring.modelLatencyHint') }}</span>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.model') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.requests') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-green-600 dark:text-green-400">Success</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-red-500 dark:text-red-400">Errors</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">Success Rate</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">AVG</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-amber-500 dark:text-amber-400">P95</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-red-500 dark:text-red-400">P99</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">TTFT</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in data.model_latencies" :key="m.model" class="border-b border-gray-50 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white font-mono">{{ m.model }}</td>
                  <td class="px-6 py-3 text-center text-sm text-gray-700 dark:text-dark-200">{{ m.request_count.toLocaleString() }}</td>
                  <td class="px-6 py-3 text-center">
                    <span class="text-sm font-semibold" :class="m.success_count > 0 ? 'text-green-600 dark:text-green-400' : 'text-gray-400'">{{ m.success_count }}</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <span :class="m.error_count > 0 ? 'text-sm font-semibold text-red-600 dark:text-red-400' : 'text-sm text-gray-400'">{{ m.error_count }}</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <div class="flex items-center justify-center gap-2">
                      <div class="h-1.5 w-16 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-700">
                        <div class="h-full rounded-full" :class="successRateColor(m.success_count, m.request_count)" :style="{ width: successRateWidth(m.success_count, m.request_count) }"></div>
                      </div>
                      <span class="text-sm" :class="successRateTextColor(m.success_count, m.request_count)">{{ successRate(m.success_count, m.request_count) }}</span>
                    </div>
                  </td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm" :class="m.p95_latency_ms > 10000 ? 'text-amber-600 dark:text-amber-400 font-semibold' : 'text-gray-700 dark:text-dark-200'">{{ formatMs(m.p95_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm" :class="m.p99_latency_ms > 30000 ? 'text-red-600 dark:text-red-400 font-semibold' : 'text-gray-700 dark:text-dark-200'">{{ formatMs(m.p99_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_first_token_ms) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>

    <footer class="border-t border-gray-100 bg-white dark:border-dark-800 dark:bg-dark-950">
      <div class="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        <span class="text-sm text-gray-500 dark:text-dark-400">&copy; {{ currentYear }} {{ siteName }}</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { monitoringAPI, type MonitoringOverview } from '@/api/admin/monitoring'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const currentYear = new Date().getFullYear()

const isDark = ref(document.documentElement.classList.contains('dark'))
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

const savedTheme = localStorage.getItem('theme')
if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
  isDark.value = true
  document.documentElement.classList.add('dark')
}

const data = ref<MonitoringOverview | null>(null)
const loading = ref(false)

interface GroupedModels {
  groupId: number
  groupName: string
  models: NonNullable<MonitoringOverview['group_models']>
}

const groupedModels = computed<GroupedModels[]>(() => {
  if (!data.value?.group_models) return []
  const map = new Map<number, GroupedModels>()
  for (const m of data.value.group_models) {
    let g = map.get(m.group_id)
    if (!g) {
      g = { groupId: m.group_id, groupName: m.group_name, models: [] }
      map.set(m.group_id, g)
    }
    g.models.push(m)
  }
  return Array.from(map.values())
})

async function refresh() {
  loading.value = true
  try {
    data.value = await monitoringAPI.getOverview()
  } catch {
    data.value = null
  } finally {
    loading.value = false
  }
}

function formatMs(ms: number): string {
  if (!ms || ms === 0) return '-'
  if (ms < 1000) return Math.round(ms) + 'ms'
  return (ms / 1000).toFixed(1) + 's'
}

function successRate(success: number, total: number): string {
  if (total === 0) return '-'
  return (success / total * 100).toFixed(1) + '%'
}

function successRateWidth(success: number, total: number): string {
  if (total === 0) return '0%'
  return Math.min(success / total * 100, 100) + '%'
}

function successRateColor(success: number, total: number): string {
  if (total === 0) return 'bg-gray-300'
  const rate = success / total
  if (rate >= 0.95) return 'bg-green-500'
  if (rate >= 0.8) return 'bg-amber-500'
  return 'bg-red-500'
}

function successRateTextColor(success: number, total: number): string {
  if (total === 0) return 'text-gray-400'
  const rate = success / total
  if (rate >= 0.95) return 'text-green-600 dark:text-green-400'
  if (rate >= 0.8) return 'text-amber-600 dark:text-amber-400'
  return 'text-red-600 dark:text-red-400'
}

onMounted(refresh)
</script>
