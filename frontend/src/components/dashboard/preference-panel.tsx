"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Trash2, Edit2 } from "lucide-react"
import { useState } from "react"

interface Preference {
  id: string
  content: string
}

export function PreferencePanel() {
  const [preferences, setPreferences] = useState<Preference[]>([
    { id: "1", content: "Tech news about AI and LLMs" },
    { id: "2", content: "Global economic trends affecting Asia" },
  ])

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Your Interests</h2>
        <Button size="sm" variant="outline">
          <Plus className="mr-2 h-4 w-4" />
          Add New
        </Button>
      </div>
      <div className="space-y-3">
        {preferences.map((pref) => (
          <Card key={pref.id} className="relative group">
            <CardContent className="p-4">
              <p className="text-sm">{pref.content}</p>
              <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                <Button variant="ghost" size="icon" className="h-6 w-6">
                  <Edit2 className="h-3 w-3" />
                </Button>
                <Button variant="ghost" size="icon" className="h-6 w-6 text-destructive hover:text-destructive">
                  <Trash2 className="h-3 w-3" />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
