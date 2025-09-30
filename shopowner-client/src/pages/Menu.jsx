import { TabHeader } from "../components/TabHeader";
import { Overlays, Tabs } from "../constants/constants";
import Searchbar from "../components/Searchbar";
import Filter from "../components/Filter";
import { FaPlus } from "react-icons/fa";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useState } from "react";
import { useOverlay } from "../store/overlayStore";
import MenuContent from "../components/menu/MenuContent";
import { useQuery } from "@tanstack/react-query";
import { fetchItemsByShopId } from "../api/api";
import LeftSidebar from "../components/LeftSidebar";

const Preview = () => {
  return (
    <div className="flex justify-center items-center gap-3 bg-white border border-gray-200 rounded-lg px-3 py-2">
      <h2>Preview</h2>
    </div>
  );
};

const ManageCategory = ({ onClick }) => {
  return (
    <div
      onClick={onClick}
      className="flex justify-center items-center gap-3 bg-white border border-gray-200 rounded-lg px-3 py-2"
    >
      <h2>Manage Categories</h2>
    </div>
  );
};

const AddItem = ({ onClick }) => (
  <div
    onClick={onClick}
    className="flex justify-center items-center gap-3 bg-gray-800 border border-gray-200 rounded-lg px-3 py-2 cursor-pointer"
  >
    <FaPlus className="text-white" />
    <h2 className="text-white">Add Items</h2>
  </div>
);

const MenuSubHeader = ({ onAddItemClick, onManageCategoryClick }) => {
  return (
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-5">
        <Searchbar />
        <Filter />
      </div>
      <div className="flex items-center gap-5">
        <Preview />
        <ManageCategory onClick={onManageCategoryClick} />
        <AddItem onClick={onAddItemClick} />
      </div>
    </div>
  );
};

const Menu = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  const { setOverlay } = useOverlay();

  const { data, isLoading, error } = useQuery({
    queryKey: ["menuItems"],
    queryFn: () => fetchItemsByShopId("8dc5b0b4-dbb6-457a-beac-b469b6c1e396"),
  });

  if (isLoading) return <p>Loading item...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const items = data.items;

  const groupedItems = items.reduce((acc, item) => {
    if (!acc[item.category]) {
      acc[item.category] = [];
    }
    acc[item.category].push(item);
    return acc;
  }, {});
  console.log(groupedItems);

  return (
    <div className="flex">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.MENU} />
      <div className="p-5 h-[100%] w-[100%] relative flex flex-col">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.MENU} />
        <MenuSubHeader
          onAddItemClick={() => setOverlay(Overlays.ADDITEM)}
          onManageCategoryClick={() => setOverlay(Overlays.MANAGECAT)}
        />
        <div className="pt-5 flex flex-col gap-2">
          {Object.entries(groupedItems).map(([category, categoryItems]) => (
            <MenuContent
              key={categoryItems.itemId}
              category={category}
              items={categoryItems}
            />
          ))}
        </div>
        <OverlaySidebar />
      </div>
    </div>
  );
};

export default Menu;
