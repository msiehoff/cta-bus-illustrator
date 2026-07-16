import { NavLink } from 'react-router-dom'
import { BusIcon } from './SVG'

interface NavItem {
  to: string
  label: string
  icon: React.ReactNode
  end?: boolean
}

const MapIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polygon points="3 6 9 3 15 6 21 3 21 18 15 21 9 18 3 21" />
    <line x1="9" y1="3" x2="9" y2="18" />
    <line x1="15" y1="6" x2="15" y2="21" />
  </svg>
)

const ChartIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
  </svg>
)

const ListIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <line x1="8" y1="6" x2="21" y2="6" />
    <line x1="8" y1="12" x2="21" y2="12" />
    <line x1="8" y1="18" x2="21" y2="18" />
    <line x1="3" y1="6" x2="3.01" y2="6" />
    <line x1="3" y1="12" x2="3.01" y2="12" />
    <line x1="3" y1="18" x2="3.01" y2="18" />
  </svg>
)

export const ridershipNavItems: NavItem[] = [
  { to: '/', label: 'Map', icon: <MapIcon />, end: true },
  { to: '/system', label: 'System', icon: <ChartIcon /> },
  { to: '/routes', label: 'Routes', icon: <ListIcon /> },
]

export const headwayNavItems: NavItem[] = [
  { to: '/headways/system', label: 'System', icon: <ChartIcon /> },
  { to: '/headways/routes', label: 'Routes', icon: <ListIcon /> },
]

/** @deprecated use ridershipNavItems */
export const navItems = ridershipNavItems

interface SidebarContentProps {
  onNavClick?: () => void
}

const NavSection = ({
  title,
  items,
  onNavClick,
}: {
  title: string
  items: NavItem[]
  onNavClick?: () => void
}) => (
  <>
    <p className="px-4 pt-2 pb-1 text-[10px] uppercase tracking-widest text-gray-600">{title}</p>
    {items.map(item => (
      <NavLink
        key={item.to}
        to={item.to}
        end={item.end}
        onClick={onNavClick}
        className={({ isActive }) =>
          `flex items-center gap-2.5 px-4 py-2 text-sm transition-colors ${
            isActive
              ? 'text-red-400 bg-red-950/40'
              : 'text-gray-400 hover:text-white hover:bg-gray-800/50'
          }`
        }
      >
        {({ isActive }) => (
          <>
            <span className={isActive ? 'text-red-500' : 'text-gray-500'}>
              {item.icon}
            </span>
            {item.label}
          </>
        )}
      </NavLink>
    ))}
  </>
)

export const SidebarContent = ({ onNavClick }: SidebarContentProps) => (
  <>
    <div className="flex items-start gap-2.5 px-4 py-3 border-b border-gray-800">
      <BusIcon size={24} className="text-red-500 shrink-0 mt-0.5" />
      <div>
        <p className="text-sm font-semibold text-white leading-tight">Chicago Transit Lab</p>
        <p className="text-xs text-gray-500 mt-0.5">Bus ridership & service</p>
      </div>
    </div>

    <nav className="flex-1 py-2">
      <NavSection title="Ridership" items={ridershipNavItems} onNavClick={onNavClick} />
      <div className="border-t border-gray-800 mt-2 pt-1">
        <NavSection title="Headway" items={headwayNavItems} onNavClick={onNavClick} />
      </div>
    </nav>
  </>
)

const Sidebar = () => (
  <aside className="hidden md:flex w-56 shrink-0 bg-gray-950 border-r border-gray-800 flex-col h-full">
    <SidebarContent />
  </aside>
)

export default Sidebar
