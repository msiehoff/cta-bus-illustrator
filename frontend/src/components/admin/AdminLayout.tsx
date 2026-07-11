import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { twMerge } from 'tailwind-merge'
import { adminLogout } from '../../lib/adminApi'

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  twMerge(
    'block px-3 py-2 rounded-md text-sm',
    isActive ? 'bg-red-950/40 text-red-400' : 'text-gray-400 hover:text-white hover:bg-gray-900',
  )

const AdminLayout = () => {
  const navigate = useNavigate()

  const handleLogout = async () => {
    await adminLogout()
    navigate('/admin/login')
  }

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      <header className="border-b border-gray-800 bg-gray-950/95 sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-wide text-gray-500">CTA Transit Lab</p>
            <h1 className="text-lg font-semibold">Admin</h1>
          </div>
          <button
            type="button"
            onClick={handleLogout}
            className="text-sm text-gray-400 hover:text-white"
          >
            Log out
          </button>
        </div>
      </header>

      <div className="max-w-6xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-[180px_1fr] gap-6">
        <nav className="space-y-1">
          <NavLink to="/admin" end className={navLinkClass}>
            Dashboard
          </NavLink>
          <NavLink to="/admin/arrivals" className={navLinkClass}>
            Arrivals
          </NavLink>
          <NavLink to="/admin/headways" className={navLinkClass}>
            Headways
          </NavLink>
          <NavLink to="/" className="block px-3 py-2 rounded-md text-sm text-gray-500 hover:text-white">
            ← Back to site
          </NavLink>
        </nav>

        <main>
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default AdminLayout
