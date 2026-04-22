import { useEffect, useState } from 'react'
import type { ComponentType } from 'react'

type AgentationModule = {
  Agentation: ComponentType
}

type DevAgentationProps = {
  isDev?: boolean
  loadAgentation?: () => Promise<AgentationModule>
}

const defaultLoadAgentation = () => import('agentation')

export function DevAgentation({
  isDev = import.meta.env.DEV,
  loadAgentation = defaultLoadAgentation,
}: DevAgentationProps) {
  const [AgentationComponent, setAgentationComponent] = useState<ComponentType | null>(null)

  useEffect(() => {
    if (!isDev) {
      return
    }

    let cancelled = false

    void loadAgentation().then((module) => {
      if (!cancelled) {
        setAgentationComponent(() => module.Agentation)
      }
    })

    return () => {
      cancelled = true
    }
  }, [isDev, loadAgentation])

  if (!isDev || !AgentationComponent) {
    return null
  }

  return <AgentationComponent />
}
