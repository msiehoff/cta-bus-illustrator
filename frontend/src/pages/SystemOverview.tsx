import { useMemo, useState } from 'react'
import { useSystemRidership } from '../hooks/useSystemRidership'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import RidershipChart, { type WindowKey, cutoffMonth } from '../components/RidershipChart'
import StatCard from '../components/StatCard'
import RoutesComparisonTable from '../components/RoutesComparisonTable'
import RecoveryBanner from '../components/RecoveryBanner'
import type { RidershipType } from '../types/api'
import {
  formatMonth,
  formatRides,
  preCovidTrend,
} from '../lib/ridershipUtils'

const SystemOverview = () => {
  const { records, loading, error } = useSystemRidership()
  const { data: comparison, loading: comparisonLoading } = useRoutesComparison()
  const [search, setSearch] = useState('')
  const [window, setWindow] = useState<WindowKey>('5y')

  const latestMonth = useMemo(() => {
    if (!records.length) return null
    return records[records.length - 1].month
  }, [records])

  const latest = useMemo(() => {
    if (!latestMonth) return null
    const rows = records.filter(r => r.month === latestMonth)
    const get = (type: RidershipType) => rows.find(r => r.type === type)?.avgRides ?? 0
    return {
      month: latestMonth,
      weekday:  get('weekday'),
      saturday: get('saturday'),
      sunday:   get('sunday'),
    }
  }, [records, latestMonth])

  const windowStart = useMemo(() => {
    if (!records.length || !latestMonth) return null
    const cutoff = cutoffMonth(window, latestMonth)
    const windowRecords = cutoff ? records.filter(r => r.month >= cutoff) : records
    if (!windowRecords.length) return null
    const firstMonth = windowRecords[0].month
    const rows = windowRecords.filter(r => r.month === firstMonth)
    const get = (type: RidershipType) => rows.find(r => r.type === type)?.avgRides ?? 0
    return {
      month: firstMonth,
      weekday:  get('weekday'),
      saturday: get('saturday'),
      sunday:   get('sunday'),
    }
  }, [records, window, latestMonth])

  const sinceLabel = windowStart ? `since ${formatMonth(windowStart.month)}` : null

  const statProps = (type: 'weekday' | 'saturday' | 'sunday') => {
    const current = latest?.[type] ?? 0
    const preCovid = latestMonth
      ? preCovidTrend(records, latestMonth, type)
      : {}
    return {
      value: formatRides(current),
      ...preCovid,
    }
  }

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">System Overview</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          CTA Bus Network · All Routes{latest ? ` · ${formatMonth(latest.month)}` : ''}
        </p>
      </div>

      {latest && (
        <div className="grid grid-cols-3 gap-3 mb-5">
          <StatCard label={`Total weekday riders · ${formatMonth(latest.month)}`} {...statProps('weekday')} />
          <StatCard label={`Total Saturday riders · ${formatMonth(latest.month)}`} {...statProps('saturday')} />
          <StatCard label={`Total Sunday riders · ${formatMonth(latest.month)}`} {...statProps('sunday')} />
        </div>
      )}

      {comparison && (
        <RecoveryBanner
          title="System recovery"
          currentMonth={comparison.currentMonth}
          benchmarkMonth={comparison.benchmarkMonth}
          current={comparison.systemCurrent}
          preCovid={comparison.systemPreCovid}
          recovery={comparison.systemRecovery}
        />
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">System ridership over time</h2>
        {loading && <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>}
        {error && <div className="text-red-400 text-sm py-10 text-center">Failed to load: {error}</div>}
        {!loading && !error && (
          <RidershipChart
            records={records}
            window={window}
            onWindowChange={setWindow}
            height={220}
          />
        )}
        {sinceLabel && (
          <p className="text-[10px] text-gray-600 mt-2">
            Chart window: {sinceLabel}. Pre-COVID benchmark uses same calendar month in 2019.
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
