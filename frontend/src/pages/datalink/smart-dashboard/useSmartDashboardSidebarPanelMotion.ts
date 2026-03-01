import { useMemo } from 'react'
import type { CSSProperties } from 'react'
import { MOTION_TOKENS, resolveIntentMotionClass } from '../../../features/datalink/motionGuidance'

type GuideStage = 'idle' | 'grid' | 'commit'

interface UseSmartDashboardSidebarPanelMotionInput {
  guideStage: GuideStage
  reducedMotion: boolean
}

interface SmartDashboardSidebarPanelMotion {
  className: string
  style: CSSProperties
}

export function useSmartDashboardSidebarPanelMotion({
  guideStage,
  reducedMotion,
}: UseSmartDashboardSidebarPanelMotionInput): SmartDashboardSidebarPanelMotion {
  return useMemo(() => {
    const className = `h-full bg-slate-900/60 backdrop-blur-xl border border-white/5 rounded-3xl overflow-hidden shadow-2xl transition-all ${resolveIntentMotionClass(guideStage === 'commit' ? 'commit' : 'idle', reducedMotion)}`
    return {
      className,
      style: {
        transitionDuration: `${MOTION_TOKENS.commitFeedbackMs}ms`,
      },
    }
  }, [guideStage, reducedMotion])
}
