export function HowItWorks() {
  const steps = [
    {
      number: "01",
      title: "Upload Once",
      description: "Drag and drop your video or select from your device. One upload is all it takes.",
      icon: (
        <svg className="w-12 h-12" viewBox="0 0 48 48" fill="none">
          <rect x="8" y="12" width="32" height="24" rx="3" stroke="url(#step-gradient-1)" strokeWidth="2.5" fill="none" />
          <path d="M24 20 L24 28 M20 24 L24 20 L28 24" stroke="url(#step-gradient-1)" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
          <defs>
            <linearGradient id="step-gradient-1" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#4F46E5" />
              <stop offset="100%" stopColor="#7C3AED" />
            </linearGradient>
          </defs>
        </svg>
      )
    },
    {
      number: "02",
      title: "Connect Platforms",
      description: "Link your social media accounts securely. Set it up once and you're ready to go.",
      icon: (
        <svg className="w-12 h-12" viewBox="0 0 48 48" fill="none">
          <circle cx="14" cy="24" r="6" stroke="url(#step-gradient-2)" strokeWidth="2.5" fill="none" />
          <circle cx="34" cy="16" r="6" stroke="url(#step-gradient-2)" strokeWidth="2.5" fill="none" />
          <circle cx="34" cy="32" r="6" stroke="url(#step-gradient-2)" strokeWidth="2.5" fill="none" />
          <path d="M20 24 L28 18 M20 24 L28 30" stroke="url(#step-gradient-2)" strokeWidth="2.5" strokeLinecap="round" />
          <defs>
            <linearGradient id="step-gradient-2" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#4F46E5" />
              <stop offset="100%" stopColor="#7C3AED" />
            </linearGradient>
          </defs>
        </svg>
      )
    },
    {
      number: "03",
      title: "Publish Everywhere",
      description: "Hit publish and watch your content go live across all platforms simultaneously.",
      icon: (
        <svg className="w-12 h-12" viewBox="0 0 48 48" fill="none">
          <circle cx="24" cy="24" r="4" fill="url(#step-gradient-3)" />
          <circle cx="24" cy="24" r="10" stroke="url(#step-gradient-3)" strokeWidth="2.5" fill="none" opacity="0.5" />
          <circle cx="24" cy="24" r="16" stroke="url(#step-gradient-3)" strokeWidth="2.5" fill="none" opacity="0.3" />
          <path d="M24 8 L24 12 M24 36 L24 40 M8 24 L12 24 M36 24 L40 24" stroke="url(#step-gradient-3)" strokeWidth="2.5" strokeLinecap="round" />
          <defs>
            <linearGradient id="step-gradient-3" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#4F46E5" />
              <stop offset="100%" stopColor="#7C3AED" />
            </linearGradient>
          </defs>
        </svg>
      )
    }
  ];

  return (
    <section id="how-it-works" className="py-24 px-6 bg-gradient-to-b from-white to-gray-50">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
            How It Works
          </h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Three simple steps to transform your content distribution workflow
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8 lg:gap-12">
          {steps.map((step, index) => (
            <div key={index} className="relative">
              {/* Connector line */}
              {index < steps.length - 1 && (
                <div className="hidden md:block absolute top-16 left-[calc(50%+40px)] w-[calc(100%-80px)] h-0.5 bg-gradient-to-r from-indigo-200 to-purple-200"></div>
              )}

              <div className="relative bg-white rounded-2xl p-8 shadow-sm hover:shadow-xl transition-shadow border border-gray-100">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-50 to-purple-50 mb-6">
                  {step.icon}
                </div>

                <div className="absolute top-6 right-6 text-5xl font-bold text-gray-100">
                  {step.number}
                </div>

                <h3 className="text-2xl font-semibold text-gray-900 mb-3">
                  {step.title}
                </h3>

                <p className="text-gray-600 leading-relaxed">
                  {step.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
