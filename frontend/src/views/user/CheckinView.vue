<template>
  <AppLayout>
    <div class="mx-auto max-w-5xl space-y-6">
      <!-- Top Stats Row -->
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-xl bg-emerald-100 p-2.5 dark:bg-emerald-900/30">
              <svg class="h-5 w-5 text-emerald-600 dark:text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('checkin.page.balance') }}</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">${{ user?.balance?.toFixed(2) || '0.00' }}</p>
            </div>
          </div>
        </div>
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-xl bg-amber-100 p-2.5 dark:bg-amber-900/30">
              <svg class="h-5 w-5 text-amber-600 dark:text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.362 5.214A8.252 8.252 0 0112 21 8.25 8.25 0 016.038 7.048 8.287 8.287 0 009 9.6a8.983 8.983 0 013.361-6.867 8.21 8.21 0 003 2.48z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18a3.75 3.75 0 00.495-7.467 5.99 5.99 0 00-1.925 3.546 5.974 5.974 0 01-2.133-1A3.75 3.75 0 0012 18z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('checkin.page.streak') }}</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">{{ checkinStore.streakDays }}<span class="text-sm font-normal text-gray-400">{{ t('checkin.page.days') }}</span></p>
            </div>
          </div>
        </div>
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-xl bg-blue-100 p-2.5 dark:bg-blue-900/30">
              <svg class="h-5 w-5 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('checkin.page.concurrency') }}</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">{{ user?.concurrency || 0 }}</p>
            </div>
          </div>
        </div>
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="rounded-xl bg-purple-100 p-2.5 dark:bg-purple-900/30">
              <svg class="h-5 w-5 text-purple-600 dark:text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('checkin.page.blindboxCount') }}</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">{{ blindboxTotal }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Checkin Action Card -->
      <div v-if="checkinStore.enabled" class="card overflow-hidden">
        <div class="relative p-6">
          <div class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex items-start gap-4">
              <div class="rounded-2xl bg-gradient-to-br from-amber-400 to-orange-500 p-3.5 shadow-lg shadow-amber-200/50 dark:shadow-amber-900/30">
                <svg class="h-7 w-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" />
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('checkin.title') }}</h3>
                <div class="mt-1 space-y-1">
                  <template v-if="checkinStore.canCheckin">
                    <p v-if="checkinStore.normalEnabled" class="text-sm text-gray-500 dark:text-dark-400">
                      {{ t('checkin.rangeHint', { min: checkinStore.status?.min_reward?.toFixed(2), max: checkinStore.status?.max_reward?.toFixed(2) }) }}
                    </p>
                    <p v-if="checkinStore.luckEnabled" class="text-sm text-purple-600 dark:text-purple-400">
                      {{ t('checkin.multiplierRange', { min: checkinStore.status?.min_multiplier?.toFixed(1), max: checkinStore.status?.max_multiplier?.toFixed(1) }) }}
                    </p>
                  </template>
                  <template v-else-if="checkinStore.todayReward !== null">
                    <p class="text-sm text-emerald-600 dark:text-emerald-400">
                      <span v-if="checkinStore.todayCheckinType === 'luck'">
                        <template v-if="(checkinStore.todayReward ?? 0) > 0">{{ t('checkin.luckSuccess', { multiplier: checkinStore.todayMultiplier?.toFixed(2) ?? '—', amount: checkinStore.todayReward?.toFixed(2) }) }}</template>
                        <template v-else-if="(checkinStore.todayReward ?? 0) < 0">{{ t('checkin.luckLoss', { multiplier: checkinStore.todayMultiplier?.toFixed(2) ?? '—', amount: Math.abs(checkinStore.todayReward ?? 0).toFixed(2) }) }}</template>
                        <template v-else>{{ t('checkin.luckEven') }}</template>
                      </span>
                      <span v-else>{{ t('checkin.todayReward', { amount: checkinStore.todayReward?.toFixed(2) }) }}</span>
                    </p>
                  </template>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3 sm:flex-col sm:items-stretch">
              <div v-if="checkinStore.streakDays > 0" class="flex items-center gap-1.5 text-sm" :class="streakColor">
                <svg :class="streakFlameClass" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 23c-3.866 0-7-2.686-7-6 0-2.972 1.637-5.342 3.276-7.097A1 1 0 019.2 10.1c-.48.907-.7 1.747-.7 2.4 0 1.933 1.567 3.5 3.5 3.5s3.5-1.567 3.5-3.5c0-1.073-.423-2.177-1.104-3.207C13.216 7.246 12 5.764 12 3.5c0-.578.077-1.16.236-1.72a1 1 0 011.69-.481C16.44 3.678 19 7.223 19 11c0 4.314-2.957 7-5.463 8.565A6.96 6.96 0 0112 23z"/>
                </svg>
                <span class="font-semibold">{{ t('checkin.streakDays', { days: checkinStore.streakDays }) }}</span>
              </div>
              <div class="flex gap-2">
                <template v-if="checkinStore.canCheckin">
                  <button v-if="checkinStore.normalEnabled" type="button" :disabled="checkinStore.loading" class="checkin-btn checkin-btn-normal" @click="checkinStore.doCheckin()">
                    <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                    {{ checkinStore.loading ? '...' : t('checkin.normalCheckin') }}
                  </button>
                  <button v-if="checkinStore.luckEnabled" type="button" :disabled="checkinStore.loading" class="checkin-btn checkin-btn-luck" @click="showLuckModal = true">
                    <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" /></svg>
                    {{ t('checkin.luckCheckin') }}
                  </button>
                </template>
                <div v-else class="flex items-center gap-1.5 rounded-lg bg-emerald-100 px-3 py-2 text-sm font-medium text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                  {{ t('checkin.checked') }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Luck Checkin Modal -->
      <BaseDialog :show="showLuckModal" :title="t('checkin.luckTitle')" width="narrow" :close-on-click-outside="true" @close="showLuckModal = false">
        <div class="mb-3 rounded-lg bg-purple-50 p-3 dark:bg-purple-900/20">
          <p class="text-xs text-purple-700 dark:text-purple-300">
            {{ t('checkin.multiplierRange', { min: checkinStore.status?.min_multiplier?.toFixed(1), max: checkinStore.status?.max_multiplier?.toFixed(1) }) }}
          </p>
        </div>
        <div class="space-y-4">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('checkin.betAmount') }}</label>
            <input v-model.number="luckBet" type="number" step="0.01" :min="0.01" :max="checkinStore.status?.balance ?? 0" class="input" :placeholder="t('checkin.betAmountPlaceholder')" @keyup.enter="submitLuck" />
          </div>
          <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>{{ t('profile.accountBalance') }}: ${{ checkinStore.status?.balance?.toFixed(2) ?? '0.00' }}</span>
            <button type="button" class="text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="luckBet = checkinStore.status?.balance ?? 0">MAX</button>
          </div>
        </div>
        <template #footer>
          <div class="flex flex-row items-center justify-end gap-3">
            <button type="button" class="btn btn-secondary" @click="showLuckModal = false">{{ t('common.cancel') }}</button>
            <button type="button" :disabled="checkinStore.loading || !luckBet || luckBet <= 0 || luckBet > (checkinStore.status?.balance ?? 0)" class="rounded-xl bg-purple-500 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-purple-600 disabled:opacity-50" @click="submitLuck">{{ checkinStore.loading ? '...' : t('checkin.luckButton') }}</button>
          </div>
        </template>
      </BaseDialog>

      <!-- Check-in Calendar -->
      <div v-if="calendarDays.length > 0" class="card p-6">
        <div class="mb-4 flex items-center gap-3">
          <div class="rounded-xl bg-emerald-100 p-2.5 dark:bg-emerald-900/30">
            <svg class="h-5 w-5 text-emerald-600 dark:text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5" />
            </svg>
          </div>
          <div>
            <h3 class="font-semibold text-gray-900 dark:text-white">{{ t('checkin.page.calendarTitle') }}</h3>
          </div>
        </div>

        <div class="grid grid-cols-7 gap-1.5 sm:gap-2">
          <div v-for="d in weekHeaders" :key="d" class="pb-1.5 text-center text-[11px] font-medium text-gray-400 dark:text-dark-500">{{ d }}</div>
          <div
            v-for="(day, i) in calendarGrid"
            :key="i"
            class="calendar-cell relative flex flex-col items-center justify-center rounded-lg py-2 text-center transition-colors"
            :class="getCalendarCellClass(day)"
          >
            <span class="text-xs font-medium" :class="day.isCurrentMonth ? 'text-gray-700 dark:text-gray-300' : 'text-gray-300 dark:text-dark-600'">{{ day.dayNum }}</span>
            <div v-if="day.checkedIn" class="mt-0.5 flex items-center justify-center">
              <div v-if="day.rewardType === 'luck'" class="h-1.5 w-1.5 rounded-full bg-purple-500"></div>
              <div v-else class="h-1.5 w-1.5 rounded-full bg-emerald-500"></div>
            </div>
          </div>
        </div>

        <div class="mt-3 flex items-center justify-center gap-5 text-xs text-gray-400 dark:text-dark-500">
          <div class="flex items-center gap-1.5"><div class="h-2 w-2 rounded-full bg-emerald-500"></div>{{ t('checkin.normalCheckin') }}</div>
          <div class="flex items-center gap-1.5"><div class="h-2 w-2 rounded-full bg-purple-500"></div>{{ t('checkin.luckCheckin') }}</div>
          <div class="flex items-center gap-1.5"><div class="h-2 w-2 rounded-full bg-gray-200 dark:bg-dark-700"></div>{{ t('checkin.page.todayNoResult') }}</div>
        </div>
      </div>

      <!-- Two Column: Blindbox History + Reward Stats -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Blindbox History (2/3 width) -->
        <div class="lg:col-span-2">
          <div class="card p-6">
            <div class="mb-4 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="rounded-xl bg-purple-100 p-2.5 dark:bg-purple-900/30">
                  <svg class="h-5 w-5 text-purple-600 dark:text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
                  </svg>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900 dark:text-white">{{ t('checkin.blindboxHistory') }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('checkin.blindboxHistoryDesc') }}</p>
                </div>
              </div>
            </div>

            <div v-if="blindboxRecords.length === 0" class="py-8 text-center">
              <svg class="mx-auto h-12 w-12 text-gray-300 dark:text-dark-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
              </svg>
              <p class="mt-2 text-sm text-gray-400 dark:text-dark-500">{{ t('checkin.page.noBlindbox') }}</p>
            </div>

            <div v-else class="space-y-2">
              <div v-for="record in blindboxRecords" :key="record.id" class="flex items-center justify-between rounded-lg border px-4 py-3 blindbox-record" :class="getRecordBorderClass(record.rarity)">
                <div class="flex items-center gap-3">
                  <span class="blindbox-rarity-badge text-xs" :class="getRarityBadgeClass(record.rarity)">{{ getRarityLabel(record.rarity) }}</span>
                  <div>
                    <span class="text-sm font-medium text-gray-900 dark:text-white">{{ record.prize_name }}</span>
                    <span v-if="record.reward_type === 'invitation_code' && record.reward_detail" class="ml-2 rounded bg-indigo-50 px-2 py-0.5 text-xs font-mono text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400">{{ record.reward_detail }}</span>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-sm font-semibold" :class="getRecordRewardColor(record)">{{ formatRecordReward(record) }}</div>
                  <div class="text-xs text-gray-400 dark:text-dark-500">{{ record.created_at }}</div>
                </div>
              </div>
            </div>

            <div v-if="blindboxTotal > blindboxRecords.length" class="mt-4 text-center">
              <button type="button" class="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="loadMoreBlindboxRecords">{{ t('common.loadMore') }}</button>
            </div>
          </div>
        </div>

        <!-- Right Sidebar: Stats & Info -->
        <div class="space-y-6">
          <!-- Today's Check-in Result -->
          <div class="card p-5">
            <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('checkin.page.todayResult') }}</h4>
            <div v-if="checkinStore.todayReward !== null && !checkinStore.canCheckin" class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('checkin.page.todayReward') }}</span>
                <span class="text-sm font-bold" :class="checkinStore.todayReward >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-500 dark:text-red-400'">
                  {{ checkinStore.todayReward >= 0 ? '+' : '' }}${{ checkinStore.todayReward?.toFixed(2) }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('checkin.checkinType') }}</span>
                <span class="rounded-full px-2.5 py-0.5 text-xs font-medium" :class="checkinStore.todayCheckinType === 'luck' ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300' : 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'">
                  {{ checkinStore.todayCheckinType === 'luck' ? t('checkin.page.todayLuck') : t('checkin.page.todayNormal') }}
                </span>
              </div>
              <div v-if="checkinStore.todayCheckinType === 'luck' && checkinStore.todayMultiplier" class="flex items-center justify-between">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('checkin.page.todayMultiplier') }}</span>
                <span class="text-sm font-semibold text-purple-600 dark:text-purple-400">{{ checkinStore.todayMultiplier?.toFixed(2) }}×</span>
              </div>
              <div v-if="checkinStore.blindboxResult" class="mt-2 rounded-lg border border-purple-200 bg-gradient-to-br from-purple-50 to-indigo-50 p-3 dark:border-purple-800/50 dark:from-purple-900/20 dark:to-indigo-900/20">
                <div class="mb-1.5 flex items-center justify-between">
                  <span class="text-xs font-medium text-purple-600 dark:text-purple-400">{{ t('checkin.page.todayBlindbox') }}</span>
                  <span class="blindbox-rarity-badge text-[10px]" :class="getRarityBadgeClass(checkinStore.blindboxResult.rarity)">{{ getRarityLabel(checkinStore.blindboxResult.rarity) }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ checkinStore.blindboxResult.prize_name }}</span>
                  <span class="text-sm font-semibold" :class="getRecordRewardColor({ rarity: checkinStore.blindboxResult.rarity })">{{ formatBlindboxReward(checkinStore.blindboxResult) }}</span>
                </div>
                <div v-if="checkinStore.blindboxResult.reward_type === 'invitation_code' && checkinStore.blindboxResult.reward_detail" class="mt-1.5 rounded bg-indigo-50 px-2 py-1 text-xs font-mono text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400">
                  {{ checkinStore.blindboxResult.reward_detail }}
                </div>
              </div>
            </div>
            <div v-else class="py-3 text-center">
              <p class="text-xs text-gray-400 dark:text-dark-500">{{ t('checkin.page.todayNoResult') }}</p>
            </div>
          </div>

          <!-- Rarity Breakdown -->
          <div v-if="blindboxRecords.length > 0" class="card p-5">
            <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('checkin.page.rarityBreakdown') }}</h4>
            <div class="space-y-2.5">
              <div v-for="r in rarityBreakdown" :key="r.key" class="flex items-center gap-2.5">
                <span class="blindbox-rarity-badge text-[10px]" :class="r.badgeClass">{{ r.label }}</span>
                <div class="h-1.5 flex-1 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                  <div class="h-full rounded-full transition-all duration-500" :class="r.barClass" :style="{ width: r.percent + '%' }"></div>
                </div>
                <span class="text-xs font-medium tabular-nums text-gray-600 dark:text-gray-400">{{ r.count }}</span>
              </div>
            </div>
          </div>

          <!-- Blindbox Trigger Info -->
          <div v-if="checkinStore.status?.blindbox_enabled" class="card p-5">
            <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('checkin.page.blindboxInfo') }}</h4>
            <div class="space-y-2 text-xs text-gray-600 dark:text-gray-400">
              <div class="flex justify-between">
                <span>{{ t('checkin.page.triggerType') }}</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ checkinStore.status.blindbox_trigger_type === 'total' ? t('checkin.page.triggerTotal') : t('checkin.page.triggerStreak') }}</span>
              </div>
              <div class="flex justify-between">
                <span>{{ t('checkin.page.triggerInterval') }}</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ checkinStore.status.blindbox_interval }}{{ t('checkin.page.days') }}</span>
              </div>
              <div v-if="nextBlindboxHint" class="mt-2 rounded-lg bg-purple-50 p-2.5 text-center dark:bg-purple-900/20">
                <span class="text-purple-700 dark:text-purple-300">{{ nextBlindboxHint }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useCheckinStore } from '@/stores/checkin'
