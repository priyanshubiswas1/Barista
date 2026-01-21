import React, { useEffect, useState } from "react";
import {
  PiStorefront,
  PiCallBell,
  PiChair,
  PiInvoice,
  PiUsersThreeLight,
} from "react-icons/pi";
import { GoGear } from "react-icons/go";
import { IoHelpCircleOutline, IoClose } from "react-icons/io5";
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
    active: "bg-indigo-200 text-indigo-700 font-medium shadow-sm",
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

const SidebarHeader = ({ collapsed, onClose, isMobile }) => (
  <div className="bg-gradient-to-br from-indigo-600 to-indigo-700 p-4 sm:p-6 relative overflow-hidden transition-all duration-300 mb-2">
    {/* Decorative circles */}
    <div className="absolute top-0 right-0 w-24 h-24 bg-white/10 rounded-full -mr-12 -mt-12" />
    <div className="absolute bottom-0 left-0 w-16 h-16 bg-white/10 rounded-full -ml-8 -mb-8" />

    <div
      className={`relative z-10 flex ${
        collapsed && !isMobile ? "justify-center" : "justify-between"
      } items-center`}
    >
      <div className={`flex ${collapsed && !isMobile ? "justify-center" : "gap-3"} items-center`}>
        <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-lg flex-shrink-0">
          <span className="text-2xl">☕</span>
        </div>
        {(!collapsed || isMobile) && (
          <div className="flex flex-col transition-opacity duration-300">
            <h3 className="text-white font-bold text-lg">Barista</h3>
            <p className="text-indigo-100 text-xs">Shop Owner Portal</p>
          </div>
        )}
      </div>
      {isMobile && onClose && (
        <button
          onClick={onClose}
          className="w-8 h-8 bg-white/10 hover:bg-white/20 rounded-lg flex items-center justify-center transition-colors"
        >
          <IoClose size={20} className="text-white" />
        </button>
      )}
    </div>
  </div>
);

const SidebarProfile = ({ collapsed, isMobile }) => (
  <div
    className={`bg-gradient-to-br from-indigo-600 to-indigo-700 ${
      collapsed && !isMobile ? "p-4" : "p-4 sm:p-5"
    } border-t border-indigo-500/30 transition-all duration-300`}
  >
    <div
      className={`flex ${collapsed && !isMobile ? "justify-center" : "gap-3"} items-center`}
    >
      <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-lg flex-shrink-0 font-bold text-indigo-700">
        {localStorage.getItem("name")?.charAt(0).toUpperCase() || "U"}
      </div>
      {(!collapsed || isMobile) && (
        <div className="flex flex-col transition-opacity duration-300 min-w-0">
          <h3 className="text-white font-semibold text-sm truncate">
            {localStorage.getItem("name")}
          </h3>
          <p className="text-indigo-100 text-xs truncate">
            {localStorage.getItem("email")}
          </p>
        </div>
      )}
    </div>
  </div>
);

const SidebarContent = ({ collapsed, tab, onNavigate, isMobile }) => {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col px-3 sm:px-5">
        {(!collapsed || isMobile) && (
          <div className="px-1.5 py-2 transition-opacity duration-300">
            <p className="text-xs text-gray-400">Primary</p>
          </div>
        )}

        <SidebarItem
          variant={tab === Tabs.MENU ? "active" : "default"}
          icon={PiStorefront}
          onClick={() => onNavigate("/menu")}
          collapsed={collapsed && !isMobile}
        >
          Menu
        </SidebarItem>

        <SidebarItem
          variant={tab === Tabs.ORDERS ? "active" : "default"}
          icon={PiCallBell}
          onClick={() => onNavigate("/orders")}
          collapsed={collapsed && !isMobile}
        >
          Orders
        </SidebarItem>

        <SidebarItem
          variant={tab === Tabs.TABLES ? "active" : "default"}
          icon={PiChair}
          onClick={() => onNavigate("/tables")}
          collapsed={collapsed && !isMobile}
        >
          Tables
        </SidebarItem>

        <SidebarItem
          variant={tab === Tabs.HISTORY ? "active" : "default"}
          icon={PiInvoice}
          onClick={() => onNavigate("/history")}
          collapsed={collapsed && !isMobile}
        >
          History
        </SidebarItem>
      </div>

      <div className="flex flex-col px-3 sm:px-5">
        {(!collapsed || isMobile) && (
          <div className="px-1.5 py-2 transition-opacity duration-300">
            <p className="text-xs text-gray-400">Others</p>
          </div>
        )}

        <SidebarItem
          variant={tab === Tabs.SETTINGS ? "active" : "default"}
          icon={GoGear}
          onClick={() => onNavigate("/settings")}
          collapsed={collapsed && !isMobile}
        >
          Settings
        </SidebarItem>

        <SidebarItem
          variant={tab === Tabs.COMMUNITY ? "active" : "default"}
          icon={PiUsersThreeLight}
          onClick={() => onNavigate("/community")}
          collapsed={collapsed && !isMobile}
        >
          Community
        </SidebarItem>

        <SidebarItem
          variant={tab === Tabs.HELP ? "active" : "default"}
          icon={IoHelpCircleOutline}
          onClick={() => onNavigate("/help")}
          collapsed={collapsed && !isMobile}
        >
          Help
        </SidebarItem>
      </div>
    </div>
  );
};

export const LeftSidebar = ({ collapsed, tab, mobileOpen, onMobileClose }) => {
  const navigate = useNavigate();
  const { setOverlay } = useOverlay();

  useEffect(() => {
    setOverlay("");
  }, [tab, setOverlay]);

  const handleNavigate = (path) => {
    navigate(path);
    if (onMobileClose) {
      onMobileClose();
    }
  };

  return (
    <>
      {/* Desktop Sidebar */}
      <aside
        className={`hidden md:flex flex-col flex-shrink-0 ${
          collapsed ? "w-20" : "w-64"
        } h-screen bg-white border-r border-gray-200 transition-all duration-300 ease-in-out`}
      >
        <div className="flex flex-col flex-grow">
          <SidebarHeader collapsed={collapsed} isMobile={false} />
          <SidebarContent
            collapsed={collapsed}
            tab={tab}
            onNavigate={handleNavigate}
            isMobile={false}
          />
        </div>
        <div className="mt-auto">
          <SidebarProfile collapsed={collapsed} isMobile={false} />
        </div>
      </aside>

      {/* Mobile Sidebar Overlay */}
      {mobileOpen && (
        <div
          className="md:hidden fixed inset-0 bg-black/60 backdrop-blur-sm z-40"
          onClick={onMobileClose}
        />
      )}

      {/* Mobile Sidebar Drawer */}
      <aside
        className={`md:hidden fixed top-0 left-0 h-screen w-72 bg-white shadow-2xl z-50 transform transition-transform duration-300 ease-in-out ${
          mobileOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex flex-col h-full">
          <SidebarHeader collapsed={false} onClose={onMobileClose} isMobile={true} />
          <div className="flex-grow overflow-y-auto">
            <SidebarContent
              collapsed={false}
              tab={tab}
              onNavigate={handleNavigate}
              isMobile={true}
            />
          </div>
          <div className="mt-auto">
            <SidebarProfile collapsed={false} isMobile={true} />
          </div>
        </div>
      </aside>
    </>
  );
};

export default LeftSidebar;
