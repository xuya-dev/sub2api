import { apiClient } from './client'

export interface PublicPricingGroup {
  id: number
  name: string
  platform: string
  rate_multiplier: number
}

export interface LiteLLMModelPricing {
  model_name: string
  input_cost_per_token: number
  output_cost_per_token: number
  cache_creation_input_token_cost?: number | null
  cache_read_input_token_cost?: number | null
  litellm_provider?: string
  mode?: string
}

export interface PublicPricingResponse {
  groups: PublicPricingGroup[]
  pricing: Record<string, LiteLLMModelPricing>
}

export async function getPublicPricing(): Promise<PublicPricingResponse> {
  const { data } = await apiClient.get<PublicPricingResponse>('/public/pricing')
  return data
}

export const pricingAPI = { getPublicPricing }
export default pricingAPI
