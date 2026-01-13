import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4 text-center">
      <div className="max-w-2xl space-y-6">
        <h1 className="text-4xl font-bold tracking-tighter sm:text-6xl md:text-7xl">
          Your Personalized <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-gray-500">
            AI News Feed
          </span>
        </h1>
        <p className="mx-auto max-w-[600px] text-muted-foreground md:text-xl">
          Curate your daily digest with precision. Powered by AI, tailored for you.
        </p>
        <div className="flex justify-center gap-4">
          <Link href="/login">
            <Button size="lg" className="h-12 px-8 text-base">
              Get Started
            </Button>
          </Link>
          <Link href="/about">
            <Button variant="outline" size="lg" className="h-12 px-8 text-base">
              Learn More
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}


