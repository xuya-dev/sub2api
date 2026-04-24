<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="flex-1 sm:max-w-64">
            <input v-model="filterQuery" type="text" :placeholder="t('admin.transfer.searchPlaceholder', '搜索用户 ID')" class="input" @input="handleSearch" />
          </div>
          <Select v-model="filters.status" :options="statusOptions" class="w-36" @change="loadTransfers" />
          <Select v-model="filters.transfer_type" :options="typeOptions" class="w-36" @change="loadTransfers" />
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <button @click="loadTransfers" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button @click="showBatchDialog = true" class="btn btn-primary">
              {{ t('admin.transfer.batchDistribute', '批量发放') }}
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="transfers"
          :loading="loading"
          :server-side-sort="true"
          default-sort-key="created_at"
          default-sort-order="desc"
          @sort="handleSort"
        >
          <template #cell-sender_id="{ value, row }">
            <div class="flex flex-col">
              <span class="text-sm text-gray-900 dark:text-white">{{ row.sender_email || '-' }}</span>
              <span class="text-xs text-gray-400 dark:text-dark-500">#{{ value }}</span>
            </div>
          </template>

          <template #cell-receiver_id="{ value, row }">
            <div class="flex flex-col">
              <span class="text-sm text-gray-900 dark:text-white">{{ row.receiver_email || '-' }}</span>
              <span class="text-xs text-gray-400 dark:text-dark-500">#{{ value }}</span>
            </div>
          </template>

          <template #cell-amount="{ value }">
            <span class="text-sm font-medium text-gray-900 dark:text-white">${{ value.toFixed(2) }}</span>
          </template>

          <template #cell-fee="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">${{ value.toFixed(2) }}</span>
          </template>

          <template #cell-transfer_type="{ value }">
            <span :class="['badge', typeBadgeClass(value)]">{{ typeLabel(value) }}</span>
          </template>

          <template #cell-status="{ value }">
            <span :class="['badge', statusBadgeClass(value)]">{{ statusLabel(value) }}</span>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center space-x-2">
              <template v-if="row.status === 'completed'">
                <button @click="confirmFreeze(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-yellow-50 hover:text-yellow-600 dark:hover:bg-yellow-900/20 dark:hover:text-yellow-400">
                  <Icon name="ban" size="sm" />
                  <span class="text-xs">{{ t('admin.transfer.freeze', '冻结') }}</span>
                </button>
                <button @click="confirmRevoke(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400">
                  <Icon name="trash" size="sm" />
                  <span class="text-xs">{{ t('admin.transfer.revoke', '撤回') }}</span>
                </button>
              </template>
              <span v-else class="text-gray-400 dark:text-dark-500">-</span>
            </div>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <ConfirmDialog
      :show="showFreezeDialog"
      :title="t('admin.transfer.freezeTitle', '冻结转账')"
      :message="t('admin.transfer.freezeConfirm', '确认冻结此笔转账？冻结后接收方余额将被扣除并退回发送方。')"
      :confirm-text="t('admin.transfer.freeze', '冻结')"
      :cancel-text="t('common.cancel')"
      danger
      @confirm="handleFreeze"
      @cancel="showFreezeDialog = false"
    />

    <ConfirmDialog
      :show="showRevokeDialog"
      :title="t('admin.transfer.revokeTitle', '撤回转账')"
      :message="t('admin.transfer.revokeConfirm', '确认撤回此笔转账？')"
      :confirm-text="t('admin.transfer.revoke', '撤回')"
      :cancel-text="t('common.cancel')"
      danger
      @confirm="handleRevoke"
      @cancel="showRevokeDialog = false"
    />

    <Teleport to="body">
      <div v-if="showBatchDialog" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/50" @click="showBatchDialog = false"></div>
        <div class="relative z-10 w-full max-w-xl rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.transfer.batchDistribute', '批量发放') }}</h2>
          <div class="max-h-96 space-y-3 overflow-y-auto">
            <div v-for="(target, i) in batchTargets" :key="i" class="rounded-lg border border-gray-200 p-3 dark:border-dark-700">
              <div class="flex items-center gap-2">
                <div class="relative flex-1">
                  <input v-model="target.query" type="text"
                    :placeholder="t('admin.transfer.searchUser', '搜索邮箱/用户名')"
                    class="input w-full"
                    @input="onBatchSearch(i)"
                    @focus="onBatchSearch(i)" />
                  <div v-if="target.results && target.results.length > 0 && !target.selected" class="absolute left-0 right-0 top-full z-10 mt-1 max-h-36 overflow-y-auto rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dark-600 dark:bg-dark-800">
                    <button v-for="u in target.results" :key="u.id" type="button" @click="selectBatchUser(i, u)"
                      class="flex w-full items-center gap-2 px-3 py-2 text-left text-sm hover:bg-gray-50 dark:hover:bg-dark-700">
                      <div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30">
                        <Icon name="user" size="xs" class="text-primary-600 dark:text-primary-400" />
                      </div>
                      <span class="flex-1 truncate text-gray-900 dark:text-white">{{ u.email }}</span>
                      <span class="text-xs text-gray-400">#{{ u.id }}</span>
                    </button>
                  </div>
                </div>
                <div v-if="target.selected" class="flex items-center gap-1">
                  <span class="rounded-md bg-primary-50 px-2 py-1 text-xs font-medium text-primary-700 dark:bg-primary-900/20 dark:text-primary-300">{{ target.selected.email }}</span>
                  <button @click="clearBatchUser(i)" class="text-gray-400 hover:text-red-500">
                    <Icon name="x" size="xs" />
                  </button>
                </div>
                <input v-model.number="target.amount" type="number" step="0.01" min="0.01"
                  :placeholder="t('admin.transfer.amount', '金额')" class="input w-28" />
                <button @click="removeBatchTarget(i)" class="text-gray-400 hover:text-red-500">
                  <Icon name="trash" size="sm" />
                </button>
              </div>
            </div>
          </div>
          <button @click="addBatchTarget" class="mt-3 text-sm font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
            + {{ t('admin.transfer.addTarget', '添加目标') }}
          </button>
          <input v-model="batchMemo" type="text" :placeholder="t('admin.transfer.memoPlaceholder', '备注（可选）')" class="input mt-3 w-full" />
          <div class="mt-4 flex gap-3">
            <button @click="handleBatch" :disabled="batchLoading" class="btn btn-primary flex-1">
              <svg v-if="batchLoading" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ t('admin.transfer.confirm', '确认发放') }}
            </button>
            <button @click="showBatchDialog = false" class="btn btn-secondary flex-1">{{ t('common.cancel') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { adminAPI } from '@/api/admin'
import { searchUsers as adminSearchUsers } from '@/api/admin/usage'
import type { TransferRecord } from '@/api/admin/transfer'
import type { Column } from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const appStore = useAppStore()

const transfers = ref<TransferRecord[]>([])
const loading = ref(false)
const filterQuery = ref('')
const filters = reactive({ status: '', transfer_type: '' })
const pagination = reactive({ total: 0, page: 1, page_size: 20 })
const sortKey = ref('created_at')
const sortOrder = ref<'asc' | 'desc'>('desc')

const showFreezeDialog = ref(false)
const showRevokeDialog = ref(false)
const freezeTarget = ref<TransferRecord | null>(null)
const revokeTarget = ref<TransferRecord | null>(null)

const showBatchDialog = ref(false)
const batchTargets = reactive<BatchTarget[]>([])
const batchMemo = ref('')
const batchLoading = ref(false)

interface BatchTarget {
  query: string
  results: { id: number; email: string }[]
  selected: { id: number; email: string } | null
  amount: number
}

const columns = computed<Column[]>(() => [
  { key: 'id', label: 'ID', sortable: true },
  { key: 'sender_id', label: t('admin.transfer.sender', '发送方') },
  { key: 'receiver_id', label: t('admin.transfer.receiver', '接收方') },
  { key: 'amount', label: t('admin.transfer.amount', '金额'), sortable: true },
  { key: 'fee', label: t('admin.transfer.fee', '手续费') },
  { key: 'transfer_type', label: t('admin.transfer.type', '类型') },
  { key: 'status', label: t('admin.transfer.status', '状态') },
  { key: 'created_at', label: t('admin.transfer.time', '时间'), sortable: true },
  { key: 'actions', label: t('admin.transfer.actions', '操作') },
])

const statusOptions = computed(() => [
  { value: '', label: t('admin.transfer.allStatus', '全部状态') },
  { value: 'completed', label: t('admin.transfer.completed', '已完成') },
  { value: 'frozen', label: t('admin.transfer.frozen', '已冻结') },
  { value: 'revoked', label: t('admin.transfer.revoked', '已撤回') },
])

const typeOptions = computed(() => [
  { value: '', label: t('admin.transfer.allType', '全部类型') },
  { value: 'direct', label: t('admin.transfer.direct', '直接转账') },
  { value: 'redpacket', label: t('admin.transfer.redpacket', '红包') },
  { value: 'batch', label: t('admin.transfer.batch', '批量发放') },
])

let searchTimer: ReturnType<typeof setTimeout> | null = null
function handleSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadTransfers()
  }, 300)
}

