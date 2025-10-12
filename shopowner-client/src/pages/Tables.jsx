import { useState } from "react";
import { TabHeader } from "../components/TabHeader";
import { Overlays, Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useOverlay } from "../store/overlayStore";
import Searchbar from "../components/Searchbar";
import Button from "../components/Button";
import { IoFilter } from "react-icons/io5";
import { FaPlus } from "react-icons/fa";
import { useQuery } from "@tanstack/react-query";
import { fetchTablesByShopId } from "../api/api";
import { FaEllipsis } from "react-icons/fa6";
import { PiQrCodeLight } from "react-icons/pi";
import QRView from "../components/overlaysidebar/QRView";

const TablesContent = ({ table, onViewQR }) => {
  return (
    <div className="bg-white border-2 border-gray-200 hover:border-gray-300 rounded-md w-[20%] p-3 flex flex-col gap-2">
      <div className="flex justify-between">
        Table: {table.tableNumber}
        <div className="text-gray-500">
          <FaEllipsis />
        </div>
      </div>
      <div className="flex justify-between text-gray-500 items-center">
        <h2>Status</h2>
        <div className="border-2 border-gray-200 rounded-md hover:border-gray-300 inline-block px-2">
          occupied
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
