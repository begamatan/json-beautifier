import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import JsonTool from '../JsonTool.vue'
import * as api from '../../api'

vi.mock('../../api')

describe('JsonTool', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders heading and core buttons', () => {
    const wrapper = mount(JsonTool)
    expect(wrapper.find('h1').text()).toBe('JSON Beautifier')
    expect(wrapper.find('[data-testid="btn-beautify"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="btn-minify"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="btn-validate"]').exists()).toBe(true)
  })

  it('calls beautify API and shows result', async () => {
    vi.mocked(api.beautify).mockResolvedValue({ result: '{\n  "a": 1\n}' })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{"a":1}')
    await wrapper.find('[data-testid="btn-beautify"]').trigger('click')
    await flushPromises()

    expect(api.beautify).toHaveBeenCalledWith('{"a":1}', 2)
    expect((wrapper.find('[data-testid="output-area"]').element as HTMLTextAreaElement).value).toBe(
      '{\n  "a": 1\n}',
    )
  })

  it('calls minify API and shows result', async () => {
    vi.mocked(api.minify).mockResolvedValue({ result: '{"a":1}' })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{\n  "a": 1\n}')
    await wrapper.find('[data-testid="btn-minify"]').trigger('click')
    await flushPromises()

    expect(api.minify).toHaveBeenCalledWith('{\n  "a": 1\n}')
    expect((wrapper.find('[data-testid="output-area"]').element as HTMLTextAreaElement).value).toBe(
      '{"a":1}',
    )
  })

  it('shows success message when JSON is valid', async () => {
    vi.mocked(api.validate).mockResolvedValue({ valid: true, message: 'JSON is valid' })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{"a":1}')
    await wrapper.find('[data-testid="btn-validate"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="success-msg"]').text()).toBe('JSON is valid')
    expect(wrapper.find('[data-testid="error-msg"]').exists()).toBe(false)
  })

  it('shows error message when JSON is invalid', async () => {
    vi.mocked(api.validate).mockResolvedValue({
      valid: false,
      message: 'input is not valid JSON',
    })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{bad}')
    await wrapper.find('[data-testid="btn-validate"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="error-msg"]').text()).toBe('input is not valid JSON')
  })

  it('shows error message on API failure', async () => {
    vi.mocked(api.beautify).mockRejectedValue({
      code: 'INVALID_JSON',
      message: 'input is not valid JSON',
    })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{bad}')
    await wrapper.find('[data-testid="btn-beautify"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="error-msg"]').text()).toBe('input is not valid JSON')
  })

  it('uses 4-space indent when selected', async () => {
    vi.mocked(api.beautify).mockResolvedValue({ result: '{\n    "a": 1\n}' })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="indent-select"]').setValue('4')
    await wrapper.find('[data-testid="input-area"]').setValue('{"a":1}')
    await wrapper.find('[data-testid="btn-beautify"]').trigger('click')
    await flushPromises()

    expect(api.beautify).toHaveBeenCalledWith('{"a":1}', 4)
  })

  it('clear button resets input, output, and messages', async () => {
    vi.mocked(api.beautify).mockResolvedValue({ result: '{}' })

    const wrapper = mount(JsonTool)
    await wrapper.find('[data-testid="input-area"]').setValue('{}')
    await wrapper.find('[data-testid="btn-beautify"]').trigger('click')
    await flushPromises()

    await wrapper.find('[data-testid="btn-clear"]').trigger('click')

    expect((wrapper.find('[data-testid="input-area"]').element as HTMLTextAreaElement).value).toBe(
      '',
    )
    expect((wrapper.find('[data-testid="output-area"]').element as HTMLTextAreaElement).value).toBe(
      '',
    )
  })

  it('copy and download buttons are disabled when output is empty', () => {
    const wrapper = mount(JsonTool)
    expect((wrapper.find('[data-testid="btn-copy"]').element as HTMLButtonElement).disabled).toBe(
      true,
    )
    expect(
      (wrapper.find('[data-testid="btn-download"]').element as HTMLButtonElement).disabled,
    ).toBe(true)
  })
})
