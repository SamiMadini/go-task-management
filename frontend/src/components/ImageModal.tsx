"use client"

import { useEffect } from "react"
import Image from "next/image"

interface ImageModalProps {
  image: {
    id: number
    fullSize: string
    alt: string
  }
  onClose: () => void
  onNext: () => void
  onPrevious: () => void
}

export default function ImageModal({ image, onClose, onNext, onPrevious }: ImageModalProps) {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose()
      if (e.key === "ArrowRight") onNext()
      if (e.key === "ArrowLeft") onPrevious()
    }

    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [onClose, onNext, onPrevious])

  return (
    <div className="fixed inset-0 bg-black bg-opacity-75 z-50 flex items-center justify-center p-4">
      <div className="relative max-w-4xl w-full h-[80vh] bg-white rounded-lg overflow-hidden">
        <div className="absolute top-4 right-4 z-10">
          <button onClick={onClose} className="bg-white rounded-full p-2 shadow-md hover:bg-gray-100">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="relative w-full h-full">
          <Image src={image.fullSize} alt={image.alt} fill className="object-contain" sizes="100vw" />
        </div>

        <div className="absolute bottom-4 left-0 right-0 flex justify-center gap-4">
          <button onClick={onPrevious} className="bg-white rounded-full p-3 shadow-md hover:bg-gray-100" aria-label="Previous image">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <button onClick={onNext} className="bg-white rounded-full p-3 shadow-md hover:bg-gray-100" aria-label="Next image">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  )
}
