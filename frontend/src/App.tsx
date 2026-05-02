import { useState } from 'react'
import { twMerge } from 'tailwind-merge'
import { BusIcon } from './components/SVG'
import RouteMap from './components/RouteMap'
import LakeShoreSimulator from './components/LakeShoreSimulator'

type Page = 'map' | 'simulator'

const PAGES: Array<{ id: Page; label: string; description: string }> = [
  {
    id: 'map',
    label: 'Bus Ridership Map',
    description: 'CTA monthly ridership averages by bus route — colored and sized by average daily ridership.',
  },
  {
    id: 'simulator',
    label: 'Lake Shore Drive Simulator',
    description: 'Compare mixed-traffic vs. bus-priority lane allocation — see how dedicating lanes to buses affects travel times and people moved.',
  },
]

const App = () => {
  const [page, setPage] = useState<Page>('map')
  const currentPage = PAGES.find((p) => p.id === page)!

  return (
    <div className="flex flex-col h-full bg-gray-950 text-white">
      <header className="flex items-start justify-between gap-4 px-6 py-3 bg-gray-900 border-b border-gray-800 shrink-0">
        <div className="flex items-start gap-3 min-w-0">
          <BusIcon size={28} className="text-red-500 shrink-0 mt-0.5" />
          <div className="min-w-0">
            <h1 className="text-xl font-bold tracking-tight">Chicago Transit Lab</h1>
            <p className="text-gray-400 text-sm mt-0.5 leading-snug max-w-2xl">
              {currentPage.description}
            </p>
          </div>
        </div>
        <nav className="flex gap-1 mt-0.5 shrink-0">
          {PAGES.map((p) => (
            <button
              key={p.id}
              onClick={() => setPage(p.id)}
              className={twMerge(
                'px-3 py-1.5 text-sm rounded-md font-medium transition-colors whitespace-nowrap',
                page === p.id
                  ? 'bg-red-600 text-white'
                  : 'text-gray-400 hover:text-white hover:bg-gray-800',
              )}
            >
              {p.label}
            </button>
          ))}
        </nav>
      </header>
      <main className="flex-1 min-h-0">
        {page === 'map' ? <RouteMap /> : <LakeShoreSimulator />}
      </main>
    </div>
  )
}

export default App
