import { apiClient } from './client'
import type { ApiResponse } from '@/types'

export interface LeaderboardEntry {
  rank: number
  username: string
  value: number
}

export interface LeaderboardData {
  items: LeaderboardEntry[]
  total: number
  page: number
  page_size: number
  pages: number
}

export async function getBalanceLeaderboard(page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<ApiResponse<LeaderboardData>>('/public/leaderboard/balance', { params: { page, page_size: pageSize } })
  return data.data
}

export async function getConsumptionLeaderboard(period: 'daily' | 'weekly' | 'monthly' = 'daily', page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<ApiResponse<LeaderboardData>>('/public/leaderboard/consumption', { params: { period, page, page_size: pageSize } })
  return data.data
}

export async function getCheckinLeaderboard(page = 1, pageSize = 10): Promise<LeaderboardData> {
  const { data } = await apiClient.get<ApiResponse<LeaderboardData>>('/public/leaderboard/checkin', { params: { page, page_size: pageSize } })
  return data.data
}

export const leaderboardAPI = {
  getBalanceLeaderboard,
  getConsumptionLeaderboard,
  getCheckinLeaderboard,
}
