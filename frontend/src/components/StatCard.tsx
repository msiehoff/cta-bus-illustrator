interface Props {
  label: string
  value: string
  trend?: string        // e.g. "+2.3% vs last month"
  trendUp?: boolean     // true = green, false = red, undefined = neutral
}

export default function StatCard({ label, value, trend, trendUp }: Props) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
      <p className="text-xs text-gray-400 mb-1">{label}</p>
      <p className="text-2xl font-semibold text-white">{value}</p>
      {trend && (
        <p className={`text-xs mt-1 ${
          trendUp === true  ? 'text-green-400' :
          trendUp === false ? 'text-red-400'   :
                              'text-gray-500'
        }`}>
          {trend}
        </p>
      )}
    </div>
  )
}
