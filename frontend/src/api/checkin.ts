import { apiClient } from './client'

export interface CheckinResult {
  reward_amount: number
  streak_days: number
  checked_at: string
  checkin_type: string
  bet_amount?: number
  multiplier?: number
}

export interface CheckinStatus {
  enabled: boolean
  luck_enabled: boolean
  can_checkin: boolean
  streak_days: number
  today_reward: number | null
  today_checkin_type?: string
  today_multiplier?: number
  min_reward: number
  max_reward: number
  min_multiplier: number
  max_multiplier: number
  balance: number
}

export async function checkin(): Promise<CheckinResult> {
  const { data } = await apiClient.post<CheckinResult>('/checkin')
  return data
}

export async function luckCheckin(betAmount: number): Promise<CheckinResult> {
  const { data } = await apiClient.post<CheckinResult>('/checkin/luck', { bet_amount: betAmount })
  return data
}

export async function getCheckinStatus(): Promise<CheckinStatus> {
  const { data } = await apiClient.get<CheckinStatus>('/checkin/status')
  return data
}

export const checkinAPI = {
  checkin,
  luckCheckin,
  getCheckinStatus
}

export default checkinAPI
