import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { checkinAPI, type CheckinStatus, type CheckinResult } from '@/api/checkin'
import { useAuthStore } from './auth'

export const useCheckinStore = defineStore('checkin', () => {
  const status = ref<CheckinStatus | null>(null)
  const loading = ref(false)
  const checkinResult = ref<CheckinResult | null>(null)

  const canCheckin = computed(() => status.value?.can_checkin ?? false)
  const enabled = computed(() => status.value?.enabled ?? false)
  const checkedInToday = computed(() => enabled.value && !canCheckin.value && status.value !== null)
  const streakDays = computed(() => status.value?.streak_days ?? 0)
  const todayReward = computed(() => status.value?.today_reward ?? null)

  async function fetchStatus() {
    try {
      status.value = await checkinAPI.getCheckinStatus()
    } catch {
      status.value = null
    }
  }

  async function doCheckin(): Promise<CheckinResult | null> {
    if (loading.value) return null
    loading.value = true
    try {
      const result = await checkinAPI.checkin()
      checkinResult.value = result

      if (status.value) {
        status.value.can_checkin = false
        status.value.streak_days = result.streak_days
        status.value.today_reward = result.reward_amount
      }

      const authStore = useAuthStore()
      await authStore.refreshUser()

      return result
    } catch {
      return null
    } finally {
      loading.value = false
    }
  }

  function $reset() {
    status.value = null
    loading.value = false
    checkinResult.value = null
  }

  return {
    status,
    loading,
    checkinResult,
    canCheckin,
    enabled,
    checkedInToday,
    streakDays,
    todayReward,
    fetchStatus,
    doCheckin,
    $reset
  }
})
