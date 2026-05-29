import { request } from './client'
import type { Unit, CreateUnitPayload, UpdateUnitPayload } from '@/types'

export async function fetchUnits(projectId: number): Promise<Unit[]> {
  const data = await request<{ units: Unit[] }>(
    `/hobby-projects/${projectId}/units`,
  )
  return data.units
}

export async function fetchUnit(projectId: number, unitId: number): Promise<Unit> {
  const data = await request<{ unit: Unit }>(
    `/hobby-projects/${projectId}/units/${unitId}`,
  )
  return data.unit
}

export async function createUnit(
  projectId: number,
  payload: CreateUnitPayload,
): Promise<Unit> {
  const data = await request<{ unit: Unit }>(
    `/hobby-projects/${projectId}/units`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
    },
  )
  return data.unit
}

export async function updateUnit(
  projectId: number,
  unitId: number,
  payload: UpdateUnitPayload,
): Promise<Unit> {
  const data = await request<{ unit: Unit }>(
    `/hobby-projects/${projectId}/units/${unitId}`,
    {
      method: 'PUT',
      body: JSON.stringify(payload),
    },
  )
  return data.unit
}

export async function deleteUnit(projectId: number, unitId: number): Promise<void> {
  await request<void>(`/hobby-projects/${projectId}/units/${unitId}`, {
    method: 'DELETE',
  })
}
