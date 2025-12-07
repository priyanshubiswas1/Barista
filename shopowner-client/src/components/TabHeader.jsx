import { SidebarSimple } from "phosphor-react";

export const TabHeader = ({ toggleSidebar, tab }) => {
  return (
    <div className="bg-gradient-to-br from-indigo-600 to-indigo-700 -mx-6 -mt-6 px-6 py-6 mb-6 relative overflow-hidden">
      {/* Decorative circles */}
      <div className="absolute top-0 right-0 w-24 h-24 bg-white/10 rounded-full -mr-12 -mt-12" />
      <div className="absolute bottom-0 left-0 w-16 h-16 bg-white/10 rounded-full -ml-8 -mb-8" />

      <div className="relative z-10 flex items-center gap-4">
        <button
          onClick={toggleSidebar}
          className="w-10 h-10 bg-white/10 hover:bg-white/20 rounded-lg flex items-center justify-center transition-colors duration-200 flex-shrink-0"
          aria-label="Toggle sidebar"
        >
          <SidebarSimple size={20} weight="regular" className="text-white" />
        </button>
        <div className="h-6 border-l border-white/30"></div>
        <h2 className="text-white font-bold text-lg">{tab}</h2>
      </div>
    </div>
  );
};
