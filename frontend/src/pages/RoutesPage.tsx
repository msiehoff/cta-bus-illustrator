import { useState } from 'react'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import { useRidershipFilters } from '../hooks/useRidershipFilters'
import RoutesComparisonTable from '../components/RoutesComparisonTable'
import RecoveryBanner from '../components/RecoveryBanner'
import RidershipFilters from '../components/RidershipFilters'
import { formatMonth } from '../lib/ridershipUtils'

const RoutesPage = () => {
  const {
    months,
    monthsLoading,
    selectedMonth,
    setSelectedMonth,
    ridershipType,
    setRidershipType,
  } = useRidershipFilters()
  const { data, loading } = useRoutesComparison(ridershipType, selectedMonth)
  const [search, setSearch] = useState('')

  return (
    <div className="p-4 sm:p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">Routes</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          CTA bus routes · {ridershipType} ridership
          {data ? ` · ${formatMonth(data.currentMonth)}` : ''}
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

      {data && (
        <RecoveryBanner
          currentMonth={data.currentMonth}
          benchmarkMonth={data.benchmarkMonth}
          current={data.systemCurrent}
          preCovid={data.systemPreCovid}
          recovery={data.systemRecovery}
          dayTypeLabel={ridershipType}
        />
      )}

      <RoutesComparisonTable
        routes={data?.routes ?? []}
        loading={loading}
        search={search}
        onSearchChange={setSearch}
      />
    </div>
  )
}

export default RoutesPage
