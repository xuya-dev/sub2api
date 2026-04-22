<template>
  <section class="border-y border-gray-100 bg-white py-20 dark:border-dark-800 dark:bg-dark-950 sm:py-24">
    <div class="mx-auto max-w-7xl px-6">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">{{ t('leaderboard.title') }}</h2>
        <p class="mt-4 text-lg text-gray-600 dark:text-dark-300">{{ t('leaderboard.subtitle') }}</p>
      </div>

      <div class="mx-auto mt-10 max-w-3xl">
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
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { leaderboardAPI, type LeaderboardEntry } from '@/api/leaderboard'

const { t } = useI18n()

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
        res = await leaderboardAPI.getBalanceLeaderboard(1, 10)
        break
      case 'consumption':
        res = await leaderboardAPI.getConsumptionLeaderboard(activePeriod.value, 1, 10)
        break
      case 'checkin':
        res = await leaderboardAPI.getCheckinLeaderboard(1, 10)
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

onMounted(() => fetchData())
</script>
