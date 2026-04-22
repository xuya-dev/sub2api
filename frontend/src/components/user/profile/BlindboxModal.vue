<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show && result" class="modal-overlay blindbox-overlay" style="z-index: 60" @click.self="handleClose">
        <div class="blindbox-container" @click.stop>
          <div class="blindbox-glow" :class="rarityGlowClass"></div>

          <div class="blindbox-card">
            <div class="blindbox-card-inner">
              <div class="blindbox-sparkles" v-if="result.rarity === 'legendary'">
                <span v-for="i in 8" :key="i" class="sparkle" :style="sparkleStyle(i)"></span>
              </div>

              <div class="blindbox-icon-wrapper" :class="rarityBgClass">
                <svg class="blindbox-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
                </svg>
              </div>

              <h3 class="blindbox-title">{{ t('checkin.blindboxTitle') }}</h3>

              <div class="blindbox-prize-name">{{ result.prize_name }}</div>

              <div class="blindbox-rarity-badge" :class="rarityBadgeClass">
                {{ rarityLabel }}
              </div>

              <div class="blindbox-reward" :class="rarityTextClass">
                {{ rewardText }}
              </div>

              <button
                type="button"
                class="blindbox-close-btn"
                :class="rarityBtnClass"
                @click="handleClose"
              >
                {{ t('checkin.normalCheckin') === t('checkin.normalCheckin') ? t('common.confirm') : t('common.confirm') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { BlindboxResult } from '@/api/checkin'

interface Props {
  show: boolean
  result: BlindboxResult | null
}

const props = defineProps<Props>()
const emit = defineEmits<{ (e: 'close'): void }>()
const { t } = useI18n()

const rarityLabel = computed(() => {
  if (!props.result) return ''
  const map: Record<string, string> = {
    common: t('checkin.blindboxCommon'),
    rare: t('checkin.blindboxRare'),
    epic: t('checkin.blindboxEpic'),
    legendary: t('checkin.blindboxLegendary'),
  }
  return map[props.result.rarity] || props.result.rarity
})

const rewardText = computed(() => {
  if (!props.result) return ''
  const v = props.result.reward_value
  switch (props.result.reward_type) {
    case 'balance':
      return t('checkin.blindboxBalanceReward', { value: v.toFixed(2) })
    case 'concurrency':
      return t('checkin.blindboxConcurrencyReward', { value: v })
    case 'subscription':
      return t('checkin.blindboxSubscriptionReward', { days: v })
    case 'invitation_code':
      return t('checkin.blindboxInvitationReward')
    default:
      return `${props.result.reward_type}: ${v}`
  }
})

const rarityGlowClass = computed(() => {
  if (!props.result) return ''
  return `glow-${props.result.rarity}`
})

const rarityBgClass = computed(() => {
  if (!props.result) return ''
  return `icon-${props.result.rarity}`
})

const rarityBadgeClass = computed(() => {
  if (!props.result) return ''
  return `badge-${props.result.rarity}`
})

const rarityTextClass = computed(() => {
  if (!props.result) return ''
  return `text-${props.result.rarity}`
})

const rarityBtnClass = computed(() => {
  if (!props.result) return ''
  return `btn-${props.result.rarity}`
})

function sparkleStyle(i: number) {
  const angle = (i * 45) * Math.PI / 180
  const r = 90 + Math.random() * 30
  const x = Math.cos(angle) * r
  const y = Math.sin(angle) * r
  const delay = i * 0.2
  const size = 4 + Math.random() * 6
  return {
    left: `calc(50% + ${x}px - ${size / 2}px)`,
    top: `calc(50% + ${y}px - ${size / 2}px)`,
    width: `${size}px`,
    height: `${size}px`,
    animationDelay: `${delay}s`,
  }
}

function handleClose() {
  emit('close')
}
</script>

<style scoped>
.blindbox-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.blindbox-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.blindbox-glow {
  position: absolute;
  width: 320px;
  height: 320px;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.4;
  animation: pulse-glow 2s ease-in-out infinite;
}

.blindbox-card {
  position: relative;
  z-index: 1;
  width: 340px;
  animation: card-enter 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.blindbox-card-inner {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 40px 32px;
  border-radius: 24px;
  background: white;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
}

:root.dark .blindbox-card-inner {
  background: #1e1e2e;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.6);
}

html.dark .blindbox-card-inner {
  background: #1e1e2e;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.6);
}

.blindbox-icon-wrapper {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: icon-bounce 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) 0.2s both;
}

.blindbox-icon {
  width: 36px;
  height: 36px;
}

.blindbox-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  text-align: center;
}

html.dark .blindbox-title {
  color: #f3f4f6;
}

.blindbox-prize-name {
  font-size: 24px;
  font-weight: 800;
  text-align: center;
  animation: prize-pop 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) 0.3s both;
}

