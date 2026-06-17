import { useEffect, useState } from 'react'
import type { GetRoutesComparisonResponse } from '../types/api'

interface UseRoutesComparisonResult {
  data: GetRoutesComparisonResponse | null
  loading: boolean
  error: string | null
}

export const useRoutesComparison = (): UseRoutesComparisonResult => {
  const [data, setData] = useState<GetRoutesComparisonResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/v1/routes/comparison')
      .then(r => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`)
        return r.json()
      })
      .then(setData)
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  return { data, loading, error }
}
