const ADMIN_API_BASE = '/api/v1/admin'

const adminFetch = (input: RequestInfo | URL, init?: RequestInit) =>
  fetch(input, { credentials: 'include', ...init })

export const getAdminSession = async () => {
  const res = await adminFetch(`${ADMIN_API_BASE}/session`)
  if (!res.ok) throw new Error('Failed to check admin session')
  return res.json()
}

export const adminLogin = async (username: string, password: string) => {
  const res = await adminFetch(`${ADMIN_API_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error ?? 'Login failed')
  }
  return res.json()
}

export const adminLogout = async () => {
  const res = await adminFetch(`${ADMIN_API_BASE}/logout`, { method: 'POST' })
  if (!res.ok) throw new Error('Logout failed')
}

export { adminFetch, ADMIN_API_BASE }
