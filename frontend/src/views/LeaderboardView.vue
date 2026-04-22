<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <PublicPageHeader active-path="/leaderboard" />

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

          <div class="relative min-h-[200px]">
            <div v-if="loading" class="absolute inset-0 z-10 flex items-center justify-center bg-white/80 backdrop-blur-sm dark:bg-dark-900/80">
              <div class="flex flex-col items-center gap-3">
                <div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
                <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('common.loading') }}</span>
              </div>
            </div>

            <div v-if="!loading && entries.length === 0" class="py-16 text-center text-sm text-gray-400 dark:text-dark-500">
              {{ t('leaderboard.empty') }}
            </div>

            <div v-else class="divide-y divide-gray-50 dark:divide-dark-800/50">
              <div
                v-for="entry in entries"
                :key="entry.rank"
                class="flex items-center gap-3 px-4 py-3 transition-colors hover:bg-gray-50 dark:hover:bg-dark-800/50 sm:px-6"
              >
                <div :class="rankClass(entry.rank)" class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold">
                  <span v-if="entry.rank <= 3">{{ ['🥇', '🥈', '🥉'][entry.rank - 1] }}</span>
                  <span v-else class="text-gray-500 dark:text-dark-400">{{ entry.rank }}</span>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ entry.username }}</p>
                  <p v-if="entry.subtitle" class="mt-0.5 truncate text-xs text-gray-400 dark:text-dark-500">{{ entry.subtitle }}</p>
                </div>
                <div class="shrink-0 text-right">
                  <template v-if="activeTab === 'checkin'">
                    <span class="text-sm font-bold text-amber-600 dark:text-amber-400">{{ entry.value }}</span>
                    <span class="text-xs text-amber-500/70 dark:text-amber-400/50"> {{ t('leaderboard.streakDays', { days: '' }).trim() }}</span>
                  </template>
                  <template v-else>
                    <span class="text-sm font-bold text-gray-900 dark:text-white">${{ entry.value.toFixed(2) }}</span>
                  </template>
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
import PublicPageHeader from '@/components/common/PublicPageHeader.vue'
import PublicPageFooter from '@/components/common/PublicPageFooter.vue'
import { leaderboardAPI, type LeaderboardEntry } from '@/api/leaderboard'

const { t } = useI18n()
const appStore = useAppStore()

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
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  fetchData()
})
</script>
