import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SegmentedProgress from '@/components/SegmentedProgress.vue'
import type { ProjectProgress } from '@/types'

function makeProgress(overrides: Partial<ProjectProgress['byStatus']> = {}): ProjectProgress {
  return {
    totalUnits: 10,
    byStatus: {
      unassembled: 0,
      assembled:   0,
      primed:      0,
      base_coated: 0,
      painted:     0,
      based:       0,
      complete:    0,
      ...overrides,
    },
  }
}

describe('SegmentedProgress', () => {
  it('renders a segment for each non-zero status', () => {
    const progress = makeProgress({ primed: 4, complete: 6 })
    const wrapper = mount(SegmentedProgress, { props: { progress } })
    const segments = wrapper.findAll('.segmented-progress__segment')
    expect(segments).toHaveLength(2)
  })

  it('renders the empty placeholder when totalUnits is 0', () => {
    const progress: ProjectProgress = { totalUnits: 0, byStatus: {} as ProjectProgress['byStatus'] }
    const wrapper = mount(SegmentedProgress, { props: { progress } })
    expect(wrapper.find('.segmented-progress__empty').exists()).toBe(true)
  })

  it('segment width is proportional to unit count', () => {
    const progress = makeProgress({ primed: 4, complete: 6 })
    const wrapper = mount(SegmentedProgress, { props: { progress } })
    const segments = wrapper.findAll('.segmented-progress__segment')

    const primedStyle   = segments.find((s) => s.classes().includes('segmented-progress__segment--primed'))
    const completeStyle = segments.find((s) => s.classes().includes('segmented-progress__segment--complete'))

    expect(primedStyle?.attributes('style')).toContain('width: 40%')
    expect(completeStyle?.attributes('style')).toContain('width: 60%')
  })
})
