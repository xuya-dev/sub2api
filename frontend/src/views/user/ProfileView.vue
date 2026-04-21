<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
        <StatCard :title="t('profile.accountBalance')" :value="formatCurrency(user?.balance || 0)" :icon="WalletIcon" icon-variant="success" />
        <StatCard :title="t('profile.concurrencyLimit')" :value="user?.concurrency || 0" :icon="BoltIcon" icon-variant="warning" />
        <StatCard :title="t('profile.memberSince')" :value="formatDate(user?.created_at || '', { year: 'numeric', month: 'long' })" :icon="CalendarIcon" icon-variant="primary" />
      </div>
      <ProfileInfoCard :user="user" />
      <div v-if="checkinStore.enabled" class="card p-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="rounded-xl bg-amber-100 p-3 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900 dark:text-dark-100">{{ t('checkin.title') }}</h3>
              <p class="text-sm text-gray-500 dark:text-dark-400">
                <template v-if="checkinStore.canCheckin">
                  <template v-if="checkinStore.normalEnabled">
                    {{ t('checkin.rangeHint', { min: checkinStore.status?.min_reward?.toFixed(2), max: checkinStore.status?.max_reward?.toFixed(2) }) }}
                  </template>
                  <template v-if="checkinStore.luckEnabled">
                    <span :class="{ 'ml-2': checkinStore.normalEnabled }">
                      {{ t('checkin.multiplierRange', { min: checkinStore.status?.min_multiplier?.toFixed(1), max: checkinStore.status?.max_multiplier?.toFixed(1) }) }}
                    </span>
                  </template>
                </template>
                <template v-else-if="checkinStore.todayReward !== null">
                  <span v-if="checkinStore.todayCheckinType === 'luck'">
                    {{ t('checkin.luckTitle') }}:
                    <template v-if="(checkinStore.todayReward ?? 0) > 0">{{ t('checkin.luckSuccess', { multiplier: checkinStore.todayMultiplier?.toFixed(2) ?? '—', amount: checkinStore.todayReward?.toFixed(2) }) }}</template>
                    <template v-else-if="(checkinStore.todayReward ?? 0) < 0">{{ t('checkin.luckLoss', { multiplier: checkinStore.todayMultiplier?.toFixed(2) ?? '—', amount: Math.abs(checkinStore.todayReward ?? 0).toFixed(2) }) }}</template>
                    <template v-else>{{ t('checkin.luckEven') }}</template>
                  </span>
                  <span v-else>{{ t('checkin.todayReward', { amount: checkinStore.todayReward?.toFixed(2) }) }}</span>
                </template>
              </p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span v-if="checkinStore.streakDays > 0" class="text-sm text-gray-500 dark:text-dark-400">
              {{ t('checkin.streakDays', { days: checkinStore.streakDays }) }}
            </span>

            <template v-if="checkinStore.canCheckin">
              <button
                v-if="checkinStore.normalEnabled"
                type="button"
                :disabled="checkinStore.loading"
                class="rounded-lg bg-amber-500 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-amber-600 disabled:opacity-50"
                @click="checkinStore.doCheckin()"
              >
                {{ checkinStore.loading ? '...' : t('checkin.normalCheckin') }}
              </button>
              <button
                v-if="checkinStore.luckEnabled"
                type="button"
                :disabled="checkinStore.loading"
                class="rounded-lg bg-purple-500 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-purple-600 disabled:opacity-50"
                @click="showProfileLuckModal = true"
              >
                {{ checkinStore.loading ? '...' : t('checkin.luckCheckin') }}
              </button>
            </template>

            <span v-else class="inline-flex items-center gap-1 rounded-lg bg-green-100 px-3 py-2 text-sm font-medium text-green-700 dark:bg-green-900/20 dark:text-green-300">
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ t('checkin.checked') }}
            </span>
          </div>
        </div>
      </div>

      <!-- Profile Luck Checkin Modal -->
      <BaseDialog :show="showProfileLuckModal" :title="t('checkin.luckTitle')" width="narrow" :close-on-click-outside="true" @close="showProfileLuckModal = false">
        <div class="mb-3 rounded-lg bg-purple-50 p-3 dark:bg-purple-900/20">
          <p class="text-xs text-purple-700 dark:text-purple-300">
            {{ t('checkin.multiplierRange', { min: checkinStore.status?.min_multiplier?.toFixed(1), max: checkinStore.status?.max_multiplier?.toFixed(1) }) }}
          </p>
        </div>
        <div class="space-y-4">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('checkin.betAmount') }}</label>
            <input v-model.number="profileLuckBet" type="number" step="0.01" :min="0.01" :max="checkinStore.status?.balance ?? 0" class="input" :placeholder="t('checkin.betAmountPlaceholder')" @keyup.enter="submitProfileLuck" />
          </div>
          <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>{{ t('profile.accountBalance') }}: ${{ checkinStore.status?.balance?.toFixed(2) ?? '0.00' }}</span>
            <button type="button" class="text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="profileLuckBet = checkinStore.status?.balance ?? 0">MAX</button>
          </div>
          <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-700">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('checkin.luckDesc', { min: checkinStore.status?.min_multiplier?.toFixed(1), max: checkinStore.status?.max_multiplier?.toFixed(1) }) }}</p>
          </div>
        </div>
        <template #footer>
          <div class="flex flex-row items-center justify-end gap-3">
            <button type="button" class="rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-700" @click="showProfileLuckModal = false">{{ t('common.cancel') }}</button>
            <button type="button" :disabled="checkinStore.loading || !profileLuckBet || profileLuckBet <= 0 || profileLuckBet > (checkinStore.status?.balance ?? 0)" class="rounded-xl bg-purple-500 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-purple-600 disabled:opacity-50" @click="submitProfileLuck">{{ checkinStore.loading ? '...' : t('checkin.luckButton') }}</button>
          </div>
        </template>
      </BaseDialog>
      <div v-if="contactInfo" class="card border-primary-200 bg-primary-50 dark:bg-primary-900/20 p-6">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-primary-100 rounded-xl text-primary-600"><Icon name="chat" size="lg" /></div>
          <div><h3 class="font-semibold text-primary-800 dark:text-primary-200">{{ t('common.contactSupport') }}</h3><p class="text-sm font-medium">{{ contactInfo }}</p></div>
        </div>
      </div>
      <ProfileEditForm :initial-username="user?.username || ''" />
      <ProfileBalanceNotifyCard
        v-if="user && balanceLowNotifyEnabled"
        :enabled="user.balance_notify_enabled ?? true"
        :threshold="user.balance_notify_threshold"
        :extra-emails="user.balance_notify_extra_emails ?? []"
        :system-default-threshold="systemDefaultThreshold"
        :user-email="user.email"
      />
      <ProfilePasswordForm />
      <ProfileTotpCard />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'; import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'; import { useCheckinStore } from '@/stores/checkin'; import { formatDate } from '@/utils/format'
