"use client"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { useRouter } from "next/navigation"

interface ConfirmDialogProps {
  title: string
  description: string | React.ReactNode
  open: boolean
  onConfirm: () => void
  onClose?: () => void
}

export function ConfirmDialog({ open, onClose, onConfirm, title, description }: ConfirmDialogProps) {
  const router = useRouter()

  const closeModal = () => {
    router.back()
  }

  return (
    <AlertDialog
      open={open}
      onOpenChange={(isOpen: boolean) => {
        if (!isOpen) {
          if (onClose) {
            onClose()
          } else {
            closeModal()
          }
        }
      }}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onClose}>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={() => onConfirm()}>Continue</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
