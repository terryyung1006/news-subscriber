import Link from "next/link"
import { Button } from "@/components/ui/button"
import { User, Settings, LogOut } from "lucide-react"

export function Header() {
  return (
    <header className="flex h-14 items-center gap-4 border-b bg-muted/40 px-6 lg:h-[60px]">
      <Link className="flex items-center gap-2 font-semibold" href="/dashboard">
        <span className="">News Subscriber</span>
      </Link>
      <div className="flex-1"></div>
      <Button variant="ghost" size="icon" className="rounded-full">
        <User className="h-5 w-5" />
        <span className="sr-only">Toggle user menu</span>
      </Button>
    </header>
  )
}
