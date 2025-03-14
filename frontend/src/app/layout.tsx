import type { Metadata } from "next"
import "../styles/globals.css"
import localFont from "next/font/local"
import MainLayout from "@/app/_components/layout/MainLayout"
import { cookies } from 'next/headers';

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
})
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
})

export const metadata: Metadata = {
  title: "Task management",
  description: "A simple task management app",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  const cookieStore = cookies();
  const isExpanded = cookieStore.get('app:sidebarExpanded') !== undefined ? cookieStore.get('app:sidebarExpanded')?.value === 'true' : true;

  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <MainLayout isExpanded={isExpanded}>{children}</MainLayout>
      </body>
    </html>
  )
}
