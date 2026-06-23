import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import Sidebar from './components/Sidebar'
import RouteMap from './components/RouteMap'
import SystemOverview from './pages/SystemOverview'
import RoutesPage from './pages/RoutesPage'
import RoutePage from './pages/RoutePage'
import { trackPageView } from './lib/analytics'

const Analytics = () => {
  const location = useLocation()

  useEffect(() => {
    trackPageView(location.pathname + location.search)
  }, [location])

  return null
}

const App = () => (
  <BrowserRouter>
    <Analytics />
    <div className="flex h-screen bg-gray-950 text-white overflow-hidden">
      <Sidebar />
      <main className="flex-1 min-w-0 overflow-y-auto">
        <Routes>
          <Route path="/" element={<RouteMap />} />
          <Route path="/system" element={<SystemOverview />} />
          <Route path="/routes" element={<RoutesPage />} />
          <Route path="/routes/:externalId" element={<RoutePage />} />
        </Routes>
      </main>
    </div>
  </BrowserRouter>
)

export default App
