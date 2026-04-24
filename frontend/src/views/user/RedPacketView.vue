<template>
  <AppLayout>
    <div class="mx-auto max-w-2xl space-y-6">
      <section
        class="relative overflow-hidden rounded-[30px] bg-gradient-to-r from-[#ff355c] via-[#ff5b46] to-[#ff8c30] px-6 py-8 shadow-[0_28px_70px_-36px_rgba(255,92,62,0.9)]"
      >
        <div class="pointer-events-none absolute inset-0">
          <div class="hero-glow hero-glow-left"></div>
          <div class="hero-glow hero-glow-right"></div>

          <svg
            class="hero-ticket hero-ticket-left"
            viewBox="0 0 72 72"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <rect x="12" y="14" width="48" height="44" rx="8" fill="rgba(255,255,255,0.12)" />
            <path d="M24 22H48V50H24V22Z" fill="rgba(255,255,255,0.14)" />
            <path d="M28 28H44" stroke="rgba(255,255,255,0.55)" stroke-width="3" stroke-linecap="round" />
            <path d="M28 36H40" stroke="rgba(255,255,255,0.45)" stroke-width="3" stroke-linecap="round" />
            <path d="M28 44H36" stroke="rgba(255,255,255,0.35)" stroke-width="3" stroke-linecap="round" />
          </svg>

          <svg
            class="hero-ticket hero-ticket-right"
            viewBox="0 0 72 72"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <rect x="12" y="14" width="48" height="44" rx="8" fill="rgba(255,255,255,0.12)" />
            <path d="M24 22H48V50H24V22Z" fill="rgba(255,255,255,0.14)" />
            <path d="M28 28H44" stroke="rgba(255,255,255,0.55)" stroke-width="3" stroke-linecap="round" />
            <path d="M28 36H40" stroke="rgba(255,255,255,0.45)" stroke-width="3" stroke-linecap="round" />
            <path d="M28 44H36" stroke="rgba(255,255,255,0.35)" stroke-width="3" stroke-linecap="round" />
          </svg>

          <svg
            class="hero-coin hero-coin-left"
            viewBox="0 0 52 52"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <ellipse cx="26" cy="26" rx="17" ry="11" transform="rotate(28 26 26)" fill="#FFE18E" />
            <ellipse cx="26" cy="26" rx="12" ry="7" transform="rotate(28 26 26)" fill="#FFD15A" />
            <path d="M26 20V32" stroke="#FF9E2F" stroke-width="2.6" stroke-linecap="round" />
            <path d="M21 26H31" stroke="#FF9E2F" stroke-width="2.6" stroke-linecap="round" />
          </svg>

          <svg
            class="hero-coin hero-coin-right"
            viewBox="0 0 52 52"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <ellipse cx="26" cy="26" rx="17" ry="11" transform="rotate(-22 26 26)" fill="#FFE18E" />
            <ellipse cx="26" cy="26" rx="12" ry="7" transform="rotate(-22 26 26)" fill="#FFD15A" />
            <path d="M26 20V32" stroke="#FF9E2F" stroke-width="2.6" stroke-linecap="round" />
            <path d="M21 26H31" stroke="#FF9E2F" stroke-width="2.6" stroke-linecap="round" />
          </svg>

          <span class="hero-spark hero-spark-a"></span>
          <span class="hero-spark hero-spark-b"></span>
          <span class="hero-spark hero-spark-c"></span>
        </div>

        <div class="relative flex flex-col items-center text-center">
          <div
            class="mb-5 flex h-24 w-24 items-center justify-center rounded-[28px] bg-white/15 shadow-[0_18px_30px_-18px_rgba(0,0,0,0.45)] ring-1 ring-white/30 backdrop-blur-sm"
          >
            <svg viewBox="0 0 120 120" class="h-20 w-20" xmlns="http://www.w3.org/2000/svg">
              <defs>
                <linearGradient id="heroPacketBody" x1="14" y1="10" x2="108" y2="108" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#FF7A6E" />
                  <stop offset="0.55" stop-color="#FF4B56" />
                  <stop offset="1" stop-color="#E82F43" />
                </linearGradient>
                <linearGradient id="heroPacketLip" x1="22" y1="12" x2="97" y2="52" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#FFB0A4" />
                  <stop offset="1" stop-color="#FF6173" />
                </linearGradient>
                <linearGradient id="heroSeal" x1="44" y1="45" x2="78" y2="79" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#FFE48D" />
                  <stop offset="1" stop-color="#FFC54D" />
                </linearGradient>
              </defs>
              <rect x="18" y="20" width="84" height="78" rx="18" fill="url(#heroPacketBody)" />
              <rect x="18" y="20" width="84" height="28" rx="18" fill="url(#heroPacketLip)" />
              <path d="M20 39C42 54 78 54 100 39" stroke="#FFCC63" stroke-width="4" stroke-linecap="round" />
              <path d="M60 40L34 54H86L60 40Z" fill="#FF6A5C" opacity="0.55" />
              <circle cx="60" cy="62" r="18" fill="url(#heroSeal)" stroke="#FFB22B" stroke-width="3" />
              <path d="M52 58H68" stroke="#FF8F20" stroke-width="3.5" stroke-linecap="round" />
              <path d="M60 50V74" stroke="#FF8F20" stroke-width="3.5" stroke-linecap="round" />
              <path d="M50 70C54 66 66 66 70 70" stroke="#FF8F20" stroke-width="3.5" stroke-linecap="round" />
            </svg>
          </div>

          <h1 class="text-3xl font-bold tracking-tight text-white sm:text-[2.25rem]">
            {{ t('redpacket.title', '红包中心') }}
          </h1>
          <p class="mt-3 text-base text-white/90 sm:text-lg">
            {{ t('redpacket.subtitle', '发红包，领红包，分享快乐') }}
          </p>
        </div>
      </section>

      <section class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <button
          type="button"
          @click="openCreate"
          :class="actionCardClass('create')"
        >
          <div class="flex items-center gap-4">
            <div class="action-illustration action-illustration-create">
              <svg viewBox="0 0 88 88" class="h-14 w-14" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="createPacketBody" x1="16" y1="14" x2="74" y2="74" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#FF8C95" />
                    <stop offset="1" stop-color="#FF4764" />
                  </linearGradient>
                </defs>
                <circle cx="44" cy="44" r="34" fill="#FFE6EA" />
                <rect x="25" y="24" width="38" height="40" rx="12" fill="url(#createPacketBody)" />
                <path d="M28 34C38 42 50 42 60 34" stroke="#FFC96A" stroke-width="3" stroke-linecap="round" />
                <circle cx="44" cy="44" r="10" fill="#FF5B72" stroke="#FFB744" stroke-width="2" />
                <path d="M44 38V50" stroke="white" stroke-width="3.5" stroke-linecap="round" />
                <path d="M38 44H50" stroke="white" stroke-width="3.5" stroke-linecap="round" />
              </svg>
            </div>

            <div class="min-w-0 flex-1 text-left">
              <p class="text-xl font-semibold text-slate-900 dark:text-white">
                {{ t('redpacket.create', '发红包') }}
              </p>
              <p class="mt-1 text-sm text-slate-500 dark:text-dark-400">
                {{ t('redpacket.createDesc', '发送余额红包给好友') }}
              </p>
            </div>
          </div>

          <span class="action-arrow">
            <svg viewBox="0 0 20 20" class="h-4 w-4" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 10H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              <path d="M10.5 4.5L16 10L10.5 15.5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </span>
        </button>

        <button
          type="button"
          @click="openClaim"
          :class="actionCardClass('claim')"
        >
          <div class="flex items-center gap-4">
            <div class="action-illustration action-illustration-claim">
              <svg viewBox="0 0 88 88" class="h-14 w-14" xmlns="http://www.w3.org/2000/svg">
                <circle cx="44" cy="44" r="34" fill="#FFF2D9" />
                <rect x="24" y="34" width="40" height="27" rx="8" fill="#FFBF3F" />
                <rect x="28" y="28" width="32" height="13" rx="6.5" fill="#FFD56D" />
                <path d="M44 28V61" stroke="#FF9E18" stroke-width="4" stroke-linecap="round" />
                <path d="M24 44H64" stroke="#FF9E18" stroke-width="4" stroke-linecap="round" />
                <path d="M44 28C44 22 37 19 33 24C29 29 35 34 44 34" stroke="#FF9E18" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" />
                <path d="M44 28C44 22 51 19 55 24C59 29 53 34 44 34" stroke="#FF9E18" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </div>

            <div class="min-w-0 flex-1 text-left">
              <p class="text-xl font-semibold text-slate-900 dark:text-white">
                {{ t('redpacket.claim', '领红包') }}
              </p>
              <p class="mt-1 text-sm text-slate-500 dark:text-dark-400">
                {{ t('redpacket.claimDesc', '输入口令领取红包') }}
              </p>
            </div>
          </div>

          <span class="action-arrow action-arrow-amber">
            <svg viewBox="0 0 20 20" class="h-4 w-4" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 10H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              <path d="M10.5 4.5L16 10L10.5 15.5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </span>
        </button>
      </section>

      <transition name="fade">
        <section
          v-if="showCreate"
          class="overflow-hidden rounded-[28px] border border-rose-100 bg-white shadow-[0_18px_40px_-28px_rgba(255,88,102,0.5)] dark:border-dark-700 dark:bg-dark-900"
        >
          <div class="border-b border-rose-100 bg-rose-50/80 px-5 py-4 dark:border-dark-700 dark:bg-dark-800/80">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-rose-100 dark:bg-rose-900/25">
                <Icon name="sparkles" size="md" class="text-rose-500" />
              </div>
              <div>
                <p class="text-base font-semibold text-slate-900 dark:text-white">
                  {{ t('redpacket.create', '发红包') }}
                </p>
                <p class="text-xs text-slate-500 dark:text-dark-400">
                  {{ t('redpacket.createPanelHint', '填写金额、份数和附言，系统会生成可分享口令') }}
                </p>
              </div>
            </div>
          </div>

          <div class="space-y-5 p-5">
            <form @submit.prevent="handleCreate" class="space-y-5">
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label class="input-label">{{ t('redpacket.totalAmount', '总金额') }}</label>
                  <div class="relative mt-1">
                    <span class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-sm text-slate-400">$</span>
                    <input
                      v-model.number="createForm.total_amount"
                      type="number"
                      step="0.01"
                      min="0.01"
                      required
                      :disabled="createLoading"
                      class="input pl-7"
                    />
                  </div>
                </div>

                <div>
                  <label class="input-label">{{ t('redpacket.count', '份数') }}</label>
                  <input
                    v-model.number="createForm.count"
                    type="number"
                    min="1"
                    max="100"
                    required
                    :disabled="createLoading"
                    class="input mt-1"
                  />
                </div>
              </div>

              <div>
                <label class="input-label">{{ t('redpacket.type', '类型') }}</label>
                <div class="mt-1 grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    @click="createForm.redpacket_type = 'equal'"
                    :class="typeToggleClass(createForm.redpacket_type === 'equal')"
                  >
                    <span class="text-base">🀄</span>
                    {{ t('redpacket.equal', '等分红包') }}
                  </button>
                  <button
                    type="button"
                    @click="createForm.redpacket_type = 'random'"
                    :class="typeToggleClass(createForm.redpacket_type === 'random')"
                  >
                    <span class="text-base">🎲</span>
                    {{ t('redpacket.random', '拼手气红包') }}
                  </button>
                </div>
              </div>

              <div>
                <label class="input-label">{{ t('redpacket.memo', '附言') }}</label>
                <input
                  v-model="createForm.memo"
                  type="text"
                  maxlength="100"
                  :placeholder="t('redpacket.memoPlaceholder', '恭喜发财，大吉大利！')"
                  :disabled="createLoading"
                  class="input mt-1 w-full"
                />
              </div>

              <p
                v-if="createError"
                class="rounded-2xl border border-rose-100 bg-rose-50 px-3 py-2 text-sm text-rose-600 dark:border-rose-900/30 dark:bg-rose-900/20 dark:text-rose-300"
              >
                {{ createError }}
              </p>

              <div class="flex flex-wrap gap-3">
                <button type="submit" :disabled="createLoading" class="btn btn-primary min-w-32 flex-1 sm:flex-none">
                  <svg v-if="createLoading" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ t('redpacket.create', '发红包') }}
                </button>
                <button type="button" @click="showCreate = false" class="btn btn-secondary min-w-24">
                  {{ t('common.cancel', '取消') }}
                </button>
              </div>
            </form>

            <div
              v-if="createdRp"
              class="rounded-[24px] border border-emerald-100 bg-emerald-50/80 p-4 dark:border-emerald-900/30 dark:bg-emerald-900/10"
            >
              <div class="flex items-start gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-white text-emerald-500 shadow-sm dark:bg-dark-800">
                  <Icon name="checkCircle" size="md" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm font-semibold text-emerald-700 dark:text-emerald-300">
                    {{ t('redpacket.created', '红包已创建！') }}
                  </p>
                  <div class="mt-3 flex items-center gap-2 rounded-2xl bg-white px-3 py-3 shadow-sm dark:bg-dark-800">
                    <code class="min-w-0 flex-1 truncate font-mono text-base font-semibold text-slate-900 dark:text-white">
                      {{ createdRp.code }}
                    </code>
                    <button
                      type="button"
                      @click="copyCode(createdRp.code)"
                      class="flex h-9 w-9 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600 transition-colors hover:bg-emerald-100 dark:bg-emerald-900/20 dark:text-emerald-300"
                    >
                      <Icon :name="copiedCode === createdRp.code ? 'checkCircle' : 'copy'" size="sm" />
                    </button>
                  </div>
                  <p class="mt-2 text-xs text-emerald-600/80 dark:text-emerald-400/70">
                    {{ t('redpacket.shareCode', '将此口令分享给好友即可领取') }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </transition>

      <transition name="fade">
        <section
          v-if="showClaim"
          class="overflow-hidden rounded-[28px] border border-amber-100 bg-white shadow-[0_18px_40px_-28px_rgba(255,181,66,0.5)] dark:border-dark-700 dark:bg-dark-900"
        >
          <div class="border-b border-amber-100 bg-amber-50/80 px-5 py-4 dark:border-dark-700 dark:bg-dark-800/80">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-amber-100 dark:bg-amber-900/25">
                <Icon name="gift" size="md" class="text-amber-500" />
              </div>
              <div>
                <p class="text-base font-semibold text-slate-900 dark:text-white">
                  {{ t('redpacket.claim', '领红包') }}
                </p>
                <p class="text-xs text-slate-500 dark:text-dark-400">
                  {{ t('redpacket.claimPanelHint', '输入好友分享的口令，系统会自动入账到你的余额') }}
                </p>
              </div>
            </div>
          </div>

          <div class="space-y-5 p-5">
            <form @submit.prevent="handleClaim" class="space-y-5">
              <div>
                <label class="input-label">{{ t('redpacket.code', '红包口令') }}</label>
                <input
                  v-model="claimCode"
                  type="text"
                  required
                  :placeholder="t('redpacket.codePlaceholder', '输入红包口令')"
                  :disabled="claimLoading"
                  class="input mt-1 w-full text-center font-mono text-lg tracking-[0.22em]"
                />
              </div>

              <div
                v-if="claimResult"
                class="rounded-[24px] border border-emerald-100 bg-gradient-to-br from-emerald-50 to-white p-6 text-center dark:border-emerald-900/30 dark:from-emerald-900/10 dark:to-dark-900"
              >
                <p class="text-xs uppercase tracking-[0.22em] text-emerald-500">
                  {{ t('redpacket.congrats', '恭喜！获得') }}
                </p>
                <p class="mt-2 text-4xl font-bold text-emerald-600 dark:text-emerald-400">
                  +${{ claimResult.amount.toFixed(2) }}
                </p>
              </div>

              <p
                v-if="claimError"
                class="rounded-2xl border border-rose-100 bg-rose-50 px-3 py-2 text-sm text-rose-600 dark:border-rose-900/30 dark:bg-rose-900/20 dark:text-rose-300"
              >
                {{ claimError }}
              </p>

              <div class="flex flex-wrap gap-3">
                <button type="submit" :disabled="claimLoading || !claimCode" class="btn btn-primary min-w-32 flex-1 sm:flex-none">
                  <svg v-if="claimLoading" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ t('redpacket.claim', '领取') }}
                </button>
                <button type="button" @click="showClaim = false; claimResult = null; claimError = ''" class="btn btn-secondary min-w-24">
                  {{ t('common.cancel', '取消') }}
                </button>
              </div>
            </form>
          </div>
        </section>
      </transition>

      <section class="overflow-hidden rounded-[30px] border border-gray-100 bg-white shadow-[0_24px_60px_-40px_rgba(15,23,42,0.2)] dark:border-dark-700 dark:bg-dark-900">
        <div class="border-b border-gray-100 px-6 py-5 dark:border-dark-700">
          <h2 class="text-xl font-semibold text-slate-900 dark:text-white">
            {{ t('redpacket.myPackets', '我的红包') }}
          </h2>
        </div>

        <div class="space-y-4 p-5">
          <div v-if="loadingPackets" class="flex items-center justify-center py-12">
            <svg class="h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>

          <template v-else-if="myPackets.length > 0">
            <article
              v-for="rp in myPackets"
              :key="rp.id"
              class="overflow-hidden rounded-[24px] border border-gray-100 bg-white shadow-[0_16px_44px_-34px_rgba(15,23,42,0.26)] transition-all hover:border-gray-200 dark:border-dark-700 dark:bg-dark-800 dark:hover:border-dark-600"
            >
              <div class="flex items-start justify-between gap-4 px-5 py-5">
                <div class="flex min-w-0 items-start gap-4">
                  <div class="packet-avatar-shell">
                    <svg viewBox="0 0 88 88" class="h-12 w-12" xmlns="http://www.w3.org/2000/svg">
                      <rect x="22" y="18" width="44" height="50" rx="14" fill="#F8485E" />
                      <rect x="22" y="18" width="44" height="16" rx="14" fill="#FFA29A" />
                      <path d="M24 29C35 36 53 36 64 29" stroke="#FFCC63" stroke-width="3" stroke-linecap="round" />
                      <circle cx="44" cy="44" r="10" fill="#FFC94E" stroke="#FFB22B" stroke-width="2" />
                      <path d="M44 39V49" stroke="#FF9021" stroke-width="2.8" stroke-linecap="round" />
                      <path d="M39 44H49" stroke="#FF9021" stroke-width="2.8" stroke-linecap="round" />
                    </svg>
                  </div>

                  <div class="min-w-0 flex-1">
                    <div class="flex flex-wrap items-center gap-2">
                      <p class="text-xl font-semibold text-slate-900 dark:text-white">
                        {{ packetTypeLabel(rp.redpacket_type) }}
                      </p>
                      <span :class="statusBadgeClass(rp.status)">
                        {{ rpStatusLabel(rp.status) }}
                      </span>
                      <span :class="packetTypeBadgeClass(rp.redpacket_type)">
                        {{ packetTypeBadgeLabel(rp.redpacket_type) }}
                      </span>
                    </div>

                    <p class="mt-2 text-sm text-slate-500 dark:text-dark-400">
                      {{ packetProgressText(rp) }}
                    </p>
                    <p v-if="rp.memo" class="mt-1 text-sm text-slate-400 dark:text-dark-500">
                      「{{ rp.memo }}」
                    </p>
                  </div>
                </div>

                <div class="flex items-start gap-3">
                  <div class="text-right">
                    <p class="text-[1.9rem] font-bold leading-none text-[#ff2b38] dark:text-[#ff6374]">
                      ${{ rp.total_amount.toFixed(2) }}
                    </p>
                    <p class="mt-2 text-sm text-slate-500 dark:text-dark-400">
                      {{ t('redpacket.remainingAmount', '剩余') }} ${{ rp.remaining_amount.toFixed(2) }}
                    </p>
                  </div>

                  <button
                    type="button"
                    @click="toggleDetail(rp)"
                    :class="[
                      'mt-1 flex h-8 w-8 items-center justify-center rounded-xl text-slate-400 transition-all hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-dark-700 dark:hover:text-dark-300',
                      detailExpandedId === rp.id ? 'rotate-180' : ''
                    ]"
                  >
                    <Icon name="chevronDown" size="sm" />
                  </button>
                </div>
              </div>

              <div class="px-5 pb-5">
                <div class="rounded-2xl bg-slate-50 px-4 py-3 dark:bg-dark-700/60">
                  <div class="flex items-center gap-2">
                    <code class="min-w-0 flex-1 truncate text-xs text-slate-600 dark:text-dark-300">
                      {{ rp.code }}
                    </code>
                    <button
                      type="button"
                      @click="copyCode(rp.code)"
                      class="flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 transition-colors hover:bg-white hover:text-primary-500 dark:hover:bg-dark-800"
                    >
                      <Icon :name="copiedCode === rp.code ? 'checkCircle' : 'copy'" size="xs" />
                    </button>
                  </div>

                  <div class="mt-3 h-2 overflow-hidden rounded-full bg-white dark:bg-dark-800">
                    <div
                      class="h-full rounded-full bg-gradient-to-r from-[#ff4a67] to-[#ff8b2f] transition-all duration-300"
                      :style="{ width: `${packetClaimedPercent(rp)}%` }"
                    ></div>
                  </div>
                </div>
              </div>

              <transition name="fade">
                <div v-if="detailExpandedId === rp.id" class="border-t border-dashed border-gray-100 bg-slate-50/70 px-5 py-4 dark:border-dark-700 dark:bg-dark-900/60">
                  <div v-if="loadingDetail" class="flex items-center justify-center py-6">
                    <svg class="h-5 w-5 animate-spin text-slate-400" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>

                  <div v-else-if="detailClaims.length > 0" class="space-y-3">
                    <p class="text-xs font-medium uppercase tracking-[0.18em] text-slate-400 dark:text-dark-500">
                      {{ t('redpacket.claimDetail', '领取详情') }} ({{ detailClaims.length }}/{{ rp.total_count }})
                    </p>

                    <div class="space-y-2">
                      <div
                        v-for="claim in detailClaims"
                        :key="claim.id"
                        class="flex items-center justify-between rounded-2xl bg-white px-4 py-3 shadow-sm dark:bg-dark-800"
                      >
                        <div class="flex min-w-0 items-center gap-3">
                          <div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary-50 text-primary-500 dark:bg-primary-900/20 dark:text-primary-300">
                            <Icon name="user" size="xs" />
                          </div>
                          <div class="min-w-0">
                            <p class="truncate text-sm font-medium text-slate-700 dark:text-dark-200">
                              {{ claim.user_email || '#' + claim.user_id }}
                            </p>
                            <p class="text-xs text-slate-400 dark:text-dark-500">
                              {{ formatDateTime(claim.created_at) }}
                            </p>
                          </div>
                        </div>

                        <span class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">
                          +${{ claim.amount.toFixed(2) }}
                        </span>
                      </div>
                    </div>

                    <div
                      v-if="rp.status === 'active' && rp.remaining_count > 0"
                      class="rounded-2xl border border-dashed border-slate-200 px-4 py-3 text-center text-sm text-slate-500 dark:border-dark-700 dark:text-dark-400"
                    >
                      还有 {{ rp.remaining_count }} 份等待领取
                    </div>
                  </div>

                  <div v-else class="rounded-2xl border border-dashed border-slate-200 px-4 py-6 text-center text-sm text-slate-400 dark:border-dark-700 dark:text-dark-500">
                    {{ t('redpacket.noClaimsYet', '暂无人领取') }}
                  </div>
                </div>
              </transition>
            </article>
          </template>

          <div v-else class="py-12 text-center">
            <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-3xl bg-slate-100 dark:bg-dark-800">
              <svg viewBox="0 0 88 88" class="h-10 w-10 opacity-60" xmlns="http://www.w3.org/2000/svg">
                <rect x="22" y="18" width="44" height="50" rx="14" fill="#FF7278" />
                <rect x="22" y="18" width="44" height="16" rx="14" fill="#FFA29A" />
                <path d="M24 29C35 36 53 36 64 29" stroke="#FFCC63" stroke-width="3" stroke-linecap="round" />
                <circle cx="44" cy="44" r="10" fill="#FFC94E" stroke="#FFB22B" stroke-width="2" />
              </svg>
            </div>
            <p class="text-sm text-slate-400 dark:text-dark-500">
              {{ t('redpacket.noPackets', '暂无红包记录') }}
            </p>
          </div>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { createRedPacket, claimRedPacket, getMyRedPackets, getRedPacketDetail } from '@/api/transfer'