import { getBlindboxRecords, getCheckinCalendar, type BlindboxRecordItem, type BlindboxResult, type CheckinCalendarDay } from '@/api/checkin'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const checkinStore = useCheckinStore()
const user = computed(() => authStore.user)

const showLuckModal = ref(false)
const luckBet = ref<number>(0)

const blindboxRecords = ref<BlindboxRecordItem[]>([])
const blindboxTotal = ref(0)
const blindboxPage = ref(1)

const calendarDays = ref<CheckinCalendarDay[]>([])

const weekHeaders = computed(() => {
  const locale = (t as any).locale?.() || 'en'
  if (locale === 'zh') return ['一', '二', '三', '四', '五', '六', '日']
  return ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
})

interface CalendarCell {
  date: string
  dayNum: number
  isCurrentMonth: boolean
  checkedIn: boolean
  rewardType: string
  rewardValue: number
  isToday: boolean
}

const calendarGrid = computed<CalendarCell[]>(() => {
  const days = calendarDays.value
  if (days.length === 0) return []

  const firstDate = new Date(days[0].date + 'T00:00:00')
  const lastDate = new Date(days[days.length - 1].date + 'T00:00:00')

  const checkedMap = new Map<string, CheckinCalendarDay>()
  for (const d of days) {
    checkedMap.set(d.date, d)
  }

  const startOfWeek = new Date(firstDate)
  const dow = startOfWeek.getDay()
  const mondayOffset = dow === 0 ? -6 : 1 - dow
  startOfWeek.setDate(startOfWeek.getDate() + mondayOffset)

  const endOfWeek = new Date(lastDate)
  const endDow = endOfWeek.getDay()
  const sundayOffset = endDow === 0 ? 0 : 7 - endDow
  endOfWeek.setDate(endOfWeek.getDate() + sundayOffset)

  const todayStr = new Date().toISOString().slice(0, 10)

  const cells: CalendarCell[] = []
  const current = new Date(startOfWeek)
  while (current <= endOfWeek) {
    const dateStr = current.toISOString().slice(0, 10)
    const calDay = checkedMap.get(dateStr)
    const firstMonth = firstDate.getFullYear() * 100 + firstDate.getMonth()
    const curMonth = current.getFullYear() * 100 + current.getMonth()

    cells.push({
      date: dateStr,
      dayNum: current.getDate(),
      isCurrentMonth: curMonth === firstMonth || curMonth === lastDate.getFullYear() * 100 + lastDate.getMonth(),
      checkedIn: calDay?.checked_in ?? false,
      rewardType: calDay?.reward_type ?? '',
      rewardValue: calDay?.reward_value ?? 0,
      isToday: dateStr === todayStr,
    })
    current.setDate(current.getDate() + 1)
  }

  return cells
})

