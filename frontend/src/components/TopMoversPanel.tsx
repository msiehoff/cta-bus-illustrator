import { useNavigate } from 'react-router-dom'
import type { RouteComparison } from '../types/api'
import { formatPct } from '../lib/ridershipUtils'

interface Props {
  routes: RouteComparison[]
}

const MoverRow = ({ route, pct }: { route: RouteComparison; pct: number }) => {
  const navigate = useNavigate()
  return (
    <button
      type="button"
      onClick={() => navigate(`/routes/${route.routeId}`, { state: { routeName: route.routeName } })}
      className="w-full flex items-center justify-between gap-2 py-1.5 text-left hover:bg-gray-800/40 rounded px-1 transition-colors"
    >
      <span className="text-white text-xs truncate">{route.routeName}</span>
      <span className={pct >= 0 ? 'text-green-400 text-xs shrink-0' : 'text-red-400 text-xs shrink-0'}>
        {formatPct(pct)}
      </span>
    </button>
  )
}

const TopMoversPanel = ({ routes }: Props) => {
  const withYoY = routes.filter(r => r.yearAgoPct != null)
  const gainers = [...withYoY].sort((a, b) => (b.yearAgoPct ?? 0) - (a.yearAgoPct ?? 0)).slice(0, 5)
  const losers = [...withYoY].sort((a, b) => (a.yearAgoPct ?? 0) - (b.yearAgoPct ?? 0)).slice(0, 5)

  if (!gainers.length) return null

  return (
    <div className="grid sm:grid-cols-2 gap-3 mb-5">
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
        <p className="text-xs text-gray-500 uppercase tracking-wide mb-2">Biggest YoY gains</p>
        <div className="space-y-0.5">
          {gainers.map(route => (
            <MoverRow key={route.routeId} route={route} pct={route.yearAgoPct!} />
          ))}
        </div>
      </div>
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
        <p className="text-xs text-gray-500 uppercase tracking-wide mb-2">Biggest YoY declines</p>
        <div className="space-y-0.5">
          {losers.map(route => (
            <MoverRow key={route.routeId} route={route} pct={route.yearAgoPct!} />
          ))}
        </div>
      </div>
    </div>
  )
}

export default TopMoversPanel