import type { RedPacketRecord, RedPacketClaimRecord } from '@/api/transfer'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const showCreate = ref(false)
const showClaim = ref(false)
const claimCode = ref('')
const claimResult = ref<RedPacketClaimRecord | null>(null)
const claimError = ref('')
const claimLoading = ref(false)
const createError = ref('')
const createLoading = ref(false)
const createdRp = ref<RedPacketRecord | null>(null)
const myPackets = ref<RedPacketRecord[]>([])
const loadingPackets = ref(false)
const copiedCode = ref('')

const detailExpandedId = ref<number | null>(null)
const detailClaims = ref<RedPacketClaimRecord[]>([])
const loadingDetail = ref(false)

const createForm = reactive({
  total_amount: 0,
  count: 1,
  redpacket_type: 'equal' as 'equal' | 'random',
  memo: '',
})

function openCreate() {
  const next = !showCreate.value
  showCreate.value = next
  showClaim.value = false
  createError.value = ''
  if (!next) {
    createdRp.value = null
  }
}

function openClaim() {
  const next = !showClaim.value
  showClaim.value = next
  showCreate.value = false
  claimError.value = ''
  if (!next) {
    claimResult.value = null
    claimCode.value = ''
  }
}

function actionCardClass(kind: 'create' | 'claim') {
  const active = kind === 'create' ? showCreate.value : showClaim.value
  const tone = kind === 'create'
    ? 'hover:border-rose-200 hover:shadow-rose-100/70 dark:hover:border-rose-800/60'
    : 'hover:border-amber-200 hover:shadow-amber-100/70 dark:hover:border-amber-800/60'
  const activeTone = kind === 'create'
    ? 'border-rose-200 shadow-[0_18px_40px_-30px_rgba(255,92,111,0.55)] dark:border-rose-800/60'
    : 'border-amber-200 shadow-[0_18px_40px_-30px_rgba(255,185,73,0.55)] dark:border-amber-800/60'
  return [
    'group flex items-center justify-between rounded-[26px] border bg-white px-5 py-5 text-left transition-all dark:bg-dark-900',
    active ? activeTone : `border-gray-100 shadow-[0_14px_36px_-32px_rgba(15,23,42,0.28)] ${tone}`,
  ]
}

