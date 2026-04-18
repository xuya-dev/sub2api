import { apiClient } from './client'

export interface CheckinResult {
  reward_amount: number
  streak_days: number
  checked_at: string
}

export interface CheckinStatus {
  can_checkin: boolean
  streak_days: number
  today_reward: number | null
  min_reward: number
  max_reward: number
}

export async function checkin(): Promise<CheckinResult> {
  const { data } = await apiClient.post<CheckinResult>('/checkin')
  return data
}

export async function getCheckinStatus(): Promise<CheckinStatus> {
  const { data } = await apiClient.get<CheckinStatus>('/checkin/status')
  return data
}

export const checkinAPI = {
  checkin,
  getCheckinStatus
}

export default checkinAPI
