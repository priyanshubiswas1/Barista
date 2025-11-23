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
    <div className="w-full px-4 pt-4 flex justify-start items-center">
      <button
        onClick={() => {
          setOrderItems([]);
          navigate(`/${shopId}/${tableId}`);
        }}
        className="bg-white border border-gray-200 rounded-md flex justify-center items-center text-xl text-gray-300 p-2"
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
    <div className="p-4 border border-gray-200 rounded-md bg-white flex flex-col gap-3">
      <div className="flex justify-between items-center">
        <h3 className="font-medium">
          Order #{order.orderItems[0].orderNumber}
        </h3>
        <span className="text-sm text-gray-600">
          {formatCurrency(subtotal)}
        </span>
      </div>

      <div className="flex flex-col gap-2">
        {items.map((it) => {
          const itemInfo = itemsMap[it.itemId] || {};
          const isVeg = itemInfo.veg ?? it.veg ?? false;
          const itemName = itemInfo.itemName ?? it.itemName ?? "Item";

          return (
            <div key={it.itemId} className="flex justify-between items-center">
              <div className="flex items-center gap-2">
                <h4 className="text-sm text-gray-800">{itemName}</h4>
                <div
                  className={`${
                    isVeg
                      ? "text-green-800 border-green-800"
                      : "text-red-800 border-red-800"
                  } border rounded-sm p-1`}
                >
                  {isVeg ? (
                    <IoEllipseSharp size={8} />
                  ) : (
                    <IoTriangleSharp size={8} />
                  )}
                </div>
                <p className="text-xs text-gray-500">Qty: {it.quantity}</p>
              </div>
              <p className="text-sm text-gray-700">
                {formatCurrency(calculateLineTotal(it))}
              </p>
            </div>
          );
        })}
      </div>

      <div className="pt-2 mt-2 border-t border-gray-100 flex justify-between text-sm">
        <span className="font-medium text-gray-700">Subtotal</span>
        <span className="font-semibold text-gray-900">
          {formatCurrency(subtotal)}
        </span>
      </div>
    </div>
  );
};

const TotalCard = ({ totalAmount }) => (
  <div className="p-4 border border-gray-300 rounded-md bg-white flex justify-between items-center font-semibold text-lg">
    <span>Total</span>
    <span>{formatCurrency(totalAmount)}</span>
  </div>
);

const Orders = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["orders"],
    queryFn: () => getOrdersBySessionId(),
  });

  const shopId = localStorage.getItem("shopId");
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
    <div className="w-screen h-screen bg-gray-50 flex flex-col">
      <div className="sticky top-0 z-10 bg-white">
        <Navbar />
        <BackButton />
      </div>

      <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
        {ordersArr.length === 0 && (
          <p className="text-gray-500 text-sm">No orders found</p>
        )}
        {ordersArr.map((order) => (
          <OrderCard key={order.orderId} order={order} itemsMap={itemsMap} />
        ))}
      </div>

      <div className="sticky bottom-0 z-10 bg-white border-t border-gray-300 p-4">
        <TotalCard totalAmount={totalAmount} />
      </div>
    </div>
  );
};

export default Orders;
