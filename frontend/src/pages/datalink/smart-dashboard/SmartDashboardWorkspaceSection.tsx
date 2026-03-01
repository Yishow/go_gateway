import type { ComponentProps, CSSProperties } from 'react'
import SmartDashboardWorkspace from './SmartDashboardWorkspace'
import SmartDashboardSidebar from './SmartDashboardSidebar'

interface SmartDashboardWorkspaceSectionProps {
  workspaceProps: ComponentProps<typeof SmartDashboardWorkspace>
  sidebarProps: ComponentProps<typeof SmartDashboardSidebar>
  sidebarContainerClassName: string
  sidebarContainerStyle: CSSProperties
}

export default function SmartDashboardWorkspaceSection({
  workspaceProps,
  sidebarProps,
  sidebarContainerClassName,
  sidebarContainerStyle,
}: SmartDashboardWorkspaceSectionProps) {
  return (
    <div className="flex-1 p-3 sm:p-4 grid grid-cols-1 xl:grid-cols-[1fr_380px] gap-4 min-h-0">
      <SmartDashboardWorkspace {...workspaceProps} />
      <div className="min-h-[300px] xl:min-h-0">
        <div className={sidebarContainerClassName} style={sidebarContainerStyle}>
          <SmartDashboardSidebar {...sidebarProps} />
        </div>
      </div>
    </div>
  )
}
