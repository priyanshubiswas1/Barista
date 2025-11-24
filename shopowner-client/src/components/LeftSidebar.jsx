import React, { useEffect } from "react";
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
import { useOverlay } from "../store/overlayStore";

const SidebarItem = ({
  children,
  variant = "default",
  icon: Icon,
  collapsed,
  ...props
}) => {
  const BASE_CLASSES =
    "flex items-center px-2 rounded-md text-sm cursor-pointer h-10 transition-all duration-300";
  const VARIANTS = {
    default: "bg-white text-gray-600 hover:bg-gray-50 hover:text-gray-800",
    active: "bg-emerald-50 text-emerald-700 border-l-4 border-emerald-600 font-medium shadow-sm",
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
  <div className={`bg-gradient-to-br from-emerald-600 to-emerald-700 ${collapsed ? 'p-5' : 'p-6'} relative overflow-hidden transition-all duration-300`}>
    {/* Decorative circles */}
    <div className="absolute top-0 right-0 w-24 h-24 bg-white/10 rounded-full -mr-12 -mt-12" />
    <div className="absolute bottom-0 left-0 w-16 h-16 bg-white/10 rounded-full -ml-8 -mb-8" />
    
    <div className={`relative z-10 flex ${collapsed ? 'justify-center' : 'gap-3'} items-center`}>
      <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-lg flex-shrink-0">
        <span className="text-2xl">☕</span>
      </div>
      {!collapsed && (
        <div className="flex flex-col transition-opacity duration-300">
          <h3 className="text-white font-bold text-lg">Barista</h3>
          <p className="text-emerald-100 text-xs">Shop Owner Portal</p>
        </div>
      )}
    </div>
  </div>
);

const SidebarProfile = ({ collapsed }) => (
  <div className={`bg-gradient-to-br from-emerald-600 to-emerald-700 ${collapsed ? 'p-4' : 'p-5'} border-t border-emerald-500/30 transition-all duration-300`}>
    <div className={`flex ${collapsed ? 'justify-center' : 'gap-3'} items-center`}>
      <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-lg flex-shrink-0 font-bold text-emerald-700">
        {localStorage.getItem("name")?.charAt(0).toUpperCase() || "U"}
      </div>
      {!collapsed && (
        <div className="flex flex-col transition-opacity duration-300 min-w-0">
          <h3 className="text-white font-semibold text-sm truncate">{localStorage.getItem("name")}</h3>
          <p className="text-emerald-100 text-xs truncate">{localStorage.getItem("email")}</p>
        </div>
      )}
    </div>
  </div>
);

export const LeftSidebar = ({ collapsed, tab }) => {
  const navigate = useNavigate();
  const { setOverlay } = useOverlay();
  useEffect(() => {
    setOverlay("");
  }, [tab, setOverlay]);
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