function handleSort(key: string, order: 'asc' | 'desc') {
  sortKey.value = key
  sortOrder.value = order
  loadTransfers()
}

function handlePageChange(page: number) {
  pagination.page = page
  loadTransfers()
}

function handlePageSizeChange(pageSize: number) {
  pagination.page_size = pageSize
  pagination.page = 1
  loadTransfers()
}

async function loadTransfers() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.page, page_size: pagination.page_size, status: filters.status, transfer_type: filters.transfer_type }
    if (filterQuery.value) {
      const uid = parseInt(filterQuery.value)
      if (uid > 0) params.user_id = uid
    }
    const res = await adminAPI.transfer.listTransfers(params)
    transfers.value = res.items || []
    pagination.total = res.total
  } catch {
    appStore.showError(t('admin.transfer.loadFailed', '加载失败'))
  } finally {
    loading.value = false
  }
}

function typeBadgeClass(type: string) {
  switch (type) {
    case 'direct': return 'badge-primary'
    case 'redpacket': return 'badge-danger'
    case 'batch': return 'badge-warning'
    default: return 'badge-gray'
  }
}

function typeLabel(type: string) {
  switch (type) {
    case 'direct': return t('admin.transfer.direct', '直接转账')
    case 'redpacket': return t('admin.transfer.redpacket', '红包')
    case 'batch': return t('admin.transfer.batch', '批量发放')
    default: return type
  }
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'completed': return 'badge-success'
    case 'frozen': return 'badge-warning'
    case 'revoked': return 'badge-danger'
    default: return 'badge-gray'
  }
}

