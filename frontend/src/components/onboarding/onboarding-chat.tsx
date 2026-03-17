"use client"

import { useState, useEffect, useRef } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { apiClient, Memory } from "@/lib/api-client"
import { Send, Bot, Sparkles } from "lucide-react"

interface Message {
  role: "user" | "assistant"
  content: string
}

interface OnboardingChatProps {
  userId: string
  onComplete: (memories: Memory[]) => void
}

export function OnboardingChat({ userId, onComplete }: OnboardingChatProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [sessionId, setSessionId] = useState<string | null>(null)
  const [isComplete, setIsComplete] = useState(false)
  const [memories, setMemories] = useState<Memory[]>([])
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  useEffect(() => {
    startOnboarding()
  }, [])

  const startOnboarding = async () => {
    try {
      setIsLoading(true)
      const response = await apiClient.onboarding.start(userId)
      setSessionId(response.session_id)
      setMessages([{ role: "assistant", content: response.first_message }])
    } catch (error) {
      console.error("Failed to start onboarding:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || !sessionId || isLoading) return

    const userMessage = input.trim()
    setInput("")
    setMessages(prev => [...prev, { role: "user", content: userMessage }])
    setIsLoading(true)

    try {
      const response = await apiClient.onboarding.sendMessage(userId, sessionId, userMessage)

      setMessages(prev => [...prev, { role: "assistant", content: response.response }])

      if (response.is_complete && response.memories) {
        setIsComplete(true)
        setMemories(response.memories)
      }
    } catch (error) {
      console.error("Failed to send message:", error)
      setMessages(prev => [...prev, {
        role: "assistant",
        content: "Sorry, something went wrong. Please try again."
      }])
    } finally {
      setIsLoading(false)
    }
  }

  const handleConfirm = () => {
    onComplete(memories)
  }

  return (
    <div className="fixed inset-0 bg-background/95 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="w-full max-w-2xl h-[80vh] flex flex-col bg-card rounded-lg border shadow-lg">
        {/* Header */}
        <div className="p-4 border-b flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-full">
            <Bot className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="font-semibold">Welcome to News Subscriber</h2>
            <p className="text-sm text-muted-foreground">Let me learn about your interests</p>
          </div>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {messages.map((message, index) => (
            <div
              key={index}
              className={`flex ${message.role === "user" ? "justify-end" : "justify-start"}`}
            >
              <div
                className={`max-w-[80%] rounded-lg p-3 ${
                  message.role === "user"
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted"
                }`}
              >
                <p className="text-sm whitespace-pre-wrap">{message.content}</p>
              </div>
            </div>
          ))}

          {isLoading && !isComplete && (
            <div className="flex justify-start">
              <div className="bg-muted rounded-lg p-3">
                <p className="text-sm text-muted-foreground">Thinking...</p>
              </div>
            </div>
          )}

          {/* Memory cards when complete */}
          {isComplete && memories.length > 0 && (
            <div className="space-y-4 mt-4">
              <div className="flex items-center gap-2 text-primary">
                <Sparkles className="h-5 w-5" />
                <span className="font-medium">Here&apos;s what I learned about you:</span>
              </div>
              <div className="space-y-3">
                {memories.map((memory, index) => (
                  <Card key={index} className="border-primary/20">
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
              <Button onClick={handleConfirm} className="w-full">
                Looks great! Start exploring
              </Button>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        {/* Input */}
        {!isComplete && (
          <form onSubmit={sendMessage} className="p-4 border-t flex gap-2">
            <Input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Tell me about your interests..."
              disabled={isLoading}
              className="flex-1"
            />
            <Button type="submit" disabled={isLoading || !input.trim()}>
              <Send className="h-4 w-4" />
            </Button>
          </form>
        )}
      </div>
    </div>
  )
}
