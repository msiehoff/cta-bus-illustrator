import { useEffect, useState } from 'react'
import type { HeadwayRouteDetailResponse } from '../types/api'

const DEFAULT_DAYS = 30

export const useHeadwayRoute = (externalId: string, days = DEFAULT_DAYS) => {
  const [data, setData] = useState<HeadwayRouteDetailResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!externalId) return
    let cancelled = false
    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const res = await fetch(
          `/api/v1/headways/routes/${encodeURIComponent(externalId)}?days=${days}`,
        )
        if (!res.ok) throw new Error(`Failed to load route headways (${res.status})`)
        const json = (await res.json()) as HeadwayRouteDetailResponse
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
  }, [externalId, days])

  return { data, loading, error }
}
