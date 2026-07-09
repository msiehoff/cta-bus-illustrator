import { useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useCorridorRidership } from '../hooks/useCorridorRidership'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import { useRidershipFilters } from '../hooks/useRidershipFilters'
import RidershipChart, { type WindowKey } from '../components/RidershipChart'
import StatCard from '../components/StatCard'
import RecoveryBanner from '../components/RecoveryBanner'
import RouteRecoveryTable from '../components/RouteRecoveryTable'
import RidershipFilters from '../components/RidershipFilters'
import RouteContextPanel from '../components/RouteContextPanel'
import RidershipDistributionChart from '../components/RidershipDistributionChart'
import SeasonalityChart from '../components/SeasonalityChart'
import WeekendShareChart from '../components/WeekendShareChart'
import {
  appendCorridorRows,
  buildCorridorRouteId,
  getCorridorBadgeLabel,
  getCorridorDisplayName,
} from '../lib/corridors'
import { getPairedRouteId } from '../lib/expressPairs'
import {
  buildRouteRecovery,
  buildSeasonality,
  buildWeekendShare,
  formatMonth,
  formatRides,
  getSnapshotForMonth,
  preCovidTrend,
} from '../lib/ridershipUtils'

const CorridorRoutePage = () => {
  const { localId = '' } = useParams<{ localId: string }>()
  const navigate = useNavigate()
  const expressId = getPairedRouteId(localId) ?? ''
  const {
    records,
    loading,
    error,
    localRecords,
    expressRecords,
  } = useCorridorRidership(localId)
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

  const corridorRouteId = buildCorridorRouteId(localId)
  const comparisonWithCorridors = useMemo(
    () => appendCorridorRows(comparison?.routes ?? []),
    [comparison?.routes],
  )
  const corridorComparison = comparisonWithCorridors.find(r => r.routeId === corridorRouteId)
  const localComparison = comparison?.routes.find(r => r.routeId === localId)
  const expressComparison = comparison?.routes.find(r => r.routeId === expressId)

  const corridorName = corridorComparison?.routeName
    ?? (localComparison ? getCorridorDisplayName(localComparison.routeName) : `${localId} Corridor`)

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

  const localLatest = selectedMonth ? getSnapshotForMonth(localRecords, selectedMonth) : null
  const expressLatest = selectedMonth ? getSnapshotForMonth(expressRecords, selectedMonth) : null

  if (!expressId) {
    return (
      <div className="p-4 sm:p-6 max-w-4xl mx-auto text-gray-400 text-sm">
        Corridor not found.
      </div>
    )
  }

  return (
    <div className="p-4 sm:p-6 max-w-4xl mx-auto">
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
        <div className="flex items-center justify-center min-w-12 h-12 px-2 rounded-full bg-blue-900 text-blue-200 text-xs font-bold shrink-0">
          {getCorridorBadgeLabel(localId)}
        </div>
        <div>
          <h1 className="text-xl font-semibold text-white">{corridorName}</h1>
          <p className="text-sm text-gray-400 mt-0.5">Express + Local</p>
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

      {latest && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-5">
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

      {localComparison && expressComparison && localLatest && expressLatest && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3 mb-5">
          <p className="text-xs text-gray-500 uppercase tracking-wide mb-2">Breakdown</p>
          <div className="flex flex-wrap items-center gap-4 text-sm">
            <div>
              <Link
                to={`/routes/${localId}`}
                state={{ routeName: localComparison.routeName }}
                className="text-gray-400 hover:text-gray-200 transition-colors"
              >
                {localComparison.routeName}
              </Link>
              <span className="text-white font-medium ml-2">
                {formatRides(localLatest[ridershipType])}
              </span>
            </div>
            <span className="text-gray-600">+</span>
            <div>
              <Link
                to={`/routes/${expressId}`}
                state={{ routeName: expressComparison.routeName }}
                className="text-gray-400 hover:text-gray-200 transition-colors"
              >
                {expressComparison.routeName}
              </Link>
              <span className="text-white font-medium ml-2">
                {formatRides(expressLatest[ridershipType])}
              </span>
            </div>
            <span className="text-gray-600">=</span>
            <div>
              <span className="text-gray-400">Combined</span>
              <span className="text-white font-semibold ml-2">
                {formatRides(latest?.[ridershipType] ?? 0)}
              </span>
            </div>
          </div>
        </div>
      )}

      {recovery && weekdayRecovery && (
        <RecoveryBanner
          title="Corridor ridership vs pre-pandemic"
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

      {comparison && selectedMonth && (
        <RouteContextPanel
          routeId={corridorRouteId}
          routeName={corridorName}
          comparison={comparisonWithCorridors}
          systemTotal={comparison.systemCurrent}
          ridershipType={ridershipType}
          selectedMonth={selectedMonth}
        />
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
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
              Combined local + express ridership. Shaded area = pandemic period (Mar 2020 – Dec 2022).
              Pre-pandemic benchmark uses the same calendar month in 2019.
            </p>
          </>
        )}
      </div>

      {comparison && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
          <h2 className="text-sm font-medium text-white mb-1">Where this corridor falls</h2>
          <p className="text-xs text-gray-500 mb-2">
            Network ridership distribution · {formatMonth(comparison.currentMonth)} · {ridershipType}
          </p>
          <RidershipDistributionChart
            routes={comparisonWithCorridors}
            highlightRouteId={corridorRouteId}
            highlightRouteName={corridorName}
            ridershipType={ridershipType}
          />
        </div>
      )}

      <div className="grid lg:grid-cols-2 gap-3 mb-5">
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Seasonality</h2>
          <p className="text-xs text-gray-500 mb-4">Avg {ridershipType} ridership by calendar month</p>
          <SeasonalityChart data={seasonality} />
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Weekend share</h2>
          <p className="text-xs text-gray-500 mb-4">(Sat + Sun) / weekday ridership over time</p>
          <WeekendShareChart data={weekendShare} />
        </div>
      </div>
    </div>
  )
}

export default CorridorRoutePage
