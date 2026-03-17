"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Brain } from "lucide-react"
import { useEffect, useState } from "react"
import { apiClient, Memory } from "@/lib/api-client"

interface MemoryPanelProps {
  userId: string
  refreshTrigger?: number
}

export function MemoryPanel({ userId, refreshTrigger }: MemoryPanelProps) {
  const [memories, setMemories] = useState<Memory[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    loadMemories()
  }, [userId, refreshTrigger])

  const loadMemories = async () => {
    try {
      setIsLoading(true)
      const data = await apiClient.memories.getAll(userId)
      setMemories(data.memories || [])
    } catch (error) {
      console.error("Failed to load memories:", error)
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <Brain className="h-5 w-5" />
          <h2 className="text-lg font-semibold">Your Memories</h2>
        </div>
        <p className="text-sm text-muted-foreground">Loading...</p>
      </div>
    )
  }

  if (memories.length === 0) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <Brain className="h-5 w-5" />
          <h2 className="text-lg font-semibold">Your Memories</h2>
        </div>
        <p className="text-sm text-muted-foreground">
          No memories yet. Complete onboarding to add some!
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2">
        <Brain className="h-5 w-5" />
        <h2 className="text-lg font-semibold">Your Memories</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Things I remember about your preferences
      </p>
      <div className="space-y-3">
        {memories.map((memory) => (
          <Card key={memory.id} className="relative group">
            <CardContent className="p-4">
              <h4 className="font-medium text-sm">{memory.title}</h4>
              <p className="text-sm text-muted-foreground mt-1">{memory.content}</p>
              <span className="text-xs bg-muted px-2 py-0.5 rounded mt-2 inline-block">
                {memory.category}
              </span>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
