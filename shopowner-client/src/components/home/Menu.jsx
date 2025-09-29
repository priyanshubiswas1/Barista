import { TabHeader } from "../TabHeader";
import { Overlays, Tabs } from "../../constants/constants";
import Searchbar from "../Searchbar";
import Filter from "../Filter";
import { FaPlus } from "react-icons/fa";
import { OverlaySidebar } from "../overlaysidebar/OverlaySidebar";
import { useState } from "react";
import { useOverlay } from "../../store/overlayStore";

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
      <div className="flex items-center">
        <Searchbar />
        <Filter />
      </div>
      <div className="flex items-center">
        <Preview />
        <ManageCategory onClick={onManageCategoryClick} />
        <AddItem onClick={onAddItemClick} />
      </div>
    </div>
  );
};

const Menu = ({ toggleSidebar }) => {
  const { setOverlay } = useOverlay();
  return (
    <div className="p-5 h-[100%] w-[100%] relative">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.MENU} />
      <MenuSubHeader
        onAddItemClick={() => setOverlay(Overlays.ADDITEM)}
        onManageCategoryClick={() => setOverlay(Overlays.MANAGECAT)}
      />
      <OverlaySidebar />
    </div>
  );
};

export default Menu;
