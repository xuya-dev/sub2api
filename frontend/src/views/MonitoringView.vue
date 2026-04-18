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
          <LocaleSwitcher />
          <button @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-50 dark:text-dark-400 dark:hover:bg-dark-800">
            <svg v-if="isDark" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" />
            </svg>
            <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
            </svg>
          </button>
          <router-link v-if="isAuthenticated" to="/dashboard"
            class="ml-1 inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100">
            {{ t('home.dashboard') }}
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
          <button @click="refresh" :disabled="loading"
            class="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:opacity-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700">
            <svg class="h-4 w-4" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
            {{ t('common.refresh') }}
          </button>
        </div>

        <!-- Summary Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.todayRequests') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data?.total_requests_today?.toLocaleString() ?? '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.avgLatency') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data ? Math.round(data.avg_latency_ms_today) + 'ms' : '-' }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.errorAccounts') }}</p>
            <p class="mt-1 text-2xl font-bold text-red-600 dark:text-red-400">{{ data?.error_accounts?.length ?? 0 }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.totalGroups') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data?.groups?.length ?? 0 }}</p>
          </div>
        </div>

        <!-- Group Health -->
        <div class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.monitoring.groupHealth') }}</h2>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.groupName') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.total') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-green-600 dark:text-green-400">{{ t('admin.monitoring.active') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-red-500 dark:text-red-400">{{ t('admin.monitoring.error') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-amber-500 dark:text-amber-400">{{ t('admin.monitoring.rateLimited') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-orange-500 dark:text-orange-400">{{ t('admin.monitoring.overload') }}</th>
                  <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.disabled') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="g in data?.groups" :key="g.group_id" class="border-b border-gray-50 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white">{{ g.group_name }}</td>
                  <td class="px-6 py-3 text-center text-sm text-gray-700 dark:text-dark-200">{{ g.total_accounts }}</td>
                  <td class="px-6 py-3 text-center">
                    <span v-if="g.active_accounts > 0" class="inline-flex items-center gap-1 text-sm font-semibold text-green-600 dark:text-green-400">{{ g.active_accounts }}</span>
                    <span v-else class="text-sm text-gray-400">0</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <span v-if="g.error_accounts > 0" class="text-sm font-semibold text-red-600 dark:text-red-400">{{ g.error_accounts }}</span>
                    <span v-else class="text-sm text-gray-400">0</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <span v-if="g.rate_limited > 0" class="text-sm font-semibold text-amber-600 dark:text-amber-400">{{ g.rate_limited }}</span>
                    <span v-else class="text-sm text-gray-400">0</span>
                  </td>
                  <td class="px-6 py-3 text-center">
                    <span v-if="g.overload > 0" class="text-sm font-semibold text-orange-600 dark:text-orange-400">{{ g.overload }}</span>
                    <span v-else class="text-sm text-gray-400">0</span>
                  </td>
                  <td class="px-6 py-3 text-center text-sm text-gray-500 dark:text-dark-400">{{ g.disabled }}</td>
                </tr>
                <tr v-if="!data?.groups?.length">
                  <td colspan="7" class="px-6 py-8 text-center text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.noData') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Model Latency -->
        <div class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.monitoring.modelLatency') }}</h2>
            <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.modelLatencyHint') }}</p>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.model') }}</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.requests') }}</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">AVG</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">P50</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-amber-500 dark:text-amber-400">P95</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-red-500 dark:text-red-400">P99</th>
                  <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">TTFT</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in data?.model_latencies" :key="m.model" class="border-b border-gray-50 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white font-mono">{{ m.model }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ m.request_count.toLocaleString() }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.p50_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm" :class="m.p95_latency_ms > 10000 ? 'text-amber-600 dark:text-amber-400 font-semibold' : 'text-gray-700 dark:text-dark-200'">{{ formatMs(m.p95_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm" :class="m.p99_latency_ms > 30000 ? 'text-red-600 dark:text-red-400 font-semibold' : 'text-gray-700 dark:text-dark-200'">{{ formatMs(m.p99_latency_ms) }}</td>
                  <td class="px-6 py-3 text-right text-sm text-gray-700 dark:text-dark-200">{{ formatMs(m.avg_first_token_ms) }}</td>
                </tr>
                <tr v-if="!data?.model_latencies?.length">
                  <td colspan="7" class="px-6 py-8 text-center text-sm text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.noData') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Error Accounts -->
        <div v-if="data?.error_accounts?.length" class="rounded-xl border border-red-200 bg-white shadow-sm dark:border-red-900/50 dark:bg-dark-900">
          <div class="border-b border-red-100 px-6 py-4 dark:border-red-900/30">
            <h2 class="text-lg font-semibold text-red-700 dark:text-red-400">{{ t('admin.monitoring.errorAccountsList') }}</h2>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.accountName') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.groupName') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.status') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.monitoring.errorMessage') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="e in data.error_accounts" :key="e.account_id" class="border-b border-gray-50 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm text-gray-700 dark:text-dark-200">{{ e.account_id }}</td>
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white">{{ e.account_name || '-' }}</td>
                  <td class="px-6 py-3 text-sm text-gray-700 dark:text-dark-200">{{ e.group_name || '-' }}</td>
                  <td class="px-6 py-3">
                    <span :class="statusClass(e.status)" class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium">{{ e.status }}</span>
                  </td>
                  <td class="max-w-md truncate px-6 py-3 text-sm text-gray-500 dark:text-dark-400" :title="e.error_message">{{ e.error_message || '-' }}</td>
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
        <button @click="refresh" :disabled="loading" class="text-sm text-gray-400 transition-colors hover:text-gray-600 dark:text-dark-500 dark:hover:text-dark-300">
          {{ t('common.refresh') }}
        </button>
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
import { monitoringAPI, type MonitoringOverview } from '@/api/admin/monitoring'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
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

function statusClass(status: string): string {
  if (status === 'error') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
  if (status === 'active') return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
  return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
}

onMounted(refresh)
</script>
