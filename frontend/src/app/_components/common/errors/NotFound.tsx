"use client"

import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import Link from "next/link"

export default function NotFoundComponent() {
  return (
    <section className="flex h-screen w-full flex-col items-center justify-center bg-black">
      <h1 className="text-9xl font-extrabold tracking-widest text-white">404</h1>
      <div className="absolute rotate-12 rounded bg-primary px-2 text-sm text-white">Page not found</div>
      <Link href="/" className={cn(buttonVariants({ variant: "default", size: "lg" }), "mt-16")}>
        Return to home
      </Link>
    </section>
  )
}
