import { apiClient } from './client'

export interface TransferRecord {
  id: number
  sender_id: number
  sender_email: string
  receiver_id: number
  receiver_email: string
  amount: number
  fee: number
  fee_rate: number
  gross_amount: number
  transfer_type: 'direct' | 'redpacket' | 'batch'
  status: 'completed' | 'frozen' | 'revoked'
  memo: string | null
  redpacket_id: number | null
  created_at: string
}

export interface TransferStats {
  total_sent: number
  total_received: number
  total_fee_paid: number
}

export interface UserSearchResult {
  id: number
  email: string
  username: string
}

export async function searchUsers(query: string): Promise<UserSearchResult[]> {
  const { data } = await apiClient.get<UserSearchResult[]>('/transfer/search-users', { params: { q: query } })
  return data
}

export interface TransferLeaderboardEntry {
  rank: number
  user_id: number
  email: string
  total_amount: number
  total_count: number
}

export async function transferBalance(receiverId: number, amount: number, memo?: string): Promise<TransferRecord> {
  const { data } = await apiClient.post<TransferRecord>('/transfer', {
    receiver_id: receiverId,
    amount,
    memo,
  })
  return data
}

export async function validateTransfer(receiverId: number, amount: number): Promise<{ fee: number; fee_rate: number }> {
  const { data } = await apiClient.post<{ fee: number; fee_rate: number }>('/transfer/validate', {
    receiver_id: receiverId,
    amount,
  })
  return data
}

export async function getTransferHistory(params: {
  role?: string
  page?: number
  page_size?: number
}): Promise<{ items: TransferRecord[]; total: number; page: number; page_size: number }> {
  const { data } = await apiClient.get('/transfer/history', { params })
  return data
}

export async function getTransferStats(): Promise<TransferStats> {
  const { data } = await apiClient.get<TransferStats>('/transfer/stats')
  return data
}

export async function getTransferLeaderboard(params: {
  period?: string
  limit?: number
}): Promise<TransferLeaderboardEntry[]> {
  const { data } = await apiClient.get<TransferLeaderboardEntry[]>('/transfer/leaderboard', { params })
  return data
}

export interface RedPacketRecord {
  id: number
  sender_id: number
  total_amount: number
  total_count: number
  remaining_amount: number
  remaining_count: number
  redpacket_type: 'equal' | 'random'
  fee: number
  fee_rate: number
  code: string
  status: 'active' | 'expired' | 'exhausted'
  memo: string | null
  expire_at: string
  created_at: string
}

export interface RedPacketClaimRecord {
  id: number
  redpacket_id: number
  user_id: number
  user_email: string
  amount: number
  transfer_id: number | null
  created_at: string
}

export async function createRedPacket(params: {
  total_amount: number
  count: number
  redpacket_type?: 'equal' | 'random'
  memo?: string
}): Promise<RedPacketRecord> {
  const { data } = await apiClient.post<RedPacketRecord>('/redpacket', params)
  return data
}

export async function claimRedPacket(code: string): Promise<RedPacketClaimRecord> {
  const { data } = await apiClient.post<RedPacketClaimRecord>('/redpacket/claim', { code })
  return data
}

export async function getRedPacketDetail(id: number): Promise<{
  redpacket: RedPacketRecord
  claims: RedPacketClaimRecord[]
}> {
  const { data } = await apiClient.get(`/redpacket/${id}`)
  return data
}

export async function getMyRedPackets(params: {
  page?: number
  page_size?: number
}): Promise<{ items: RedPacketRecord[]; total: number; page: number; page_size: number }> {
  const { data } = await apiClient.get('/redpacket/my', { params })
  return data
}
