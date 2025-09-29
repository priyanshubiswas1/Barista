import { useState } from "react";
import LeftSidebar from "../components/LeftSidebar.jsx";
import Menu from "../components/home/Menu.jsx";
import Orders from "../components/home/Orders.jsx";
import Tables from "../components/home/Tables.jsx";
import History from "../components/home/History.jsx";
import Settings from "../components/home/Settings.jsx";
import Community from "../components/home/Community.jsx";
import Help from "../components/home/Help.jsx";
import { Tabs } from "../constants/constants.js";
import { useTabs } from "../store/tabStore.js";

const Home = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };

  const tabComponents = {
    [Tabs.MENU]: <Menu toggleSidebar={toggleSidebar} />,
    [Tabs.ORDERS]: <Orders toggleSidebar={toggleSidebar} />,
    [Tabs.TABLES]: <Tables toggleSidebar={toggleSidebar} />,
    [Tabs.HISTORY]: <History toggleSidebar={toggleSidebar} />,
    [Tabs.SETTINGS]: <Settings toggleSidebar={toggleSidebar} />,
    [Tabs.COMMUNITY]: <Community toggleSidebar={toggleSidebar} />,
    [Tabs.HELP]: <Help toggleSidebar={toggleSidebar} />,
  };

  return (
    <div className="flex">
      <LeftSidebar collapsed={isSidebarCollapesed} />
      {tabComponents[useTabs((state) => state.tab)]}
    </div>
  );
};

export default Home;