function typeToggleClass(active: boolean) {
  return [
    'flex items-center justify-center gap-2 rounded-2xl border px-3 py-3 text-sm font-medium transition-all',
    active
      ? 'border-rose-200 bg-rose-50 text-rose-700 dark:border-rose-800/60 dark:bg-rose-900/20 dark:text-rose-300'
      : 'border-gray-200 bg-white text-slate-600 hover:bg-slate-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300 dark:hover:bg-dark-700',
  ]
}

function rpStatusStyle(status: string) {
  switch (status) {
    case 'active':
      return {
        badge: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
      }
    case 'exhausted':
      return {
        badge: 'bg-slate-100 text-slate-600 dark:bg-dark-700 dark:text-dark-300',
      }
    case 'expired':
      return {
        badge: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
      }
    default:
      return {
        badge: 'bg-slate-100 text-slate-600 dark:bg-dark-700 dark:text-dark-300',
      }
  }
}

function rpStatusLabel(status: string) {
  switch (status) {
    case 'active':
      return t('redpacket.statusActive', '进行中')
    case 'exhausted':
      return t('redpacket.statusExhausted', '已领完')
    case 'expired':
      return t('redpacket.statusExpired', '已过期')
    default:
      return status
  }
}

function statusBadgeClass(status: string) {
  return ['inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold', rpStatusStyle(status).badge]
}