.blindbox-rarity-badge {
  padding: 4px 16px;
  border-radius: 9999px;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.blindbox-reward {
  font-size: 20px;
  font-weight: 700;
  text-align: center;
}

.blindbox-close-btn {
  margin-top: 8px;
  padding: 10px 40px;
  border-radius: 14px;
  font-size: 15px;
  font-weight: 600;
  color: white;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
}

.blindbox-close-btn:hover {
  opacity: 0.9;
}

/* Glow colors */
.glow-common { background-color: #9ca3af; }
.glow-rare { background-color: #3b82f6; }
.glow-epic { background-color: #8b5cf6; }
.glow-legendary { background-color: #f59e0b; }

/* Icon backgrounds */
.icon-common { background-color: #f3f4f6; color: #6b7280; }
.icon-rare { background-color: #dbeafe; color: #2563eb; }
.icon-epic { background-color: #ede9fe; color: #7c3aed; }
.icon-legendary { background-color: #fef3c7; color: #d97706; }

html.dark .icon-common { background-color: #374151; color: #9ca3af; }
html.dark .icon-rare { background-color: #1e3a5f; color: #60a5fa; }
html.dark .icon-epic { background-color: #2d1b69; color: #a78bfa; }
html.dark .icon-legendary { background-color: #451a03; color: #fbbf24; }

/* Badge styles */
.badge-common { background-color: #f3f4f6; color: #6b7280; }
.badge-rare { background-color: #dbeafe; color: #2563eb; }
.badge-epic { background-color: #ede9fe; color: #7c3aed; }
.badge-legendary { background-color: #fef3c7; color: #d97706; }

html.dark .badge-common { background-color: #374151; color: #9ca3af; }
html.dark .badge-rare { background-color: #1e3a5f; color: #60a5fa; }
html.dark .badge-epic { background-color: #2d1b69; color: #a78bfa; }
html.dark .badge-legendary { background-color: #451a03; color: #fbbf24; }

/* Reward text colors */
.text-common { color: #6b7280; }
.text-rare { color: #2563eb; }
.text-epic { color: #7c3aed; }
.text-legendary { color: #d97706; }

html.dark .text-common { color: #9ca3af; }
html.dark .text-rare { color: #60a5fa; }
html.dark .text-epic { color: #a78bfa; }
html.dark .text-legendary { color: #fbbf24; }

/* Button colors */
.btn-common { background-color: #6b7280; }
.btn-rare { background-color: #2563eb; }
.btn-epic { background-color: #7c3aed; }
.btn-legendary { background-color: #d97706; }

/* Sparkles for legendary */
.blindbox-sparkles {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.sparkle {
  position: absolute;
  width: 6px;
  height: 6px;
  background: #fbbf24;
  border-radius: 50%;
  animation: sparkle-float 2s ease-in-out infinite;
  box-shadow: 0 0 6px #fbbf24;
}

/* Animations */
@keyframes card-enter {
  from {
    opacity: 0;
    transform: scale(0.5) rotateY(180deg);
  }
  to {
    opacity: 1;
    transform: scale(1) rotateY(0deg);
  }
}

@keyframes icon-bounce {
  from {
    opacity: 0;
    transform: scale(0);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes prize-pop {
  from {
    opacity: 0;
    transform: scale(0.5);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes pulse-glow {
  0%, 100% {
    opacity: 0.3;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.1);
  }
}

@keyframes sparkle-float {
  0%, 100% {
    opacity: 0;
    transform: translateY(0) scale(0);
  }
  25% {
    opacity: 1;
    transform: translateY(-10px) scale(1);
  }
  75% {
    opacity: 0.5;
    transform: translateY(-20px) scale(0.5);
  }
}

/* Modal transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
