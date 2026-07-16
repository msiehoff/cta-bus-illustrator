import type { ReactNode } from 'react'
import { twMerge } from 'tailwind-merge'

interface Props {
  hint: string
  children: ReactNode
  className?: string
  /** Prefer left/right on narrow layouts so the panel isn’t clipped. */
  align?: 'left' | 'center' | 'right'
  /** Table headers inside overflow containers should use "above". */
  placement?: 'below' | 'above'
}

/** Dotted-underline label with an immediate CSS hover tooltip (not native title). */
const HintLabel = ({
  hint,
  children,
  className,
  align = 'left',
  placement = 'below',
}: Props) => (
  <span className={twMerge('relative inline-flex group/hint', className)}>
    <span
      tabIndex={0}
      className="underline decoration-dotted decoration-gray-600 underline-offset-2 cursor-help outline-none"
    >
      {children}
    </span>
    <span
      role="tooltip"
      className={twMerge(
        'pointer-events-none absolute z-50 w-56 max-w-[70vw]',
        'rounded-md border border-gray-700 bg-gray-900 px-2.5 py-2 text-[11px] leading-snug text-gray-300 shadow-lg',
        'opacity-0 invisible group-hover/hint:opacity-100 group-hover/hint:visible',
        'group-focus-within/hint:opacity-100 group-focus-within/hint:visible',
        'transition-opacity duration-100',
        placement === 'below' ? 'top-full mt-1.5' : 'bottom-full mb-1.5',
        align === 'center' && 'left-1/2 -translate-x-1/2',
        align === 'right' && 'right-0',
        align === 'left' && 'left-0',
      )}
    >
      {hint}
    </span>
  </span>
)

export default HintLabel
