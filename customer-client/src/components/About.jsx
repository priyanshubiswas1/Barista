const About = () => {
  return (
    <div className="w-full flex justify-center items-center px-4 py-4">
      <div className="w-full max-w-md bg-gradient-to-br from-emerald-700 to-emerald-800 text-white rounded-2xl p-6 shadow-xl relative overflow-hidden">
        <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
        <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />
        <div className="relative z-10">
          <div className="bg-amber-500 inline-block rounded-lg px-3 py-1 text-xs text-black font-semibold mb-3 shadow-md">
            About Us
          </div>
          <h2 className="text-xl font-bold mb-2">Welcome</h2>
          <p className="text-emerald-50 text-sm leading-relaxed">
            Experience the finest dining with our curated menu. Scan, order, and
            enjoy!
          </p>
        </div>
      </div>
    </div>
  );
};

export default About;
