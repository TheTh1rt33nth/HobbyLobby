export type HobbyStatus =
  | 'unassembled'
  | 'assembled'
  | 'primed'
  | 'base_coated'
  | 'painted'
  | 'based'
  | 'complete'

// Ordered from least to most progress (used for sorting / display order)
export const HOBBY_STATUS_ORDER: HobbyStatus[] = [
  'unassembled',
  'assembled',
  'primed',
  'base_coated',
  'painted',
  'based',
  'complete',
]

export const HOBBY_STATUS_LABELS: Record<HobbyStatus, string> = {
  unassembled: 'UNASSEMBLED',
  assembled:   'ASSEMBLED',
  primed:      'PRIMED',
  base_coated: 'BASE COATED',
  painted:     'PAINTED',
  based:       'BASED',
  complete:    'COMPLETE',
}

// API models

export interface User {
  id: number
  username: string
  email: string
  createdAt: string
  updatedAt: string
}

export interface ProjectProgress {
  totalUnits: number
  byStatus: Record<HobbyStatus, number>
}

export interface HobbyProject {
  id: number
  userId: number
  name: string
  description: string | null
  gameSystem: string | null
  faction: string | null
  progress?: ProjectProgress
  createdAt: string
  updatedAt: string
}

export interface Unit {
  id: number
  projectId: number
  name: string
  quantity: number
  status: HobbyStatus
  notes: string | null
  paintProfileId: number | null
  createdAt: string
  updatedAt: string
}

export interface PaintStep {
  id: number
  paintProfileId: number
  stepOrder: number
  paintName: string
  brand: string | null
  paintType: string | null
  applicationMethod: string | null
  colorHex: string | null
  notes: string | null
}

export interface PaintProfile {
  id: number
  userId: number
  name: string
  description: string | null
  targetArea: string | null
  steps?: PaintStep[]
  createdAt: string
  updatedAt: string
}

// Mutation Payloads

export interface RegisterPayload {
  username: string
  email: string
  password: string
}

export interface LoginPayload {
  username: string
  password: string
}

export interface CreateProjectPayload {
  name: string
  description?: string | null
  gameSystem?: string | null
  faction?: string | null
}

export type UpdateProjectPayload = CreateProjectPayload

export interface CreateUnitPayload {
  name: string
  quantity: number
  status: HobbyStatus
  notes?: string | null
  paintProfileId?: number | null
}

export type UpdateUnitPayload = CreateUnitPayload

export interface CreatePaintProfilePayload {
  name: string
  description?: string | null
  targetArea?: string | null
}

export type UpdatePaintProfilePayload = CreatePaintProfilePayload

export interface CreatePaintStepPayload {
  stepOrder: number
  paintName: string
  brand?: string | null
  paintType?: string | null
  applicationMethod?: string | null
  colorHex?: string | null
  notes?: string | null
}

export type UpdatePaintStepPayload = CreatePaintStepPayload
