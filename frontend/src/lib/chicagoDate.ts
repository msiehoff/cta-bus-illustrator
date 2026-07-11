/** Format an ISO timestamp as YYYY-MM-DD in America/Chicago (for date inputs). */
export const toChicagoServiceDate = (iso: string) =>
  new Intl.DateTimeFormat('en-CA', {
    timeZone: 'America/Chicago',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date(iso))
