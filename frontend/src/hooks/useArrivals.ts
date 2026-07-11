import { useEffect, useState } from 'react'
import type { ListArrivalsResponse } from '../types/api'
import { adminFetch, ADMIN_API_BASE } from '../lib/adminApi'

export interface ArrivalFilters {
  route?: string
  direction?: string
  stop?: string
  vehicle?: string
  date?: string
  sort?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export const useArrivals = (filters: ArrivalFilters) => {
  const [data, setData] = useState<ListArrivalsResponse | null>(null)
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
        if (filters.sort) params.set('sort', filters.sort)
        if (filters.limit != null) params.set('limit', String(filters.limit))
        if (filters.offset != null) params.set('offset', String(filters.offset))

        const res = await adminFetch(`${ADMIN_API_BASE}/arrivals?${params}`)
        if (!res.ok) throw new Error('Failed to load arrivals')
        const json = await res.json()
        if (!cancelled) {
          setData(json)
          setError(null)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load arrivals')
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    const id = window.setInterval(load, 10000)
    return () => {
      cancelled = true
      window.clearInterval(id)
    }
  }, [
    filters.route,
    filters.direction,
    filters.stop,
    filters.vehicle,
    filters.date,
    filters.sort,
    filters.limit,
    filters.offset,
  ])

  return { data, loading, error }
}
