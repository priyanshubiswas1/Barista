import { useState, useRef, useEffect } from "react";
import { IoFilter } from "react-icons/io5";
import { IoCalendarClearOutline } from "react-icons/io5";
import { IoIosArrowDown } from "react-icons/io";
import { IoChevronBack, IoChevronForward } from "react-icons/io5";
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
        {/* <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button> */}
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

const SalesSummaryCard = ({ filter, totalAmount }) => {
  return (
    <div className="bg-white border-2 p-4 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col gap-y-4 w-65 shadow-sm hover:shadow-md transition-all duration-200">
      <div className="flex justify-between gap-10 text-xs font-semibold items-center">
        <h2 className="text-gray-600">Total Sales</h2>
        <h2 className="bg-emerald-100 border-emerald-200 rounded-lg px-3 py-1 text-emerald-700">
          {filter}
        </h2>
      </div>
      <div className="flex justify-between items-center">
        <h1 className="font-bold text-2xl text-gray-900">
          ₹{totalAmount || 0}
        </h1>
      </div>
    </div>
  );
};
const SessionsSummaryCard = ({ filter, totalSessions }) => {
  return (
    <div className="bg-white border-2 p-4 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col gap-y-4 w-65 shadow-sm hover:shadow-md transition-all duration-200">
      <div className="flex justify-between gap-10 text-xs font-semibold items-center">
        <h2 className="text-gray-600">Total Sessions</h2>
        <h2 className="bg-emerald-100 border-emerald-200 rounded-lg px-3 py-1 text-emerald-700">
          {filter}
        </h2>
      </div>
      <div className="flex justify-between items-center">
        <h1 className="font-bold text-2xl text-gray-900">
          {totalSessions || 0}
        </h1>
      </div>
    </div>
  );
};

const HistoryContent = ({
  data,
  onPageChange,
  currentPage,
  limit,
  totalSessions,
}) => {
  const sessionSummaries = data?.sessionSummaries || [];
  const totalPages = Math.ceil((totalSessions || 0) / limit);

  const calculateSessionTime = (createdAt, updatedAt) => {
    const createdDate = new Date(createdAt);
    const updatedDate = new Date(updatedAt);
    const timeDiff = updatedDate - createdDate;

    const hours = Math.floor(timeDiff / (1000 * 60 * 60));
    const minutes = Math.floor((timeDiff % (1000 * 60 * 60)) / (1000 * 60));

    return `${hours}h${minutes}m`;
  };

  return (
    <div className="flex flex-col p-2">
      <div className="flex bg-gray-100 p-2 text-sm font-medium text-gray-600">
        <div className="flex-2">Order ID</div>
        <div className="flex-1">Items</div>
        <div className="flex-1">Session</div>
        <div className="flex-1">Amount</div>
      </div>
      <div className="w-full h-full flex flex-col gap-2">
        {sessionSummaries.map((session, index) => (
          <div
            key={index}
            className="flex bg-white p-2 text-sm text-gray-700 border-b border-gray-200 rounded-md"
          >
            <div className="flex-2">{session.sessionNumber}</div>
            <div className="flex-1">{session.orderedItems.length}</div>
            <div className="flex-1">
              {calculateSessionTime(session.createdAt, session.updatedAt)}
            </div>
            <div className="flex-1">${session.totalAmount.toFixed(2)}</div>
          </div>
        ))}
      </div>
      {totalPages > 1 && (
        <div className="flex justify-between items-center mt-4 pt-4 border-t border-gray-200">
          <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium ${
              currentPage === 1
                ? "bg-gray-100 text-gray-400 cursor-not-allowed"
                : "bg-gray-200 text-gray-700 hover:bg-gray-300"
            }`}
          >
            <IoChevronBack size={16} />
            Previous
          </button>
          <span className="text-sm text-gray-600">
            Page {currentPage} of {totalPages}
          </span>
          <button
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium ${
              currentPage >= totalPages
                ? "bg-gray-100 text-gray-400 cursor-not-allowed"
                : "bg-gray-200 text-gray-700 hover:bg-gray-300"
            }`}
          >
            Next
            <IoChevronForward size={16} />
          </button>
        </div>
      )}
    </div>
  );
};

const History = () => {
  const [isSidebarCollapesed, setIsSidebarCollapsed] = useState(false);
  const [selectedFilter, setSelectedFilter] = useState("Today");
  const [currentPage, setCurrentPage] = useState(1);
  const limit = 5;

  const toggleSidebar = () => {
    setIsSidebarCollapsed((prev) => !prev);
  };

  const mapFilterToTimeRange = (filter) => {
    const filterMap = {
      Today: "today",
      "Last Week": "last-week",
      "Last Month": "last-month",
      "Last 3 Months": "last-three-months",
    };
    return filterMap[filter] || "today";
  };

  const handlePageChange = (newPage) => {
    if (newPage >= 1) {
      setCurrentPage(newPage);
    }
  };

  const handleFilterChange = (filter) => {
    setSelectedFilter(filter);
    setCurrentPage(1);
  };

  const offset = (currentPage - 1) * limit;
  const payload = {
    shopId: localStorage.getItem("shopId"),
    offset: offset,
    limit: limit,
    timeRange: mapFilterToTimeRange(selectedFilter),
  };

  const { data, isLoading, error } = useQuery({
    queryKey: ["history", selectedFilter, currentPage, limit],
    queryFn: () => getSessionSummary(payload),
  });

  if (isLoading) return <p>Loading history...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.HISTORY} />
      <div className="p-5 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HISTORY} />
        <HistorySubHeader
          selectedFilter={selectedFilter}
          setSelectedFilter={handleFilterChange}
        />
        <div className="flex gap-3 py-3">
          <SalesSummaryCard
            filter={selectedFilter}
            totalAmount={data?.totalAmount || 0}
          />
          <SessionsSummaryCard
            filter={selectedFilter}
            totalSessions={data?.totalSessions || 0}
          />
        </div>
        <HistoryContent
          data={data}
          onPageChange={handlePageChange}
          currentPage={currentPage}
          limit={limit}
          totalSessions={data?.totalSessions || 0}
        />
      </div>
    </div>
  );
};

export default History;
