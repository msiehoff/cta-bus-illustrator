import { useState, useMemo } from 'react'
import { twMerge } from 'tailwind-merge'
import { BusIcon, CarIcon, UsersIcon, ClockIcon } from './SVG'

// ─── Constants ───────────────────────────────────────────────────────────────

const CORRIDOR_MILES = 10
const TOTAL_LANES = 4
const LANE_CAP = 1800        // vehicles / lane / hr (theoretical max flow)
const CAR_OCC = 1.2          // people per car
const BUS_OCC = 45           // people per bus
const CAR_FREE_FLOW = 45     // mph, unimpeded
const BUS_MIXED_SPEED = 12   // mph, stuck in mixed traffic with stops
const BUS_LANE_CAP = 60      // buses / dedicated lane / hr (one every ~60 s)
const BPR_ALPHA = 0.15       // Bureau of Public Roads congestion parameters
const BPR_BETA = 4
const MIN_CAR_SPEED = 5      // mph floor (gridlock)

// ─── Types ───────────────────────────────────────────────────────────────────

interface ScenarioResult {
  busRiders: number
  carRiders: number
  busVehicles: number
  carVehicles: number
  totalVehicles: number
  carLanes: number
  busLanes: number
  vc: number
  carSpeed: number
  carTravelTime: number
  busTravelTime: number
  avgTravelTime: number
  totalDemand: number
  personMinutes: number
}

// ─── Simulation engine ───────────────────────────────────────────────────────

const computeScenario = (
  totalDemand: number,
  transitSharePct: number,
  busLanesCount: number,
  busSpeed: number,
): ScenarioResult => {
  const transitShare = transitSharePct / 100
  const busRiders = Math.round(totalDemand * transitShare)
  const carRiders = totalDemand - busRiders

  const busVehicles = Math.round(busRiders / BUS_OCC)
  const carVehicles = Math.round(carRiders / CAR_OCC)
  const totalVehicles = carVehicles + busVehicles

  const carLanes = Math.max(TOTAL_LANES - busLanesCount, 1)
  const vehiclesPerCarLane = carVehicles / carLanes
  const vc = vehiclesPerCarLane / LANE_CAP

  const rawCarSpeed = CAR_FREE_FLOW / (1 + BPR_ALPHA * Math.pow(Math.max(vc, 0), BPR_BETA))
  const carSpeed = Math.max(rawCarSpeed, MIN_CAR_SPEED)

  const carTravelTime = (CORRIDOR_MILES / carSpeed) * 60
  const busTravelTime = (CORRIDOR_MILES / busSpeed) * 60
  const avgTravelTime = (carRiders * carTravelTime + busRiders * busTravelTime) / totalDemand
  const personMinutes = avgTravelTime * totalDemand

  return {
    busRiders, carRiders, busVehicles, carVehicles, totalVehicles,
    carLanes, busLanes: busLanesCount, vc,
    carSpeed, carTravelTime, busTravelTime, avgTravelTime, totalDemand, personMinutes,
  }
}

// ─── Congestion theming ──────────────────────────────────────────────────────

const vcTheme = (vc: number) => {
  if (vc < 0.6)  return { bg: 'bg-green-950',  border: 'border-green-800',  text: 'text-green-400',  label: 'Free Flow' }
  if (vc < 0.85) return { bg: 'bg-yellow-950', border: 'border-yellow-800', text: 'text-yellow-400', label: 'Moderate' }
  if (vc < 1.05) return { bg: 'bg-orange-950', border: 'border-orange-800', text: 'text-orange-400', label: 'Heavy' }
  return           { bg: 'bg-red-950',    border: 'border-red-800',    text: 'text-red-400',    label: 'Gridlock' }
}

// ─── Sub-components ──────────────────────────────────────────────────────────

interface SliderProps {
  label: string
  value: number
  min: number
  max: number
  step?: number
  onChange: (v: number) => void
  format?: (v: number) => string
}

