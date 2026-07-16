import { useEffect, useState } from 'react'
import type { HeadwayRoutesListResponse } from '../types/api'

const DEFAULT_DAYS = 30

export const useHeadwayRoutes = (days = DEFAULT_DAYS) => {
  const [data, setData] = useState<HeadwayRoutesListResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const res = await fetch(`/api/v1/headways/routes?days=${days}`)
        if (!res.ok) throw new Error(`Failed to load headways (${res.status})`)
        const json = (await res.json()) as HeadwayRoutesListResponse
        if (!cancelled) setData(json)
      } catch (err) {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load')
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    void load()
    return () => {
      cancelled = true
    }
  }, [days])

  return { data, loading, error }
}
