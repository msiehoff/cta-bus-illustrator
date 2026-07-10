import { useEffect, useState } from 'react'
import type { PipelineStatusResponse } from '../types/api'
import { adminFetch, ADMIN_API_BASE } from '../lib/adminApi'

export const usePipelineStatus = (pollMs = 5000) => {
  const [data, setData] = useState<PipelineStatusResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      try {
        const res = await adminFetch(`${ADMIN_API_BASE}/pipeline/status`)
        if (!res.ok) throw new Error('Failed to load pipeline status')
        const json = await res.json()
        if (!cancelled) {
          setData(json)
          setError(null)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load pipeline status')
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    const id = window.setInterval(load, pollMs)
    return () => {
      cancelled = true
      window.clearInterval(id)
    }
  }, [pollMs])

  return { data, loading, error }
}
