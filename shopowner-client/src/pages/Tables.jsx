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
        className="cursor-pointer text-gray-500"
      >
        <FaEllipsis />
      </div>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute right-0 top-6 z-50 bg-white border border-gray-200 rounded-md shadow-lg min-w-[120px] overflow-hidden"
        >
          <button
            onClick={handleEdit}
            className="w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 text-left transition-colors"
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2 text-sm text-red-700 bg-red-100 hover:bg-red-200 text-left transition-colors"
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
    <div className="bg-white border-2 border-gray-200 hover:border-gray-300 rounded-md w-[20%] p-3 flex flex-col gap-2">
      <div className="flex justify-between">
        Table: {table.tableNumber}
        <TableActions table={table} />
      </div>
      <div className="flex justify-between text-gray-500 items-center">
        <h2>Status</h2>
        <div className="border-2 border-gray-200 rounded-md hover:border-gray-300 inline-block px-2">
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
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-5">
        <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button>
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

  console.log(data);
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
      <div className="p-5 flex flex-col w-full">
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
