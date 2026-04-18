import { apiClient } from '../client'

export interface GroupHealth {
  group_id: number
  group_name: string
  total_accounts: number
  active_accounts: number
  error_accounts: number
  rate_limited: number
  overload: number
  disabled: number
}

export interface ModelLatency {
  model: string
  request_count: number
  avg_latency_ms: number
  p50_latency_ms: number
  p95_latency_ms: number
  p99_latency_ms: number
  avg_first_token_ms: number
}

export interface ErrorAccount {
  account_id: number
  account_name: string
  group_name: string
  status: string
  error_message: string
  rate_limited_at?: string
  overload_until?: string
}

export interface MonitoringOverview {
  groups: GroupHealth[]
  model_latencies: ModelLatency[]
  error_accounts: ErrorAccount[]
  total_requests_today: number
  avg_latency_ms_today: number
}

export async function getMonitoringOverview(): Promise<MonitoringOverview> {
  const { data } = await apiClient.get<MonitoringOverview>('/monitoring/overview')
  return data
}

export const monitoringAPI = {
  getOverview: getMonitoringOverview,
}

export default monitoringAPI
