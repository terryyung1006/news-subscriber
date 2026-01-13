"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { RefreshCw, FileText, MessageSquare } from "lucide-react"
import { useState } from "react"
import { ChatInterface } from "./chat-interface"
import { cn } from "@/lib/utils"

export function ReportPanel() {
  const [activeTab, setActiveTab] = useState<"report" | "chat">("report")

  return (
    <div className="h-full flex flex-col space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex space-x-2 bg-muted/50 p-1 rounded-lg">
          <button
            onClick={() => setActiveTab("report")}
            className={cn(
              "flex items-center px-3 py-1.5 text-sm font-medium rounded-md transition-colors",
              activeTab === "report"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <FileText className="mr-2 h-4 w-4" />
            Report
          </button>
          <button
            onClick={() => setActiveTab("chat")}
            className={cn(
              "flex items-center px-3 py-1.5 text-sm font-medium rounded-md transition-colors",
              activeTab === "chat"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <MessageSquare className="mr-2 h-4 w-4" />
            Chat
          </button>
        </div>
        {activeTab === "report" && (
          <Button size="sm">
            <RefreshCw className="mr-2 h-4 w-4" />
            Generate Now
          </Button>
        )}
      </div>

      <Card className="flex-1 overflow-hidden">
        <CardContent className="h-full p-6">
          {activeTab === "report" ? (
            <div className="h-full overflow-y-auto">
              <div className="prose dark:prose-invert max-w-none">
                <h2 className="mt-0">Today's Digest</h2>
                <p>
                  Your personalized report will appear here. It will be generated based on the interests you've defined on the left.
                </p>
                <h3>AI Updates</h3>
                <ul>
                  <li>Google releases new Gemini model with enhanced coding capabilities.</li>
                  <li>OpenAI announces partnership with major news publishers.</li>
                </ul>
                <h3>Economic Trends</h3>
                <ul>
                  <li>Asian markets show resilience amidst global inflation concerns.</li>
                  <li>Tech sector drives growth in the region.</li>
                </ul>
              </div>
            </div>
          ) : (
            <ChatInterface />
          )}
        </CardContent>
      </Card>
    </div>
  )
}
