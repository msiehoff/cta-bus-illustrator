import { useState } from 'react'
import { twMerge } from 'tailwind-merge'
import { runHeadwayJob, useHeadwayJobRuns } from '../../hooks/useHeadwayJobs'
import type { HeadwayJobRun } from '../../types/api'

const formatTime = (value?: string) => {
  if (!value) return '—'
  return new Date(value).toLocaleString()
}

const yesterdayChicago = () => {
  const fmt = new Intl.DateTimeFormat('en-CA', {
    timeZone: 'America/Chicago',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
  const parts = fmt.formatToParts(new Date(Date.now() - 24 * 60 * 60 * 1000))
  const y = parts.find(p => p.type === 'year')?.value
  const m = parts.find(p => p.type === 'month')?.value
  const d = parts.find(p => p.type === 'day')?.value
  return `${y}-${m}-${d}`
}

const statusClass = (status: string) => {
  switch (status) {
    case 'success':
      return 'text-emerald-400'
    case 'failed':
      return 'text-red-400'
    case 'running':
      return 'text-amber-400'
    default:
      return 'text-gray-400'
  }
}

const AdminHeadwayJobs = () => {
  const { data, loading, error, reload } = useHeadwayJobRuns()
  const [serviceDate, setServiceDate] = useState(yesterdayChicago)
  const [running, setRunning] = useState(false)
  const [runError, setRunError] = useState<string | null>(null)
  const [lastRun, setLastRun] = useState<HeadwayJobRun | null>(null)

  const handleRun = async () => {
    setRunning(true)
    setRunError(null)
    try {
      const result = await runHeadwayJob(serviceDate)
      setLastRun(result)
      await reload()
    } catch (err) {
      setRunError(err instanceof Error ? err.message : 'Run failed')
    } finally {
      setRunning(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold">Headway Jobs</h2>
        <p className="text-sm text-gray-400 mt-1">
          Compute observed headways for a Chicago service date. Idempotent per day.
        </p>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg p-4 space-y-4">
        <h3 className="text-sm font-medium text-gray-300">Run rollup</h3>
        <div className="flex flex-wrap items-end gap-3">
          <label className="text-sm text-gray-400">
            Service date
            <input
              type="date"
              value={serviceDate}
              onChange={e => setServiceDate(e.target.value)}
              className="mt-1 block rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white [color-scheme:dark]"
            />
          </label>
          <button
            type="button"
            onClick={handleRun}
            disabled={running || !serviceDate}
            className={twMerge(
              'rounded-md px-4 py-2 text-sm font-medium',
              running
                ? 'bg-gray-800 text-gray-500'
                : 'bg-red-900/60 text-red-100 hover:bg-red-900',
            )}
          >
            {running ? 'Running…' : 'Run headway job'}
          </button>
        </div>
        {runError && <p className="text-sm text-red-400">{runError}</p>}
        {lastRun && (
          <p className="text-sm text-gray-400">
            Last trigger: <span className={statusClass(lastRun.status)}>{lastRun.status}</span>
            {' · '}
            {lastRun.arrivalsProcessed} arrivals → {lastRun.headwaysWritten} headways
          </p>
        )}
      </div>

      <div>
        <h3 className="text-sm font-medium text-gray-300 mb-3">Recent runs</h3>
        {loading && !data && <p className="text-gray-400">Loading…</p>}
        {error && <p className="text-red-400">{error}</p>}
        {data && data.runs.length === 0 && (
          <p className="text-sm text-gray-500">No headway jobs yet.</p>
        )}
        {data && data.runs.length > 0 && (
          <div className="overflow-x-auto rounded-lg border border-gray-800">
            <table className="min-w-full text-sm">
              <thead className="bg-gray-900 text-left text-gray-400">
                <tr>
                  <th className="px-3 py-2 font-medium">Service date</th>
                  <th className="px-3 py-2 font-medium">Status</th>
                  <th className="px-3 py-2 font-medium">Trigger</th>
                  <th className="px-3 py-2 font-medium">Arrivals</th>
                  <th className="px-3 py-2 font-medium">Headways</th>
                  <th className="px-3 py-2 font-medium">Started</th>
                  <th className="px-3 py-2 font-medium">Finished</th>
                </tr>
              </thead>
              <tbody>
                {data.runs.map(run => (
                  <tr key={run.id} className="border-t border-gray-800">
                    <td className="px-3 py-2">{run.serviceDate}</td>
                    <td className={twMerge('px-3 py-2', statusClass(run.status))}>
                      {run.status}
                      {run.errorMessage ? (
                        <span className="block text-xs text-red-400/80 max-w-xs truncate" title={run.errorMessage}>
                          {run.errorMessage}
                        </span>
                      ) : null}
                    </td>
                    <td className="px-3 py-2 text-gray-400">{run.triggeredBy}</td>
                    <td className="px-3 py-2">{run.arrivalsProcessed}</td>
                    <td className="px-3 py-2">{run.headwaysWritten}</td>
                    <td className="px-3 py-2 text-gray-400">{formatTime(run.startedAt)}</td>
                    <td className="px-3 py-2 text-gray-400">{formatTime(run.finishedAt)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default AdminHeadwayJobs
