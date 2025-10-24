import { TabHeader } from "../components/TabHeader";
import { Overlays, Tabs } from "../constants/constants";
import Searchbar from "../components/Searchbar";
import { IoFilter } from "react-icons/io5";
import { Plus } from "phosphor-react";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useState } from "react";
import { useOverlay } from "../store/overlayStore";
import MenuContent from "../components/menu/MenuContent";
import { useQuery } from "@tanstack/react-query";
import { fetchItemsByShopId } from "../api/api";
import LeftSidebar from "../components/LeftSidebar";
import MenuNoItems from "../components/menu/MenuNoItems";
import Button from "../components/Button";

const MenuSubHeader = ({ onAddItemClick, onManageCategoryClick }) => {
  return (
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-3">
        <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button>
      </div>
      <div className="flex items-center gap-3">
        <Button variant="outline">Preview</Button>
        <Button onClick={onManageCategoryClick} variant="outline">
          Manage Categories
        </Button>
        <Button onClick={onAddItemClick} variant="primary" icon={Plus}>
          Add New Items
        </Button>
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
    queryFn: () => fetchItemsByShopId(localStorage.getItem("shopId")),
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
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.MENU} />
      <div className="p-5 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.MENU} />
        <MenuSubHeader
          onAddItemClick={() => setOverlay(Overlays.ADDITEM)}
          onManageCategoryClick={() => setOverlay(Overlays.MANAGECAT)}
        />
        {items.length > 0 ? (
          <div className="pt-5 flex-grow flex flex-col gap-4">
            {Object.entries(groupedItems).map(
              ([category, categoryItems], index) => (
                <MenuContent
                  key={index}
                  category={category}
                  items={categoryItems}
                />
              )
            )}
          </div>
        ) : (
          <div className="pt-5 flex-grow">
            <MenuNoItems />
          </div>
        )}
        <OverlaySidebar />
      </div>
    </div>
  );
};

export default Menu;
