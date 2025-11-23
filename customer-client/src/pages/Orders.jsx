import React, { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import Navbar from "../components/Navbar";
import { IoEllipseSharp, IoTriangleSharp } from "react-icons/io5";
import { IoIosArrowBack } from "react-icons/io";
import { getOrdersBySessionId, fetchItemsByShopId } from "../api/api";
import { useAuthStore } from "../store/authStore";

const formatCurrency = (amount) => {
  if (typeof amount !== "number") return "₹ 0";
  return `₹ ${amount.toFixed(2)}`;
};

const calculateLineTotal = (item) =>
  Number(item.priceAtOrderTime || 0) * Number(item.quantity || 0);

const BackButton = () => {
  const shopId = localStorage.getItem("shopId");
  const tableId = localStorage.getItem("tableId");
  const navigate = useNavigate();
  const { setOrderItems } = useAuthStore();

  return (
    <div className="w-full max-w-3xl mx-auto px-4 py-3 flex justify-start items-center">
      <button
        onClick={() => {
          setOrderItems([]);
          navigate(`/${shopId}/${tableId}`);
        }}
        className="bg-white border-2 border-gray-300 rounded-lg flex justify-center items-center text-xl text-gray-600 p-2 hover:border-emerald-600 hover:text-emerald-600 hover:bg-emerald-50 transition-all duration-200"
      >
        <IoIosArrowBack />
      </button>
    </div>
  );
};

const OrderCard = ({ order, itemsMap }) => {
  const items = Array.isArray(order.orderItems)
    ? order.orderItems
    : Object.values(order.orderItems || {});
  const subtotal = useMemo(
    () =>
      items.reduce((acc, it) => {
        return acc + calculateLineTotal(it);
      }, 0),
    [items]
  );

  return (
    <div className="p-4 border border-gray-200 rounded-2xl bg-white flex flex-col gap-3 shadow-md hover:shadow-lg transition-all duration-200 hover:border-emerald-200">
      <div className="flex justify-between items-center">
        <h3 className="font-bold text-gray-900 text-base">
          Order #{order.orderItems[0].orderNumber}
        </h3>
      </div>

      <div className="flex flex-col gap-3">
        {items.map((it) => {
          const itemInfo = itemsMap[it.itemId] || {};
          const isVeg = itemInfo.veg ?? it.veg ?? false;
          const itemName = itemInfo.itemName ?? it.itemName ?? "Item";

          return (
            <div key={it.itemId} className="flex justify-between items-center py-1">
              <div className="flex items-center gap-2">
                <h4 className="text-sm text-gray-800 font-medium">{itemName}</h4>
                <div
                  className={`${
                    isVeg
                      ? "text-green-700 border-green-700"
                      : "text-red-700 border-red-700"
                  } border-2 rounded p-1`}
                >
                  {isVeg ? (
                    <IoEllipseSharp size={8} />
                  ) : (
                    <IoTriangleSharp size={8} />
                  )}
                </div>
                <p className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full">Qty: {it.quantity}</p>
              </div>
              <p className="text-sm text-gray-900 font-semibold">
                {formatCurrency(calculateLineTotal(it))}
              </p>
            </div>
          );
        })}
      </div>

      <div className="pt-2 mt-2 border-t-2 border-gray-100 flex justify-between text-sm">
        <span className="font-semibold text-gray-700">Subtotal</span>
        <span className="font-bold text-emerald-700 text-base">
          {formatCurrency(subtotal)}
        </span>
      </div>
    </div>
  );
};

const TotalCard = ({ totalAmount }) => (
  <div className="p-6 border-2 border-emerald-200 rounded-2xl bg-gradient-to-br from-emerald-50 to-white flex justify-between items-center shadow-lg">
    <span className="font-bold text-emerald-900 text-lg">Grand Total</span>
    <span className="font-bold text-emerald-700 text-3xl">
      {formatCurrency(totalAmount)}
    </span>
  </div>
);

const Orders = () => {
  const navigate = useNavigate();
  const shopId = localStorage.getItem("shopId");
  const tableId = localStorage.getItem("tableId");
  const sessionToken = localStorage.getItem("session-token");

  if (!sessionToken) {
    return (
      <div className="w-screen h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex flex-col">
        <div className="sticky top-0 z-10 bg-white shadow-sm">
          <Navbar />
          <BackButton />
        </div>
        <div className="flex-1 flex flex-col items-center justify-center p-4 max-w-md mx-auto">
          <div className="bg-white rounded-2xl shadow-xl p-8 text-center border border-gray-200">
            <div className="w-16 h-16 bg-gray-100 rounded-full mx-auto mb-4 flex items-center justify-center">
              <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <p className="text-gray-700 text-base mb-2 font-semibold">No orders yet</p>
            <p className="text-gray-500 text-sm mb-6">Place your first order to see it here!</p>
            <button
              onClick={() => navigate(`/${shopId}/${tableId}`)}
              className="px-8 py-3 bg-gradient-to-r from-emerald-600 to-emerald-700 text-white rounded-xl font-semibold shadow-lg hover:shadow-xl hover:from-emerald-700 hover:to-emerald-800 transition-all duration-200 active:scale-95"
            >
              Browse Menu
            </button>
          </div>
        </div>
      </div>
    );
  }

  const { data, isLoading, error } = useQuery({
    queryKey: ["orders"],
    queryFn: () => getOrdersBySessionId(),
  });

  const { data: itemsData } = useQuery({
    queryKey: ["items", shopId],
    queryFn: () => fetchItemsByShopId(shopId),
    enabled: !!shopId,
  });

  const itemsArray = Array.isArray(itemsData?.items)
    ? itemsData.items
    : Array.isArray(itemsData)
    ? itemsData
    : [];

  const itemsMap = useMemo(() => {
    const map = {};
    itemsArray.forEach((it) => {
      const id = it.itemId ?? it.id;
      map[id] = {
        itemName: it.itemName ?? it.name ?? it.title ?? "Item",
        veg: it.veg ?? false,
      };
    });
    return map;
  }, [itemsArray]);

  if (isLoading) return <p>Loading orders...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const orderItemsArray = Array.isArray(data?.orderItems)
    ? data.orderItems
    : [];

  const groupedOrders = orderItemsArray.reduce((acc, item) => {
    const { orderId } = item;
    if (!acc[orderId]) {
      acc[orderId] = {
        orderId,
        orderItems: [],
      };
    }
    acc[orderId].orderItems.push(item);
    return acc;
  }, {});

  const ordersArr = Object.values(groupedOrders);

  const totalAmount = ordersArr.reduce((acc, order) => {
    const subtotal = order.orderItems.reduce((sum, it) => {
      return sum + calculateLineTotal(it);
    }, 0);
    return acc + subtotal;
  }, 0);

  return (
    <div className="w-screen h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex flex-col">
      <div className="sticky top-0 z-10 bg-white shadow-sm">
        <Navbar />
        <BackButton />
      </div>

      <div className="flex-1 overflow-y-auto max-w-3xl mx-auto w-full p-4 flex flex-col gap-4">
        {ordersArr.length === 0 && (
          <div className="flex-1 flex items-center justify-center">
            <div className="bg-white rounded-2xl shadow-lg p-8 text-center border border-gray-200">
              <p className="text-gray-600 font-medium">No orders found</p>
              <p className="text-gray-400 text-sm mt-1">Your orders will appear here</p>
            </div>
          </div>
        )}
        {ordersArr.map((order) => (
          <OrderCard key={order.orderId} order={order} itemsMap={itemsMap} />
        ))}
      </div>

      <div className="sticky bottom-0 z-10 bg-white border-t border-gray-200 p-4 shadow-lg">
        <div className="max-w-3xl mx-auto">
          <TotalCard totalAmount={totalAmount} />
        </div>
      </div>
    </div>
  );
};

export default Orders;
