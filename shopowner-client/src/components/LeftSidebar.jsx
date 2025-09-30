import React from "react";
import {
  PiStorefront,
  PiCallBell,
  PiChair,
  PiInvoice,
  PiUsersThreeLight,
} from "react-icons/pi";
import { GoGear } from "react-icons/go";
import { IoHelpCircleOutline } from "react-icons/io5";
import { Tabs } from "../constants/constants";
import { useTabs } from "../store/tabStore";
import { useNavigate } from "react-router-dom";

const SidebarItem = ({
  children,
  variant = "default",
  icon: Icon,
  collapsed,
  ...props
}) => {
  const BASE_CLASSES =
    "flex items-center px-2 rounded-md text-sm cursor-pointer h-10 transition-colors duration-300";
  const VARIANTS = {
    default: "bg-white text-gray-600 hover:bg-gray-100 hover:text-gray-800",
    active: "bg-gray-100 text-emerald-700",
  };
  const variantClasses = VARIANTS[variant] || VARIANTS.default;

  return (
    <button
      className={`${BASE_CLASSES} ${variantClasses} ${
        collapsed ? "justify-center" : "gap-2"
      }`}
      {...props}
    >
      <span className="w-6 h-6 flex items-center justify-center">
        <Icon size={18} color="currentColor" />
      </span>
      {!collapsed && <span className="whitespace-nowrap">{children}</span>}
    </button>
  );
};

const SidebarHeader = ({ collapsed }) => (
  <div className="flex gap-3 p-5 items-center h-22 transition-all duration-300">
    <div className="w-8 h-8 bg-gray-800 rounded-sm"></div>
    {!collapsed && (
      <div className="flex flex-col transition-opacity duration-300">
        <h3 className="text-gray-800">My Cafe App</h3>
        <p className="text-gray-800">v101</p>
      </div>
    )}
  </div>
);

const SidebarProfile = ({ collapsed }) => (
  <div className="flex gap-3 p-5 items-center h-22 transition-all duration-300">
    <div className="w-8 h-8 bg-gray-800 rounded-full"></div>
    {!collapsed && (
      <div className="flex flex-col transition-opacity duration-300">
        <h3 className="text-gray-800">John Doe</h3>
        <p className="text-gray-800">johndoe@gmail.com</p>
      </div>
    )}
  </div>
);

export const LeftSidebar = ({ collapsed, tab }) => {
  const navigate = useNavigate();
  return (
    <aside
      className={`flex flex-col flex-shrink-0 ${
        collapsed ? "w-20" : "w-64"
      } h-screen bg-white border-r border-gray-200 transition-all duration-300 ease-in-out`}
    >
      <div className="flex flex-col flex-grow">
        <SidebarHeader collapsed={collapsed} />

        <div className="flex flex-col gap-4">
          <div className="flex flex-col px-5">
            {!collapsed && (
              <div className="px-1.5 py-2 transition-opacity duration-300">
                <p className="text-xs text-gray-400">Primary</p>
              </div>
            )}

            <SidebarItem
              variant={tab === Tabs.MENU ? "active" : "default"}
              icon={PiStorefront}
              onClick={() => navigate("/menu")}
              collapsed={collapsed}
            >
              Menu
            </SidebarItem>

            <SidebarItem
              variant={tab === Tabs.ORDERS ? "active" : "default"}
              icon={PiCallBell}
              onClick={() => navigate("/orders")}
              collapsed={collapsed}
            >
              Orders
            </SidebarItem>

            <SidebarItem
              variant={tab === Tabs.TABLES ? "active" : "default"}
              icon={PiChair}
              onClick={() => navigate("/tables")}
              collapsed={collapsed}
            >
              Tables
            </SidebarItem>

            <SidebarItem
              variant={tab === Tabs.HISTORY ? "active" : "default"}
              icon={PiInvoice}
              onClick={() => navigate("/history")}
              collapsed={collapsed}
            >
              History
            </SidebarItem>
          </div>

          <div className="flex flex-col px-5">
            {!collapsed && (
              <div className="px-1.5 py-2 transition-opacity duration-300">
                <p className="text-xs text-gray-400">Others</p>
              </div>
            )}

            <SidebarItem
              variant={tab === Tabs.SETTINGS ? "active" : "default"}
              icon={GoGear}
              onClick={() => navigate("/settings")}
              collapsed={collapsed}
            >
              Settings
            </SidebarItem>

            <SidebarItem
              variant={tab === Tabs.COMMUNITY ? "active" : "default"}
              icon={PiUsersThreeLight}
              onClick={() => navigate("/community")}
              collapsed={collapsed}
            >
              Community
            </SidebarItem>

            <SidebarItem
              variant={tab === Tabs.HELP ? "active" : "default"}
              icon={IoHelpCircleOutline}
              onClick={() => navigate("/help")}
              collapsed={collapsed}
            >
              Help
            </SidebarItem>
          </div>
        </div>
      </div>

      <div className="mt-auto">
        <SidebarProfile collapsed={collapsed} />
      </div>
    </aside>
  );
};

export default LeftSidebar;