function packetTypeLabel(type: 'equal' | 'random') {
  return type === 'equal' ? t('redpacket.equal', '等分红包') : t('redpacket.random', '拼手气')
}

function packetTypeBadgeClass(type: 'equal' | 'random') {
  return [
    'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold',
    type === 'equal'
      ? 'bg-rose-50 text-rose-600 dark:bg-rose-900/20 dark:text-rose-300'
      : 'bg-orange-50 text-orange-600 dark:bg-orange-900/20 dark:text-orange-300',
  ]
}

function packetTypeBadgeLabel(type: 'equal' | 'random') {
  return type === 'equal' ? t('redpacket.equalBadge', '均分') : t('redpacket.randomBadge', '限领红包')
}

function packetClaimedPercent(rp: RedPacketRecord) {
  if (rp.total_count <= 0) {
    return 0
  }
  return Math.max(0, Math.min(100, ((rp.total_count - rp.remaining_count) / rp.total_count) * 100))
}

function packetProgressText(rp: RedPacketRecord) {
  return `${rp.total_count}${t('redpacket.copies', '份')} · ${t('redpacket.remaining', '剩余')} ${rp.remaining_count}${t('redpacket.copies', '份')}`
}

async function copyCode(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    copiedCode.value = code
    setTimeout(() => {
      copiedCode.value = ''
    }, 2000)
  } catch {}
}

