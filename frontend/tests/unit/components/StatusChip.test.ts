import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusChip from '@/components/StatusChip.vue'
import { HOBBY_STATUS_LABELS, HOBBY_STATUS_ORDER } from '@/types'
import type { HobbyStatus } from '@/types'

describe('StatusChip', () => {
  it.each(HOBBY_STATUS_ORDER)('renders correct label for status: %s', (status) => {
    const wrapper = mount(StatusChip, { props: { status } })
    expect(wrapper.text()).toBe(HOBBY_STATUS_LABELS[status as HobbyStatus])
  })

  it.each(HOBBY_STATUS_ORDER)('applies the correct CSS class for status: %s', (status) => {
    const wrapper = mount(StatusChip, { props: { status } })
    expect(wrapper.find('.chip').classes()).toContain(`chip--${status}`)
  })
})
