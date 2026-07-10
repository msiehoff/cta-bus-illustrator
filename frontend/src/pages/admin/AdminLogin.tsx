import { FormEvent, useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import { adminLogin } from '../../lib/adminApi'
import { useAdminSession } from '../../hooks/useAdminSession'

const AdminLogin = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { session, loading, refresh } = useAdminSession()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const from = (location.state as { from?: string } | null)?.from ?? '/admin'

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center px-4">
        <p className="text-gray-400">Loading…</p>
      </div>
    )
  }

  if (session?.authenticated) {
    return <Navigate to={from} replace />
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setError(null)
    try {
      await adminLogin(username, password)
      await refresh()
      navigate(from, { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center px-4">
      <div className="w-full max-w-sm bg-gray-900 border border-gray-800 rounded-lg p-6">
        <p className="text-xs uppercase tracking-wide text-gray-500 mb-1">CTA Transit Lab</p>
        <h1 className="text-xl font-semibold mb-6">Admin Login</h1>

        <form onSubmit={handleSubmit} className="space-y-4">
          <label className="block">
            <span className="text-sm text-gray-400">Username</span>
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              autoComplete="username"
              className="mt-1 w-full rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white"
              required
            />
          </label>

          <label className="block">
            <span className="text-sm text-gray-400">Password</span>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              autoComplete="current-password"
              className="mt-1 w-full rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white"
              required
            />
          </label>

          {error && <p className="text-sm text-red-400">{error}</p>}

          <button
            type="submit"
            disabled={submitting}
            className="w-full rounded-md bg-red-600 hover:bg-red-500 disabled:opacity-60 px-4 py-2 text-sm font-medium"
          >
            {submitting ? 'Signing in…' : 'Sign in'}
          </button>
        </form>
      </div>
    </div>
  )
}

export default AdminLogin
