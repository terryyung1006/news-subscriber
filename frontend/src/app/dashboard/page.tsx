"use client"

import { useState, useEffect } from "react"
import { Header } from "@/components/dashboard/header"
import { MemoryPanel } from "@/components/dashboard/memory-panel"
import { ReportPanel } from "@/components/dashboard/report-panel"
import { OnboardingChat } from "@/components/onboarding/onboarding-chat"
import { apiClient, Memory } from "@/lib/api-client"

export default function DashboardPage() {
  const [showOnboarding, setShowOnboarding] = useState(false)
  const [isCheckingStatus, setIsCheckingStatus] = useState(true)
  const [userId, setUserId] = useState<string | null>(null)
  const [refreshMemories, setRefreshMemories] = useState(0)

  useEffect(() => {
    const storedUserId = localStorage.getItem("user_id")
    if (storedUserId) {
      setUserId(storedUserId)
      checkOnboardingStatus(storedUserId)
    } else {
      setIsCheckingStatus(false)
    }
  }, [])

  const checkOnboardingStatus = async (uid: string) => {
    try {
      const status = await apiClient.onboarding.getStatus(uid)
      setShowOnboarding(status.needs_onboarding)
    } catch (error) {
      console.error("Failed to check onboarding status:", error)
    } finally {
      setIsCheckingStatus(false)
    }
  }

  const handleOnboardingComplete = (memories: Memory[]) => {
    setShowOnboarding(false)
    setRefreshMemories(prev => prev + 1)
  }

  if (isCheckingStatus) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex flex-1 flex-col md:flex-row overflow-hidden">
        <aside className="w-full md:w-1/3 lg:w-1/4 border-r bg-muted/10 p-4 overflow-y-auto">
          {userId && <MemoryPanel userId={userId} refreshTrigger={refreshMemories} />}
        </aside>
        <section className="flex-1 p-4 overflow-y-auto">
          <ReportPanel />
        </section>
      </main>

      {showOnboarding && userId && (
        <OnboardingChat
          userId={userId}
          onComplete={handleOnboardingComplete}
        />
      )}
    </div>
  )
}
