import { request } from './client'
import type {
  HobbyProject,
  PaintProfile,
  CreateProjectPayload,
  UpdateProjectPayload,
} from '@/types'

export async function fetchProjects(): Promise<HobbyProject[]> {
  const data = await request<{ hobbyProjects: HobbyProject[] }>('/hobby-projects')
  return data.hobbyProjects
}

export async function fetchProject(projectId: number): Promise<HobbyProject> {
  const data = await request<{ hobbyProject: HobbyProject }>(
    `/hobby-projects/${projectId}`,
  )
  return data.hobbyProject
}

export async function createProject(payload: CreateProjectPayload): Promise<HobbyProject> {
  const data = await request<{ hobbyProject: HobbyProject }>('/hobby-projects', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  return data.hobbyProject
}

export async function updateProject(
  projectId: number,
  payload: UpdateProjectPayload,
): Promise<HobbyProject> {
  const data = await request<{ hobbyProject: HobbyProject }>(
    `/hobby-projects/${projectId}`,
    {
      method: 'PUT',
      body: JSON.stringify(payload),
    },
  )
  return data.hobbyProject
}

export async function deleteProject(projectId: number): Promise<void> {
  await request<void>(`/hobby-projects/${projectId}`, { method: 'DELETE' })
}

export async function fetchProjectPaintProfiles(
  projectId: number,
): Promise<PaintProfile[]> {
  const data = await request<{ paintProfiles: PaintProfile[] }>(
    `/hobby-projects/${projectId}/paint-profiles`,
  )
  return data.paintProfiles
}

export async function assignPaintProfile(
  projectId: number,
  profileId: number,
): Promise<void> {
  await request<void>(`/hobby-projects/${projectId}/paint-profiles/${profileId}`, {
    method: 'POST',
  })
}

export async function unassignPaintProfile(
  projectId: number,
  profileId: number,
): Promise<void> {
  await request<void>(`/hobby-projects/${projectId}/paint-profiles/${profileId}`, {
    method: 'DELETE',
  })
}
