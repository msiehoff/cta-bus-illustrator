import StatCard from '../../components/StatCard'
import { usePipelineStatus } from '../../hooks/usePipelineStatus'

const formatTime = (value?: string) => {
  if (!value) return '—'
  return new Date(value).toLocaleString()
}

const statusLabel = (running: boolean, enabled: boolean) => {
  if (!enabled) return 'Disabled'
  if (running) return 'Running'
  return 'Stopped'
}

const AdminDashboard = () => {
  const { data, loading, error } = usePipelineStatus()

  if (loading && !data) {
    return <p className="text-gray-400">Loading pipeline status…</p>
  }

  if (error) {
    return <p className="text-red-400">{error}</p>
  }

  if (!data) return null

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold">Pipeline Dashboard</h2>
        <p className="text-sm text-gray-400 mt-1">
          Live status refreshes every 5 seconds.
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label="Pipeline"
          value={statusLabel(data.running, data.enabled)}
          trend={data.enabled ? `Poll every ${data.pollInterval}` : 'Set PIPELINE_ENABLED=true'}
          trendUp={data.running ? true : undefined}
        />
        <StatCard
          label="Routes monitored"
          value={String(data.routeCount)}
        />
        <StatCard
          label="Last poll pings"
          value={String(data.lastPingCount)}
          trend={data.lastPollAt ? `Last poll ${formatTime(data.lastPollAt)}` : 'No polls yet'}
        />
        <StatCard
          label="Arrivals recorded"
          value={String(data.arrivalCount)}
          trend={data.startedAt ? `Started ${formatTime(data.startedAt)}` : undefined}
        />
      </div>

      {data.lastError && (
        <div className="rounded-lg border border-red-900/60 bg-red-950/20 px-4 py-3">
          <p className="text-sm font-medium text-red-400">Last poll error</p>
          <p className="text-sm text-red-300 mt-1">{data.lastError}</p>
        </div>
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
        <h3 className="text-sm font-medium text-gray-300 mb-3">Monitored routes</h3>
        {data.routes.length === 0 ? (
          <p className="text-sm text-gray-500">No routes configured.</p>
        ) : (
          <div className="flex flex-wrap gap-2">
            {data.routes.map(routeId => (
              <span
                key={routeId}
                className="text-xs rounded-full bg-gray-950 border border-gray-800 px-2 py-1 text-gray-300"
              >
                {routeId}
              </span>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default AdminDashboard
