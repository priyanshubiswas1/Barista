import { useNavigate } from "react-router-dom";
import { useState } from "react";

const Home = () => {
  const navigate = useNavigate();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const features = [
    {
      image: "/images/menu_feature.png",
      title: "Menu Management",
      description:
        "Create and manage your café menu with ease. Add items, set prices, and organize categories effortlessly.",
    },
    {
      image: "/images/table_feature.png",
      title: "Table Management",
      description:
        "Organize tables and generate unique QR codes for seamless customer ordering experience.",
    },
    {
      image: "/images/order_feature.png",
      title: "Order Tracking",
      description:
        "Real-time order management system. Track orders from placement to completion.",
    },
    {
      image: "/images/analytics_feature.png",
      title: "Analytics",
      description:
        "Track sales, monitor trends, and gain valuable customer insights to grow your business.",
    },
  ];

  const stats = [
    { number: "Easy", label: "Menu Setup" },
    { number: "Simple", label: "QR Ordering" },
    { number: "99.9%", label: "Uptime" },
    { number: "24/7", label: "Support" },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-indigo-50/30 to-gray-100 overflow-hidden">
      {/* Navigation Bar */}
      <nav className="w-full px-4 sm:px-6 py-3 sm:py-4 flex justify-between items-center bg-white/80 backdrop-blur-md shadow-sm fixed top-0 z-50 border-b border-gray-100/50">
        <div className="flex items-center gap-2 sm:gap-3">
          <div className="bg-gradient-to-br from-indigo-600 to-purple-600 rounded-lg sm:rounded-xl h-9 w-9 sm:h-11 sm:w-11 flex items-center justify-center shadow-lg shadow-indigo-200/50">
            <span className="text-lg sm:text-xl">☕</span>
          </div>
          <span className="font-bold text-xl sm:text-2xl bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">
            Barista
          </span>
        </div>

        {/* Desktop Navigation */}
        <div className="hidden sm:flex gap-3">
          <button
            onClick={() => navigate("/login")}
            className="px-4 sm:px-6 py-2 sm:py-2.5 text-indigo-700 font-semibold hover:bg-indigo-50 rounded-xl transition-all duration-300 cursor-pointer border border-transparent hover:border-indigo-100"
          >
            Sign In
          </button>
          <button
            onClick={() => navigate("/register")}
            className="px-4 sm:px-6 py-2 sm:py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl shadow-lg shadow-indigo-300/50 hover:shadow-xl hover:shadow-indigo-400/50 hover:from-indigo-700 hover:to-purple-700 transition-all duration-300 cursor-pointer transform hover:scale-105"
          >
            Get Started
          </button>
        </div>

        {/* Mobile Menu Button */}
        <button
          className="sm:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        >
          {mobileMenuOpen ? (
            <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          ) : (
            <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          )}
        </button>
      </nav>

      {/* Mobile Menu Dropdown */}
      {mobileMenuOpen && (
        <div className="fixed top-14 left-0 right-0 bg-white/95 backdrop-blur-md shadow-lg z-40 sm:hidden border-b border-gray-100">
          <div className="flex flex-col p-4 gap-3">
            <button
              onClick={() => {
                navigate("/login");
                setMobileMenuOpen(false);
              }}
              className="w-full py-3 text-indigo-700 font-semibold hover:bg-indigo-50 rounded-xl transition-all duration-300 cursor-pointer border border-indigo-100"
            >
              Sign In
            </button>
            <button
              onClick={() => {
                navigate("/register");
                setMobileMenuOpen(false);
              }}
              className="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl shadow-lg cursor-pointer"
            >
              Get Started
            </button>
          </div>
        </div>
      )}

      {/* Hero Section */}
      <section className="pt-24 sm:pt-32 pb-10 sm:pb-16 px-4 sm:px-6 relative">
        {/* Animated Background Elements - Hidden on mobile for performance */}
        <div className="hidden sm:block absolute top-20 right-10 w-96 h-96 bg-gradient-to-br from-indigo-200/50 to-purple-200/40 rounded-full blur-3xl animate-pulse" />
        <div className="hidden sm:block absolute bottom-10 left-10 w-80 h-80 bg-gradient-to-br from-purple-200/40 to-pink-200/30 rounded-full blur-3xl" />
        <div className="hidden sm:block absolute top-60 left-1/4 w-64 h-64 bg-gradient-to-br from-indigo-100/60 to-blue-100/50 rounded-full blur-2xl" />

        <div className="max-w-7xl mx-auto relative z-10">
          <div className="grid lg:grid-cols-2 gap-8 lg:gap-12 items-center">
            {/* Left Content */}
            <div className="text-center lg:text-left order-2 lg:order-1">
              <div className="inline-flex items-center gap-2 bg-white/90 backdrop-blur-sm rounded-full px-3 sm:px-5 py-2 sm:py-2.5 shadow-lg shadow-gray-200/50 mb-6 sm:mb-8 border border-gray-100">
                <span className="text-xl sm:text-2xl animate-bounce">☕</span>
                <span className="text-gray-700 font-semibold text-sm sm:text-base">
                  Café Management Made Simple
                </span>
              </div>

              <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-bold text-gray-900 mb-4 sm:mb-6 leading-tight">
                Manage Your
                <br />
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-purple-600 to-indigo-600 animate-gradient">
                  Café Like a Pro
                </span>
              </h1>

              <p className="text-base sm:text-lg md:text-xl text-gray-600 mb-8 sm:mb-10 leading-relaxed max-w-xl mx-auto lg:mx-0">
                Barista helps café owners streamline operations with powerful
                tools for menu management, order tracking, and customer
                engagement.
              </p>

              <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center lg:justify-start">
                <button
                  onClick={() => navigate("/register")}
                  className="group px-6 sm:px-8 py-3 sm:py-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-bold rounded-xl sm:rounded-2xl shadow-xl shadow-indigo-300/50 hover:shadow-2xl hover:shadow-indigo-400/60 hover:from-indigo-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105 cursor-pointer flex items-center justify-center gap-2"
                >
                  Start Free Trial
                  <svg
                    className="w-5 h-5 group-hover:translate-x-1 transition-transform"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M17 8l4 4m0 0l-4 4m4-4H3"
                    />
                  </svg>
                </button>
                <button
                  onClick={() => navigate("/login")}
                  className="px-6 sm:px-8 py-3 sm:py-4 bg-white text-indigo-700 font-bold rounded-xl sm:rounded-2xl shadow-lg shadow-gray-200/50 border-2 border-indigo-100 hover:border-indigo-200 hover:shadow-xl transition-all duration-300 cursor-pointer"
                >
                  Sign In
                </button>
              </div>
            </div>

            {/* Right Content - Hero Illustration */}
            <div className="relative order-1 lg:order-2">
              <div className="absolute inset-0 bg-gradient-to-br from-indigo-100 to-purple-100 rounded-2xl sm:rounded-3xl transform rotate-2 sm:rotate-3 scale-105" />
              <div className="relative bg-white rounded-2xl sm:rounded-3xl shadow-xl sm:shadow-2xl shadow-indigo-200/50 overflow-hidden border border-gray-100">
                <img
                  src="/images/hero_illustration.png"
                  alt="Café Management"
                  className="w-full h-auto transform hover:scale-105 transition-transform duration-500"
                />
              </div>
              {/* Floating Elements - Hidden on small screens */}
              <div className="hidden md:flex absolute -top-4 -right-4 bg-white rounded-2xl shadow-xl p-3 sm:p-4 border border-gray-100 animate-float">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse" />
                  <span className="text-xs sm:text-sm font-semibold text-gray-700">
                    Live Orders
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-8 sm:py-12 px-4 sm:px-6">
        <div className="max-w-6xl mx-auto">
          <div className="bg-white/70 backdrop-blur-sm rounded-2xl sm:rounded-3xl shadow-xl shadow-gray-200/50 border border-gray-100 p-4 sm:p-8">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-8">
              {stats.map((stat, index) => (
                <div key={index} className="text-center">
                  <div className="text-2xl sm:text-4xl md:text-5xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-1 sm:mb-2">
                    {stat.number}
                  </div>
                  <div className="text-gray-600 font-medium text-xs sm:text-base">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-12 sm:py-20 px-4 sm:px-6 relative">
        <div className="hidden sm:block absolute top-0 left-1/2 w-96 h-96 bg-gradient-to-br from-purple-100/40 to-indigo-100/30 rounded-full blur-3xl -translate-x-1/2" />

        <div className="max-w-7xl mx-auto relative z-10">
          <div className="text-center mb-10 sm:mb-16">
            <span className="inline-block bg-indigo-100 text-indigo-700 font-semibold px-3 sm:px-4 py-1.5 sm:py-2 rounded-full text-xs sm:text-sm mb-3 sm:mb-4">
              Powerful Features
            </span>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold text-gray-900 mb-3 sm:mb-4">
              Everything You Need
            </h2>
            <p className="text-gray-600 text-sm sm:text-lg max-w-2xl mx-auto px-4">
              Powerful features designed specifically for café owners to
              streamline operations
            </p>
          </div>

          <div className="grid sm:grid-cols-2 gap-4 sm:gap-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className="group bg-white/80 backdrop-blur-sm rounded-2xl sm:rounded-3xl overflow-hidden shadow-lg sm:shadow-xl shadow-gray-200/50 hover:shadow-2xl hover:shadow-indigo-200/50 transition-all duration-500 transform hover:-translate-y-2 border border-gray-100"
              >
                <div className="p-4 sm:p-6 flex flex-col sm:flex-row gap-4 sm:gap-6 items-center">
                  <div className="w-20 h-20 sm:w-32 sm:h-32 flex-shrink-0 bg-gradient-to-br from-indigo-50 to-purple-50 rounded-xl sm:rounded-2xl overflow-hidden border border-indigo-100/50">
                    <img
                      src={feature.image}
                      alt={feature.title}
                      className="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-500"
                    />
                  </div>
                  <div className="text-center sm:text-left">
                    <h3 className="font-bold text-gray-900 text-lg sm:text-xl mb-1 sm:mb-2">
                      {feature.title}
                    </h3>
                    <p className="text-gray-600 text-sm sm:text-base">{feature.description}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-12 sm:py-20 px-4 sm:px-6 bg-gradient-to-b from-white/50 to-transparent">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-10 sm:mb-16">
            <span className="inline-block bg-purple-100 text-purple-700 font-semibold px-3 sm:px-4 py-1.5 sm:py-2 rounded-full text-xs sm:text-sm mb-3 sm:mb-4">
              Simple Process
            </span>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold text-gray-900 mb-3 sm:mb-4">
              How It Works
            </h2>
          </div>

          <div className="grid sm:grid-cols-3 gap-6 sm:gap-8">
            {[
              {
                step: "01",
                title: "Sign Up",
                description: "Create your account in less than 2 minutes",
                icon: "✨",
              },
              {
                step: "02",
                title: "Setup Menu",
                description: "Add your menu items, categories, and prices",
                icon: "📝",
              },
              {
                step: "03",
                title: "Start Serving",
                description: "Generate QR codes and start taking orders",
                icon: "🚀",
              },
            ].map((item, index) => (
              <div key={index} className="relative">
                <div className="bg-white rounded-2xl sm:rounded-3xl p-6 sm:p-8 shadow-lg sm:shadow-xl shadow-gray-200/50 border border-gray-100 text-center hover:shadow-2xl transition-all duration-300">
                  <div className="absolute -top-3 sm:-top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-xs sm:text-sm font-bold px-3 sm:px-4 py-1 rounded-full">
                    Step {item.step}
                  </div>
                  <div className="text-3xl sm:text-4xl mb-3 sm:mb-4 mt-3 sm:mt-4">{item.icon}</div>
                  <h3 className="text-lg sm:text-xl font-bold text-gray-900 mb-1 sm:mb-2">
                    {item.title}
                  </h3>
                  <p className="text-gray-600 text-sm sm:text-base">{item.description}</p>
                </div>
                {index < 2 && (
                  <div className="hidden md:block absolute top-1/2 -right-4 transform translate-x-1/2">
                    <svg
                      className="w-8 h-8 text-indigo-300"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M17 8l4 4m0 0l-4 4m4-4H3"
                      />
                    </svg>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-12 sm:py-20 px-4 sm:px-6">
        <div className="max-w-4xl mx-auto">
          <div className="bg-gradient-to-br from-indigo-600 via-purple-600 to-indigo-700 rounded-2xl sm:rounded-[2.5rem] p-8 sm:p-12 md:p-16 text-center relative overflow-hidden shadow-2xl shadow-indigo-300/50">
            {/* Decorative circles */}
            <div className="absolute top-0 right-0 w-32 sm:w-48 h-32 sm:h-48 bg-white/10 rounded-full -mr-16 sm:-mr-24 -mt-16 sm:-mt-24" />
            <div className="absolute bottom-0 left-0 w-28 sm:w-40 h-28 sm:h-40 bg-white/10 rounded-full -ml-14 sm:-ml-20 -mb-14 sm:-mb-20" />
            <div className="hidden sm:block absolute top-1/2 right-10 w-28 h-28 bg-white/5 rounded-full" />
            <div className="hidden sm:block absolute top-1/4 left-10 w-20 h-20 bg-white/5 rounded-full" />

            <div className="relative z-10">
              <div className="bg-white/20 backdrop-blur-sm rounded-full h-14 w-14 sm:h-20 sm:w-20 mx-auto mb-6 sm:mb-8 flex items-center justify-center shadow-xl ring-4 ring-white/30">
                <span className="text-2xl sm:text-4xl">☕</span>
              </div>
              <h2 className="text-2xl sm:text-4xl md:text-5xl font-bold text-white mb-4 sm:mb-6">
                Ready to Transform Your Café?
              </h2>
              <p className="text-indigo-100 text-sm sm:text-lg mb-8 sm:mb-10 max-w-lg mx-auto px-4">
                Join thousands of café owners already using Barista to grow
                their business and delight customers.
              </p>
              <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center">
                <button
                  onClick={() => navigate("/register")}
                  className="px-8 sm:px-10 py-3 sm:py-4 bg-white text-indigo-700 font-bold rounded-xl sm:rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105 cursor-pointer"
                >
                  Get Started for Free
                </button>
                <button
                  onClick={() => navigate("/login")}
                  className="px-8 sm:px-10 py-3 sm:py-4 bg-white/10 backdrop-blur-sm text-white font-bold rounded-xl sm:rounded-2xl border-2 border-white/30 hover:bg-white/20 transition-all duration-300 cursor-pointer"
                >
                  Sign In
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-6 sm:py-10 px-4 sm:px-6 border-t border-gray-200/50 bg-white/70 backdrop-blur-sm">
        <div className="max-w-6xl mx-auto flex flex-col gap-4 sm:gap-6 items-center">
          <div className="flex items-center gap-2 sm:gap-3">
            <div className="bg-gradient-to-br from-indigo-600 to-purple-600 rounded-lg sm:rounded-xl h-8 w-8 sm:h-10 sm:w-10 flex items-center justify-center shadow-lg">
              <span className="text-sm sm:text-lg">☕</span>
            </div>
            <span className="font-bold text-lg sm:text-xl text-gray-800">Barista</span>
          </div>
          <div className="flex gap-6 sm:gap-8">
            <button
              onClick={() => navigate("/login")}
              className="text-gray-600 hover:text-indigo-600 font-medium text-sm sm:text-base transition-colors cursor-pointer"
            >
              Sign In
            </button>
            <button
              onClick={() => navigate("/register")}
              className="text-gray-600 hover:text-indigo-600 font-medium text-sm sm:text-base transition-colors cursor-pointer"
            >
              Sign Up
            </button>
          </div>
          <p className="text-gray-500 text-xs sm:text-sm text-center">
            © 2026 Barista. All rights reserved. Made with ❤️ for café owners.
          </p>
        </div>
      </footer>

      {/* Custom CSS for animations */}
      <style>{`
        @keyframes float {
          0%, 100% { transform: translateY(0px); }
          50% { transform: translateY(-10px); }
        }
        @keyframes float-delayed {
          0%, 100% { transform: translateY(0px); }
          50% { transform: translateY(-8px); }
        }
        @keyframes gradient {
          0% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
          100% { background-position: 0% 50%; }
        }
        .animate-float {
          animation: float 3s ease-in-out infinite;
        }
        .animate-float-delayed {
          animation: float-delayed 3s ease-in-out infinite 0.5s;
        }
        .animate-gradient {
          background-size: 200% 200%;
          animation: gradient 3s ease infinite;
        }
      `}</style>
    </div>
  );
};

export default Home;
