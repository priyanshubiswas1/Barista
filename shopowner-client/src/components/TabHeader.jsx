// import { BsLayoutSidebar } from "react-icons/bs";
import { SidebarSimple } from "phosphor-react";

export const TabHeader = ({ toggleSidebar, tab }) => {
  return (
    <div className="flex items-center gap-4">
      <div className="cursor-pointer" onClick={toggleSidebar}>
        <SidebarSimple size={16} />
      </div>
      <div className="h-3.5 border-l border-gray-300"></div>
      <p className="text-gray-800">{tab}</p>
    </div>
  );
};
