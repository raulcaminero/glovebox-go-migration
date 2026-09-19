export type RoutingMode = 'auto' | 'go' | 'legacy';

export interface Policyholder {
  id: string;
  full_name?: string;
  fullName?: string;
  email: string;
  phone?: string;
  created_at?: string;
}

export interface InspectionResult {
  endpoint: string;
  status: string;
  ok: boolean;
  routedTo: string;
  targetURL: string;
  latencyMs: number;
  body: string;
}

export interface AdrInfo {
  id: string;
  title: string;
  summary: string;
  status: string;
}
