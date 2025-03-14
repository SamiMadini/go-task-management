"use client"

import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet"
import { useRouter } from "next/navigation"

const SHEET_SIDES = ["right"] as const

type SheetRightProps = {
  open: boolean
  title: string
  description: string
  overlay: boolean
  children: React.ReactNode
  isModal?: boolean
  setIsOpen?: (open: boolean) => void
}

export function SheetRight({ open, title, description, overlay, children, isModal, setIsOpen }: SheetRightProps) {
  const router = useRouter()

  const closeSheet = () => {
    router.back()
  }

  return (
    <div className="flex">
      {SHEET_SIDES.map((side) => (
        <Sheet
          key={side}
          open={open}
          onOpenChange={(isOpen) => {
            if (!isOpen) {
              if (setIsOpen) {
                setIsOpen(false)
              } else {
                closeSheet()
              }
            }
          }}
          modal={isModal !== undefined ? isModal : true}
        >
          <SheetContent side={side} overlay={overlay} className="overflow-y-scroll">
            <SheetHeader>
              <SheetTitle>{title}</SheetTitle>
              <SheetDescription>{description}</SheetDescription>
            </SheetHeader>
            {children}
          </SheetContent>
        </Sheet>
      ))}
    </div>
  )
}
