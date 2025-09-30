import { useState } from "react";
import { TabHeader } from "../components/TabHeader";
import { Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";

const Help = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  return (
    <div className="flex">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.HELP} />
      <div className="p-5 h-[100%] w-[100%] relative flex flex-col">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HELP} />
      </div>
    </div>
  );
};

export default Help;
