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
import CircleLoading from "../components/CircleLoading";

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
    <div className="w-full flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 pt-2 pb-4 sm:pb-6">
      <div className="flex items-center gap-3">
        {/* <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button> */}
      </div>

      <div className="flex items-center gap-3 w-full sm:w-auto">
        <div className="relative" ref={dropdownRef}>
          <button
            type="button"
            onClick={toggleDropdown}
            className="rounded-md shadow-xs disabled:shadow-none text-sm font-medium cursor-pointer disabled:cursor-not-allowed transition-colors duration-120 ease-in-out
                 bg-gray-200 text-gray-500 hover:text-gray-700 disabled:bg-gray-100 disabled:text-gray-200
                 w-full sm:w-48 h-10 sm:h-8 flex items-center gap-2 justify-center"
          >
            <IoCalendarClearOutline size={16} color="currentColor" />
            {selectedFilter}
            <IoIosArrowDown size={16} color="currentColor" />
          </button>

          {isDropdownOpen && (
            <div className="absolute right-0 mt-2 w-full sm:w-48 bg-white border border-gray-200 rounded-md shadow-md z-10">
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
    <div className="bg-white border-2 p-4 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col gap-y-4 w-full sm:w-64 shadow-sm hover:shadow-md transition-all duration-200">
      <div className="flex justify-between gap-10 text-xs font-semibold items-center">
        <h2 className="text-gray-600">Total Sales</h2>
        <h2 className="bg-indigo-100 border-indigo-200 rounded-lg px-3 py-1 text-indigo-700">
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
    <div className="bg-white border-2 p-4 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col gap-y-4 w-full sm:w-64 shadow-sm hover:shadow-md transition-all duration-200">
      <div className="flex justify-between gap-10 text-xs font-semibold items-center">
        <h2 className="text-gray-600">Total Sessions</h2>
        <h2 className="bg-indigo-100 border-indigo-200 rounded-lg px-3 py-1 text-indigo-700">
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
    <div className="flex flex-col flex-1 overflow-hidden">
      <div className="hidden sm:flex bg-gray-100 p-2 text-sm font-medium text-gray-600 flex-shrink-0 rounded-md">
        <div className="flex-[2]">Order ID</div>
        <div className="flex-1">Items</div>
        <div className="flex-1">Session</div>
        <div className="flex-1 text-right">Amount</div>
      </div>
      <div className="w-full flex-1 overflow-y-auto min-h-0 flex flex-col gap-2">
        {sessionSummaries.map((session, index) => (
          <div
            key={index}
            className="flex flex-col sm:flex-row bg-white p-3 sm:p-2 text-sm text-gray-700 border border-gray-200 rounded-xl sm:rounded-md"
          >
            <div className="flex items-center justify-between sm:contents gap-2">
              <div className="sm:flex-[2] font-semibold text-gray-900">
                {session.sessionNumber}
              </div>
              <div className="sm:hidden text-gray-500 text-xs">
                {calculateSessionTime(session.createdAt, session.updatedAt)}
              </div>
            </div>
            <div className="mt-2 sm:mt-0 flex items-center justify-between sm:contents gap-2">
              <div className="sm:flex-1 text-gray-600">
                <span className="sm:hidden text-gray-500">Items: </span>
                {session.orderedItems.length}
              </div>
              <div className="hidden sm:block sm:flex-1">
                {calculateSessionTime(session.createdAt, session.updatedAt)}
              </div>
              <div className="sm:flex-1 sm:text-right font-semibold text-gray-900">
                ₹{session.totalAmount.toFixed(2)}
              </div>
            </div>
          </div>
        ))}
      </div>
      {totalPages > 1 && (
        <div className="flex flex-col sm:flex-row justify-between items-center mt-4 pt-4 border-t border-gray-200 flex-shrink-0 gap-3">
          <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className={`w-full sm:w-auto flex items-center justify-center gap-2 px-4 py-2 rounded-md text-sm font-medium ${
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
            className={`w-full sm:w-auto flex items-center justify-center gap-2 px-4 py-2 rounded-md text-sm font-medium ${
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
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
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

  if (isLoading) return <CircleLoading text="Loading history..." />;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div className="flex h-screen">
      <LeftSidebar
        collapsed={isSidebarCollapesed}
        tab={Tabs.HISTORY}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div className="px-4 sm:px-6 pt-4 sm:pt-6 flex flex-col flex-1 overflow-hidden">
        <div className="flex-shrink-0">
          <TabHeader
            toggleSidebar={toggleSidebar}
            tab={Tabs.HISTORY}
            onMobileMenuClick={() => setMobileMenuOpen(true)}
          />
        </div>
        <div className="flex-shrink-0">
          <HistorySubHeader
            selectedFilter={selectedFilter}
            setSelectedFilter={handleFilterChange}
          />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pb-4 flex-shrink-0">
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