function statusLabel(status: string) {
  switch (status) {
    case 'completed': return t('admin.transfer.completed', '已完成')
    case 'frozen': return t('admin.transfer.frozen', '已冻结')
    case 'revoked': return t('admin.transfer.revoked', '已撤回')
    default: return status
  }
}

function confirmFreeze(row: TransferRecord) {
  freezeTarget.value = row
  showFreezeDialog.value = true
}

function confirmRevoke(row: TransferRecord) {
  revokeTarget.value = row
  showRevokeDialog.value = true
}

async function handleFreeze() {
  if (!freezeTarget.value) return
  try {
    await adminAPI.transfer.freezeTransfer(freezeTarget.value.id)
    appStore.showSuccess(t('admin.transfer.freezeSuccess', '冻结成功'))
    loadTransfers()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || t('admin.transfer.freezeFailed', '冻结失败'))
  } finally {
    showFreezeDialog.value = false
    freezeTarget.value = null
  }
}

async function handleRevoke() {
  if (!revokeTarget.value) return
  try {
    await adminAPI.transfer.revokeTransfer(revokeTarget.value.id, t('admin.transfer.adminRevoke', '管理员撤回'))
    appStore.showSuccess(t('admin.transfer.revokeSuccess', '撤回成功'))
    loadTransfers()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || t('admin.transfer.revokeFailed', '撤回失败'))
  } finally {
    showRevokeDialog.value = false
    revokeTarget.value = null
  }
}

function addBatchTarget() {
  batchTargets.push({ query: '', results: [], selected: null, amount: 0 })
}

function removeBatchTarget(i: number) {
  batchTargets.splice(i, 1)
}

function selectBatchUser(i: number, u: { id: number; email: string }) {
  const target = batchTargets[i]
  target.selected = u
  target.query = u.email
  target.results = []
}

function clearBatchUser(i: number) {
  const target = batchTargets[i]
  target.selected = null
  target.query = ''
  target.results = []
}

let batchSearchTimers: Map<number, ReturnType<typeof setTimeout>> = new Map()
async function onBatchSearch(i: number) {
  const target = batchTargets[i]
  target.selected = null
  const old = batchSearchTimers.get(i)
  if (old) clearTimeout(old)
  if (!target.query || target.query.length < 1) {
    target.results = []
    return
  }
  const timer = setTimeout(async () => {
    try {
      target.results = await adminSearchUsers(target.query)
    } catch {
      target.results = []
    }
  }, 300)
  batchSearchTimers.set(i, timer)
}

async function handleBatch() {
  const valid = batchTargets.filter(t => t.selected && t.amount > 0).map(t => ({ user_id: t.selected!.id, amount: t.amount }))
  if (valid.length === 0) {
    appStore.showError(t('admin.transfer.noValidTargets', '请添加有效的发放目标'))
    return
  }
  batchLoading.value = true
  try {
    await adminAPI.transfer.batchDistribute(valid, batchMemo.value || undefined)
    appStore.showSuccess(t('admin.transfer.batchSuccess', '批量发放成功'))
    showBatchDialog.value = false
    batchTargets.splice(0, batchTargets.length)
    batchMemo.value = ''
    loadTransfers()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || t('admin.transfer.batchFailed', '批量发放失败'))
  } finally {
    batchLoading.value = false
  }
}

onMounted(loadTransfers)
</script>
