import React from "react";
import { PiStorefront, PiCallBell, PiChair, PiInvoice } from "react-icons/pi";
import { GoGear } from "react-icons/go";
import { PiUsersThreeLight } from "react-icons/pi";
import { IoHelpCircleOutline } from "react-icons/io5";
import { Tabs } from "../constants/constants";
import { useTabs } from "../store/tabStore";

const SidebarItem = ({
  children,
  variant = "default",
  icon: Icon,
  ...props
}) => {
  const BASE_CLASSES =
    "flex gap-2 items-center px-2 py-1.5 rounded-md text-sm cursor-pointer transition-colors duration-120 ease-in-out";

  const VARIANTS = {
    default: "bg-white text-gray-600 hover:bg-gray-100 hover:text-gray-800",
    active: "bg-gray-100 text-emerald-700",
  };

  const variantClasses = VARIANTS[variant] || VARIANTS.default;

  return (
    <button className={`${BASE_CLASSES} ${variantClasses}`} {...props}>
      {Icon && (
        <span className="w-4 h-4 flex items-center justify-center">
          <Icon size={16} color="currentColor" />
        </span>
      )}
      {children}
    </button>
  );
};

const SidebarHeader = () => {
  return (
    <div className="flex gap-3 p-5 items-center">
      <div className="w-8 h-8 bg-gray-800 rounded-sm"></div>
      <div className="flex flex-col">
        <h3 className="text-gray-800">My Cafe App</h3>
        <p className="text-gray-800">v101</p>
      </div>
    </div>
  );
};

const SidebarProfile = () => {
  return (
    <div className="flex gap-3 p-5 items-center">
      <div className="w-8 h-8 bg-gray-800 rounded-full"></div>
      <div className="flex flex-col">
        <h3 className="text-gray-800">John Doe</h3>
        <p className="text-gray-800">johndoe@gmail.com</p>
      </div>
    </div>
  );
};

export const LeftSidebar = ({ collapsed }) => {
  const { tab, setTab } = useTabs();
  return (
    <aside
      className={`flex flex-col flex-shrink-0 justify-between ${
        collapsed ? "translate-x-[-100%]" : "translate-x-0"
      } h-screen bg-white border-r border-gray-200 transition-all duration-50 ease-in-out`}
    >
      {!collapsed && (
        <div className="flex flex-col">
          <SidebarHeader />
          <div className="flex flex-col gap-4">
            <div className="flex flex-col px-5">
              <div className="px-1.5 py-2">
                <p className="text-xs text-gray-400">Primary</p>
              </div>
              <SidebarItem
                variant={tab === Tabs.MENU ? "active" : "default"}
                icon={PiStorefront}
                onClick={() => setTab(Tabs.MENU)}
              >
                Menu
              </SidebarItem>
              <SidebarItem
                variant={tab === Tabs.ORDERS ? "active" : "default"}
                icon={PiCallBell}
                onClick={() => setTab(Tabs.ORDERS)}
              >
                Orders
              </SidebarItem>
              <SidebarItem
                variant={tab === Tabs.TABLES ? "active" : "default"}
                icon={PiChair}
                onClick={() => setTab(Tabs.TABLES)}
              >
                Tables
              </SidebarItem>
              <SidebarItem
                variant={tab === Tabs.HISTORY ? "active" : "default"}
                icon={PiInvoice}
                onClick={() => setTab(Tabs.HISTORY)}
              >
                History
              </SidebarItem>
            </div>
            <div className="flex flex-col px-5">
              <div className="px-1.5 py-2">
                <p className="text-xs text-gray-400">Others</p>
              </div>
              <SidebarItem
                variant={tab === Tabs.SETTINGS ? "active" : "default"}
                icon={GoGear}
                onClick={() => setTab(Tabs.SETTINGS)}
              >
                Settings
              </SidebarItem>
              <SidebarItem
                variant={tab === Tabs.COMMUNITY ? "active" : "default"}
                icon={PiUsersThreeLight}
                onClick={() => setTab(Tabs.COMMUNITY)}
              >
                Community
              </SidebarItem>
              <SidebarItem
                variant={tab === Tabs.HELP ? "active" : "default"}
                icon={IoHelpCircleOutline}
                onClick={() => setTab(Tabs.HELP)}
              >
                Help
              </SidebarItem>
            </div>
          </div>
        </div>
      )}
      {!collapsed && <SidebarProfile />}
    </aside>
  );
};

export default LeftSidebar;
