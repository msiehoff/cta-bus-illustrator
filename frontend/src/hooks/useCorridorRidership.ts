import { useMemo } from 'react'
import type { RidershipDataPoint } from '../types/api'
import { mergeRidershipRecords } from '../lib/corridors'
import { getPairedRouteId } from '../lib/expressPairs'
import { useRouteRidership } from './useRouteRidership'

interface UseCorridorRidershipResult {
  records: RidershipDataPoint[]
  loading: boolean
  error: string | null
  localId: string
  expressId: string
  localRecords: RidershipDataPoint[]
  expressRecords: RidershipDataPoint[]
}

export const useCorridorRidership = (localId: string): UseCorridorRidershipResult => {
  const expressId = getPairedRouteId(localId) ?? ''
  const local = useRouteRidership(localId)
  const express = useRouteRidership(expressId)

  const records = useMemo(
    () => mergeRidershipRecords(local.records, express.records),
    [local.records, express.records],
  )

  return {
    records,
    loading: local.loading || express.loading,
    error: local.error ?? express.error,
    localId,
    expressId,
    localRecords: local.records,
    expressRecords: express.records,
  }
}