function getCalendarCellClass(day: CalendarCell): string {
  const classes: string[] = []
  if (day.isToday) classes.push('calendar-cell-today')
  if (day.checkedIn) {
    if (day.rewardType === 'luck') classes.push('calendar-cell-luck')
    else classes.push('calendar-cell-checked')
  } else if (day.isCurrentMonth) {
    classes.push('calendar-cell-missed')
  }
  return classes.join(' ')
}

const streakColor = computed(() => {
  const d = checkinStore.streakDays
  if (d >= 30) return 'text-orange-500 dark:text-orange-400'
  if (d >= 7) return 'text-amber-500 dark:text-amber-400'
  return 'text-gray-500 dark:text-dark-400'
})

const streakFlameClass = computed(() => {
  const d = checkinStore.streakDays
  if (d >= 100) return 'h-6 w-6 animate-pulse'
  if (d >= 30) return 'h-5 w-5'
  if (d >= 7) return 'h-4 w-4'
  return 'h-3.5 w-3.5'
})

const nextBlindboxHint = computed(() => {
  const s = checkinStore.status
  if (!s?.blindbox_enabled || !s.blindbox_interval) return ''
  const interval = s.blindbox_interval
  if (s.blindbox_trigger_type === 'total') return ''
  const streak = s.streak_days || 0
  const next = interval - (streak % interval)
  if (next === interval && streak > 0) return t('checkin.page.blindboxNextCycle', { interval })
  return t('checkin.page.blindboxNextIn', { days: next })
})

