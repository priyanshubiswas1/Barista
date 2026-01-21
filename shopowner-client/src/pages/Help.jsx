import { useState } from "react";
import { TabHeader } from "../components/TabHeader";
import { Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";

const Help = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  return (
    <div className="flex h-screen">
      <LeftSidebar
        collapsed={isSidebarCollapesed}
        tab={Tabs.HELP}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div className="px-4 sm:px-6 pt-4 sm:pt-6 flex flex-col flex-1 overflow-hidden">
        <TabHeader
          toggleSidebar={toggleSidebar}
          tab={Tabs.HELP}
          onMobileMenuClick={() => setMobileMenuOpen(true)}
        />
        <div className="mt-2 sm:mt-4 text-gray-600 text-sm">
          Help page (responsive shell wired).
        </div>
      </div>
    </div>
  );
};

export default Help;
