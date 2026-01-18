export function WhoItsFor() {
  const audiences = [
    {
      title: "Content Creators",
      description: "Share your creativity across all platforms without the hassle of multiple uploads.",
      emoji: "🎨",
      gradient: "from-pink-500 to-rose-500"
    },
    {
      title: "Influencers",
      description: "Maximize your reach and engagement by being everywhere your audience is.",
      emoji: "⭐",
      gradient: "from-purple-500 to-indigo-500"
    },
    {
      title: "Small Businesses",
      description: "Promote your brand consistently across all social channels with minimal effort.",
      emoji: "🚀",
      gradient: "from-blue-500 to-cyan-500"
    },
    {
      title: "Social Media Managers",
      description: "Manage multiple clients and accounts efficiently from a single dashboard.",
      emoji: "📊",
      gradient: "from-emerald-500 to-teal-500"
    },
    {
      title: "Marketing Agencies",
      description: "Scale your content operations and deliver better results for your clients.",
      emoji: "💼",
      gradient: "from-amber-500 to-orange-500"
    },
    {
      title: "Educators & Coaches",
      description: "Share your knowledge and lessons with students across all platforms simultaneously.",
      emoji: "📚",
      gradient: "from-violet-500 to-purple-500"
    }
  ];

  return (
    <section className="py-24 px-6 bg-gradient-to-b from-gray-50 to-white">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
            Built For Everyone
          </h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Whether you're just starting out or managing multiple accounts, elci.io scales with you
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {audiences.map((audience, index) => (
            <div
              key={index}
              className="relative group overflow-hidden rounded-2xl bg-white p-8 border-2 border-gray-100 hover:border-transparent transition-all duration-300"
            >
              {/* Gradient background on hover */}
              <div className={`absolute inset-0 bg-gradient-to-br ${audience.gradient} opacity-0 group-hover:opacity-5 transition-opacity duration-300`}></div>

              <div className="relative">
                <div className="text-4xl mb-4">{audience.emoji}</div>

                <h3 className="text-xl font-semibold text-gray-900 mb-3">
                  {audience.title}
                </h3>

                <p className="text-gray-600 leading-relaxed">
                  {audience.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
