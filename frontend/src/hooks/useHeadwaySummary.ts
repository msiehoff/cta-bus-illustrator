import { useEffect, useState } from 'react'
import type { HeadwaySummaryResponse } from '../types/api'
import { adminFetch, ADMIN_API_BASE } from '../lib/adminApi'
import type { HeadwayFilters } from './useHeadways'

export const useHeadwaySummary = (filters: Omit<HeadwayFilters, 'sort' | 'limit' | 'offset'>) => {
  const [data, setData] = useState<HeadwaySummaryResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      try {
        const params = new URLSearchParams()
        if (filters.route) params.set('route', filters.route)
        if (filters.direction) params.set('direction', filters.direction)
        if (filters.stop) params.set('stop', filters.stop)
        if (filters.vehicle) params.set('vehicle', filters.vehicle)
        if (filters.date) params.set('date', filters.date)

        const res = await adminFetch(`${ADMIN_API_BASE}/headways/summary?${params}`)
        if (!res.ok) throw new Error('Failed to load headway summary')
        const json = await res.json()
        if (!cancelled) {
          setData(json)
          setError(null)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load headway summary')
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    void load()
    return () => {
      cancelled = true
    }
  }, [filters.route, filters.direction, filters.stop, filters.vehicle, filters.date])

  return { data, loading, error }
}
