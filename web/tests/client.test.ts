import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, login, me, setUnauthorizedHandler } from '@/api/client'
import type { User } from '@/api/types'

// The client against a mocked fetch: envelope unwrapping on success, the
// ApiError shape on failures, and the single 401 interceptor.

function stubFetch(status: number, body: unknown) {
  const mock = vi.fn(() =>
    Promise.resolve({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(body),
    }),
  )
  vi.stubGlobal('fetch', mock)
  return mock
}

const demo: User = {
  id: 1,
  username: 'demo',
  display_name: '',
  avatar: '',
  is_admin: false,
  created_at: 1700000000000,
}

afterEach(() => {
  vi.unstubAllGlobals()
  setUnauthorizedHandler(undefined)
})

describe('client', () => {
  it('unwraps the {data} envelope', async () => {
    stubFetch(200, { data: demo })
    await expect(me()).resolves.toEqual(demo)
  })

  it('posts credentials as JSON', async () => {
    const mock = stubFetch(200, { data: demo })
    await login({ username: 'demo', password: 'demo1234' })
    expect(mock).toHaveBeenCalledWith('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: 'demo', password: 'demo1234' }),
    })
  })

  it('throws a typed ApiError on 404', async () => {
    stubFetch(404, { error: { code: 'not_found', message: 'course es-from-en does not exist' } })
    const err = await me().then(
      () => null,
      (e: unknown) => e,
    )
    expect(err).toBeInstanceOf(ApiError)
    const apiErr = err as ApiError
    expect(apiErr.code).toBe('not_found')
    expect(apiErr.message).toBe('course es-from-en does not exist')
  })

  it('throws on 401 and fires the unauthorized handler once', async () => {
    stubFetch(401, { error: { code: 'unauthorized', message: 'authentication required' } })
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    const err = await me().then(
      () => null,
      (e: unknown) => e,
    )
    expect(err).toBeInstanceOf(ApiError)
    expect((err as ApiError).code).toBe('unauthorized')
    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('does not fire the unauthorized handler on other errors', async () => {
    stubFetch(404, { error: { code: 'not_found', message: 'nope' } })
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    await expect(me()).rejects.toBeInstanceOf(ApiError)
    expect(handler).not.toHaveBeenCalled()
  })
})
