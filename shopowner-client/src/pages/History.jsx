import { useState } from "react";
import { TabHeader } from "../components/TabHeader";
import { Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";

const History = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  return (
    <div className="flex">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.HISTORY} />
      <div className="p-5 h-[100%] w-[100%] relative flex flex-col">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HISTORY} />
      </div>
    </div>
  );
};

export default History;
