import { useState } from "react";
import { IoFilter } from "react-icons/io5";
import { FaPlus } from "react-icons/fa";
import { Overlays, Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar.jsx";
import { TabHeader } from "../components/TabHeader.jsx";
import Searchbar from "../components/Searchbar.jsx";
import Button from "../components/Button.jsx";
import { useOverlay } from "../store/overlayStore.js";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";

const OrderSubHeader = ({ onAddNewOrderClick }) => {
  return (
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-5">
        <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button>
      </div>
      <div className="flex items-center gap-5">
        <Button onClick={onAddNewOrderClick} variant="primary" icon={FaPlus}>
          Add New Order
        </Button>
      </div>
    </div>
  );
};

const OrderContent = () => {
  return <div className="bg-black w-full h-[50%]"></div>;
};

const Orders = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  const { setOverlay } = useOverlay();
  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.ORDERS} />
      <div className="p-5 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.ORDERS} />
        <OrderSubHeader
          onAddNewOrderClick={() => setOverlay(Overlays.CREATEORDER)}
        />
        <div className="pt-5 flex-grow">
          <OrderContent />
        </div>
      </div>
      <OverlaySidebar />
    </div>
  );
};

export default Orders;
