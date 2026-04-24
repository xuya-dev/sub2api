<template>
  <AppLayout>
    <div class="mx-auto max-w-2xl space-y-6">
      <!-- Hero Banner -->
      <div class="card overflow-hidden">
        <div class="relative bg-gradient-to-br from-red-500 via-red-500 to-orange-500 px-6 py-10 text-center">
          <div class="absolute inset-0 opacity-10" style="background-image: url(&quot;data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E&quot;);"></div>
          <div class="relative">
            <div class="mb-3 inline-flex h-20 w-20 items-center justify-center rounded-3xl bg-white/20 shadow-lg shadow-red-500/30 backdrop-blur-sm">
              <span class="text-4xl">🧧</span>
            </div>
            <p class="text-2xl font-bold text-white">{{ t('redpacket.title', '红包中心') }}</p>
            <p class="mt-2 text-sm text-red-100/90">{{ t('redpacket.subtitle', '发红包、领红包，分享快乐') }}</p>
          </div>
        </div>
      </div>

      <!-- Action Cards -->
      <div class="grid grid-cols-2 gap-4">
        <button @click="openCreate"
          class="group card flex flex-col items-center gap-3 p-5 transition-all hover:shadow-md hover:shadow-red-100 dark:hover:shadow-red-900/10">
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-red-100 transition-transform group-hover:scale-110 dark:bg-red-900/30">
            <Icon name="plus" size="lg" class="text-red-600 dark:text-red-400" />
          </div>
          <div class="text-center">
            <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('redpacket.create', '发红包') }}</p>
            <p class="mt-0.5 text-xs text-gray-400 dark:text-dark-500">{{ t('redpacket.createDesc', '发送余额红包给好友') }}</p>
          </div>
        </button>
        <button @click="openClaim"
          class="group card flex flex-col items-center gap-3 p-5 transition-all hover:shadow-md hover:shadow-amber-100 dark:hover:shadow-amber-900/10">
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 transition-transform group-hover:scale-110 dark:bg-amber-900/30">
            <Icon name="gift" size="lg" class="text-amber-600 dark:text-amber-400" />
          </div>
          <div class="text-center">
            <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('redpacket.claim', '领红包') }}</p>
            <p class="mt-0.5 text-xs text-gray-400 dark:text-dark-500">{{ t('redpacket.claimDesc', '输入口令领取红包') }}</p>
          </div>
        </button>
      </div>

      <!-- Create Red Packet -->
      <transition name="fade">
        <div v-if="showCreate" class="card overflow-hidden">
          <div class="bg-gradient-to-r from-red-50 to-orange-50 px-5 py-4 dark:from-red-900/20 dark:to-orange-900/20">
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-red-100 dark:bg-red-900/40">
                <Icon name="sparkles" size="md" class="text-red-600 dark:text-red-400" />
              </div>
              <h3 class="text-sm font-semibold text-red-800 dark:text-red-300">{{ t('redpacket.create', '发红包') }}</h3>
            </div>
          </div>
          <div class="p-5">
            <form @submit.prevent="handleCreate" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="input-label">{{ t('redpacket.totalAmount', '总金额') }}</label>
                  <div class="relative mt-1">
                    <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <span class="text-sm text-gray-400">$</span>
                    </div>
                    <input v-model.number="createForm.total_amount" type="number" step="0.01" min="0.01" required
                      :disabled="createLoading" class="input pl-7" />
                  </div>
                </div>
                <div>
                  <label class="input-label">{{ t('redpacket.count', '份数') }}</label>
                  <input v-model.number="createForm.count" type="number" min="1" max="100" required
                    :disabled="createLoading" class="input mt-1" />
                </div>
              </div>

              <div>
                <label class="input-label">{{ t('redpacket.type', '类型') }}</label>
                <div class="mt-1 grid grid-cols-2 gap-2">
                  <button type="button" @click="createForm.redpacket_type = 'equal'"
                    :class="['flex items-center justify-center gap-2 rounded-lg border px-3 py-2.5 text-sm font-medium transition-all',
                      createForm.redpacket_type === 'equal'
                        ? 'border-red-300 bg-red-50 text-red-700 dark:border-red-700 dark:bg-red-900/30 dark:text-red-300'
                        : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300 dark:hover:bg-dark-700']">
                    🀄 {{ t('redpacket.equal', '等分红包') }}
                  </button>
                  <button type="button" @click="createForm.redpacket_type = 'random'"
                    :class="['flex items-center justify-center gap-2 rounded-lg border px-3 py-2.5 text-sm font-medium transition-all',
                      createForm.redpacket_type === 'random'
                        ? 'border-red-300 bg-red-50 text-red-700 dark:border-red-700 dark:bg-red-900/30 dark:text-red-300'
                        : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300 dark:hover:bg-dark-700']">
                    🎲 {{ t('redpacket.random', '拼手气红包') }}
                  </button>
                </div>
              </div>

              <div>
                <label class="input-label">{{ t('redpacket.memo', '附言') }}</label>
                <input v-model="createForm.memo" type="text" maxlength="100"
                  :placeholder="t('redpacket.memoPlaceholder', '恭喜发财，大吉大利！')"
                  :disabled="createLoading" class="input mt-1 w-full" />
              </div>

              <p v-if="createError" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">{{ createError }}</p>

              <div class="flex gap-2">
                <button type="submit" :disabled="createLoading" class="btn btn-primary flex-1">
                  <svg v-if="createLoading" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ t('redpacket.create', '发红包') }}
                </button>
                <button type="button" @click="showCreate = false" class="btn btn-secondary">{{ t('common.cancel', '取消') }}</button>
              </div>
            </form>

            <!-- Success Result -->
            <div v-if="createdRp" class="mt-4 rounded-xl border border-emerald-200 bg-gradient-to-r from-emerald-50 to-teal-50 p-4 dark:border-emerald-800/50 dark:from-emerald-900/20 dark:to-teal-900/20">
              <div class="flex items-start gap-3">
                <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-emerald-100 dark:bg-emerald-900/40">
                  <Icon name="checkCircle" size="sm" class="text-emerald-600 dark:text-emerald-400" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm font-medium text-emerald-700 dark:text-emerald-300">{{ t('redpacket.created', '红包已创建！') }}</p>
                  <div class="mt-2 flex items-center gap-2 rounded-lg bg-white/60 px-3 py-2 dark:bg-dark-800/60">
                    <code class="flex-1 font-mono text-lg font-bold text-emerald-800 select-all dark:text-emerald-200">{{ createdRp.code }}</code>
                    <button @click="copyCode(createdRp.code)" class="flex-shrink-0 text-emerald-500 hover:text-emerald-700 dark:hover:text-emerald-300">
                      <Icon :name="copiedCode === createdRp.code ? 'checkCircle' : 'copy'" size="sm" />
                    </button>
                  </div>
                  <p class="mt-1 text-xs text-emerald-600/70 dark:text-emerald-400/50">{{ t('redpacket.shareCode', '将此口令分享给好友即可领取') }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <!-- Claim Red Packet -->
      <transition name="fade">
        <div v-if="showClaim" class="card overflow-hidden">
          <div class="bg-gradient-to-r from-amber-50 to-yellow-50 px-5 py-4 dark:from-amber-900/20 dark:to-yellow-900/20">
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-100 dark:bg-amber-900/40">
                <Icon name="gift" size="md" class="text-amber-600 dark:text-amber-400" />
              </div>
              <h3 class="text-sm font-semibold text-amber-800 dark:text-amber-300">{{ t('redpacket.claim', '领红包') }}</h3>
            </div>
          </div>
          <div class="p-5">
            <form @submit.prevent="handleClaim" class="space-y-4">
              <div>
                <label class="input-label">{{ t('redpacket.code', '红包口令') }}</label>
                <input v-model="claimCode" type="text" required
                  :placeholder="t('redpacket.codePlaceholder', '输入红包口令')"
                  :disabled="claimLoading" class="input mt-1 w-full text-center font-mono text-lg tracking-wider" />
              </div>

              <transition name="fade">
                <div v-if="claimResult" class="rounded-xl border border-emerald-200 bg-gradient-to-br from-emerald-50 to-teal-50 p-6 text-center dark:border-emerald-800/50 dark:from-emerald-900/20 dark:to-teal-900/20">
                  <p class="text-xs text-emerald-600/70 dark:text-emerald-400/50">{{ t('redpacket.congrats', '恭喜！获得') }}</p>
                  <p class="mt-1 text-3xl font-bold text-emerald-600 dark:text-emerald-400">+${{ claimResult.amount.toFixed(2) }}</p>
                </div>
              </transition>

              <p v-if="claimError" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">{{ claimError }}</p>

              <div class="flex gap-2">
                <button type="submit" :disabled="claimLoading || !claimCode" class="btn btn-primary flex-1">
                  <svg v-if="claimLoading" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ t('redpacket.claim', '领取') }}
                </button>
                <button type="button" @click="showClaim = false; claimResult = null; claimError = ''" class="btn btn-secondary">{{ t('common.cancel', '取消') }}</button>
              </div>
            </form>
          </div>
        </div>
      </transition>

      <!-- My Red Packets -->
      <div class="card">
        <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('redpacket.myPackets', '我的红包') }}</h2>
        </div>
        <div class="p-5">
          <div v-if="loadingPackets" class="flex items-center justify-center py-12">
            <svg class="h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>

          <div v-else-if="myPackets.length > 0" class="space-y-3">
            <div v-for="rp in myPackets" :key="rp.id" class="rounded-xl border border-gray-100 bg-gradient-to-r from-white to-gray-50/50 transition-colors hover:border-gray-200 dark:border-dark-700 dark:from-dark-800 dark:to-dark-800/50 dark:hover:border-dark-600">
              <!-- Main Row -->
              <div class="flex items-center justify-between p-4">
                <div class="flex items-start gap-3">
                  <div :class="['flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl', rpStatusStyle(rp.status).bg]">
                    <span v-if="rp.status === 'active'" class="text-lg">🧧</span>
                    <Icon v-else name="gift" size="md" :class="rpStatusStyle(rp.status).icon" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <p class="text-sm font-medium text-gray-900 dark:text-white">
                        {{ rp.redpacket_type === 'equal' ? t('redpacket.equal', '等分') : t('redpacket.random', '拼手气') }}
                      </p>
                      <span :class="['rounded-full px-2 py-0.5 text-[10px] font-semibold', rpStatusStyle(rp.status).badge]">
                        {{ rpStatusLabel(rp.status) }}
                      </span>
                    </div>
                    <p class="mt-0.5 text-xs text-gray-400 dark:text-dark-500">
                      {{ rp.total_count }}{{ t('redpacket.copies', '份') }}
                      <template v-if="rp.status === 'active'">
                        · {{ t('redpacket.remaining', '剩余') }} {{ rp.remaining_count }}{{ t('redpacket.copies', '份') }}
                      </template>
                    </p>
                    <p v-if="rp.memo" class="mt-1 text-xs text-gray-500 dark:text-dark-400">「{{ rp.memo }}」</p>
                  </div>
                </div>
                <div class="flex items-center gap-3">
                  <div class="text-right">
                    <p class="text-base font-bold text-red-600 dark:text-red-400">${{ rp.total_amount.toFixed(2) }}</p>
                    <p v-if="rp.status === 'active' && rp.remaining_amount > 0" class="text-xs text-gray-400 dark:text-dark-500">
                      {{ t('redpacket.remainingAmount', '剩余') }} ${{ rp.remaining_amount.toFixed(2) }}
                    </p>
                  </div>
                  <button @click="toggleDetail(rp)"
                    :class="['flex h-8 w-8 items-center justify-center rounded-lg text-gray-400 transition-all hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dark-700 dark:hover:text-dark-300',
                      detailExpandedId === rp.id ? 'rotate-180' : '']">
                    <Icon name="chevronDown" size="sm" />
                  </button>
                </div>
              </div>

              <!-- Code Bar (active only) -->
              <div v-if="rp.status === 'active' && rp.code" class="mx-4 mb-2 flex items-center gap-2 rounded-lg bg-gray-50 px-3 py-1.5 dark:bg-dark-700/50">
                <code class="flex-1 text-xs text-gray-600 select-all dark:text-dark-300">{{ rp.code }}</code>
                <button @click="copyCode(rp.code)" class="text-gray-400 hover:text-primary-500 dark:hover:text-primary-400">
                  <Icon :name="copiedCode === rp.code ? 'checkCircle' : 'copy'" size="xs" />
                </button>
              </div>

              <!-- Detail: Claims -->
              <transition name="fade">
                <div v-if="detailExpandedId === rp.id" class="border-t border-gray-100 dark:border-dark-700">
                  <div v-if="loadingDetail" class="flex items-center justify-center py-6">
                    <svg class="h-5 w-5 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>
                  <div v-else-if="detailClaims.length > 0" class="p-4">
                    <p class="mb-3 text-xs font-medium text-gray-500 dark:text-dark-400">
                      {{ t('redpacket.claimDetail', '领取详情') }}
                      ({{ detailClaims.length }}/{{ rp.total_count }})
                    </p>
                    <div class="space-y-2">
                      <div v-for="claim in detailClaims" :key="claim.id"
                        class="flex items-center justify-between rounded-lg bg-gray-50 px-3 py-2 dark:bg-dark-700/50">
                        <div class="flex items-center gap-2">
                          <div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30">
                            <Icon name="user" size="xs" class="text-primary-600 dark:text-primary-400" />
                          </div>
                          <span class="text-sm text-gray-700 dark:text-gray-300">{{ claim.user_email || '#' + claim.user_id }}</span>
                        </div>
                        <div class="flex items-center gap-3">
                          <span class="text-sm font-medium text-emerald-600 dark:text-emerald-400">+${{ claim.amount.toFixed(2) }}</span>
                          <span class="text-xs text-gray-400 dark:text-dark-500">{{ formatDateTime(claim.created_at) }}</span>
                        </div>
                      </div>
                    </div>
                    <div v-if="rp.status === 'active' && rp.remaining_count > 0" class="mt-2 rounded-lg bg-gray-50 px-3 py-2 text-center dark:bg-dark-700/50">
                      <span class="text-xs text-gray-400 dark:text-dark-500">
                        还有 {{ rp.remaining_count }} 份等待领取
                      </span>
                    </div>
                  </div>
                  <div v-else class="py-6 text-center">
                    <p class="text-xs text-gray-400 dark:text-dark-500">{{ t('redpacket.noClaimsYet', '暂无人领取') }}</p>
                  </div>
                </div>
              </transition>
            </div>
          </div>

          <div v-else class="py-12 text-center">
            <div class="mb-3 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800">
              <span class="text-3xl opacity-40">🧧</span>
            </div>
            <p class="text-sm text-gray-400 dark:text-dark-500">{{ t('redpacket.noPackets', '暂无红包记录') }}</p>
          </div>
        </div>
      </div>
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
  showCreate.value = true
  showClaim.value = false
  createError.value = ''
  createdRp.value = null
}

