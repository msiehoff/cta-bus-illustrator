import { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom'
import type { GetRoutesComparisonResponse } from '../types/api'

export const useRouteName = (externalId: string): string => {
  const location = useLocation()
  const fromState = (location.state as { routeName?: string } | null)?.routeName
  const [name, setName] = useState(fromState ?? '')

  useEffect(() => {
    if (fromState) {
      setName(fromState)
      return
    }

    fetch('/api/v1/routes/comparison')
      .then(r => r.json())
      .then((data: GetRoutesComparisonResponse) => {
        const route = data.routes.find(r => r.routeId === externalId)
        if (route) setName(route.routeName)
      })
      .catch(() => {})
  }, [externalId, fromState])

  return name || externalId
}
