import { apiClient } from './client'

export interface BlindboxResult {
  prize_name: string
  rarity: string
  reward_type: string
  reward_value: number
  subscription_days?: number
  reward_detail?: string
}

export interface CheckinResult {
  reward_amount: number
  streak_days: number
  checked_at: string
  checkin_type: string
  bet_amount?: number
  multiplier?: number
  blindbox?: BlindboxResult
}

export interface CheckinStatus {
  enabled: boolean
  luck_enabled: boolean
  blindbox_enabled: boolean
  blindbox_trigger_type?: string
  blindbox_interval?: number
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

export interface CheckinCalendarDay {
  date: string
  checked_in: boolean
  reward_type?: string
  reward_value?: number
  streak_days?: number
}

export interface CheckinCalendar {
  days: CheckinCalendarDay[]
}

export async function getCheckinCalendar(): Promise<CheckinCalendar> {
  const { data } = await apiClient.get<CheckinCalendar>('/checkin/calendar')
  return data
}

export const checkinAPI = {
  checkin,
  luckCheckin,
  getCheckinStatus,
  getCalendar: getCheckinCalendar,
  getBlindboxRecords,
}

export default checkinAPI

export interface BlindboxRecordItem {
  id: number
  prize_name: string
  rarity: string
  reward_type: string
  reward_value: number
  reward_detail?: string
  subscription_days?: number
  streak_days: number
  created_at: string
}

export interface BlindboxRecordList {
  items: BlindboxRecordItem[]
  total: number
}

export async function getBlindboxRecords(page = 1, pageSize = 20): Promise<BlindboxRecordList> {
  const { data } = await apiClient.get<BlindboxRecordList>('/checkin/blindbox/records', { params: { page, page_size: pageSize } })
  return data
}
