import { useState } from 'react'
import { useRoutesComparison } from '../hooks/useRoutesComparison'
import RoutesComparisonTable from '../components/RoutesComparisonTable'
import RecoveryBanner from '../components/RecoveryBanner'
import { formatMonth } from '../lib/ridershipUtils'

const RoutesPage = () => {
  const { data, loading } = useRoutesComparison()
  const [search, setSearch] = useState('')

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">Routes</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          CTA bus routes · weekday ridership
          {data ? ` · ${formatMonth(data.currentMonth)}` : ''}
        </p>
      </div>

      {data && (
        <RecoveryBanner
          title="System recovery"
          currentMonth={data.currentMonth}
          benchmarkMonth={data.benchmarkMonth}
          current={data.systemCurrent}
          preCovid={data.systemPreCovid}
          recovery={data.systemRecovery}
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
