import { useEffect, useState } from 'react'
import type { RidershipType } from '../types/api'
import { useRidershipMonths } from './useRidershipMonths'

interface UseRidershipFiltersResult {
  months: string[]
  monthsLoading: boolean
  selectedMonth: string | null
  setSelectedMonth: (month: string) => void
  ridershipType: RidershipType
  setRidershipType: (type: RidershipType) => void
}

export const useRidershipFilters = (): UseRidershipFiltersResult => {
  const { months, loading: monthsLoading } = useRidershipMonths()
  const [selectedMonth, setSelectedMonth] = useState<string | null>(null)
  const [ridershipType, setRidershipType] = useState<RidershipType>('weekday')

  useEffect(() => {
    if (months.length && !selectedMonth) {
      setSelectedMonth(months[0])
    }
  }, [months, selectedMonth])

  return {
    months,
    monthsLoading,
    selectedMonth,
    setSelectedMonth,
    ridershipType,
    setRidershipType,
  }
}
