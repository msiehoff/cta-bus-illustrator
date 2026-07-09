import { useEffect, useState } from 'react'
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import Sidebar, { SidebarContent } from './components/Sidebar'
import MobileHeader from './components/MobileHeader'
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

const AppShell = () => {
  const location = useLocation()
  const [mobileNavOpen, setMobileNavOpen] = useState(false)

  useEffect(() => {
    setMobileNavOpen(false)
  }, [location.pathname])

  useEffect(() => {
    document.body.style.overflow = mobileNavOpen ? 'hidden' : ''
    return () => {
      document.body.style.overflow = ''
    }
  }, [mobileNavOpen])

  return (
    <div className="flex h-screen bg-gray-950 text-white overflow-hidden">
      <Sidebar />

      {mobileNavOpen && (
        <>
          <button
            type="button"
            aria-label="Close menu"
            className="fixed inset-0 z-40 bg-black/60 md:hidden"
            onClick={() => setMobileNavOpen(false)}
          />
          <aside className="fixed inset-y-0 left-0 z-50 w-56 bg-gray-950 border-r border-gray-800 flex flex-col md:hidden shadow-2xl">
            <SidebarContent onNavClick={() => setMobileNavOpen(false)} />
          </aside>
        </>
      )}

      <div className="flex flex-1 flex-col min-w-0 min-h-0">
        <MobileHeader
          open={mobileNavOpen}
          onToggle={() => setMobileNavOpen(open => !open)}
        />
        <main className="flex-1 min-h-0 overflow-y-auto">
          <Routes>
            <Route path="/" element={<RouteMap />} />
            <Route path="/system" element={<SystemOverview />} />
            <Route path="/routes" element={<RoutesPage />} />
            <Route path="/routes/:externalId" element={<RoutePage />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}

const App = () => (
  <BrowserRouter>
    <Analytics />
    <AppShell />
  </BrowserRouter>
)

export default App
