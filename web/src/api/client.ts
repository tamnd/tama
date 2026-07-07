import type { Course, Credentials, ErrorCode, LogoutResult, User } from './types'

// The hand-kept API client. Every call goes through request(), which owns
// the {data}/{error} envelope, so endpoint functions stay one-liners.

/** The {error:{code,message}} envelope as a throwable. */
export class ApiError extends Error {
  readonly code: ErrorCode

  constructor(code: ErrorCode, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
  }
}

type ErrorEnvelope = { error?: { code?: ErrorCode; message?: string } }

// The one 401 interceptor: the app points this at the router once, so an
// expired session lands on /login without per-call handling.
let onUnauthorized: (() => void) | undefined

export function setUnauthorizedHandler(handler: (() => void) | undefined): void {
  onUnauthorized = handler
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(path, {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  if (!res.ok) {
    let code: ErrorCode = 'internal'
    let message = `unexpected response ${String(res.status)}`
    try {
      const envelope = (await res.json()) as ErrorEnvelope
      if (envelope.error?.code) code = envelope.error.code
      if (envelope.error?.message) message = envelope.error.message
    } catch {
      // A non-JSON body keeps the fallback message.
    }
    if (res.status === 401) onUnauthorized?.()
    throw new ApiError(code, message)
  }

  const envelope = (await res.json()) as { data: T }
  return envelope.data
}

export function register(creds: Credentials): Promise<User> {
  return request<User>('POST', '/api/auth/register', creds)
}

export function login(creds: Credentials): Promise<User> {
  return request<User>('POST', '/api/auth/login', creds)
}

export function logout(): Promise<LogoutResult> {
  return request<LogoutResult>('POST', '/api/auth/logout')
}

export function me(): Promise<User> {
  return request<User>('GET', '/api/me')
}

export function catalog(from: string): Promise<Course[]> {
  return request<Course[]>('GET', `/api/catalog?from=${encodeURIComponent(from)}`)
}
