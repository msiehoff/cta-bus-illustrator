export const formatHeadwayMinutes = (mins: number) => {
  if (!Number.isFinite(mins)) return '—'
  if (mins >= 10) return mins.toFixed(0)
  return mins.toFixed(1)
}

export const formatHeadwayCV = (cv: number) => {
  if (!Number.isFinite(cv)) return '—'
  return cv.toFixed(2)
}

export const formatHeadwayPeriod = (start?: string, end?: string, daysWithData?: number) => {
  if (!start || !end) {
    return daysWithData ? `${daysWithData} day${daysWithData === 1 ? '' : 's'} of data` : 'No data yet'
  }
  if (start === end) return start
  return `${start} → ${end}`
}
