<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <header class="relative z-20 border-b border-gray-100 bg-white/80 backdrop-blur-xl dark:border-dark-800 dark:bg-dark-950/80">
      <nav class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-8 w-8 overflow-hidden rounded-lg">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-bold text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-2">
          <router-link to="/leaderboard"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-gray-900 bg-gray-100 dark:text-white dark:bg-dark-800 sm:flex">
            {{ t('leaderboard.title') }}
          </router-link>
          <router-link to="/key-usage"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('home.keyUsage') }}
          </router-link>
          <router-link to="/monitoring"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('admin.monitoring.title') }}
          </router-link>
          <router-link to="/pricing"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('pricing.title') }}
          </router-link>
          <LocaleSwitcher />
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

    <main class="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 sm:py-8">
      <div class="space-y-6">
        <div>
          <h1 class="text-xl font-bold text-gray-900 dark:text-white sm:text-2xl">{{ t('leaderboard.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('leaderboard.subtitle') }}</p>
        </div>

        <div class="card overflow-hidden p-0">
          <div class="flex border-b border-gray-200 dark:border-dark-700">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              @click="activeTab = tab.key"
              :class="[
                'flex-1 px-4 py-3 text-sm font-medium transition-colors',
                activeTab === tab.key
                  ? 'border-b-2 border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'
              ]"
            >
              {{ tab.label }}
            </button>
          </div>

          <div v-if="activeTab === 'consumption'" class="flex gap-1 border-b border-gray-100 px-4 py-2 dark:border-dark-800">
            <button
              v-for="p in periods"
              :key="p.key"
              @click="activePeriod = p.key"
              :class="[
                'rounded-md px-3 py-1 text-xs font-medium transition-colors',
                activePeriod === p.key
                  ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                  : 'text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800'
              ]"
            >
              {{ p.label }}
            </button>
          </div>

          <div class="p-4">
            <div v-if="loading" class="flex items-center justify-center py-12">
              <div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <div v-else-if="entries.length === 0" class="py-12 text-center text-sm text-gray-400 dark:text-dark-500">
              {{ t('leaderboard.empty') }}
            </div>
            <div v-else class="space-y-1">
              <div
                v-for="entry in entries"
                :key="entry.rank"
                class="flex items-center gap-3 rounded-lg px-3 py-2.5 transition-colors hover:bg-gray-50 dark:hover:bg-dark-800/50"
              >
                <div :class="rankClass(entry.rank)" class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-bold">
                  <span v-if="entry.rank <= 3">{{ ['🥇', '🥈', '🥉'][entry.rank - 1] }}</span>
                  <span v-else>{{ entry.rank }}</span>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ entry.username }}</p>
                </div>
                <div class="shrink-0 text-right">
                  <span v-if="activeTab === 'checkin'" class="text-sm font-semibold text-amber-600 dark:text-amber-400">
                    {{ t('leaderboard.streakDays', { days: entry.value }) }}
                  </span>
                  <span v-else class="text-sm font-semibold text-gray-900 dark:text-white">
                    ${{ entry.value.toFixed(2) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <PublicPageFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import PublicPageFooter from '@/components/common/PublicPageFooter.vue'
import { leaderboardAPI, type LeaderboardEntry } from '@/api/leaderboard'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

const isDark = ref(document.documentElement.classList.contains('dark'))

type TabKey = 'balance' | 'consumption' | 'checkin'
type PeriodKey = 'daily' | 'weekly' | 'monthly'

const activeTab = ref<TabKey>('balance')
const activePeriod = ref<PeriodKey>('daily')
const entries = ref<LeaderboardEntry[]>([])
const loading = ref(false)

const tabs = computed(() => [
  { key: 'balance' as TabKey, label: t('leaderboard.tabs.balance') },
  { key: 'consumption' as TabKey, label: t('leaderboard.tabs.consumption') },
  { key: 'checkin' as TabKey, label: t('leaderboard.tabs.checkin') },
])

const periods = computed(() => [
  { key: 'daily' as PeriodKey, label: t('leaderboard.periods.daily') },
  { key: 'weekly' as PeriodKey, label: t('leaderboard.periods.weekly') },
  { key: 'monthly' as PeriodKey, label: t('leaderboard.periods.monthly') },
])

function rankClass(rank: number): string {
  if (rank === 1) return 'bg-amber-100 dark:bg-amber-900/30'
  if (rank === 2) return 'bg-gray-200 dark:bg-dark-700'
  if (rank === 3) return 'bg-orange-100 dark:bg-orange-900/30'
  return 'bg-gray-100 dark:bg-dark-800'
}

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

async function fetchData() {
  loading.value = true
  try {
    let res
    switch (activeTab.value) {
      case 'balance':
        res = await leaderboardAPI.getBalanceLeaderboard(1, 20)
        break
      case 'consumption':
        res = await leaderboardAPI.getConsumptionLeaderboard(activePeriod.value, 1, 20)
        break
      case 'checkin':
        res = await leaderboardAPI.getCheckinLeaderboard(1, 20)
        break
    }
    entries.value = res.items || []
  } catch {
    entries.value = []
  } finally {
    loading.value = false
  }
}

watch([activeTab, activePeriod], () => fetchData())

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  authStore.checkAuth()
  fetchData()
})
</script>
