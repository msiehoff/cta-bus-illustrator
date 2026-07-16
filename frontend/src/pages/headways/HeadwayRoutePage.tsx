import { Link, useNavigate, useParams } from 'react-router-dom'
import { useHeadwayRoute } from '../../hooks/useHeadwayRoute'
import { useHeadwayRoutes } from '../../hooks/useHeadwayRoutes'
import { useHeadwaySystem } from '../../hooks/useHeadwaySystem'
import { useRouteName } from '../../hooks/useRouteName'
import StatCard from '../../components/StatCard'
import HeadwayDailyChart from '../../components/HeadwayDailyChart'
import HeadwayContextPanel from '../../components/HeadwayContextPanel'
import HeadwayDistributionChart from '../../components/HeadwayDistributionChart'
import {
  formatHeadwayCV,
  formatHeadwayMinutes,
  formatHeadwayPeriod,
  describeHeadwayConsistency,
  HEADWAY_CONSISTENCY_LABEL,
  HEADWAY_CONSISTENCY_TOOLTIP,
} from '../../lib/headwayUtils'

const PERIOD_DAYS = 30

const HeadwayRoutePage = () => {
  const { externalId = '' } = useParams<{ externalId: string }>()
  const navigate = useNavigate()
  const { data, loading, error } = useHeadwayRoute(externalId, PERIOD_DAYS)
  const { data: routesData } = useHeadwayRoutes(PERIOD_DAYS)
  const { data: systemData } = useHeadwaySystem(PERIOD_DAYS)
  const fallbackName = useRouteName(externalId)

  const route = data?.route
  const routeName = route?.routeName || fallbackName
  const periodLabel = route
    ? formatHeadwayPeriod(route.periodStart, route.periodEnd, route.daysWithData)
    : null
  const networkPeriod = systemData?.period
  const allRoutes = routesData?.routes ?? []
  const consistency = route ? describeHeadwayConsistency(route.cv) : null

  return (
    <div className="p-4 sm:p-6 max-w-4xl mx-auto">
      <button
        type="button"
        onClick={() => navigate(-1)}
        className="flex items-center gap-1.5 text-xs text-gray-500 hover:text-gray-300 mb-4 transition-colors"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <line x1="19" y1="12" x2="5" y2="12" />
          <polyline points="12 19 5 12 12 5" />
        </svg>
        Back
      </button>

      <div className="flex items-start justify-between gap-4 mb-5">
        <div className="flex items-center gap-4">
          <div className="flex items-center justify-center w-12 h-12 rounded-full bg-blue-900 text-blue-200 text-lg font-bold shrink-0">
            {externalId}
          </div>
          <div>
            <h1 className="text-xl font-semibold text-white">{routeName}</h1>
            <p className="text-sm text-gray-400 mt-0.5">
              Headways · both directions
              {periodLabel ? ` · ${periodLabel}` : ''}
            </p>
          </div>
        </div>
        <Link
          to={`/routes/${encodeURIComponent(externalId)}`}
          className="text-xs text-gray-500 hover:text-red-400 transition-colors shrink-0 mt-1"
        >
          Ridership for this route →
        </Link>
      </div>

      {error && (
        <div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 rounded-lg px-4 py-3">
          {error}
        </div>
      )}

      {loading && (
        <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>
      )}

      {!loading && route && route.daysWithData > 0 && (
        <div className="mb-5">
          <p className="text-xs text-gray-500 mb-2">
            Across {route.daysWithData} available service day
            {route.daysWithData === 1 ? '' : 's'}
            {periodLabel ? ` · ${periodLabel}` : ''}
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <StatCard
              label="Median headway"
              value={`${formatHeadwayMinutes(route.medianMinutes)} min`}
            />
            <StatCard
              label="Avg wait"
              value={`${formatHeadwayMinutes(route.avgWaitMinutes)} min`}
            />
            <StatCard
              label={HEADWAY_CONSISTENCY_LABEL}
              labelHint={HEADWAY_CONSISTENCY_TOOLTIP}
              value={formatHeadwayCV(route.cv)}
              trend={consistency?.label}
              trendUp={consistency?.trendUp}
            />
          </div>
        </div>
      )}

      {!loading && (!route || route.daysWithData === 0) && !error && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-10 text-center text-sm text-gray-500 mb-5">
          No headway summaries for this route yet.
        </div>
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-1">Daily series</h2>
        <p className="text-xs text-gray-500 mb-4">
          Equal-stop · both directions · up to last 30 service days with data
        </p>
        {loading ? (
          <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>
        ) : (
          <HeadwayDailyChart series={data?.series ?? []} />
        )}
      </div>

      {!loading && route && route.daysWithData > 0 && networkPeriod && networkPeriod.daysWithData > 0 && (
        <HeadwayContextPanel
          routeId={externalId}
          routes={allRoutes}
          routeMedian={route.medianMinutes}
          routeWait={route.avgWaitMinutes}
          networkMedian={networkPeriod.medianMinutes}
          networkWait={networkPeriod.avgWaitMinutes}
          periodLabel={periodLabel}
        />
      )}

      {allRoutes.length > 0 && route && route.daysWithData > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
          <h2 className="text-sm font-medium text-white mb-1">Where this route falls</h2>
          <p className="text-xs text-gray-500 mb-2">
            Network median headway distribution
            {periodLabel ? ` · ${periodLabel}` : ''}
          </p>
          <HeadwayDistributionChart
            routes={allRoutes}
            highlightRouteId={externalId}
            highlightRouteName={routeName}
          />
        </div>
      )}
    </div>
  )
}

export default HeadwayRoutePage
