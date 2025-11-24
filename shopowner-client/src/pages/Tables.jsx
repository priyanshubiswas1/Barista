import { useState, useEffect, useRef } from "react";
import { TabHeader } from "../components/TabHeader";
import { Overlays, Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useOverlay } from "../store/overlayStore";
import Searchbar from "../components/Searchbar";
import Button from "../components/Button";
import { IoFilter } from "react-icons/io5";
import { FaPlus } from "react-icons/fa";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { deleteTableByShopId, fetchTablesByShopId } from "../api/api";
import { FaEllipsis } from "react-icons/fa6";
import { PiQrCodeLight } from "react-icons/pi";
import QRView from "../components/overlaysidebar/QRView";
import { toast } from "react-toastify";

const TableActions = ({ table }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);
  const buttonRef = useRef(null);
  const { setOverlay } = useOverlay();
  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    mutationFn: deleteTableByShopId,
    onSuccess: () => {
      toast("Table deleted successfully");
      queryClient.invalidateQueries({ queryKey: ["tables"] });
    },
    onError: (error) => {
      toast.error("Failed to delete table: " + error.message);
    },
  });

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  const handleEdit = () => {
    setIsOpen(false);
    setOverlay(Overlays.UPDATETABLE, table);
  };

  const handleDelete = () => {
    setIsOpen(false);
    deleteMutation.mutate({
      tableId: table.tableId,
      shopId: localStorage.getItem("shopId"),
    });
  };

  return (
    <div className="relative">
      <div
        ref={buttonRef}
        onClick={() => setIsOpen(!isOpen)}
        className="cursor-pointer text-gray-500 hover:text-gray-700 p-1 rounded hover:bg-gray-100 transition-colors"
      >
        <FaEllipsis />
      </div>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute right-0 top-8 z-50 bg-white border border-gray-200 rounded-lg shadow-xl min-w-[120px] overflow-hidden"
        >
          <button
            onClick={handleEdit}
            className="w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 text-left transition-colors font-medium"
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2.5 text-sm text-red-600 bg-red-50 hover:bg-red-100 text-left transition-colors font-medium"
          >
            Delete
          </button>
        </div>
      )}
    </div>
  );
};

const TablesContent = ({ table, onViewQR }) => {
  return (
    <div className="bg-white border-2 border-gray-200 hover:border-gray-300 rounded-xl w-[20%] min-w-[250px] p-4 flex flex-col gap-3 shadow-sm hover:shadow-md transition-all duration-200">
      <div className="flex justify-between items-center">
        <span className="font-bold text-gray-900">
          Table: {table.tableNumber}
        </span>
        <TableActions table={table} />
      </div>
      <div className="flex justify-between text-gray-600 items-center text-sm">
        <h2 className="font-medium">Status</h2>
        <div className="border-2 border-gray-300 rounded-lg hover:border-gray-400 inline-block px-3 py-1 font-medium text-gray-700 transition-colors">
          {table.status}
        </div>
      </div>
      <Button variant="subtle" icon={PiQrCodeLight} onClick={onViewQR}>
        View QR
      </Button>
    </div>
  );
};

const TableSubHeader = ({ onAddNewTableClick }) => {
  return (
    <div className="w-[100%] flex justify-between items-center pt-2">
      <div className="flex items-center gap-5">
        {/* <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button> */}
      </div>
      <div className="flex items-center gap-5">
        <Button onClick={onAddNewTableClick} variant="primary" icon={FaPlus}>
          Add New Table
        </Button>
      </div>
    </div>
  );
};

const Tables = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const { setOverlay } = useOverlay();
  const [selectedTable, setSelectedTable] = useState(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ["tables"],
    queryFn: () => fetchTablesByShopId(localStorage.getItem("shopId")),
  });

  if (isLoading) return <p>Loading tables...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const tables = data.tables;

  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };

  const openQRModal = (table) => {
    setSelectedTable(table);
  };

  const closeQRModal = () => {
    setSelectedTable(null);
  };

  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.TABLES} />
      <div className="p-6 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.TABLES} />
        <TableSubHeader
          onAddNewTableClick={() => setOverlay(Overlays.CREATETABLE)}
        />
        <div className="pt-5 flex-grow">
          {tables.map((table) => (
            <TablesContent
              key={table.tableId}
              table={table}
              onViewQR={() => openQRModal(table)}
            />
          ))}
        </div>
      </div>
      <OverlaySidebar />
      {selectedTable && <QRView table={selectedTable} onClose={closeQRModal} />}
    </div>
  );
};

export default Tables;