async function toggleDetail(rp: RedPacketRecord) {
  if (detailExpandedId.value === rp.id) {
    detailExpandedId.value = null
    detailClaims.value = []
    return
  }
  detailExpandedId.value = rp.id
  loadingDetail.value = true
  detailClaims.value = []
  try {
    const res = await getRedPacketDetail(rp.id)
    detailClaims.value = res.claims || []
  } catch {
    detailClaims.value = []
  } finally {
    loadingDetail.value = false
  }
}

async function loadMyPackets() {
  loadingPackets.value = true
  try {
    const res = await getMyRedPackets({ page: 1, page_size: 20 })
    myPackets.value = res.items || []
  } catch {} finally {
    loadingPackets.value = false
  }
}

async function handleCreate() {
  createError.value = ''
  createLoading.value = true
  try {
    createdRp.value = await createRedPacket({
      total_amount: createForm.total_amount,
      count: createForm.count,
      redpacket_type: createForm.redpacket_type,
      memo: createForm.memo || undefined,
    })
    appStore.showSuccess(t('redpacket.created', '红包创建成功'))
    loadMyPackets().catch(() => {})
    authStore.refreshUser().catch(() => {})
  } catch (e: any) {
    createError.value = e?.response?.data?.error || t('redpacket.createFailed', '创建失败')
  } finally {
    createLoading.value = false
  }
}

