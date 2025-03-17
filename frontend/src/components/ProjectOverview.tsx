export default function ProjectOverview({
  projectInfo,
}: {
  projectInfo: {
    name: string
    description: string
    services: string[]
    architecture: string
  }
}) {
  return (
    <section className="mb-12">
      <h2 className="text-2xl font-semibold mb-4">Project Overview</h2>
      <p className="mb-4">{projectInfo.description}</p>

      <div className="mt-6">
        <h3 className="text-xl font-medium mb-2">Services</h3>
        <ul className="list-disc pl-5">
          {projectInfo.services.map((service, index) => (
            <li key={index}>{service}</li>
          ))}
        </ul>
      </div>
    </section>
  )
}
