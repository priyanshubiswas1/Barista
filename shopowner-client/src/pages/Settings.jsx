import { useState } from "react";
import { TabHeader } from "../components/TabHeader";
import { Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";

const Settings = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.SETTINGS} />
      <div className="p-6 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.SETTINGS} />
      </div>
    </div>
  );
};

export default Settings;