import { authAPI } from '@/api'; import AppLayout from '@/components/layout/AppLayout.vue'
import StatCard from '@/components/common/StatCard.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ProfileInfoCard from '@/components/user/profile/ProfileInfoCard.vue'
import ProfileEditForm from '@/components/user/profile/ProfileEditForm.vue'
import ProfileBalanceNotifyCard from '@/components/user/profile/ProfileBalanceNotifyCard.vue'
import ProfilePasswordForm from '@/components/user/profile/ProfilePasswordForm.vue'
import ProfileTotpCard from '@/components/user/profile/ProfileTotpCard.vue'
import { Icon } from '@/components/icons'

const { t } = useI18n(); const authStore = useAuthStore(); const checkinStore = useCheckinStore(); const user = computed(() => authStore.user)
const contactInfo = ref('')
const balanceLowNotifyEnabled = ref(false)
const systemDefaultThreshold = ref(0)
const showProfileLuckModal = ref(false)
const profileLuckBet = ref<number>(0)

async function submitProfileLuck() {
  if (!profileLuckBet.value || profileLuckBet.value <= 0) return
  const result = await checkinStore.doLuckCheckin(profileLuckBet.value)
  if (result) {
    showProfileLuckModal.value = false
    profileLuckBet.value = 0
  }
}

const WalletIcon = { render: () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [h('path', { d: 'M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12' })]) }
const BoltIcon = { render: () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [h('path', { d: 'm3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z' })]) }
const CalendarIcon = { render: () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [h('path', { d: 'M6.75 3v2.25M17.25 3v2.25' })]) }

onMounted(async () => { try { const s = await authAPI.getPublicSettings(); contactInfo.value = s.contact_info || ''; balanceLowNotifyEnabled.value = s.balance_low_notify_enabled ?? false; systemDefaultThreshold.value = s.balance_low_notify_threshold ?? 0 } catch (error) { console.error('Failed to load settings:', error) } })
const formatCurrency = (v: number) => `$${v.toFixed(2)}`
</script>