function openClaim() {
  showClaim.value = true
  showCreate.value = false
  claimError.value = ''
  claimResult.value = null
  claimCode.value = ''
}

function rpStatusStyle(status: string) {
  switch (status) {
    case 'active': return {
      bg: 'bg-red-100 dark:bg-red-900/30',
      icon: 'text-red-600 dark:text-red-400',
      badge: 'bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-300',
    }
    case 'exhausted': return {
      bg: 'bg-gray-100 dark:bg-dark-800',
      icon: 'text-gray-500 dark:text-dark-400',
      badge: 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400',
    }
    case 'expired': return {
      bg: 'bg-gray-100 dark:bg-dark-800',
      icon: 'text-gray-400 dark:text-dark-500',
      badge: 'bg-orange-100 text-orange-600 dark:bg-orange-900/50 dark:text-orange-300',
    }
    default: return {
      bg: 'bg-gray-100 dark:bg-dark-800',
      icon: 'text-gray-500 dark:text-dark-400',
      badge: 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400',
    }
  }
}

function rpStatusLabel(status: string) {
  switch (status) {
    case 'active': return t('redpacket.statusActive', '进行中')
    case 'exhausted': return t('redpacket.statusExhausted', '已领完')
    case 'expired': return t('redpacket.statusExpired', '已过期')
    default: return status
  }
}

async function copyCode(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    copiedCode.value = code
    setTimeout(() => { copiedCode.value = '' }, 2000)
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
  transform: translateY(-6px);
}
</style>
