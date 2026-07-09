import { useMemo, useState } from 'react'
import { useSystemRidership } from '../hooks/useSystemRidership'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import { useRidershipFilters } from '../hooks/useRidershipFilters'
import RidershipChart, { type WindowKey, cutoffMonth } from '../components/RidershipChart'
import StatCard from '../components/StatCard'
import RoutesComparisonTable from '../components/RoutesComparisonTable'
import RecoveryBanner from '../components/RecoveryBanner'
import RidershipFilters from '../components/RidershipFilters'
import RecoveryDistributionChart from '../components/RecoveryDistributionChart'
import RecoveryScatterChart from '../components/RecoveryScatterChart'
import SeasonalityChart from '../components/SeasonalityChart'
import TopMoversPanel from '../components/TopMoversPanel'
import RidershipDistributionChart from '../components/RidershipDistributionChart'
import {
  buildRecoveryDistribution,
  buildRecoveryScatter,
  buildSeasonality,
  formatMonth,
  formatRides,
  getSnapshotForMonth,
  preCovidTrend,
} from '../lib/ridershipUtils'

const SystemOverview = () => {
  const {
    months,
    monthsLoading,
    selectedMonth,
    setSelectedMonth,
    ridershipType,
    setRidershipType,
  } = useRidershipFilters()
  const { records, loading, error } = useSystemRidership()
  const { data: comparison, loading: comparisonLoading } = useRoutesComparison(ridershipType, selectedMonth)
  const [search, setSearch] = useState('')
  const [window, setWindow] = useState<WindowKey>('5y')

  const latest = useMemo(
    () => (selectedMonth ? getSnapshotForMonth(records, selectedMonth) : null),
    [records, selectedMonth],
  )

  const windowStart = useMemo(() => {
    if (!records.length || !selectedMonth) return null
    const cutoff = cutoffMonth(window, selectedMonth)
    const windowRecords = cutoff ? records.filter(r => r.month >= cutoff) : records
    if (!windowRecords.length) return null
    const firstMonth = windowRecords[0].month
    return getSnapshotForMonth(records, firstMonth)
  }, [records, window, selectedMonth])

  const sinceLabel = windowStart ? `since ${formatMonth(windowStart.month)}` : null

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

  const distribution = useMemo(
    () => buildRecoveryDistribution(comparison?.routes ?? []),
    [comparison?.routes],
  )

  const scatter = useMemo(
    () => buildRecoveryScatter(comparison?.routes ?? []),
    [comparison?.routes],
  )

  const seasonality = useMemo(
    () => buildSeasonality(records, ridershipType),
    [records, ridershipType],
  )

  return (
    <div className="p-4 sm:p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">System Overview</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          CTA Bus Network · All Routes{latest ? ` · ${formatMonth(latest.month)}` : ''}
        </p>
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
          <StatCard label={`Total weekday riders · ${formatMonth(latest.month)}`} {...statProps('weekday')} />
          <StatCard label={`Total Saturday riders · ${formatMonth(latest.month)}`} {...statProps('saturday')} />
          <StatCard label={`Total Sunday riders · ${formatMonth(latest.month)}`} {...statProps('sunday')} />
        </div>
      )}

      {comparison && (
        <RecoveryBanner
          currentMonth={comparison.currentMonth}
          benchmarkMonth={comparison.benchmarkMonth}
          current={comparison.systemCurrent}
          preCovid={comparison.systemPreCovid}
          recovery={comparison.systemRecovery}
          dayTypeLabel={ridershipType}
        />
      )}

      {comparison && <TopMoversPanel routes={comparison.routes} />}

      {comparison && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
          <h2 className="text-sm font-medium text-white mb-1">Route ridership distribution</h2>
          <p className="text-xs text-gray-500 mb-2">
            How avg daily {ridershipType} ridership is spread across the network · {formatMonth(comparison.currentMonth)}
            {' · '}y-axis = routes per ridership band
          </p>
          <RidershipDistributionChart
            routes={comparison.routes}
            ridershipType={ridershipType}
          />
        </div>
      )}

      <div className="grid lg:grid-cols-2 gap-3 mb-5">
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Pre-pandemic recovery distribution</h2>
          <p className="text-xs text-gray-500 mb-4">Routes grouped by ridership vs same month in 2019</p>
          <RecoveryDistributionChart data={distribution} />
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4">
          <h2 className="text-sm font-medium text-white mb-1">Ridership vs pre-pandemic</h2>
          <p className="text-xs text-gray-500 mb-4">Each dot is a route · click to open detail</p>
          <RecoveryScatterChart data={scatter} />
        </div>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-1">Seasonality</h2>
        <p className="text-xs text-gray-500 mb-4">
          Average {ridershipType} ridership by calendar month across all years
        </p>
        <SeasonalityChart data={seasonality} />
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 sm:px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">System ridership over time</h2>
        {loading && <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>}
        {error && <div className="text-red-400 text-sm py-10 text-center">Failed to load: {error}</div>}
        {!loading && !error && (
          <RidershipChart
            records={records}
            window={window}
            onWindowChange={setWindow}
            height={220}
            highlightType={ridershipType}
          />
        )}
        {sinceLabel && (
          <p className="text-[10px] text-gray-600 mt-2">
            Chart window: {sinceLabel}. Pre-pandemic benchmark = same calendar month in 2019.
          </p>
        )}
      </div>

      <div>
        <h2 className="text-sm font-medium text-white mb-3">Routes</h2>
        <RoutesComparisonTable
          routes={comparison?.routes ?? []}
          loading={comparisonLoading}
          showHeader
          search={search}
          onSearchChange={setSearch}
        />
      </div>
    </div>
  )
}

export default SystemOverview
