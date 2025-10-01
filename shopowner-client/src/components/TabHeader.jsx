import { BsLayoutSidebar } from "react-icons/bs";

export const TabHeader = ({ toggleSidebar, tab }) => {
  return (
    <div className="flex items-center gap-10 h-10">
      <div className="" onClick={toggleSidebar}>
        <BsLayoutSidebar className="text-xl" />
      </div>
      <div className="h-[70%] border-l border-gray-300"></div>
      <h2 className="text-xl text-gray-800">{tab}</h2>
    </div>
  );
};
