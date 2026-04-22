<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('admin.settings.checkin.blindboxTitle') }}
          </h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t('admin.settings.checkin.blindboxDescription') }}
          </p>
        </div>
        <Toggle :modelValue="enabled" @update:modelValue="$emit('update:enabled', $event)" />
      </div>
    </div>
    <template v-if="enabled">
      <div class="space-y-4 p-6">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('admin.settings.checkin.blindboxTriggerType') }}
            </label>
            <select :value="triggerType" @change="$emit('update:triggerType', ($event.target as HTMLSelectElement).value)" class="input">
              <option value="streak">{{ t('admin.settings.checkin.blindboxTriggerStreak') }}</option>
              <option value="total">{{ t('admin.settings.checkin.blindboxTriggerTotal') }}</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('admin.settings.checkin.blindboxInterval') }}
            </label>
            <input :value="interval" @input="$emit('update:interval', parseInt(($event.target as HTMLInputElement).value) || 1)" type="number" min="1" class="input" placeholder="7" />
            <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.checkin.blindboxIntervalHint') }}
            </p>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium opacity-0">&nbsp;</label>
            <button type="button" @click="showCreate = true" class="btn btn-primary">
              <Icon name="plus" size="md" class="mr-1" />
              {{ t('admin.blindbox.createItem') }}
            </button>
          </div>
        </div>

        <!-- Stats -->
        <div v-if="stats" class="grid grid-cols-3 gap-3">
          <div class="rounded-lg bg-gray-50 p-3 text-center dark:bg-dark-800">
            <p class="text-lg font-bold text-gray-900 dark:text-white">{{ stats.total_items }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.totalItems') }}</p>
          </div>
          <div class="rounded-lg bg-gray-50 p-3 text-center dark:bg-dark-800">
            <p class="text-lg font-bold text-green-600 dark:text-green-400">{{ stats.enabled_items }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.enabledItems') }}</p>
          </div>
          <div class="rounded-lg bg-gray-50 p-3 text-center dark:bg-dark-800">
            <p class="text-lg font-bold text-purple-600 dark:text-purple-400">{{ stats.total_draws }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.totalDraws') }}</p>
          </div>
        </div>

        <!-- Prize Items Table -->
        <div class="overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
          <div v-if="loading" class="flex items-center justify-center py-8">
            <div class="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
          </div>
          <div v-else-if="items.length === 0" class="py-8 text-center text-sm text-gray-400 dark:text-gray-500">
            {{ t('admin.blindbox.empty') }}
          </div>
          <table v-else class="w-full min-w-[640px] text-sm">
            <thead class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
              <tr>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colName') }}</th>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colRarity') }}</th>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colRewardType') }}</th>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colReward') }}</th>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colWeight') }}</th>
                <th class="px-3 py-2 text-left font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colStatus') }}</th>
                <th class="px-3 py-2 text-right font-medium text-gray-500 dark:text-gray-400">{{ t('admin.blindbox.colActions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50 dark:divide-dark-800/50">
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
                <td class="px-3 py-2 font-medium text-gray-900 dark:text-white">{{ item.name }}</td>
                <td class="px-3 py-2">
                  <span :class="rarityClass(item.rarity)" class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium">
                    {{ rarityLabel(item.rarity) }}
                  </span>
                </td>
                <td class="px-3 py-2 text-gray-600 dark:text-gray-300">{{ rewardTypeLabel(item.reward_type) }}</td>
                <td class="px-3 py-2 text-gray-600 dark:text-gray-300">{{ formatReward(item) }}</td>
                <td class="px-3 py-2 text-gray-600 dark:text-gray-300">{{ item.weight }}</td>
                <td class="px-3 py-2">
                  <span :class="item.is_enabled ? 'text-green-600 dark:text-green-400' : 'text-gray-400 dark:text-gray-500'" class="text-xs font-medium">
                    {{ item.is_enabled ? t('common.enabled') : t('common.disabled') }}
                  </span>
                </td>
                <td class="px-3 py-2 text-right">
                  <button type="button" @click="editItem(item)" class="mr-2 text-primary-600 hover:text-primary-700 dark:text-primary-400">{{ t('common.edit') }}</button>
                  <button type="button" @click="deleteItem(item)" class="text-red-500 hover:text-red-600">{{ t('common.delete') }}</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Create/Edit Modal -->
    <div v-if="showCreate || editingItem" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="closeModal">
      <div class="card mx-4 w-full max-w-lg space-y-4 p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ editingItem ? t('admin.blindbox.editItem') : t('admin.blindbox.createItem') }}
        </h3>
        <div class="space-y-3">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.colName') }}</label>
            <input v-model="form.name" class="input" :placeholder="t('admin.blindbox.namePlaceholder')" />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.colRarity') }}</label>
              <select v-model="form.rarity" class="input">
                <option value="common">{{ t('checkin.blindboxCommon') }}</option>
                <option value="rare">{{ t('checkin.blindboxRare') }}</option>
                <option value="epic">{{ t('checkin.blindboxEpic') }}</option>
                <option value="legendary">{{ t('checkin.blindboxLegendary') }}</option>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.colRewardType') }}</label>
              <select v-model="form.reward_type" class="input">
                <option value="balance">{{ t('admin.blindbox.rewardBalance') }}</option>
                <option value="concurrency">{{ t('admin.blindbox.rewardConcurrency') }}</option>
                <option value="subscription">{{ t('admin.blindbox.rewardSubscription') }}</option>
                <option value="invitation_code">{{ t('admin.blindbox.rewardInvitation') }}</option>
              </select>
            </div>
          </div>
          <div v-if="form.reward_type === 'balance'" class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.minValue') }}</label>
              <input v-model.number="form.reward_value" type="number" step="0.01" min="0" class="input" />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.maxValue') }}</label>
              <input v-model.number="form.reward_value_max" type="number" step="0.01" min="0" class="input" />
            </div>
          </div>
          <div v-if="form.reward_type === 'concurrency'">
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.concurrencyValue') }}</label>
            <input v-model.number="form.reward_value" type="number" min="1" class="input" />
          </div>
          <div v-if="form.reward_type === 'subscription'">
            <div class="mb-3">
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.subscriptionGroup') }}</label>
              <select v-model.number="form.subscription_id" class="input">
                <option :value="0" disabled>{{ t('admin.blindbox.selectGroup') }}</option>
                <option v-for="g in subscriptionGroups" :key="g.id" :value="g.id">{{ g.name }}</option>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.subscriptionDays') }}</label>
              <input v-model.number="form.subscription_days" type="number" min="1" class="input" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.blindbox.colWeight') }}</label>
              <input v-model.number="form.weight" type="number" min="1" class="input" />
              <p class="mt-1 text-xs text-gray-400">{{ t('admin.blindbox.weightHint') }}</p>
            </div>
            <div class="flex items-end">
              <label class="flex items-center gap-2 text-sm">
                <input v-model="form.is_enabled" type="checkbox" class="rounded border-gray-300 dark:border-dark-600" />
                {{ t('common.enabled') }}
              </label>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-3 pt-2">
          <button type="button" @click="closeModal" class="btn btn-secondary">{{ t('common.cancel') }}</button>
          <button type="button" @click="saveItem" class="btn btn-primary" :disabled="saving">
            {{ saving ? t('common.saving') : t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import Toggle from '@/components/common/Toggle.vue'
import { blindboxAPI, type PrizeItem } from '@/api/admin/blindbox'
import * as groupAPI from '@/api/admin/groups'
import { useAppStore } from '@/stores/app'

const props = defineProps<{
  enabled: boolean
  triggerType: string
  interval: number
}>()

defineEmits<{
  'update:enabled': [value: boolean]
  'update:triggerType': [value: string]
  'update:interval': [value: number]
}>()

const { t } = useI18n()
const appStore = useAppStore()
const loading = ref(false)
const saving = ref(false)
const items = ref<PrizeItem[]>([])
const stats = ref<{ total_items: number; enabled_items: number; total_draws: number } | null>(null)
const showCreate = ref(false)
const editingItem = ref<PrizeItem | null>(null)
const subscriptionGroups = ref<{ id: number; name: string }[]>([])

const defaultForm = { name: '', rarity: 'common', reward_type: 'balance', reward_value: 0, reward_value_max: 0, subscription_id: 0, subscription_days: 0, weight: 100, is_enabled: true }
const form = ref({ ...defaultForm })

function rarityClass(rarity: string): string {
  switch (rarity) {
    case 'legendary': return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
    case 'epic': return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'
    case 'rare': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
    default: return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
  }
}

function rarityLabel(rarity: string): string {
  const map: Record<string, string> = { common: t('checkin.blindboxCommon'), rare: t('checkin.blindboxRare'), epic: t('checkin.blindboxEpic'), legendary: t('checkin.blindboxLegendary') }
  return map[rarity] || rarity
}

function rewardTypeLabel(type: string): string {
  const map: Record<string, string> = { balance: t('admin.blindbox.rewardBalance'), concurrency: t('admin.blindbox.rewardConcurrency'), subscription: t('admin.blindbox.rewardSubscription'), invitation_code: t('admin.blindbox.rewardInvitation') }
  return map[type] || type
}

function formatReward(item: PrizeItem): string {
  switch (item.reward_type) {
    case 'balance': return item.reward_value_max > item.reward_value ? `$${item.reward_value}~$${item.reward_value_max}` : `$${item.reward_value}`
    case 'concurrency': return `+${item.reward_value}`
    case 'subscription': {
      const days = item.subscription_days
      const group = subscriptionGroups.value.find(g => g.id === item.subscription_id)
      return group ? `${group.name} ${days}${t('admin.blindbox.days')}` : t('checkin.blindboxSubscriptionReward', { days })
    }
    case 'invitation_code': return '×1'
    default: return `${item.reward_value}`
  }
}

function editItem(item: PrizeItem) {
  editingItem.value = item
  form.value = { name: item.name, rarity: item.rarity, reward_type: item.reward_type, reward_value: item.reward_value, reward_value_max: item.reward_value_max, subscription_id: item.subscription_id || 0, subscription_days: item.subscription_days, weight: item.weight, is_enabled: item.is_enabled }
}

function closeModal() {
  showCreate.value = false
  editingItem.value = null
  form.value = { ...defaultForm }
}

async function saveItem() {
  saving.value = true
  try {
    const f = form.value
    const payload = {
      name: f.name,
      rarity: f.rarity,
      reward_type: f.reward_type,
      reward_value: f.reward_value,
      reward_value_max: f.reward_value_max,
      subscription_id: f.reward_type === 'subscription' ? (f.subscription_id || null) : null,
      subscription_days: f.reward_type === 'subscription' ? f.subscription_days : 0,
      weight: f.weight,
      is_enabled: f.is_enabled,
    }
    if (editingItem.value) {
      await blindboxAPI.updatePrizeItem(editingItem.value.id, payload)
    } else {
      await blindboxAPI.createPrizeItem(payload)
    }
    closeModal()
    await loadData()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.detail || t('common.error'))
  } finally {
    saving.value = false
  }
}

async function deleteItem(item: PrizeItem) {
  if (!confirm(t('admin.blindbox.confirmDelete'))) return
  try {
    await blindboxAPI.deletePrizeItem(item.id)
    await loadData()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.detail || t('common.error'))
  }
}

async function loadData() {
  loading.value = true
  try {
    const [prizeItems, prizeStats] = await Promise.all([blindboxAPI.listPrizeItems(), blindboxAPI.getBlindboxStats()])
    items.value = prizeItems
    stats.value = prizeStats
  } catch {
    items.value = []
  } finally {
    loading.value = false
  }
}

async function loadSubscriptionGroups() {
  try {
    const groups = await groupAPI.getAll()
    subscriptionGroups.value = groups
      .filter(g => g.subscription_type === 'subscription' && g.status === 'active')
      .map(g => ({ id: g.id, name: g.name }))
  } catch {
    subscriptionGroups.value = []
  }
}

watch(() => props.enabled, (val) => {
  if (val && items.value.length === 0) loadData()
})

watch([showCreate, editingItem, () => form.value.reward_type], ([isCreating, isEditing, rewardType]) => {
  if ((isCreating || isEditing) && rewardType === 'subscription' && subscriptionGroups.value.length === 0) {
    loadSubscriptionGroups()
  }
})

onMounted(async () => {
  await loadSubscriptionGroups()
  if (props.enabled) loadData()
})
</script>
