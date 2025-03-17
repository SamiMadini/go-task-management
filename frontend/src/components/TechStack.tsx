export default function TechStack({
  languages,
  database,
  technologies,
}: {
  languages: string[]
  database: string
  technologies: string[]
}) {
  return (
    <section className="mb-12">
      <h2 className="text-2xl font-semibold mb-4">Technology Stack</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-xl font-medium mb-2">Languages</h3>
          <div className="flex flex-wrap gap-2">
            {languages.map((lang, index) => (
              <span key={index} className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                {lang}
              </span>
            ))}
          </div>
        </div>

        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-xl font-medium mb-2">Database</h3>
          <span className="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm">{database}</span>
        </div>
      </div>

      <div className="mt-6 bg-gray-50 p-4 rounded-lg">
        <h3 className="text-xl font-medium mb-2">Technologies</h3>
        <div className="flex flex-wrap gap-2">
          {technologies.map((tech, index) => (
            <span key={index} className="px-3 py-1 bg-purple-100 text-purple-800 rounded-full text-sm">
              {tech}
            </span>
          ))}
        </div>
      </div>
    </section>
  )
}
