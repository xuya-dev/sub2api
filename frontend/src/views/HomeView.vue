<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe v-if="isHomeContentUrl" :src="homeContent.trim()" class="h-screen w-full border-0" allowfullscreen></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page - Professional Business Style -->
  <div v-else class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <PublicPageHeader />

    <!-- Hero Section -->
    <section class="relative overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-b from-primary-50/40 to-transparent dark:from-primary-950/20 dark:to-transparent"></div>
      <div class="relative mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:py-40">
        <div class="mx-auto max-w-3xl text-center">
          <div class="mb-6 inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50 px-4 py-1.5 dark:border-primary-800 dark:bg-primary-950/40">
            <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
            <span class="text-sm font-medium text-primary-700 dark:text-primary-300">{{ t('home.tags.subscriptionToApi') }}</span>
          </div>
          <h1 class="text-4xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-5xl lg:text-6xl">
            {{ t('home.heroSubtitle') }}
          </h1>
          <p class="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.heroDescription') }}
          </p>
          <div class="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
            <router-link :to="isAuthenticated ? dashboardPath : '/login'"
              class="inline-flex items-center gap-2 rounded-lg bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-500/25 transition-all hover:bg-primary-700 hover:shadow-xl hover:shadow-primary-500/30">
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
              <Icon name="arrowRight" size="md" :stroke-width="2" />
            </router-link>
            <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-8 py-3.5 text-base font-semibold text-gray-700 transition-all hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700">
              {{ t('home.viewDocs') }}
            </a>
          </div>
        </div>

        <!-- Terminal Preview -->
        <div class="mx-auto mt-16 max-w-2xl">
          <div class="terminal-window">
            <div class="terminal-header">
              <div class="terminal-buttons">
                <span class="btn-close"></span>
                <span class="btn-minimize"></span>
                <span class="btn-maximize"></span>
              </div>
              <span class="terminal-title">terminal</span>
            </div>
            <div class="terminal-body">
              <div class="code-line line-1">
                <span class="code-prompt">$</span>
                <span class="code-cmd">curl</span>
                <span class="code-flag">-X POST</span>
                <span class="code-url">/v1/messages</span>
              </div>
              <div class="code-line line-2">
                <span class="code-comment"># Routing to upstream...</span>
              </div>
              <div class="code-line line-3">
                <span class="code-success">200 OK</span>
                <span class="code-response">{ "content": "Hello!" }</span>
              </div>
              <div class="code-line line-4">
                <span class="code-prompt">$</span>
                <span class="cursor"></span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Supported Providers -->
    <section class="border-y border-gray-100 bg-gray-50/50 py-12 dark:border-dark-800 dark:bg-dark-900/50">
      <div class="mx-auto max-w-7xl px-6">
        <p class="mb-8 text-center text-sm font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">
          {{ t('home.providers.title') }}
        </p>
        <div class="flex flex-wrap items-center justify-center gap-10">
          <div v-for="p in providers" :key="p.name" class="flex flex-col items-center gap-2">
            <div class="flex h-10 w-10 items-center justify-center">
              <ModelIcon :model="p.model" size="32px" />
            </div>
            <span class="text-sm font-semibold text-gray-700 dark:text-dark-200">{{ p.name }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Pain Points -->
    <section class="py-20 sm:py-24">
      <div class="mx-auto max-w-7xl px-6">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">{{ t('home.painPoints.title') }}</h2>
        </div>
        <div class="mx-auto mt-14 grid max-w-4xl gap-6 sm:grid-cols-2">
          <div v-for="(item, key) in painPointItems" :key="key"
            class="flex gap-4 rounded-xl border border-gray-100 bg-white p-6 shadow-sm transition-shadow hover:shadow-md dark:border-dark-800 dark:bg-dark-900">
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-red-50 dark:bg-red-950/30">
              <svg class="h-5 w-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900 dark:text-white">{{ item.title }}</h3>
              <p class="mt-1 text-sm leading-relaxed text-gray-500 dark:text-dark-400">{{ item.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features / Solutions -->
    <section class="border-y border-gray-100 bg-gray-50/50 py-20 sm:py-24 dark:border-dark-800 dark:bg-dark-900/50">
      <div class="mx-auto max-w-7xl px-6">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">{{ t('home.solutions.title') }}</h2>
          <p class="mt-4 text-lg text-gray-600 dark:text-dark-300">{{ t('home.solutions.subtitle') }}</p>
        </div>
        <div class="mx-auto mt-14 grid max-w-5xl gap-8 sm:grid-cols-3">
          <div v-for="(feature, idx) in featureItems" :key="idx"
            class="group rounded-2xl border border-gray-100 bg-white p-8 shadow-sm transition-all hover:border-primary-200 hover:shadow-lg dark:border-dark-800 dark:bg-dark-900 dark:hover:border-primary-800">
            <div :class="feature.bgColor" class="mb-5 flex h-12 w-12 items-center justify-center rounded-xl transition-transform group-hover:scale-110">
              <Icon :name="feature.icon" size="lg" class="text-white" />
            </div>
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ feature.title }}</h3>
            <p class="mt-2 text-sm leading-relaxed text-gray-500 dark:text-dark-400">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Comparison Table -->
    <section class="py-20 sm:py-24">
      <div class="mx-auto max-w-7xl px-6">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">{{ t('home.comparison.title') }}</h2>
        </div>
        <div class="mx-auto mt-14 max-w-3xl overflow-hidden rounded-2xl border border-gray-200 dark:border-dark-700">
          <table class="w-full">
            <thead>
              <tr class="bg-gray-50 dark:bg-dark-800">
                <th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-white">{{ t('home.comparison.headers.feature') }}</th>
                <th class="px-6 py-4 text-center text-sm font-semibold text-gray-500 dark:text-dark-400">{{ t('home.comparison.headers.official') }}</th>
                <th class="bg-primary-600 px-6 py-4 text-center text-sm font-semibold text-white">{{ t('home.comparison.headers.us') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, key) in comparisonItems" :key="key"
                class="border-t border-gray-100 dark:border-dark-800">
                <td class="px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">{{ item.feature }}</td>
                <td class="px-6 py-4 text-center text-sm text-gray-500 dark:text-dark-400">{{ item.official }}</td>
                <td class="bg-primary-50/50 px-6 py-4 text-center text-sm font-medium text-primary-700 dark:bg-primary-950/20 dark:text-primary-300">{{ item.us }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="border-t border-gray-100 bg-gray-50/50 py-20 dark:border-dark-800 dark:bg-dark-900/50">
      <div class="mx-auto max-w-7xl px-6">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">{{ t('home.cta.title') }}</h2>
          <p class="mt-4 text-lg text-gray-600 dark:text-dark-300">{{ t('home.cta.description') }}</p>
          <div class="mt-8">
            <router-link :to="isAuthenticated ? dashboardPath : '/login'"
              class="inline-flex items-center gap-2 rounded-lg bg-primary-600 px-10 py-4 text-lg font-semibold text-white shadow-lg shadow-primary-500/25 transition-all hover:bg-primary-700 hover:shadow-xl">
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.cta.button') }}
              <Icon name="arrowRight" size="md" :stroke-width="2" />
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <PublicPageFooter />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import PublicPageHeader from '@/components/common/PublicPageHeader.vue'
import PublicPageFooter from '@/components/common/PublicPageFooter.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

const providers = [
  { name: 'Claude', model: 'claude-sonnet-4-5' },
  { name: 'OpenAI', model: 'gpt-4o' },
  { name: 'Gemini', model: 'gemini-2.5-pro' },
  { name: 'DeepSeek', model: 'deepseek-chat' },
  { name: 'GLM', model: 'glm-4' },
  { name: 'Qwen', model: 'qwen-max' },
  { name: 'Kimi', model: 'kimi' },
  { name: 'MiniMax', model: 'minimax' },
  { name: 'Yi', model: 'yi-lightning' },
  { name: 'Mistral', model: 'mistral-large' },
  { name: 'Llama', model: 'llama-4' },
]

const featureItems = computed(() => [
  {
    title: t('home.features.unifiedGateway'),
    desc: t('home.features.unifiedGatewayDesc'),
    icon: 'server' as const,
    bgColor: 'bg-gradient-to-br from-blue-500 to-blue-600',
  },
  {
    title: t('home.features.multiAccount'),
    desc: t('home.features.multiAccountDesc'),
    icon: 'shield' as const,
    bgColor: 'bg-gradient-to-br from-primary-500 to-primary-600',
  },
  {
    title: t('home.features.balanceQuota'),
    desc: t('home.features.balanceQuotaDesc'),
    icon: 'chart' as const,
    bgColor: 'bg-gradient-to-br from-purple-500 to-purple-600',
  },
])

const painPointItems = computed(() => ({
  expensive: { title: t('home.painPoints.items.expensive.title'), desc: t('home.painPoints.items.expensive.desc') },
  complex: { title: t('home.painPoints.items.complex.title'), desc: t('home.painPoints.items.complex.desc') },
  unstable: { title: t('home.painPoints.items.unstable.title'), desc: t('home.painPoints.items.unstable.desc') },
  noControl: { title: t('home.painPoints.items.noControl.title'), desc: t('home.painPoints.items.noControl.desc') },
}))

const comparisonItems = computed(() => ({
  pricing: { feature: t('home.comparison.items.pricing.feature'), official: t('home.comparison.items.pricing.official'), us: t('home.comparison.items.pricing.us') },
  models: { feature: t('home.comparison.items.models.feature'), official: t('home.comparison.items.models.official'), us: t('home.comparison.items.models.us') },
  management: { feature: t('home.comparison.items.management.feature'), official: t('home.comparison.items.management.official'), us: t('home.comparison.items.management.us') },
  stability: { feature: t('home.comparison.items.stability.feature'), official: t('home.comparison.items.stability.official'), us: t('home.comparison.items.stability.us') },
  control: { feature: t('home.comparison.items.control.feature'), official: t('home.comparison.items.control.official'), us: t('home.comparison.items.control.us') },
}))
</script>

<style scoped>
.terminal-window {
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

:deep(.dark) .terminal-window {
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.06);
}

.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close { background: #ef4444; }
.btn-minimize { background: #eab308; }
.btn-maximize { background: #22c55e; }

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 { animation-delay: 0.3s; }
.line-2 { animation-delay: 1s; }
.line-3 { animation-delay: 1.8s; }
.line-4 { animation-delay: 2.5s; }

@keyframes line-appear {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}

.code-prompt { color: #22c55e; font-weight: bold; }
.code-cmd { color: #38bdf8; }
.code-flag { color: #a78bfa; }
.code-url { color: #14b8a6; }
.code-comment { color: #64748b; font-style: italic; }
.code-success { color: #22c55e; background: rgba(34, 197, 94, 0.15); padding: 2px 8px; border-radius: 4px; font-weight: 600; }
.code-response { color: #fbbf24; }

.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
