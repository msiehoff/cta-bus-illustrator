import { useEffect, useState } from 'react'
import type { GetRidershipResponse, RidershipDataPoint } from '../types/api'

interface UseRouteRidershipResult {
  records: RidershipDataPoint[]
  loading: boolean
  error: string | null
}

export function useRouteRidership(externalId: string): UseRouteRidershipResult {
  const [records, setRecords] = useState<RidershipDataPoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!externalId) return
    let cancelled = false
    setLoading(true)
    setError(null)

    fetch(`/api/v1/ridership/routes/${encodeURIComponent(externalId)}`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json() as Promise<GetRidershipResponse>
      })
      .then(data => {
        if (!cancelled) setRecords(data.records)
      })
      .catch(err => {
        if (!cancelled) setError(err.message)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => { cancelled = true }
  }, [externalId])

  return { records, loading, error }
}
