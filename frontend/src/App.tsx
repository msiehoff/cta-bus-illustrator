import { BusIcon } from './components/SVG'
import RouteMap from './components/RouteMap'

const App = () => (
  <div className="flex flex-col h-full bg-gray-950 text-white">
    <header className="flex items-center gap-3 px-6 py-4 bg-gray-900 border-b border-gray-800 shrink-0">
      <BusIcon size={28} className="text-red-500" />
      <h1 className="text-xl font-bold tracking-tight">CTA Bus Illustrator</h1>
    </header>
    <main className="flex-1 min-h-0">
      <RouteMap />
    </main>
  </div>
)

export default App
