import { useCallback, useEffect, useState } from 'react'
import type { AdminSessionResponse } from '../types/api'
import { getAdminSession } from '../lib/adminApi'

export const useAdminSession = () => {
  const [session, setSession] = useState<AdminSessionResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getAdminSession()
      setSession(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load session')
      setSession({ authenticated: false })
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  return { session, loading, error, refresh }
}
