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
import CircleLoading from "../components/CircleLoading";

const MenuSubHeader = ({ onAddItemClick, onManageCategoryClick }) => {
  return (
    <div className="w-full flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 pt-2 pb-4 sm:pb-6">
      <div className="flex items-center gap-3">
        {/* <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button> */}
      </div>
      <div className="flex flex-wrap items-center gap-2 sm:gap-3 w-full sm:w-auto">
        {/* <Button variant="outline">Preview</Button> */}
        <Button onClick={onManageCategoryClick} variant="outline">
          <span className="hidden sm:inline">Manage Categories</span>
          <span className="sm:hidden">Categories</span>
        </Button>
        <Button onClick={onAddItemClick} variant="primary" icon={Plus}>
          <span className="hidden sm:inline">Add New Items</span>
          <span className="sm:hidden">Add Item</span>
        </Button>
      </div>
    </div>
  );
};

const Menu = () => {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarCollapsed((prev) => !prev);
  };

  const { setOverlay } = useOverlay();

  const { data, isLoading, error } = useQuery({
    queryKey: ["menuItems"],
    queryFn: () => fetchItemsByShopId(localStorage.getItem("shopId")),
  });

  if (isLoading) return <CircleLoading text="Loading items..." />;
  if (error) return <p>Error: {error.message}</p>;

  const items = data.items;

  const groupedItems = items.reduce((acc, item) => {
    if (!acc[item.category]) {
      acc[item.category] = [];
    }
    acc[item.category].push(item);
    return acc;
  }, {});

  Object.keys(groupedItems).forEach((category) => {
    groupedItems[category].sort((a, b) => {
      const nameA = (a.itemName || "").toLowerCase();
      const nameB = (b.itemName || "").toLowerCase();
      return nameA.localeCompare(nameB);
    });
  });

  const sortedCategoryEntries = Object.entries(groupedItems).sort((a, b) => {
    const categoryA = (a[0] || "").toLowerCase();
    const categoryB = (b[0] || "").toLowerCase();
    return categoryA.localeCompare(categoryB);
  });

  return (
    <div className="flex h-screen">
      <LeftSidebar
        collapsed={isSidebarCollapsed}
        tab={Tabs.MENU}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div className="flex flex-col w-full overflow-hidden">
        <div className="px-4 sm:px-6 pt-4 sm:pt-6">
          <TabHeader
            toggleSidebar={toggleSidebar}
            tab={Tabs.MENU}
            onMobileMenuClick={() => setMobileMenuOpen(true)}
          />
          <MenuSubHeader
            onAddItemClick={() => setOverlay(Overlays.ADDITEM)}
            onManageCategoryClick={() => setOverlay(Overlays.MANAGECAT)}
          />
        </div>
        {items.length > 0 ? (
          <div className="px-4 sm:px-6 pb-4 sm:pb-6 pt-2 flex-grow flex flex-col gap-3 sm:gap-4 overflow-y-auto">
            {sortedCategoryEntries.map(([category, categoryItems]) => (
              <MenuContent
                key={category}
                category={category}
                items={categoryItems}
              />
            ))}
          </div>
        ) : (
          <div className="px-4 sm:px-6 pb-4 sm:pb-6 pt-5 flex-grow overflow-y-auto">
            <MenuNoItems />
          </div>
        )}
        <OverlaySidebar />
      </div>
    </div>
  );
};

export default Menu;