async function submitLuck() {
  if (!luckBet.value || luckBet.value <= 0) return
  const result = await checkinStore.doLuckCheckin(luckBet.value)
  if (result) {
    showLuckModal.value = false
    luckBet.value = 0
    fetchBlindboxRecords()
  }
}

function getRarityBadgeClass(rarity: string) {
  const map: Record<string, string> = { common: 'badge-common', rare: 'badge-rare', epic: 'badge-epic', legendary: 'badge-legendary' }
  return map[rarity] || 'badge-common'
}

function getRarityLabel(rarity: string) {
  const map: Record<string, string> = { common: t('checkin.blindboxCommon'), rare: t('checkin.blindboxRare'), epic: t('checkin.blindboxEpic'), legendary: t('checkin.blindboxLegendary') }
  return map[rarity] || rarity
}

function getRecordBorderClass(rarity: string) {
  const map: Record<string, string> = {
    common: 'border-gray-200 dark:border-dark-700',
    rare: 'border-blue-200 dark:border-blue-900/50',
    epic: 'border-purple-200 dark:border-purple-900/50',
    legendary: 'border-amber-200 dark:border-amber-900/50',
  }
  return map[rarity] || map.common
}

function getRecordRewardColor(record: { rarity: string }) {
  const map: Record<string, string> = {
    common: 'text-gray-600 dark:text-gray-400',
    rare: 'text-blue-600 dark:text-blue-400',
    epic: 'text-purple-600 dark:text-purple-400',
    legendary: 'text-amber-600 dark:text-amber-400',
  }
  return map[record.rarity] || map.common
}

