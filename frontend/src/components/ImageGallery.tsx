"use client"

import { useState } from "react"
import Image from "next/image"
import ImageModal from "./ImageModal"

interface GalleryImage {
  id: number
  thumbnail: string
  fullSize: string
  alt: string
}

interface ImageGalleryProps {
  images: GalleryImage[]
}

export default function ImageGallery({ images }: ImageGalleryProps) {
  const [modalOpen, setModalOpen] = useState(false)
  const [currentImageIndex, setCurrentImageIndex] = useState(0)

  const openModal = (index: number) => {
    setCurrentImageIndex(index)
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
  }

  const goToNextImage = () => {
    setCurrentImageIndex((prev) => (prev + 1) % images.length)
  }

  const goToPrevImage = () => {
    setCurrentImageIndex((prev) => (prev - 1 + images.length) % images.length)
  }

  return (
    <div>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {images.map((image, index) => (
          <div key={image.id} className="relative overflow-hidden rounded-lg cursor-pointer aspect-square" onClick={() => openModal(index)}>
            <Image
              src={image.thumbnail}
              alt={image.alt}
              fill
              className="object-cover hover:scale-105 transition-transform duration-300"
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            />
          </div>
        ))}
      </div>

      {modalOpen && <ImageModal image={images[currentImageIndex]} onClose={closeModal} onNext={goToNextImage} onPrevious={goToPrevImage} />}
    </div>
  )
}
