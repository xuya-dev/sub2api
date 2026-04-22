import { apiClient } from './client'

export interface LeaderboardEntry {
  rank: number
  username: string
  value: number
  extra_int?: number
  extra_int2?: number
  extra_float?: number
  extra_date?: string
}

export interface LeaderboardData {
  items: LeaderboardEntry[]
  total: number
  page: number
  page_size: number
  pages: number
}

export async function getBalanceLeaderboard(page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<LeaderboardData>('/public/leaderboard/balance', { params: { page, page_size: pageSize } })
  return data
}

export async function getConsumptionLeaderboard(period: 'daily' | 'weekly' | 'monthly' = 'daily', page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<LeaderboardData>('/public/leaderboard/consumption', { params: { period, page, page_size: pageSize } })
  return data
}

export async function getCheckinLeaderboard(page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<LeaderboardData>('/public/leaderboard/checkin', { params: { page, page_size: pageSize } })
  return data
}

export const leaderboardAPI = {
  getBalanceLeaderboard,
  getConsumptionLeaderboard,
  getCheckinLeaderboard,
}
