// Color stops for the ridership heat ramp: yellow (low) → orange → deep red (high).
// Each entry is [avgRides threshold, hex color].
// These stops drive both the MapLibre line-color expression and the legend gradient,
// so editing here keeps both in sync automatically.
export const RIDERSHIP_STOPS: [number, string][] = [
  [0,      '#fef9c3'],  // yellow-100  — very low
  [5000,   '#fde047'],  // yellow-300  — low
  [10000,  '#fb923c'],  // orange-400  — medium
  [20000,  '#ef4444'],  // red-500     — high
  [30000,  '#7f1d1d'],  // red-900     — very high
]

// Color shown for routes that have no ridership data for the selected month/type.
export const NO_DATA_COLOR = '#cbd5e1'  // slate-300

// MapLibre expression: gray for routes with no data, color ramp otherwise.
export const ridershipColorExpression = [
  'case',
  ['!', ['has', 'avgRides']],
  NO_DATA_COLOR,
  [
    'interpolate', ['linear'],
    ['get', 'avgRides'],
    ...RIDERSHIP_STOPS.flat(),
  ],
]

// CSS linear-gradient string for the legend bar.
export const legendGradient =
  `linear-gradient(to right, ${RIDERSHIP_STOPS.map(([, color]) => color).join(', ')})`

// Human-readable tick labels aligned to the stop thresholds.
export const LEGEND_TICKS = RIDERSHIP_STOPS.map(([value]) =>
  value === 0 ? '0' : value >= 1000 ? `${value / 1000}k` : String(value)
)
