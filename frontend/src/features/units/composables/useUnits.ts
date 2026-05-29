import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import {
  fetchUnits,
  createUnit,
  updateUnit,
  deleteUnit,
} from '@/api/units'
import type { CreateUnitPayload, UpdateUnitPayload } from '@/types'

export function useUnits(projectId: number) {
  const queryClient = useQueryClient()
  const key = () => ['units', projectId]

  const unitsQuery = useQuery({
    queryKey: key(),
    queryFn: () => fetchUnits(projectId),
  })

  const createMutation = useMutation({
    mutationFn: (payload: CreateUnitPayload) => createUnit(projectId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: key() })
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ unitId, payload }: { unitId: number; payload: UpdateUnitPayload }) =>
      updateUnit(projectId, unitId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: key() })
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (unitId: number) => deleteUnit(projectId, unitId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: key() })
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
    },
  })

  return { unitsQuery, createMutation, updateMutation, deleteMutation }
}
