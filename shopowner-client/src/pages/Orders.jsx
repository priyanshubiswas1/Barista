import { useState } from "react";
import { IoFilter } from "react-icons/io5";
import { Plus } from "phosphor-react";
import { FaEllipsis } from "react-icons/fa6";
import { LuCirclePlus } from "react-icons/lu";
import { FaArrowRightLong } from "react-icons/fa6";
import { Overlays, Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar.jsx";
import { TabHeader } from "../components/TabHeader.jsx";
import Searchbar from "../components/Searchbar.jsx";
import Button from "../components/Button.jsx";
import { useOverlay } from "../store/overlayStore.js";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useQuery } from "@tanstack/react-query";
import { fetchActiveSessions, fetchOrdersBySessionId } from "../api/api.js";

const OrderSubHeader = ({ onAddNewOrderClick }) => {
  return (
    <div className="w-[100%] flex justify-between items-center pt-10">
      <div className="flex items-center gap-3">
        <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button>
      </div>
      <div className="flex items-center gap-3">
        <Button onClick={onAddNewOrderClick} variant="primary" icon={Plus}>
          Add New Order
        </Button>
      </div>
    </div>
  );
};

const OrderCard = ({ orderId, orderItems, index }) => {
  const totalQuantity = orderItems.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);
  const totalPrice = orderItems.reduce((sum, item) => {
    return sum + item.priceAtOrderTime * item.quantity;
  }, 0);
  return (
    <div className="bg-white border-2 p-3 border-gray-200 hover:border-gray-300 rounded-xl">
      <div className="flex justify-between items-center">
        <h1 className="font-semibold text-xl">Order {index}</h1>
        <FaEllipsis className="text-gray-500" />
      </div>
      <div className="pt-3">
        <button className="bg-blue-400 text-white font-semibold rounded-md flex items-center justify-center px-4 py-0.5">
          new
        </button>
        <div className="pt-3 font-medium">
          <div className="flex justify-between">
            <h2 className="text-gray-500">Items</h2>
            <h2 className="text-gray-700">{totalQuantity}</h2>
          </div>
          <div className="flex justify-between">
            <h2 className="text-gray-500">Subtotal</h2>
            <h2 className="text-gray-700">${totalPrice}</h2>
          </div>
        </div>
      </div>
    </div>
  );
};

const TableCard = ({ session }) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["order-items", session.tableId],
    queryFn: () => fetchOrdersBySessionId(session.sessionId),
  });
  if (isLoading) return <p>Loading tables...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const orderItems = data.orderItems;
  const groupedOrders = orderItems.reduce((acc, item) => {
    const { orderId } = item;
    if (!acc[orderId]) {
      acc[orderId] = [];
    }
    acc[orderId].push(item);
    return acc;
  }, {});
  console.log(groupedOrders);

  return (
    <div className="bg-white w-[25%] p-3 border-2 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col">
      <h1 className="text-gray-500 text-xl px-3 font-semibold">
        Table {session.tableNumber}
      </h1>
      <div className="pt-3 flex flex-col gap-4">
        {Object.entries(groupedOrders).map(([orderId, items], index) => (
          <OrderCard
            key={orderId}
            orderId={orderId}
            orderItems={items}
            index={index}
          />
        ))}
      </div>
      <div className="flex flex-col pt-3 gap-3">
        <button className="flex items-center justify-center border-1 border-gray-200 rounded-xl hover:border-gray-300 text-gray-500 text-xl gap-3 py-3.5">
          <LuCirclePlus />
          <h2 className="font-medium">Add New Order</h2>
        </button>
        <button className="flex items-center justify-center bg-gray-200 rounded-xl hover:bg-gray-300 text-gray-500 text-xl gap-3 py-3.5">
          <FaArrowRightLong />
          <h2 className="font-medium">Generate Bill</h2>
        </button>
      </div>
    </div>
  );
};

const OrderContent = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["sessions"],
    queryFn: () => fetchActiveSessions(localStorage.getItem("shopId")),
    select: (data) => ({
      ...data,
      sessions: [...data.sessions].sort((a, b) =>
        a.tableNumber.localeCompare(b.tableNumber, undefined, { numeric: true })
      ),
    }),
  });
  if (isLoading) return <p>Loading tables...</p>;
  if (error) return <p>Error: {error.message}</p>;

  console.log(data);
  const sessions = data.sessions;
  return (
    <div className="bg-gray-100 w-full h-full overflow-x-auto overflow-y-auto flex">
      {sessions.map((session) => (
        <TableCard session={session} key={session.sessionId} />
      ))}
    </div>
  );
};

const Orders = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  const { setOverlay } = useOverlay();
  return (
    <div className="flex h-screen">
      <LeftSidebar collapsed={isSidebarCollapesed} tab={Tabs.ORDERS} />
      <div className="p-5 flex flex-col w-full">
        <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.ORDERS} />
        <OrderSubHeader
          onAddNewOrderClick={() => setOverlay(Overlays.CREATEORDER)}
        />
        <div className="pt-5 flex-grow">
          <OrderContent />
        </div>
      </div>
      <OverlaySidebar />
    </div>
  );
};

export default Orders;