function formatRecordReward(record: BlindboxRecordItem) {
  switch (record.reward_type) {
    case 'balance': return `+$${record.reward_value.toFixed(2)}`
    case 'concurrency': return `+${Math.round(record.reward_value)}`
    case 'subscription': return `${record.subscription_days || Math.round(record.reward_value)}${t('checkin.blindboxDays')}`
    case 'invitation_code': return '×1'
    default: return `${record.reward_value}`
  }
}

function formatBlindboxReward(result: BlindboxResult) {
  switch (result.reward_type) {
    case 'balance': return `+$${result.reward_value.toFixed(2)}`
    case 'concurrency': return `+${Math.round(result.reward_value)}`
    case 'subscription': return `${result.subscription_days || Math.round(result.reward_value)}${t('checkin.blindboxDays')}`
    case 'invitation_code': return '×1'
    default: return `${result.reward_value}`
  }
}

const rarityBreakdown = computed(() => {
  const total = blindboxTotal.value || blindboxRecords.value.length
  const counts: Record<string, number> = { common: 0, rare: 0, epic: 0, legendary: 0 }
  blindboxRecords.value.forEach(r => { if (counts[r.rarity] !== undefined) counts[r.rarity]++ })

  const barClasses: Record<string, string> = {
    common: 'bg-gray-400 dark:bg-gray-500',
    rare: 'bg-blue-500',
    epic: 'bg-purple-500',
    legendary: 'bg-amber-500',
  }

  return (['common', 'rare', 'epic', 'legendary'] as const).map(key => ({
    key,
    label: getRarityLabel(key),
    badgeClass: getRarityBadgeClass(key),
    barClass: barClasses[key],
    count: counts[key],
    percent: total > 0 ? (counts[key] / total) * 100 : 0,
  }))
})