async function handleClaim() {
  claimError.value = ''
  claimResult.value = null
  claimLoading.value = true
  try {
    claimResult.value = await claimRedPacket(claimCode.value)
    appStore.showSuccess(t('redpacket.claimSuccess', '领取成功！'))
    loadMyPackets().catch(() => {})
    authStore.refreshUser().catch(() => {})
  } catch (e: any) {
    claimError.value = e?.response?.data?.error || t('redpacket.claimFailed', '领取失败')
  } finally {
    claimLoading.value = false
  }
}

onMounted(loadMyPackets)
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.hero-glow {
  position: absolute;
  border-radius: 9999px;
  filter: blur(50px);
  opacity: 0.4;
}

.hero-glow-left {
  left: -2rem;
  top: -1rem;
  height: 9rem;
  width: 9rem;
  background: rgba(255, 255, 255, 0.22);
}

.hero-glow-right {
  bottom: -2rem;
  right: -1rem;
  height: 10rem;
  width: 10rem;
  background: rgba(255, 195, 120, 0.28);
}

.hero-ticket,
.hero-coin,
.hero-spark {
  position: absolute;
  animation: heroFloat 6s ease-in-out infinite;
}

.hero-ticket-left {
  left: 1.3rem;
  bottom: 3.5rem;
  height: 3.75rem;
  width: 3.75rem;
  transform: rotate(-18deg);
  opacity: 0.65;
}

