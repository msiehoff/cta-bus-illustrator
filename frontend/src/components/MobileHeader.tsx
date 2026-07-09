import { MenuIcon } from './SVG'

interface Props {
  open: boolean
  onToggle: () => void
}

const MobileHeader = ({ open, onToggle }: Props) => (
  <header className="md:hidden shrink-0 flex items-center gap-3 px-3 py-2.5 bg-gray-950 border-b border-gray-800">
    <button
      type="button"
      onClick={onToggle}
      aria-label={open ? 'Close menu' : 'Open menu'}
      aria-expanded={open}
      className="flex items-center justify-center w-10 h-10 rounded-lg text-gray-300 hover:text-white hover:bg-gray-800 transition-colors"
    >
      <MenuIcon size={22} />
    </button>
    <div className="min-w-0">
      <p className="text-sm font-semibold text-white truncate">Chicago Transit Lab</p>
      <p className="text-xs text-gray-500 truncate">Bus Ridership</p>
    </div>
  </header>
)

export default MobileHeader
