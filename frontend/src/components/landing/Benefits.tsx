export function Benefits() {
  const benefits = [
    {
      title: "Save Hours Every Week",
      description: "Stop uploading the same video multiple times. Reclaim your time and focus on creating great content.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <circle cx="16" cy="16" r="12" stroke="currentColor" strokeWidth="2" fill="none" />
          <path d="M16 8 L16 16 L22 16" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
      color: "from-indigo-500 to-purple-500"
    },
    {
      title: "Consistent Posting",
      description: "Maintain a unified presence across all platforms. Same quality, same timing, everywhere.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <rect x="4" y="8" width="8" height="8" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <rect x="4" y="20" width="8" height="8" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <rect x="16" y="8" width="8" height="8" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <rect x="16" y="20" width="8" height="8" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <circle cx="28" cy="12" r="2" fill="currentColor" />
          <circle cx="28" cy="24" r="2" fill="currentColor" />
        </svg>
      ),
      color: "from-purple-500 to-pink-500"
    },
    {
      title: "No Manual Re-uploads",
      description: "Eliminate the tedious task of re-uploading videos. One click does it all, automatically.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <path d="M8 24 L24 24 M16 20 L16 8 M12 12 L16 8 L20 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
          <circle cx="8" cy="16" r="2" fill="currentColor" opacity="0.5" />
          <circle cx="24" cy="16" r="2" fill="currentColor" opacity="0.5" />
        </svg>
      ),
      color: "from-blue-500 to-indigo-500"
    },
    {
      title: "Reach More Audiences",
      description: "Maximize your visibility by being present on every major platform where your audience lives.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <circle cx="16" cy="16" r="3" fill="currentColor" />
          <circle cx="8" cy="8" r="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <circle cx="24" cy="8" r="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <circle cx="8" cy="24" r="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <circle cx="24" cy="24" r="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <path d="M16 16 L8 8 M16 16 L24 8 M16 16 L8 24 M16 16 L24 24" stroke="currentColor" strokeWidth="1.5" opacity="0.3" />
        </svg>
      ),
      color: "from-emerald-500 to-teal-500"
    },
    {
      title: "Analytics Dashboard",
      description: "Track performance across all platforms in one place. Make data-driven decisions effortlessly.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <rect x="6" y="18" width="4" height="8" rx="1" fill="currentColor" />
          <rect x="14" y="12" width="4" height="14" rx="1" fill="currentColor" />
          <rect x="22" y="8" width="4" height="18" rx="1" fill="currentColor" />
          <path d="M6 8 L14 12 L22 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" opacity="0.3" />
        </svg>
      ),
      color: "from-amber-500 to-orange-500"
    },
    {
      title: "Schedule in Advance",
      description: "Plan your content calendar weeks ahead. Set it and forget it with automated publishing.",
      icon: (
        <svg className="w-8 h-8" viewBox="0 0 32 32" fill="none">
          <rect x="6" y="8" width="20" height="18" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
          <path d="M6 12 L26 12 M12 6 L12 10 M20 6 L20 10" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          <circle cx="12" cy="17" r="1" fill="currentColor" />
          <circle cx="16" cy="17" r="1" fill="currentColor" />
          <circle cx="20" cy="17" r="1" fill="currentColor" />
          <circle cx="12" cy="21" r="1" fill="currentColor" />
          <circle cx="16" cy="21" r="1" fill="currentColor" />
        </svg>
      ),
      color: "from-rose-500 to-pink-500"
    }
  ];

  return (
    <section id="benefits" className="py-24 px-6 bg-white">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
            Why Creators Love elci.io
          </h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Everything you need to supercharge your content distribution
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {benefits.map((benefit, index) => (
            <div
              key={index}
              className="group p-8 rounded-2xl border-2 border-gray-100 hover:border-transparent hover:shadow-2xl transition-all duration-300"
            >
              <div className={`inline-flex items-center justify-center w-14 h-14 rounded-xl bg-gradient-to-br ${benefit.color} text-white mb-5 group-hover:scale-110 transition-transform`}>
                {benefit.icon}
              </div>

              <h3 className="text-xl font-semibold text-gray-900 mb-3">
                {benefit.title}
              </h3>

              <p className="text-gray-600 leading-relaxed">
                {benefit.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
