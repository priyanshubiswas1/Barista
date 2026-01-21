import { useEffect, useState } from "react";
import { IoFilter } from "react-icons/io5";
import { Plus, DotsThree } from "phosphor-react";
import { LuCirclePlus } from "react-icons/lu";
import { FaArrowRightLong } from "react-icons/fa6";
import { Overlays, Tabs } from "../constants/constants";
import LeftSidebar from "../components/LeftSidebar.jsx";
import CircleLoading from "../components/CircleLoading";
import { TabHeader } from "../components/TabHeader.jsx";
import Searchbar from "../components/Searchbar.jsx";
import Button from "../components/Button.jsx";
import { useOverlay } from "../store/overlayStore.js";
import { OverlaySidebar } from "../components/overlaysidebar/OverlaySidebar";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { fetchActiveSessions, fetchOrdersBySessionId } from "../api/api.js";

const getStatusClasses = (status) => {
  switch (status) {
    case "new":
      return "bg-blue-400";
    case "ready":
      return "bg-purple-400";
    case "completed":
      return "bg-teal-400";
    case "cancelled":
      return "bg-red-400";
    default:
      return "bg-gray-400";
  }
};

const OrderSubHeader = ({ onAddNewOrderClick }) => {
  return (
    <div className="w-full flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 pt-2 pb-4 sm:pb-6 flex-shrink-0">
      <div className="flex items-center gap-3">
        {/* <Searchbar />
        <Button variant="outline" icon={IoFilter}>
          Filter
        </Button> */}
      </div>
      <div className="flex items-center gap-3 w-full sm:w-auto">
        <Button onClick={onAddNewOrderClick} variant="primary" icon={Plus}>
          Add New Order
        </Button>
      </div>
    </div>
  );
};

const OrderCard = ({ orderId, orderItems, index, onClick }) => {
  const totalQuantity = orderItems.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);
  const totalPrice = orderItems.reduce((sum, item) => {
    return sum + item.priceAtOrderTime * item.quantity;
  }, 0);
  const orderStatus = orderItems[0]?.orderStatus || "new";
  const statusBgColor = getStatusClasses(orderStatus);

  return (
    <button
      type="button"
      onClick={onClick}
      className="bg-white border-2 p-4 border-gray-200 hover:border-gray-300 rounded-xl text-left w-full shadow-sm hover:shadow-md transition-all duration-200"
    >
      <div className="flex justify-between items-center">
        <h1 className="font-bold text-xl text-gray-900">Order {orderId}</h1>
        <DotsThree className="text-gray-500" />
      </div>
      <div className="pt-3">
        <div
          className={`${statusBgColor} text-white font-semibold rounded-lg flex items-center justify-center px-4 py-1 w-fit capitalize shadow-sm`}
        >
          {orderStatus}
        </div>
        <div className="pt-3 font-medium">
          <div className="flex justify-between">
            <h2 className="text-gray-500">Items</h2>
            <h2 className="text-gray-800 font-semibold">{totalQuantity}</h2>
          </div>
          <div className="flex justify-between">
            <h2 className="text-gray-500">Subtotal</h2>
            <h2 className="text-gray-800 font-semibold">₹{totalPrice}</h2>
          </div>
        </div>
      </div>
    </button>
  );
};

