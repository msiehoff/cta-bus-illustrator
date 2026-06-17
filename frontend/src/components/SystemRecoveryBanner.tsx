import {
  formatMonth,
  formatRecoveryPct,
  formatRides,
  recoveryBarColorClass,
  recoveryColorClass,
} from '../lib/ridershipUtils'

interface Props {
  currentMonth: string
  benchmarkMonth: string
  systemCurrent: number
  systemPreCovid?: number
  systemRecovery?: number
}

const SystemRecoveryBanner = ({
  currentMonth,
  benchmarkMonth,
  systemCurrent,
  systemPreCovid,
  systemRecovery,
}: Props) => {
  if (systemRecovery == null || systemPreCovid == null) return null

  const barWidth = Math.min(systemRecovery, 100)

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <p className="text-xs text-gray-500 uppercase tracking-wide">System recovery</p>
          <p className="text-2xl font-semibold text-white mt-0.5">
            <span className={recoveryColorClass(systemRecovery)}>
              {formatRecoveryPct(systemRecovery)}
            </span>
            <span className="text-gray-400 text-base font-normal ml-2">
              of {formatMonth(benchmarkMonth)} weekday ridership
            </span>
          </p>
          <p className="text-xs text-gray-500 mt-1">
            {formatRides(systemCurrent)} avg weekday riders · {formatMonth(currentMonth)}
            {' · '}
            {formatRides(systemPreCovid)} pre-COVID
          </p>
        </div>
        <div className="sm:w-48 shrink-0">
          <div className="h-2 bg-gray-800 rounded-full overflow-hidden">
            <div
              className={`h-full rounded-full transition-all ${recoveryBarColorClass(systemRecovery)}`}
              style={{ width: `${barWidth}%` }}
            />
          </div>
          <p className="text-[10px] text-gray-600 mt-1 text-right">100% = pre-COVID baseline</p>
        </div>
      </div>
    </div>
  )
}

export default SystemRecoveryBanner
