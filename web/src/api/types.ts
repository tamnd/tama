// Payload types for the API. Each one mirrors a type in pkg/api and is kept
// in lockstep by review; the Go side is the source of truth.

/** UserPayload in pkg/api/auth.go. */
export type User = {
  id: number
  username: string
  display_name: string
  avatar: string
  is_admin: boolean
  created_at: number
}

/** credentials in pkg/api/auth.go; register and login share it. */
export type Credentials = {
  username: string
  password: string
}

/** POST /api/auth/logout response body. */
export type LogoutResult = {
  ok: boolean
}

/** Language in pkg/course/catalog.go. */
export type Language = {
  code: string
  name: string
  native: string
  rtl?: boolean
}

/** Course in pkg/course/catalog.go; GET /api/catalog returns a list. */
export type Course = {
  id: string
  base: Language
  target: Language
}

/** The closed code set from pkg/api/errors.go; the HTTP status always
 * matches the code, so callers switch on the code alone. */
export type ErrorCode =
  | 'bad_request'
  | 'unauthorized'
  | 'forbidden'
  | 'not_found'
  | 'conflict'
  | 'rate_limited'
  | 'internal'
