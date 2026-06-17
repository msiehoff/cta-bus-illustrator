import { useMemo, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useRouteRidership } from '../hooks/useRouteRidership'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import { useRouteName } from '../hooks/useRouteName'
import { useRidershipFilters } from '../hooks/useRidershipFilters'
import RidershipChart, { type WindowKey } from '../components/RidershipChart'
import StatCard from '../components/StatCard'
import RecoveryBanner from '../components/RecoveryBanner'
import RouteRecoveryTable from '../components/RouteRecoveryTable'
import RidershipFilters from '../components/RidershipFilters'
import RouteContextPanel from '../components/RouteContextPanel'
import SeasonalityChart from '../components/SeasonalityChart'
import WeekendShareChart from '../components/WeekendShareChart'
import {
  buildRouteRecovery,
  buildSeasonality,
  buildWeekendShare,
  formatMonth,
  formatRides,
  getSnapshotForMonth,
  preCovidTrend,
} from '../lib/ridershipUtils'

const RoutePage = () => {
  const { externalId = '' } = useParams<{ externalId: string }>()
  const navigate = useNavigate()
  const { records, loading, error } = useRouteRidership(externalId)
  const {
    months,
    monthsLoading,
    selectedMonth,
    setSelectedMonth,
    ridershipType,
    setRidershipType,
  } = useRidershipFilters()
  const { data: comparison } = useRoutesComparison(ridershipType, selectedMonth)
  const [window, setWindow] = useState<WindowKey>('5y')

  const routeName = useRouteName(externalId)

  const latest = useMemo(
    () => (selectedMonth ? getSnapshotForMonth(records, selectedMonth) : null),
    [records, selectedMonth],
  )

  const recovery = useMemo(
    () => (selectedMonth ? buildRouteRecovery(records, selectedMonth) : null),
    [records, selectedMonth],
  )

  const weekdayRecovery = recovery?.rows.find(r => r.type === ridershipType)

  const statProps = (type: 'weekday' | 'saturday' | 'sunday') => {
    const current = latest?.[type] ?? 0
    const preCovid = selectedMonth
      ? preCovidTrend(records, selectedMonth, type)
      : {}
    return {
      value: formatRides(current),
      ...preCovid,
    }
  }

  const seasonality = useMemo(
    () => buildSeasonality(records, ridershipType),
    [records, ridershipType],
  )

  const weekendShare = useMemo(
    () => buildWeekendShare(records),
    [records],
  )

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-1.5 text-xs text-gray-500 hover:text-gray-300 mb-4 transition-colors"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <line x1="19" y1="12" x2="5" y2="12" />
          <polyline points="12 19 5 12 12 5" />
        </svg>
        Back
      </button>

      <div className="flex items-center gap-4 mb-5">
        <div className="flex items-center justify-center w-12 h-12 rounded-full bg-blue-900 text-blue-200 text-lg font-bold shrink-0">
          {externalId}
        </div>
        <div>
          <h1 className="text-xl font-semibold text-white">{routeName}</h1>
          {latest && (
            <p className="text-sm text-gray-400 mt-0.5">
              {formatRides(latest[ridershipType])} avg {ridershipType} riders · {formatMonth(latest.month)}
            </p>
          )}
        </div>
      </div>

      <RidershipFilters
        months={months}
        selectedMonth={selectedMonth}
        ridershipType={ridershipType}
        monthsLoading={monthsLoading}
        onMonthChange={setSelectedMonth}
        onTypeChange={setRidershipType}
      />

      {comparison && selectedMonth && (
        <RouteContextPanel
          routeId={externalId}
          routeName={routeName}
          comparison={comparison.routes}
          systemTotal={comparison.systemCurrent}
          ridershipType={ridershipType}
          selectedMonth={selectedMonth}
        />
      )}

      {latest && (
        <div className="grid grid-cols-3 gap-3 mb-5">
          <StatCard
            label={`Weekday riders · ${formatMonth(latest.month)}`}
            {...statProps('weekday')}
          />
          <StatCard
            label={`Saturday riders · ${formatMonth(latest.month)}`}
            {...statProps('saturday')}
          />
          <StatCard
            label={`Sunday riders · ${formatMonth(latest.month)}`}
            {...statProps('sunday')}
          />
        </div>
      )}

      {recovery && weekdayRecovery && (
        <RecoveryBanner
          title="Route ridership vs pre-pandemic"
          currentMonth={recovery.currentMonth}
          benchmarkMonth={recovery.benchmarkMonth}
          current={weekdayRecovery.current}
          preCovid={weekdayRecovery.preCovid2019 ?? undefined}
          recovery={weekdayRecovery.recoveryPct ?? undefined}
          dayTypeLabel={ridershipType}
        />
      )}

      {recovery && (
        <RouteRecoveryTable
          currentMonth={recovery.currentMonth}
          benchmarkMonth={recovery.benchmarkMonth}
          yearAgoMonth={recovery.yearAgoMonth}
          fiveYearsAgoMonth={recovery.fiveYearsAgoMonth}
          rows={recovery.rows}
        />
      )}

      <div className="grid lg:grid-cols-2 gap-3 mb-5">
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Seasonality</h2>
          <p className="text-xs text-gray-500 mb-4">Avg {ridershipType} ridership by calendar month</p>
          <SeasonalityChart data={seasonality} />
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Weekend share</h2>
          <p className="text-xs text-gray-500 mb-4">(Sat + Sun) / weekday ridership over time</p>
          <WeekendShareChart data={weekendShare} />
        </div>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">Ridership over time</h2>
        {loading && <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>}
        {error && <div className="text-red-400 text-sm py-10 text-center">Failed to load: {error}</div>}
        {!loading && !error && (
          <>
            <RidershipChart
              records={records}
              window={window}
              onWindowChange={setWindow}
              height={220}
              highlightType={ridershipType}
            />
            <p className="text-[10px] text-gray-600 mt-2">
              Shaded area = pandemic period (Mar 2020 – Dec 2022).
              Pre-pandemic benchmark uses the same calendar month in 2019.
            </p>
          </>
        )}
      </div>

      <div className="bg-gray-900 border border-dashed border-gray-700 rounded-lg px-5 py-4 opacity-60 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-gray-800 flex items-center justify-center shrink-0">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#6b7280" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
        </div>
        <div>
          <span className="text-[10px] bg-gray-800 text-gray-500 px-2 py-0.5 rounded-full">Coming soon</span>
          <p className="text-sm font-medium text-gray-400 mt-1">Headway Metrics</p>
          <p className="text-xs text-gray-600 mt-0.5">Track frequency, reliability, and schedule adherence for this route</p>
        </div>
      </div>
    </div>
  )
}

export default RoutePage
