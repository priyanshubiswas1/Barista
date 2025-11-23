import { useMemo } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import Navbar from "../components/Navbar";
import { IoEllipseSharp, IoTriangleSharp } from "react-icons/io5";
import { IoIosArrowBack } from "react-icons/io";
import { getOrdersBySessionId, fetchItemsByShopId } from "../api/api";
import { useAuthStore } from "../store/authStore";

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

const formatCurrency = (amount) => {
  if (typeof amount !== "number") return "₹ 0";
  return `₹ ${amount.toFixed(2)}`;
};

const calculateLineTotal = (item) =>
  Number(item.priceAtOrderTime || 0) * Number(item.quantity || 0);

const OrderDetails = ({ orderItems, itemsMap }) => {
  const total = useMemo(
    () => orderItems.reduce((acc, item) => acc + calculateLineTotal(item), 0),
    [orderItems]
  );

  return (
    <div className="w-full max-w-3xl mx-auto px-4 pt-4">
      <div className="p-4 border-2 border-gray-200 hover:border-emerald-200 rounded-2xl bg-white flex flex-col gap-3 shadow-md hover:shadow-lg transition-all duration-200">
        <h2 className="text-base font-bold text-gray-900">Order Items</h2>

        {orderItems.length === 0 && (
          <p className="text-gray-500 text-sm py-2">No items in order</p>
        )}

        {orderItems.map((item) => {
          const itemInfo = itemsMap[item.itemId] || {};
          const displayName = itemInfo.itemName ?? item.itemName ?? "Item";
          const isVeg = itemInfo.veg ?? item.veg ?? false;

          return (
            <div
              key={item.itemId}
              className="flex justify-between items-center py-2 border-b border-gray-100 last:border-0"
            >
              <div className="flex items-center gap-2">
                <h2 className="text-sm text-gray-800 font-bold">
                  {displayName}
                </h2>
                <div
                  className={`${
                    isVeg
                      ? "text-green-700 border-green-700"
                      : "text-red-700 border-red-700"
                  } border-2 rounded p-0.5`}
                >
                  {isVeg ? (
                    <IoEllipseSharp size={6} />
                  ) : (
                    <IoTriangleSharp size={6} />
                  )}
                </div>
                <p className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full">Qty: {item.quantity}</p>
              </div>
              <p className="text-sm text-gray-900 font-semibold">
                {formatCurrency(calculateLineTotal(item))}
              </p>
            </div>
          );
        })}

        {orderItems.length > 0 && (
          <div className="pt-2 mt-2 border-t-2 border-gray-100 flex justify-between text-sm">
            <span className="font-semibold text-gray-700">Subtotal</span>
            <span className="font-bold text-emerald-700 text-base">
              {formatCurrency(total)}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

const OrderSummary = () => {
  const { orderid } = useParams();

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

  if (isLoading) return <p>Loading menu...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const orders = Array.isArray(data?.orderItems)
    ? data.orderItems
    : Object.values(data?.orderItems || {});

  const filteredItems = orders.filter(
    (item) => String(item.orderId) === String(orderid)
  );

  return (
    <div className="w-screen h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="sticky top-0 z-10 bg-white shadow-sm">
        <Navbar />
      </div>
      <BackButton />
      <OrderDetails orderItems={filteredItems} itemsMap={itemsMap} />
    </div>
  );
};

export default OrderSummary;