const Slider = ({ label, value, min, max, step = 1, onChange, format = String }: SliderProps) => (
  <label className="flex flex-col gap-1.5 cursor-pointer">
    <div className="flex justify-between items-baseline">
      <span className="text-xs text-gray-400">{label}</span>
      <span className="text-xs font-bold text-white tabular-nums">{format(value)}</span>
    </div>
    <input
      type="range"
      min={min}
      max={max}
      step={step}
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
      className="w-full h-1.5 rounded-full appearance-none cursor-pointer accent-red-500 bg-gray-700"
    />
    <div className="flex justify-between text-xs text-gray-600">
      <span>{format(min)}</span>
      <span>{format(max)}</span>
    </div>
  </label>
)

interface LaneDiagramProps {
  s: ScenarioResult
  busSpeed: number
}

const LaneDiagram = ({ s, busSpeed }: LaneDiagramProps) => {
  const carTheme = vcTheme(s.vc)
  const carDensity = Math.min(Math.max(Math.round(s.vc * 7), s.vc > 0 ? 1 : 0), 7)
  const busDensity = Math.min(
    Math.max(Math.round((s.busVehicles / Math.max(s.busLanes, 1) / BUS_LANE_CAP) * 7), s.busVehicles > 0 ? 1 : 0),
    7,
  )

  return (
    <div className="space-y-1.5">
      {/* Road label */}
      <div className="text-xs text-gray-600 uppercase tracking-widest px-1 mb-2">Lake Shore Dr — {TOTAL_LANES} lanes</div>

      {/* Bus lanes */}
      {s.busLanes > 0 && Array.from({ length: s.busLanes }).map((_, i) => (
        <div
          key={`bus-${i}`}
          className="flex items-center gap-2 px-3 h-11 rounded-lg bg-amber-950 border border-amber-700"
        >
          <BusIcon size={14} className="text-amber-500 shrink-0" />
          <span className="text-xs font-bold text-amber-400 w-20 shrink-0">BUS ONLY</span>
          <div className="flex-1 flex items-center gap-1.5 overflow-hidden">
            {Array.from({ length: busDensity }).map((_, j) => (
              <BusIcon key={j} size={13} className="text-amber-300 shrink-0" />
            ))}
          </div>
          <span className="text-xs font-semibold text-amber-300 tabular-nums shrink-0">{Math.round(busSpeed)} mph</span>
        </div>
      ))}

      {/* Car / mixed lanes */}
      {Array.from({ length: s.carLanes }).map((_, i) => (
        <div
          key={`car-${i}`}
          className={twMerge('flex items-center gap-2 px-3 h-11 rounded-lg border', carTheme.bg, carTheme.border)}
        >
          <CarIcon size={14} className={twMerge(carTheme.text, 'shrink-0')} />
          <span className={twMerge('text-xs font-bold w-20 shrink-0', carTheme.text)}>
            {s.busLanes > 0 ? `CAR ${i + 1}` : `LANE ${i + 1}`}
          </span>
          <div className="flex-1 flex items-center gap-1.5 overflow-hidden">
            {Array.from({ length: carDensity }).map((_, j) => (
              <CarIcon key={j} size={13} className={twMerge(carTheme.text, 'shrink-0')} />
            ))}
          </div>
          <span className={twMerge('text-xs font-semibold tabular-nums shrink-0', carTheme.text)}>
            {Math.round(s.carSpeed)} mph
          </span>
        </div>
      ))}

      {/* Congestion badge */}
      <div className={twMerge('text-xs font-medium px-2 py-1 rounded w-fit mt-1', carTheme.text)}>
        Car traffic: {carTheme.label} ({Math.round(s.vc * 100)}% of capacity)
      </div>
    </div>
  )
}

interface MetricRowProps {
  icon: React.ReactNode
  label: string
  value: string
  color?: string
  sub?: React.ReactNode
}

