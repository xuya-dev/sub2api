import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { checkinAPI, type CheckinStatus, type CheckinResult, type BlindboxResult } from '@/api/checkin'
import { useAuthStore } from './auth'

export const useCheckinStore = defineStore('checkin', () => {
  const status = ref<CheckinStatus | null>(null)
  const loading = ref(false)
  const checkinResult = ref<CheckinResult | null>(null)
  const blindboxResult = ref<BlindboxResult | null>(null)

  const canCheckin = computed(() => status.value?.can_checkin ?? false)
  const enabled = computed(() => (status.value?.enabled ?? false) || (status.value?.luck_enabled ?? false))
  const normalEnabled = computed(() => status.value?.enabled ?? false)
  const luckEnabled = computed(() => status.value?.luck_enabled ?? false)
  const checkedInToday = computed(() => enabled.value && !canCheckin.value && status.value !== null)
  const streakDays = computed(() => status.value?.streak_days ?? 0)
  const todayReward = computed(() => status.value?.today_reward ?? null)
  const todayCheckinType = computed(() => status.value?.today_checkin_type ?? null)
  const todayMultiplier = computed(() => status.value?.today_multiplier ?? null)

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
      blindboxResult.value = result.blindbox ?? null

      if (status.value) {
        status.value.can_checkin = false
        status.value.streak_days = result.streak_days
        status.value.today_reward = result.reward_amount
        status.value.today_checkin_type = result.checkin_type
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

  async function doLuckCheckin(betAmount: number): Promise<CheckinResult | null> {
    if (loading.value) return null
    loading.value = true
    try {
      const result = await checkinAPI.luckCheckin(betAmount)
      checkinResult.value = result
      blindboxResult.value = result.blindbox ?? null

      if (status.value) {
        status.value.can_checkin = false
        status.value.streak_days = result.streak_days
        status.value.today_reward = result.reward_amount
        status.value.today_checkin_type = result.checkin_type
        if (result.multiplier !== undefined) {
          status.value.today_multiplier = result.multiplier
        }
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

  function clearBlindboxResult() {
    blindboxResult.value = null
  }

  function $reset() {
    status.value = null
    loading.value = false
    checkinResult.value = null
    blindboxResult.value = null
  }

  return {
    status,
    loading,
    checkinResult,
    blindboxResult,
    canCheckin,
    enabled,
    normalEnabled,
    luckEnabled,
    checkedInToday,
    streakDays,
    todayReward,
    todayCheckinType,
    todayMultiplier,
    fetchStatus,
    doCheckin,
    doLuckCheckin,
    clearBlindboxResult,
    $reset
  }
})
