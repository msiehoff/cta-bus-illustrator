import { Link } from 'react-router-dom'
import { useHeadwaySystem } from '../../hooks/useHeadwaySystem'
import StatCard from '../../components/StatCard'
import HeadwayDailyChart from '../../components/HeadwayDailyChart'
import {
  formatHeadwayCV,
  formatHeadwayMinutes,
  formatHeadwayPeriod,
} from '../../lib/headwayUtils'

const HeadwaySystemPage = () => {
  const { data, loading, error } = useHeadwaySystem(30)
  const period = data?.period
  const periodLabel = period
    ? formatHeadwayPeriod(period.periodStart, period.periodEnd, period.daysWithData)
    : null

  return (
    <div className="p-4 sm:p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">Headway · System</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          Network service frequency
          {periodLabel
            ? ` · up to last 30 available days (${periodLabel})`
            : ' · up to last 30 available service days'}
        </p>
      </div>

      {error && (
        <div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 rounded-lg px-4 py-3">
          {error}
        </div>
      )}

      {loading && (
        <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>
      )}

      {!loading && period && period.daysWithData > 0 && (
        <div className="mb-5">
          <p className="text-xs text-gray-500 mb-2">
            Across {period.daysWithData} available service day
            {period.daysWithData === 1 ? '' : 's'}
            {periodLabel ? ` · ${periodLabel}` : ''}
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <StatCard
              label="Network median headway"
              value={`${formatHeadwayMinutes(period.medianMinutes)} min`}
            />
            <StatCard
              label="Network avg wait"
              value={`${formatHeadwayMinutes(period.avgWaitMinutes)} min`}
            />
            <StatCard
              label="Network CV"
              value={formatHeadwayCV(period.cv)}
            />
          </div>
        </div>
      )}

      {!loading && (!period || period.daysWithData === 0) && !error && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-10 text-center text-sm text-gray-500 mb-5">
          No system headway summaries yet. Run a daily headway job first.
        </div>
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-1">Network daily trend</h2>
        <p className="text-xs text-gray-500 mb-4">
          Equal-stop · all routes · up to last 30 service days
        </p>
        {!loading && <HeadwayDailyChart series={data?.series ?? []} />}
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-1">Shortest median headways</h2>
        <p className="text-xs text-gray-500 mb-4">
          Lowest period medians
          {period?.daysWithData
            ? ` · across ${period.daysWithData} available service day${period.daysWithData === 1 ? '' : 's'}`
            : ''}
          {periodLabel ? ` · ${periodLabel}` : ''}
        </p>
        {!loading && !(data?.shortestHeadways?.length) && (
          <p className="text-sm text-gray-500 py-4 text-center">No routes to rank yet.</p>
        )}
        {!loading && !!data?.shortestHeadways?.length && (
          <ul className="divide-y divide-gray-800">
            {data.shortestHeadways.map(route => (
              <li key={route.routeId}>
                <Link
                  to={`/headways/routes/${encodeURIComponent(route.routeId)}`}
                  className="flex items-center justify-between gap-3 py-2.5 hover:bg-gray-800/40 -mx-1 px-1 rounded transition-colors"
                >
                  <div>
                    <span className="text-white font-medium">{route.routeId}</span>
                    {route.routeName && (
                      <span className="text-gray-500 ml-2 text-sm">{route.routeName}</span>
                    )}
                  </div>
                  <span className="text-sm text-gray-300 tabular-nums">
                    {formatHeadwayMinutes(route.medianMinutes)} min
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        )}
        <p className="mt-3">
          <Link
            to="/headways/routes"
            className="text-xs text-gray-500 hover:text-red-400 transition-colors"
          >
            View all routes →
          </Link>
        </p>
      </div>
    </div>
  )
}

export default HeadwaySystemPage