const MetricRow = ({ icon, label, value, color = 'text-white', sub }: MetricRowProps) => (
  <div className="flex items-center gap-3 py-3 border-b border-gray-800 last:border-0">
    <div className="text-gray-500 shrink-0">{icon}</div>
    <span className="text-sm text-gray-400 flex-1">{label}</span>
    <div className="text-right">
      <span className={twMerge('text-sm font-bold tabular-nums', color)}>{value}</span>
      {sub && <div className="text-xs text-gray-500 mt-0.5">{sub}</div>}
    </div>
  </div>
)

const Delta = ({ val, suffix = '', invert = false }: { val: number; suffix?: string; invert?: boolean }) => {
  const improved = invert ? val > 0 : val < 0
  const neutral = Math.abs(val) < 0.5
  if (neutral) return <span className="text-gray-600">no change</span>
  const display = `${val > 0 ? '+' : ''}${Math.round(val)}${suffix}`
  return (
    <span className={twMerge('font-semibold', improved ? 'text-green-400' : 'text-red-400')}>
      {display}
    </span>
  )
}

// ─── Main component ──────────────────────────────────────────────────────────

const LakeShoreSimulator = () => {
  const [totalDemand,  setTotalDemand]  = useState(10000)
  const [transitShare, setTransitShare] = useState(15)
  const [busLanes,     setBusLanes]     = useState(1)
  const [busSpeed,     setBusSpeed]     = useState(25)
  const [modeShift,    setModeShift]    = useState(5)

  const current = useMemo(
    () => computeScenario(totalDemand, transitShare, 0, BUS_MIXED_SPEED),
    [totalDemand, transitShare],
  )

  const proposed = useMemo(
    () => computeScenario(totalDemand, Math.min(transitShare + modeShift, 90), busLanes, busSpeed),
    [totalDemand, transitShare, busLanes, busSpeed, modeShift],
  )

  const carDeltaMin     = proposed.carTravelTime - current.carTravelTime
  const busSavingMin    = current.busTravelTime  - proposed.busTravelTime
  const personMinSaved  = Math.round(current.personMinutes - proposed.personMinutes)
  const moreOnBus       = proposed.busRiders - current.busRiders
  const fewerVehicles   = current.totalVehicles - proposed.totalVehicles

  const fmt      = (n: number) => n.toLocaleString()
  const fmtMph   = (v: number) => `${Math.round(v)} mph`
  const fmtMin   = (t: number) => `${Math.round(t)} min`
  const fmtPct   = (p: number) => `${p}%`
  const fmtK     = (n: number) => Math.abs(n) >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(Math.round(n))

  return (
    <div className="flex flex-col h-full bg-gray-950 text-white overflow-auto">

      {/* ── Controls ─────────────────────────────────────────────────────── */}
      <div className="bg-gray-900 border-b border-gray-800 px-6 py-4 shrink-0">
        <p className="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-4">Simulation Controls</p>
        <div className="grid grid-cols-2 lg:grid-cols-5 gap-x-8 gap-y-5">
          <Slider
            label="Peak Hour Demand"
            value={totalDemand}
            min={5000} max={15000} step={500}
            onChange={setTotalDemand}
            format={(v) => `${(v / 1000).toFixed(1)}k people`}
          />
          <Slider
            label="Current Bus Mode Share"
            value={transitShare}
            min={5} max={40}
            onChange={setTransitShare}
            format={fmtPct}
          />
          <Slider
            label="Dedicated Bus Lanes"
            value={busLanes}
            min={1} max={2}
            onChange={setBusLanes}
            format={(v) => `${v} lane${v > 1 ? 's' : ''}`}
          />
          <Slider
            label="Bus Speed (dedicated lane)"
            value={busSpeed}
            min={15} max={35}
            onChange={setBusSpeed}
            format={fmtMph}
          />
          <Slider
            label="Mode Shift to Bus"
            value={modeShift}
            min={0} max={25}
            onChange={setModeShift}
            format={fmtPct}
          />
        </div>
      </div>

      {/* ── Two scenario panels ───────────────────────────────────────────── */}
      <div className="flex-1 grid grid-cols-2 divide-x divide-gray-800 min-h-0">

        {/* Current */}
        <div className="p-6 flex flex-col gap-5 overflow-auto">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="w-2.5 h-2.5 rounded-full bg-gray-500 shrink-0" />
              <h3 className="text-base font-bold text-gray-100">Current Conditions</h3>
            </div>
            <p className="text-xs text-gray-500 ml-4.5">
              All {TOTAL_LANES} lanes mixed traffic — cars and buses compete for space
            </p>
          </div>

          <LaneDiagram s={current} busSpeed={BUS_MIXED_SPEED} />

          <div className="rounded-xl bg-gray-900 border border-gray-800 px-4 divide-y divide-gray-800">
            <MetricRow
              icon={<UsersIcon size={15} />}
              label="People moving / hr"
              value={fmt(current.totalDemand)}
            />
            <MetricRow
              icon={<CarIcon size={15} />}
              label="Car travel time"
              value={fmtMin(current.carTravelTime)}
              color={current.vc > 0.9 ? 'text-orange-400' : 'text-white'}
            />
            <MetricRow
              icon={<BusIcon size={15} />}
              label="Bus travel time"
              value={fmtMin(current.busTravelTime)}
              color="text-amber-400"
            />
            <MetricRow
              icon={<CarIcon size={15} />}
              label="Vehicles on road"
              value={fmt(current.totalVehicles)}
            />
            <MetricRow
              icon={<ClockIcon size={15} />}
              label="Avg travel time"
              value={fmtMin(current.avgTravelTime)}
            />
          </div>
        </div>

        {/* Proposed */}
        <div className="p-6 flex flex-col gap-5 overflow-auto">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="w-2.5 h-2.5 rounded-full bg-blue-400 shrink-0" />
              <h3 className="text-base font-bold text-blue-100">Proposed: Bus Priority</h3>
            </div>
            <p className="text-xs text-gray-500 ml-4.5">
              {busLanes} dedicated bus lane{busLanes > 1 ? 's' : ''}, {TOTAL_LANES - busLanes} car lanes
              {modeShift > 0 ? ` — +${modeShift}% mode shift to bus` : ''}
            </p>
          </div>

          <LaneDiagram s={proposed} busSpeed={busSpeed} />

          <div className="rounded-xl bg-gray-900 border border-gray-800 px-4 divide-y divide-gray-800">
            <MetricRow
              icon={<UsersIcon size={15} />}
              label="People moving / hr"
              value={fmt(proposed.totalDemand)}
              sub={moreOnBus > 0 ? <><Delta val={moreOnBus} suffix="" invert /> more by bus</> : undefined}
            />
            <MetricRow
              icon={<CarIcon size={15} />}
              label="Car travel time"
              value={fmtMin(proposed.carTravelTime)}
              color={carDeltaMin > 5 ? 'text-red-400' : carDeltaMin > 2 ? 'text-orange-300' : 'text-green-400'}
              sub={<Delta val={carDeltaMin} suffix=" min" />}
            />
            <MetricRow
              icon={<BusIcon size={15} />}
              label="Bus travel time"
              value={fmtMin(proposed.busTravelTime)}
              color="text-green-400"
              sub={busSavingMin > 0 ? <span className="text-green-400">−{Math.round(busSavingMin)} min saved</span> : undefined}
            />
            <MetricRow
              icon={<CarIcon size={15} />}
              label="Vehicles on road"
              value={fmt(proposed.totalVehicles)}
              sub={fewerVehicles !== 0 ? <Delta val={-fewerVehicles} suffix="" /> : undefined}
            />
            <MetricRow
              icon={<ClockIcon size={15} />}
              label="Avg travel time"
              value={fmtMin(proposed.avgTravelTime)}
              color={proposed.avgTravelTime < current.avgTravelTime ? 'text-green-400' : 'text-orange-400'}
              sub={<Delta val={proposed.avgTravelTime - current.avgTravelTime} suffix=" min" />}
            />
          </div>
        </div>
      </div>

      {/* ── Impact summary bar ────────────────────────────────────────────── */}
      <div className="bg-gray-900 border-t border-gray-800 px-6 py-4 shrink-0">
        <p className="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-3">Impact Summary</p>
        <div className="grid grid-cols-3 gap-4">

          {/* Bus time savings */}
          <div className={twMerge(
            'rounded-xl p-4 border',
            busSavingMin >= 10 ? 'bg-green-950 border-green-800' : 'bg-gray-800 border-gray-700',
          )}>
            <div className={twMerge('text-3xl font-extrabold tabular-nums', busSavingMin >= 0 ? 'text-green-300' : 'text-red-300')}>
              {busSavingMin >= 0 ? '−' : '+'}{Math.abs(Math.round(busSavingMin))} min
            </div>
            <div className="text-xs text-gray-400 mt-1">
              Bus riders save <span className="text-green-400 font-semibold">{Math.round(busSavingMin / current.busTravelTime * 100)}%</span> per trip
              &nbsp;({fmtMin(current.busTravelTime)} → {fmtMin(proposed.busTravelTime)})
            </div>
          </div>

          {/* Car impact */}
          <div className={twMerge(
            'rounded-xl p-4 border',
            carDeltaMin <= 3 ? 'bg-gray-800 border-gray-700' :
            carDeltaMin <= 6 ? 'bg-orange-950 border-orange-800' : 'bg-red-950 border-red-800',
          )}>
            <div className={twMerge(
              'text-3xl font-extrabold tabular-nums',
              carDeltaMin <= 3 ? 'text-gray-300' :
              carDeltaMin <= 6 ? 'text-orange-300' : 'text-red-300',
            )}>
              {carDeltaMin > 0 ? '+' : ''}{Math.round(carDeltaMin)} min
            </div>
            <div className="text-xs text-gray-400 mt-1">
              Change in car travel time
              {carDeltaMin <= 3
                ? ' — minimal impact on drivers'
                : carDeltaMin <= 6
                ? ' — moderate increase for drivers'
                : ' — significant increase; consider more mode shift'}
            </div>
          </div>

          {/* Person-minutes */}
          <div className={twMerge(
            'rounded-xl p-4 border',
            personMinSaved > 0 ? 'bg-blue-950 border-blue-800' : 'bg-gray-800 border-gray-700',
          )}>
            <div className={twMerge('text-3xl font-extrabold tabular-nums', personMinSaved > 0 ? 'text-blue-300' : 'text-red-300')}>
              {personMinSaved > 0 ? '+' : ''}{fmtK(personMinSaved)}
            </div>
            <div className="text-xs text-gray-400 mt-1">
              {personMinSaved > 0 ? 'Person-minutes saved' : 'Additional person-minutes'} per peak hour
              {personMinSaved > 0 && (
                <> — equivalent to <span className="text-blue-300 font-semibold">{fmtK(personMinSaved / 60)} person-hours</span></>
              )}
            </div>
          </div>

        </div>

        {/* Assumptions footnote */}
        <p className="text-xs text-gray-600 mt-3">
          Assumptions: {CORRIDOR_MILES}-mile corridor · {CAR_OCC} people/car · {BUS_OCC} people/bus · {LANE_CAP} veh/lane/hr capacity · BPR congestion model (α={BPR_ALPHA}, β={BPR_BETA}) · Bus mixed-traffic speed {BUS_MIXED_SPEED} mph
        </p>
      </div>

    </div>
  )
}

export default LakeShoreSimulator
