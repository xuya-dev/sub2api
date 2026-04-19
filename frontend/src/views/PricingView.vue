<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <header class="relative z-20 border-b border-gray-100 bg-white/80 backdrop-blur-xl dark:border-dark-800 dark:bg-dark-950/80">
      <nav class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-8 w-8 overflow-hidden rounded-lg">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-bold text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-2">
          <router-link to="/key-usage"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('home.keyUsage') }}
          </router-link>
          <router-link to="/monitoring"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('admin.monitoring.title') }}
          </router-link>
          <router-link to="/pricing"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-gray-900 bg-gray-100 dark:text-white dark:bg-dark-800 sm:flex">
            {{ t('pricing.title') }}
          </router-link>
          <LocaleSwitcher />
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white sm:flex">
            {{ t('home.docs') }}
          </a>
          <button @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-50 dark:text-dark-400 dark:hover:bg-dark-800">
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link v-if="isAuthenticated" :to="dashboardPath"
            class="ml-1 inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100">
            {{ t('home.dashboard') }}
            <Icon name="arrowRight" size="xs" :stroke-width="2" />
          </router-link>
          <router-link v-else to="/login"
            class="ml-1 inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100">
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="mx-auto w-full max-w-7xl flex-1 px-6 py-8">
      <div class="space-y-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.description') }}</p>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('pricing.perMillionTokens') }}</span>
            <button @click="refresh" :disabled="loading"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:opacity-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700">
              <svg class="h-4 w-4" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
              </svg>
              {{ t('common.refresh') }}
            </button>
          </div>
        </div>

        <div v-for="group in data?.groups" :key="group.id" class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900 overflow-hidden">
          <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ group.name }}</span>
            <span v-if="group.platform" class="rounded-full bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700 dark:bg-blue-950 dark:text-blue-300">{{ group.platform }}</span>
            <span v-if="group.rate_multiplier !== 1" class="rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-950 dark:text-amber-300">×{{ group.rate_multiplier.toFixed(2) }}</span>
            <span class="rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ group.models.length }} {{ group.models.length === 1 ? t('pricing.model') : t('pricing.models') }}
            </span>
          </div>
          <div v-if="group.models.length > 0" class="overflow-x-auto">
            <table class="w-full" style="table-layout: fixed;">
              <colgroup>
                <col style="width: 28%;">
                <col style="width: 14%;">
                <col style="width: 14%;">
                <col style="width: 14%;">
                <col style="width: 14%;">
                <col style="width: 16%;">
              </colgroup>
              <thead>
                <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('pricing.modelName') }}</th>
                  <th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('pricing.inputPrice') }} ($)</th>
                  <th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('pricing.outputPrice') }} ($)</th>
                  <th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-blue-600 dark:text-blue-400">{{ t('pricing.effectiveInput') }} ($)</th>
                  <th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-blue-600 dark:text-blue-400">{{ t('pricing.effectiveOutput') }} ($)</th>
                  <th class="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('pricing.requests7d') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in group.models" :key="m.model_name" class="border-b border-gray-50 last:border-b-0 dark:border-dark-800">
                  <td class="px-6 py-3 text-sm font-medium text-gray-900 dark:text-white font-mono truncate" :title="m.model_name">{{ m.model_name }}</td>
                  <td class="px-4 py-3 text-right text-sm text-gray-600 dark:text-dark-300 tabular-nums">{{ formatPrice(m.input_cost_per_million) }}</td>
                  <td class="px-4 py-3 text-right text-sm text-gray-600 dark:text-dark-300 tabular-nums">{{ formatPrice(m.output_cost_per_million) }}</td>
                  <td class="px-4 py-3 text-right text-sm font-medium tabular-nums" :class="group.rate_multiplier !== 1 ? 'text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-dark-200'">{{ formatPrice(m.effective_input) }}</td>
                  <td class="px-4 py-3 text-right text-sm font-medium tabular-nums" :class="group.rate_multiplier !== 1 ? 'text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-dark-200'">{{ formatPrice(m.effective_output) }}</td>
                  <td class="px-4 py-3 text-center text-sm tabular-nums" :class="m.request_count > 0 ? 'text-gray-700 dark:text-dark-200' : 'text-gray-400'">{{ m.request_count.toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="py-12 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.noModelsInGroup') }}</p>
          </div>
        </div>

        <div v-if="data && (!data.groups || data.groups.length === 0)" class="rounded-xl border border-gray-200 bg-white py-16 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.noData') }}</p>
        </div>
      </div>
    </main>

    <footer class="border-t border-gray-100 bg-white dark:border-dark-800 dark:bg-dark-950">
      <div class="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        <span class="text-sm text-gray-500 dark:text-dark-400">&copy; {{ currentYear }} {{ siteName }}</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { getPublicPricing, type PublicPricingResponse } from '@/api/pricing'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const currentYear = new Date().getFullYear()

const isDark = ref(document.documentElement.classList.contains('dark'))
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

const savedTheme = localStorage.getItem('theme')
if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
  isDark.value = true
  document.documentElement.classList.add('dark')
}

const data = ref<PublicPricingResponse | null>(null)
const loading = ref(false)

async function refresh() {
  loading.value = true
  try {
    data.value = await getPublicPricing()
  } catch {
    data.value = null
  } finally {
    loading.value = false
  }
}

function formatPrice(value: number): string {
  if (!value || value === 0) return '-'
  if (value >= 1) return '$' + value.toFixed(2)
  if (value >= 0.01) return '$' + value.toFixed(4)
  return '$' + value.toFixed(6)
}

onMounted(refresh)
</script>
