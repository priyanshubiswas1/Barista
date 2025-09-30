import { useState } from "react";
import { Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";
import { TabHeader } from "../components/TabHeader";

const Orders = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  return (
    <div className="flex">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.ORDERS} />
      <div className="p-5 h-[100%] w-[100%] relative flex flex-col">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.ORDERS} />
      </div>
    </div>
  );
};

export default Orders;
