import { useState } from 'react'
import { useHeadwayRoutes } from '../../hooks/useHeadwayRoutes'
import HeadwayRoutesTable from '../../components/HeadwayRoutesTable'
import { formatHeadwayPeriod } from '../../lib/headwayUtils'

const HeadwayRoutesPage = () => {
  const { data, loading, error } = useHeadwayRoutes(30)
  const [search, setSearch] = useState('')

  const periodLabel = data
    ? formatHeadwayPeriod(
        data.period.periodStart,
        data.period.periodEnd,
        data.period.daysWithData,
      )
    : null

  return (
    <div className="p-4 sm:p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">Headway · Routes</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          Both directions · equal-stop weighting
          {periodLabel ? ` · ${periodLabel}` : ''}
          {data?.period.daysWithData
            ? ` · ${data.period.daysWithData} service day${data.period.daysWithData === 1 ? '' : 's'}`
            : ''}
        </p>
      </div>

      {error && (
        <div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 rounded-lg px-4 py-3">
          {error}
        </div>
      )}

      <HeadwayRoutesTable
        routes={data?.routes ?? []}
        loading={loading}
        search={search}
        onSearchChange={setSearch}
      />

      <p className="text-[10px] text-gray-600 mt-3">
        Period uses up to the last 30 available service days with headway summaries.
        Median, wait, and consistency are observation-weighted across those days.
      </p>
    </div>
  )
}

export default HeadwayRoutesPage
