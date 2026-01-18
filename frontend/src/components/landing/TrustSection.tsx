export function TrustSection() {
  return (
    <section className="py-24 px-6 bg-white">
      <div className="max-w-7xl mx-auto">
        <div className="grid lg:grid-cols-2 gap-16 items-center">
          {/* Left - UI Mockup */}
          <div className="relative">
            {/* Decorative background */}
            <div className="absolute inset-0 bg-gradient-to-br from-indigo-100 to-purple-100 rounded-3xl transform rotate-3"></div>

            {/* Main mockup */}
            <div className="relative bg-white rounded-2xl shadow-2xl overflow-hidden border border-gray-200">
              {/* Browser bar */}
              <div className="bg-gray-100 px-4 py-3 flex items-center gap-2 border-b border-gray-200">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-red-400"></div>
                  <div className="w-3 h-3 rounded-full bg-yellow-400"></div>
                  <div className="w-3 h-3 rounded-full bg-green-400"></div>
                </div>
                <div className="flex-1 ml-4">
                  <div className="bg-white rounded px-3 py-1 text-xs text-gray-500 inline-block">
                    app.elci.io/dashboard
                  </div>
                </div>
              </div>

              {/* Dashboard content */}
              <div className="p-6 space-y-4">
                {/* Header */}
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold text-gray-900">Upload Video</h3>
                  <div className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium">
                    Ready
                  </div>
                </div>

                {/* Upload area */}
                <div className="border-2 border-dashed border-indigo-200 rounded-xl p-8 bg-gradient-to-br from-indigo-50 to-purple-50 hover:border-indigo-300 transition-colors cursor-pointer">
                  <div className="text-center">
                    <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 text-white mb-3">
                      <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                      </svg>
                    </div>
                    <p className="font-medium text-gray-700 mb-1">Drop your video here</p>
                    <p className="text-sm text-gray-500">or click to browse</p>
                  </div>
                </div>

                {/* Platform selection */}
                <div>
                  <p className="text-sm font-medium text-gray-700 mb-3">Select platforms</p>
                  <div className="grid grid-cols-2 gap-2">
                    {['TikTok', 'Instagram', 'YouTube', 'Twitter'].map((platform, i) => (
                      <div
                        key={platform}
                        className={`flex items-center gap-2 px-3 py-2 rounded-lg border-2 ${
                          i < 3
                            ? 'border-indigo-500 bg-indigo-50'
                            : 'border-gray-200 bg-white'
                        } transition-colors`}
                      >
                        <div className={`w-4 h-4 rounded ${
                          i < 3
                            ? 'bg-gradient-to-br from-indigo-500 to-purple-500'
                            : 'bg-gray-200'
                        } flex items-center justify-center`}>
                          {i < 3 && (
                            <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                            </svg>
                          )}
                        </div>
                        <span className="text-sm font-medium text-gray-700">{platform}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Publish button */}
                <button className="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg font-medium hover:shadow-lg transition-shadow">
                  Publish to 3 Platforms
                </button>
              </div>
            </div>
          </div>

          {/* Right - Content */}
          <div className="space-y-8">
            <div>
              <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
                Simple. Fast. Reliable.
              </h2>
              <p className="text-xl text-gray-600 leading-relaxed">
                Our intuitive interface makes cross-platform publishing effortless. No technical knowledge required.
              </p>
            </div>

            <div className="space-y-6">
              {[
                {
                  title: "Lightning Fast",
                  description: "Upload and publish in seconds, not hours. Our optimized infrastructure ensures quick delivery."
                },
                {
                  title: "Secure & Private",
                  description: "Your content and credentials are encrypted and protected. We never share your data."
                },
                {
                  title: "Always Available",
                  description: "99.9% uptime guaranteed. Publish whenever inspiration strikes, day or night."
                }
              ].map((item, index) => (
                <div key={index} className="flex gap-4">
                  <div className="flex-shrink-0">
                    <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center">
                      <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                    </div>
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 mb-1">
                      {item.title}
                    </h3>
                    <p className="text-gray-600">
                      {item.description}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