const TableCard = ({ session }) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["order-items", session.tableId],
    queryFn: () => fetchOrdersBySessionId(session.sessionId),
    staleTime: 0,
  });

  const { setOverlay } = useOverlay();

  if (isLoading) return <CircleLoading text="Loading tables..." />;
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

  const handleOrderClick = (orderId, items) => {
    const firstItem = items[0] || {};
    setOverlay(Overlays.ORDERDETAILS, {
      orderId: items[0].orderId,
      orderNumber: items[0].orderNumber,
      sessionId: session.sessionId,
      tableId: session.tableId,
      tableNumber: session.tableNumber,
      orderItems: items,
      status: firstItem.orderStatus || firstItem.status || "new",
      createdAt:
        firstItem.createdAt ||
        firstItem.orderCreatedAt ||
        firstItem.placedAt ||
        null,
    });
  };

  return (
    <div className="bg-white w-full sm:w-[25%] sm:min-w-[300px] h-auto sm:h-full overflow-hidden p-4 border-2 border-gray-200 hover:border-gray-300 rounded-xl flex flex-col shadow-sm hover:shadow-md transition-all duration-200 min-w-0">
      <h1 className="text-gray-600 text-xl px-3 font-bold">
        Table {session.tableNumber}
      </h1>
      <div className="pt-3 flex flex-col gap-4 flex-1 overflow-y-auto min-h-0 pr-2">
        {Object.entries(groupedOrders).map(([orderId, items], index) => (
          <OrderCard
            key={orderId}
            orderId={items[0].orderNumber}
            orderItems={items}
            index={index}
            onClick={() => handleOrderClick(orderId, items)}
          />
        ))}
      </div>
      <div className="flex flex-col pt-3 gap-3">
        <button
          onClick={() => setOverlay(Overlays.CREATEORDER)}
          className="flex items-center justify-center border-2 border-gray-200 rounded-xl hover:border-gray-300 hover:bg-gray-50 text-gray-600 text-lg gap-3 py-3 transition-all duration-200 font-medium"
        >
          <LuCirclePlus />
          <h2>Add New Order</h2>
        </button>
        <button
          onClick={() =>
            setOverlay(Overlays.GENERATEBILL, {
              sessionId: session.sessionId,
              orderIds: Object.keys(groupedOrders),
              orderNumbers: Array.from(
                new Set(
                  Object.values(groupedOrders).flatMap((items) =>
                    items.map((item) => item.orderNumber),
                  ),
                ),
              ),
              tableId: session.tableId,
              tableNumber: session.tableNumber,
              createdAt: session.createdAt,
              orderItems: orderItems,
            })
          }
          className="flex items-center justify-center bg-gray-700 hover:bg-gray-800 rounded-xl text-white text-lg gap-3 py-3 transition-all duration-200 font-medium shadow-sm"
        >
          <FaArrowRightLong />
          <h2>Generate Bill</h2>
        </button>
      </div>
    </div>
  );
};

const OrderContent = () => {
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["sessions"],
    queryFn: () => fetchActiveSessions(localStorage.getItem("shopId")),
    select: (data) => ({
      ...data,
      sessions: [...data.sessions].sort((a, b) =>
        a.tableNumber.localeCompare(b.tableNumber, undefined, {
          numeric: true,
        }),
      ),
    }),
  });

  // WebSocket connection
  useEffect(() => {
    const wsBaseUrl = import.meta.env.VITE_WS_URL || "ws://localhost:8081/ws";
    const shopId = localStorage.getItem("shopId");

    if (!shopId) {
      console.error("No shopId found in localStorage");
      return;
    }

    const ws = new WebSocket(`${wsBaseUrl}?id=${shopId}`);

    ws.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.onmessage = async (event) => {
      try {
        const message = JSON.parse(event.data);

        if (message.type === "order_placed") {
          await queryClient.refetchQueries({ queryKey: ["order-items"] });
          await queryClient.refetchQueries({ queryKey: ["sessions"] });
        }
      } catch (err) {
        console.error("Error parsing WebSocket message:", err);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => {
      ws.close();
    };
  }, [queryClient]);

  if (isLoading) return <CircleLoading text="Loading orders..." />;
  if (error) return <p>Error: {error.message}</p>;

  const sessions = data.sessions;
  return (
    <div className="bg-gray-100 w-full h-full flex flex-col sm:flex-row sm:overflow-x-auto overflow-y-auto sm:overflow-y-hidden gap-4 sm:gap-6 pb-4 sm:pr-6">
      {sessions.map((session) => (
        <TableCard session={session} key={session.sessionId} />
      ))}
    </div>
  );
};

const Orders = () => {
  const [isSidebarCollapesed, setisSidebarCollapesed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const toggleSidebar = () => {
    setisSidebarCollapesed((prev) => !prev);
  };
  const { setOverlay } = useOverlay();
  return (
    <div className="flex h-screen">
      <LeftSidebar
        collapsed={isSidebarCollapesed}
        tab={Tabs.ORDERS}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div className="px-4 sm:px-6 pt-4 sm:pt-6 flex flex-col flex-1 overflow-hidden">
        <div className="flex-shrink-0">
          <TabHeader
            toggleSidebar={toggleSidebar}
            tab={Tabs.ORDERS}
            onMobileMenuClick={() => setMobileMenuOpen(true)}
          />
        </div>
        <OrderSubHeader
          onAddNewOrderClick={() => setOverlay(Overlays.CREATEORDER)}
        />
        <div className="pt-2 flex-grow min-w-0 overflow-hidden min-h-0">
          <OrderContent />
        </div>
      </div>
      <OverlaySidebar />
    </div>
  );
};

export default Orders;
