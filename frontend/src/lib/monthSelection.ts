const MONTH_NAMES = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
] as const

export const getAvailableYears = (months: string[]): number[] =>
  [...new Set(months.map(m => Number(m.split('-')[0])))].sort((a, b) => b - a)

export const getAvailableMonthsInYear = (months: string[], year: number): string[] =>
  months
    .filter(m => Number(m.split('-')[0]) === year)
    .sort((a, b) => a.localeCompare(b))

export const parseMonth = (month: string | null): { year: number; monthNum: string } | null => {
  if (!month) return null
  const [year, monthNum] = month.split('-')
  return { year: Number(year), monthNum }
}

export const formatMonthOption = (month: string): string => {
  const parsed = parseMonth(month)
  if (!parsed) return month
  return MONTH_NAMES[Number(parsed.monthNum) - 1] ?? month
}

export const getAdjacentMonth = (
  months: string[],
  current: string | null,
  direction: 'newer' | 'older',
): string | null => {
  if (!current) return null
  const index = months.indexOf(current)
  if (index < 0) return null
  const nextIndex = direction === 'newer' ? index - 1 : index + 1
  return months[nextIndex] ?? null
}

export const resolveMonthForYear = (
  months: string[],
  year: number,
  preferredMonthNum?: string,
): string | null => {
  const inYear = getAvailableMonthsInYear(months, year)
  if (!inYear.length) return null
  if (preferredMonthNum) {
    const match = inYear.find(m => m.split('-')[1] === preferredMonthNum)
    if (match) return match
  }
  return inYear[inYear.length - 1]
}
