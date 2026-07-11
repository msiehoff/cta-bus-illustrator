import { useEffect, useState } from 'react'
import type { ListHeadwaySummariesResponse } from '../types/api'
import { adminFetch, ADMIN_API_BASE } from '../lib/adminApi'

export interface HeadwaySummaryListFilters {
  date?: string
  grain?: string
  method?: string
  route?: string
  direction?: string
  stop?: string
  sort?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export const useHeadwaySummaries = (filters: HeadwaySummaryListFilters) => {
  const [data, setData] = useState<ListHeadwaySummariesResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      try {
        const params = new URLSearchParams()
        if (filters.date) params.set('date', filters.date)
        if (filters.grain) params.set('grain', filters.grain)
        if (filters.method) params.set('method', filters.method)
        if (filters.route) params.set('route', filters.route)
        if (filters.direction) params.set('direction', filters.direction)
        if (filters.stop) params.set('stop', filters.stop)
        if (filters.sort) params.set('sort', filters.sort)
        if (filters.limit != null) params.set('limit', String(filters.limit))
        if (filters.offset != null) params.set('offset', String(filters.offset))

        const res = await adminFetch(`${ADMIN_API_BASE}/headway-summaries?${params}`)
        if (!res.ok) throw new Error('Failed to load headway summaries')
        const json = await res.json()
        if (!cancelled) {
          setData(json)
          setError(null)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load headway summaries')
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    void load()
    return () => {
      cancelled = true
    }
  }, [
    filters.date,
    filters.grain,
    filters.method,
    filters.route,
    filters.direction,
    filters.stop,
    filters.sort,
    filters.limit,
    filters.offset,
  ])

  return { data, loading, error }
}
