import { useEffect, useState } from 'react'

interface UseRidershipMonthsResult {
  months: string[]
  loading: boolean
  error: string | null
}

export const useRidershipMonths = (): UseRidershipMonthsResult => {
  const [months, setMonths] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    fetch('/api/v1/ridership/months')
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json()
      })
      .then(data => {
        if (!cancelled) setMonths(data.months ?? [])
      })
      .catch(err => {
        if (!cancelled) setError(err.message)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => { cancelled = true }
  }, [])

  return { months, loading, error }
}
