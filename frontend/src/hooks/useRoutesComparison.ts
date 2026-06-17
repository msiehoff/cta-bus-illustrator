import { useEffect, useState } from 'react'
import type { GetRoutesComparisonResponse, RidershipType } from '../types/api'

interface UseRoutesComparisonResult {
  data: GetRoutesComparisonResponse | null
  loading: boolean
  error: string | null
}

export const useRoutesComparison = (
  ridershipType: RidershipType = 'weekday',
  month: string | null = null,
): UseRoutesComparisonResult => {
  const [data, setData] = useState<GetRoutesComparisonResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!month) {
      setLoading(true)
      return
    }

    let cancelled = false
    setLoading(true)
    setError(null)

    const params = new URLSearchParams({ type: ridershipType, month })
    fetch(`/api/v1/routes/comparison?${params}`)
      .then(r => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`)
        return r.json()
      })
      .then(next => {
        if (!cancelled) setData(next)
      })
      .catch(err => {
        if (!cancelled) setError(err.message)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => { cancelled = true }
  }, [ridershipType, month])

  return { data, loading, error }
}
