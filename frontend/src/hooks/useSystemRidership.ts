import { useEffect, useState } from 'react'
import type { GetRidershipResponse, RidershipDataPoint } from '../types/api'

interface UseSystemRidershipResult {
  records: RidershipDataPoint[]
  loading: boolean
  error: string | null
}

export function useSystemRidership(): UseSystemRidershipResult {
  const [records, setRecords] = useState<RidershipDataPoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)

    fetch('/api/v1/ridership/system')
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
  }, [])

  return { records, loading, error }
}
