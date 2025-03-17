import ProjectOverview from "../components/ProjectOverview"
import TechStack from "../components/TechStack"
import ImageGallery from "../components/ImageGallery"

interface ProjectInfo {
  name: string
  description: string
  services: string[]
  languages: string[]
  database: string
  architecture: string
  technologies: string[]
}

interface GalleryImage {
  id: number
  thumbnail: string
  fullSize: string
  alt: string
}

export default async function Page() {
  const projectInfo: ProjectInfo = {
    name: "Task Management System",
    description: "A simple task management with events history.",
    services: ["Task Management Frontend", "API gateway", "Notification Service", "Email service"],
    languages: ["Golang", "TypeScript"],
    database: "PostgreSQL",
    architecture: "Microservices with API Gateway",
    technologies: ["Docker", "gRPC", "Next.js", "Tailwind CSS"],
  }

  const galleryImages: GalleryImage[] = [
    { id: 1, thumbnail: "/images/thumbnail1.png", fullSize: "/images/full1.png", alt: "System Architecture Diagram" },
    { id: 2, thumbnail: "/images/thumbnail2.png", fullSize: "/images/full2.png", alt: "Database Schema" },
    { id: 3, thumbnail: "/images/thumbnail3.png", fullSize: "/images/full3.png", alt: "User Interface" },
  ]

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-center mb-12">Task Management System</h1>

      <ProjectOverview projectInfo={projectInfo} />

      <TechStack languages={projectInfo.languages} database={projectInfo.database} technologies={projectInfo.technologies} />

      <section className="mt-12">
        <h2 className="text-2xl font-semibold mb-6">Project Gallery</h2>
        <p className="mb-4">Click on any image to view in full size.</p>
        <ImageGallery images={galleryImages} />
      </section>
    </div>
  )
}
