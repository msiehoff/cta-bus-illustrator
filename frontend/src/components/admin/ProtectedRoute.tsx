import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAdminSession } from '../../hooks/useAdminSession'

const ProtectedRoute = () => {
  const location = useLocation()
  const { session, loading } = useAdminSession()

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center">
        <p className="text-gray-400">Checking session…</p>
      </div>
    )
  }

  if (!session?.authenticated) {
    return <Navigate to="/admin/login" replace state={{ from: location.pathname }} />
  }

  return <Outlet />
}

export default ProtectedRoute
