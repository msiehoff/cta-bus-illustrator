import { twMerge } from 'tailwind-merge'
import {
  HEADWAY_WINDOW_OPTIONS,
  type HeadwayWindowKey,
} from '../lib/headwayUtils'

type Props = {
  value: HeadwayWindowKey
  onChange: (key: HeadwayWindowKey) => void
  className?: string
}

const HeadwayWindowFilter = ({ value, onChange, className }: Props) => (
  <div className={twMerge('flex flex-wrap gap-1', className)} role="group" aria-label="Time window">
    {HEADWAY_WINDOW_OPTIONS.map(opt => (
      <button
        key={opt.key}
        type="button"
        onClick={() => onChange(opt.key)}
        className={twMerge(
          'px-2.5 py-1.5 rounded-md text-xs transition-colors',
          value === opt.key
            ? 'bg-gray-700 text-white'
            : 'text-gray-500 hover:text-gray-300 hover:bg-gray-800',
        )}
      >
        {opt.label}
      </button>
    ))}
  </div>
)

export default HeadwayWindowFilter