.hero-ticket-right {
  right: 1.5rem;
  top: 4.1rem;
  height: 4.2rem;
  width: 4.2rem;
  transform: rotate(20deg);
  opacity: 0.55;
  animation-delay: 1.2s;
}

.hero-coin-left {
  left: 8rem;
  bottom: 2rem;
  height: 2.75rem;
  width: 2.75rem;
  opacity: 0.85;
  animation-delay: 0.7s;
}

.hero-coin-right {
  right: 11rem;
  top: 4rem;
  height: 2.8rem;
  width: 2.8rem;
  opacity: 0.9;
  animation-delay: 1.5s;
}

.hero-spark {
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.65);
}

.hero-spark-a {
  left: 5rem;
  top: 3.1rem;
  height: 0.55rem;
  width: 0.55rem;
}

.hero-spark-b {
  right: 5rem;
  top: 2.3rem;
  height: 0.6rem;
  width: 0.6rem;
  animation-delay: 1.3s;
}

.hero-spark-c {
  right: 13rem;
  bottom: 2.8rem;
  height: 0.45rem;
  width: 0.45rem;
  animation-delay: 2.2s;
}

.action-illustration {
  display: flex;
  height: 6rem;
  width: 6rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
}

.action-illustration-create {
  background: linear-gradient(135deg, #fff0f2 0%, #ffe0e7 100%);
}

.action-illustration-claim {
  background: linear-gradient(135deg, #fff6e1 0%, #ffedd0 100%);
}

.action-arrow {
  display: inline-flex;
  height: 3rem;
  width: 3rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  background: #fff0f2;
  color: #ff4f67;
  transition: transform 0.2s ease;
}

.group:hover .action-arrow {
  transform: translateX(3px);
}

.action-arrow-amber {
  background: #fff4dc;
  color: #f59e0b;
}

.packet-avatar-shell {
  display: flex;
  height: 4.75rem;
  width: 4.75rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 1.5rem;
  background: linear-gradient(135deg, #fff1f2 0%, #ffe4e8 100%);
}

@keyframes heroFloat {
  0%,
  100% {
    transform: translateY(0) rotate(var(--hero-rotate, 0deg));
  }
  50% {
    transform: translateY(-4px) rotate(var(--hero-rotate, 0deg));
  }
}

.hero-ticket-left {
  --hero-rotate: -18deg;
}

.hero-ticket-right {
  --hero-rotate: 20deg;
}

.hero-coin-left {
  --hero-rotate: 28deg;
}

.hero-coin-right {
  --hero-rotate: -22deg;
}
</style>
