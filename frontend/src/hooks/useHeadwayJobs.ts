import { useCallback, useEffect, useState } from 'react'
import type { ListHeadwayJobRunsResponse } from '../types/api'
import { adminFetch, ADMIN_API_BASE } from '../lib/adminApi'

export const useHeadwayJobRuns = (limit = 50) => {
  const [data, setData] = useState<ListHeadwayJobRunsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const reload = useCallback(async () => {
    setLoading(true)
    try {
      const res = await adminFetch(`${ADMIN_API_BASE}/headways/runs?limit=${limit}`)
      if (!res.ok) throw new Error('Failed to load headway job runs')
      const json = await res.json()
      setData(json)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load headway job runs')
    } finally {
      setLoading(false)
    }
  }, [limit])

  useEffect(() => {
    void reload()
  }, [reload])

  return { data, loading, error, reload }
}

export const runHeadwayJob = async (serviceDate?: string) => {
  const res = await adminFetch(`${ADMIN_API_BASE}/headways/run`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(serviceDate ? { service_date: serviceDate } : {}),
  })
  const body = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(body.errorMessage || body.error || 'Headway run failed')
  }
  return body
}
