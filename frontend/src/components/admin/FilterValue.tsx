import type { ReactNode } from 'react'
import { twMerge } from 'tailwind-merge'

type FilterValueProps = {
  value: string
  active?: boolean
  onSelect: (value: string) => void
  className?: string
  children?: ReactNode
}

/** Clickable table cell value that applies an admin list filter. */
const FilterValue = ({ value, active, onSelect, className, children }: FilterValueProps) => {
  if (!value) return <span className={className}>—</span>

  return (
    <button
      type="button"
      title={`Filter by ${value}`}
      onClick={() => onSelect(value)}
      className={twMerge(
        'text-left hover:underline decoration-red-400/60 underline-offset-2',
        active ? 'text-red-400' : 'hover:text-red-300',
        className,
      )}
    >
      {children ?? value}
    </button>
  )
}

export default FilterValue
