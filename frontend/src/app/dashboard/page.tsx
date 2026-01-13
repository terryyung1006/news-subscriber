import { Header } from "@/components/dashboard/header"
import { PreferencePanel } from "@/components/dashboard/preference-panel"
import { ReportPanel } from "@/components/dashboard/report-panel"

export default function DashboardPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex flex-1 flex-col md:flex-row overflow-hidden">
        <aside className="w-full md:w-1/3 lg:w-1/4 border-r bg-muted/10 p-4 overflow-y-auto">
          <PreferencePanel />
        </aside>
        <section className="flex-1 p-4 overflow-y-auto">
          <ReportPanel />
        </section>
      </main>
    </div>
  )
}
