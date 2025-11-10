import { useState, useRef, useEffect } from "react";
import { IoFilter } from "react-icons/io5";
import { IoCalendarClearOutline } from "react-icons/io5";
import { IoIosArrowDown } from "react-icons/io";
import Button from "../components/Button";
import Searchbar from "../components/Searchbar";
import { TabHeader } from "../components/TabHeader";
import LeftSidebar from "../components/LeftSidebar";
import { Tabs } from "../constants/constants";
import { useQuery } from "@tanstack/react-query";
import { getSessionSummary } from "../api/api";

const HistorySubHeader = ({ selectedFilter, setSelectedFilter }) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);

  const toggleDropdown = () => setIsDropdownOpen((prev) => !prev);

  const handleOptionClick = (option) => {
    setSelectedFilter(option);
    setIsDropdownOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-3">
        <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative" ref={dropdownRef}>
          <button
            type="button"
            onClick={toggleDropdown}
            className="rounded-md shadow-xs disabled:shadow-none text-sm font-medium cursor-pointer disabled:cursor-not-allowed transition-colors duration-120 ease-in-out
                 bg-gray-200 text-gray-500 hover:text-gray-700 disabled:bg-gray-100 disabled:text-gray-200
                 w-48 h-8 flex items-center gap-2 justify-center"
          >
            <IoCalendarClearOutline size={16} color="currentColor" />
            {selectedFilter}
            <IoIosArrowDown size={16} color="currentColor" />
          </button>

          {isDropdownOpen && (
            <div className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-md z-10">
              {["Today", "Last Week", "Last Month", "Last 3 Months"].map(
                (option) => (
                  <button
                    key={option}
                    onClick={() => handleOptionClick(option)}
                    className={`w-full text-left px-4 py-2 text-sm hover:bg-gray-100 ${
                      selectedFilter === option ? "bg-gray-100 font-medium" : ""
                    }`}
                    style={{
                      minHeight: "32px",
                      lineHeight: "32px",
                    }}
                  >
                    {option}
                  </button>
                )
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const SalesSummaryCard = ({ filter }) => {
  return (
    <div className="bg-white border-2 p-3 border-gray-200 hover:border-gray-300 rounded-md flex flex-col gap-y-5 w-65">
      <div className="flex justify-between gap-10 text-xs font-medium items-center">
        <h2>Total Sales</h2>
        <h2 className="bg-gray-200 border-gray-200 rounded-xl px-2 py-0.5 text-gray-800">
          {filter}
        </h2>
      </div>
      <div className="flex justify-between items-center">
        <h1 className="font-semibold">₹4000</h1>
        <h2 className="text-green-800 text-xs">+30%</h2>
      </div>
    </div>
  );
};
const SessionsSummaryCard = ({ filter }) => {
  return (
    <div className="bg-white border-2 p-3 border-gray-200 hover:border-gray-300 rounded-md flex flex-col gap-y-5 w-65">
      <div className="flex justify-between gap-10 text-xs font-medium items-center">
        <h2>Total Sessions</h2>
        <h2 className="bg-gray-200 border-gray-200 rounded-xl px-2 py-0.5 text-gray-800">
          {filter}
        </h2>
      </div>
      <div className="flex justify-between items-center">
        <h1 className="font-semibold">₹4000</h1>
        <h2 className="text-green-800 text-xs">+30%</h2>
      </div>
    </div>
  );
};

const HistoryContent = () => {
  const payload = {
    shopId: localStorage.getItem("shopId"),
    offset: 1,
    limit: 5,
  };
  const { data, isLoading, error } = useQuery({
    queryKey: ["history"],
    queryFn: () => getSessionSummary(payload),
  });
  if (isLoading) return <p>Loading history...</p>;
  if (error) return <p>Error: {error.message}</p>;

  console.log(data);
  // const sessions = data.sessions;
  const orderData = [
    { orderId: "002", items: 5, session: "1h45m", amount: "120.00" },
    { orderId: "002", items: 5, session: "1h45m", amount: "120.00" },
    { orderId: "002", items: 5, session: "1h45m", amount: "120.00" },
    { orderId: "002", items: 5, session: "1h45m", amount: "120.00" },
    { orderId: "002", items: 5, session: "1h45m", amount: "120.00" },
  ];

  return (
    <div className="flex flex-col p-2">
      <div className="flex bg-gray-100 p-2 text-sm font-medium text-gray-600">
        <div className="flex-2">Order ID</div>
        <div className="flex-1">Items</div>
        <div className="flex-1">Session</div>
        <div className="flex-1">Amount</div>
      </div>
      <div className="w-full h-full flex flex-col gap-2">
        {orderData.map((order, index) => (
          <div
            key={index}
            className="flex bg-white p-2 text-sm text-gray-700 border-b border-gray-200 rounded-md"
          >
            <div className="flex-2">{order.orderId}</div>
            <div className="flex-1">{order.items}</div>
            <div className="flex-1">{order.session}</div>
            <div className="flex-1">{order.amount}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

const History = () => {
  const [isSidebarCollapesed, setIsSidebarCollapsed] = useState(false);
  const [selectedFilter, setSelectedFilter] = useState("Today");

  const toggleSidebar = () => {
    setIsSidebarCollapsed((prev) => !prev);
  };

  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.HISTORY} />
      <div className="p-5 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HISTORY} />
        <HistorySubHeader
          selectedFilter={selectedFilter}
          setSelectedFilter={setSelectedFilter}
        />
        <div className="flex gap-3 py-3">
          <SalesSummaryCard filter={selectedFilter} />
          <SessionsSummaryCard filter={selectedFilter} />
        </div>
        <HistoryContent />
      </div>
    </div>
  );
};

export default History;
