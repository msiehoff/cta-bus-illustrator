import { BusIcon } from './components/SVG'
import RouteMap from './components/RouteMap'

const App = () => (
  <div className="flex flex-col h-full bg-gray-950 text-white">
    <header className="flex items-start gap-3 px-6 py-3 bg-gray-900 border-b border-gray-800 shrink-0">
      <BusIcon size={28} className="text-red-500 shrink-0 mt-0.5" />
      <div className="min-w-0">
        <h1 className="text-xl font-bold tracking-tight">Chicago Transit Lab</h1>
        <p className="text-gray-400 text-sm mt-0.5 leading-snug max-w-2xl">
          Bus routes are colored and sized by average daily ridership—use the map and filters to explore where demand is highest.
        </p>
      </div>
    </header>
    <main className="flex-1 min-h-0">
      <RouteMap />
    </main>
  </div>
)

export default App