async function fetchCalendar() {
  try {
    const result = await getCheckinCalendar()
    calendarDays.value = result.days || []
  } catch { /* noop */ }
}

async function fetchBlindboxRecords() {
  try {
    const result = await getBlindboxRecords(blindboxPage.value, 20)
    blindboxRecords.value = result.items || []
    blindboxTotal.value = result.total || 0
  } catch { /* noop */ }
}

async function loadMoreBlindboxRecords() {
  blindboxPage.value++
  try {
    const result = await getBlindboxRecords(blindboxPage.value, 20)
    blindboxRecords.value = [...blindboxRecords.value, ...(result.items || [])]
    blindboxTotal.value = result.total || 0
  } catch { /* noop */ }
}

onMounted(() => {
  checkinStore.fetchStatus()
  fetchCalendar()
  fetchBlindboxRecords()
})
</script>

<style scoped>
.blindbox-rarity-badge {
  padding: 2px 10px;
  border-radius: 9999px;
  font-weight: 600;
  letter-spacing: 0.3px;
}
.badge-common { background-color: #f3f4f6; color: #6b7280; }
.badge-rare { background-color: #dbeafe; color: #2563eb; }
.badge-epic { background-color: #ede9fe; color: #7c3aed; }
.badge-legendary { background-color: #fef3c7; color: #d97706; }
html.dark .badge-common { background-color: #374151; color: #9ca3af; }
html.dark .badge-rare { background-color: #1e3a5f; color: #60a5fa; }
html.dark .badge-epic { background-color: #2d1b69; color: #a78bfa; }
html.dark .badge-legendary { background-color: #451a03; color: #fbbf24; }

.checkin-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 1.25rem;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}
.checkin-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.checkin-btn-normal { background: linear-gradient(135deg, #f59e0b, #d97706); }
.checkin-btn-normal:hover:not(:disabled) { background: linear-gradient(135deg, #d97706, #b45309); }
.checkin-btn-luck { background: linear-gradient(135deg, #8b5cf6, #7c3aed); }
.checkin-btn-luck:hover:not(:disabled) { background: linear-gradient(135deg, #7c3aed, #6d28d9); }

.calendar-cell { min-height: 40px; }
.calendar-cell-checked { background-color: #ecfdf5; border: 1px solid #a7f3d0; }
.calendar-cell-luck { background-color: #f5f3ff; border: 1px solid #c4b5fd; }
.calendar-cell-missed { background-color: transparent; border: 1px solid #f3f4f6; }
.calendar-cell-today { box-shadow: 0 0 0 2px #3b82f6; z-index: 1; }
html.dark .calendar-cell-checked { background-color: rgba(16, 185, 129, 0.1); border-color: rgba(16, 185, 129, 0.3); }
html.dark .calendar-cell-luck { background-color: rgba(139, 92, 246, 0.1); border-color: rgba(139, 92, 246, 0.3); }
html.dark .calendar-cell-missed { border-color: #374151; }
html.dark .calendar-cell-today { box-shadow: 0 0 0 2px #60a5fa; }
</style>
