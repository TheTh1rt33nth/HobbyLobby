// SECURITY NOTE: Backend paint-step endpoints authorize access to the parent
// profile but do not clearly verify that stepId belongs to the given profileId.
// This is a known authorization gap — flagged for backend team review.
// Do not treat update/delete step routes as fully safe until verified.

import { request } from './client'
import type {
  PaintProfile,
  PaintStep,
  CreatePaintProfilePayload,
  UpdatePaintProfilePayload,
  CreatePaintStepPayload,
  UpdatePaintStepPayload,
} from '@/types'

export async function fetchPaintProfiles(): Promise<PaintProfile[]> {
  const data = await request<{ paintProfiles: PaintProfile[] }>('/paint-profiles')
  return data.paintProfiles
}

export async function fetchPaintProfile(profileId: number): Promise<PaintProfile> {
  const data = await request<{ paintProfile: PaintProfile }>(
    `/paint-profiles/${profileId}`,
  )
  return data.paintProfile
}

export async function createPaintProfile(
  payload: CreatePaintProfilePayload,
): Promise<PaintProfile> {
  const data = await request<{ paintProfile: PaintProfile }>('/paint-profiles', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  return data.paintProfile
}

export async function updatePaintProfile(
  profileId: number,
  payload: UpdatePaintProfilePayload,
): Promise<PaintProfile> {
  const data = await request<{ paintProfile: PaintProfile }>(
    `/paint-profiles/${profileId}`,
    {
      method: 'PUT',
      body: JSON.stringify(payload),
    },
  )
  return data.paintProfile
}

export async function deletePaintProfile(profileId: number): Promise<void> {
  await request<void>(`/paint-profiles/${profileId}`, { method: 'DELETE' })
}

export async function createPaintStep(
  profileId: number,
  payload: CreatePaintStepPayload,
): Promise<PaintStep> {
  const data = await request<{ paintStep: PaintStep }>(
    `/paint-profiles/${profileId}/steps`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
    },
  )
  return data.paintStep
}

export async function updatePaintStep(
  profileId: number,
  stepId: number,
  payload: UpdatePaintStepPayload,
): Promise<PaintStep> {
  const data = await request<{ paintStep: PaintStep }>(
    `/paint-profiles/${profileId}/steps/${stepId}`,
    {
      method: 'PUT',
      body: JSON.stringify(payload),
    },
  )
  return data.paintStep
}

export async function deletePaintStep(
  profileId: number,
  stepId: number,
): Promise<void> {
  await request<void>(`/paint-profiles/${profileId}/steps/${stepId}`, {
    method: 'DELETE',
  })
}
