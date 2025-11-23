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
    <div className="w-full px-4 pt-4">
      <div className="p-3 border border-gray-200 hover:border-gray-300 rounded-md bg-white flex flex-col gap-3">
        <h2 className="text-sm font-medium">Order Items</h2>

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
              className="flex justify-between items-center"
            >
              <div className="flex items-center gap-2">
                <h2 className="text-sm text-gray-800 font-medium">
                  {displayName}
                </h2>
                <div
                  className={`${
                    isVeg
                      ? "text-green-800 border-green-800"
                      : "text-red-800 border-red-800"
                  } border rounded-sm p-1`}
                >
                  {isVeg ? (
                    <IoEllipseSharp size={6} />
                  ) : (
                    <IoTriangleSharp size={6} />
                  )}
                </div>
                <p className="text-xs text-gray-500">Qty: {item.quantity}</p>
              </div>
              <p className="text-sm text-gray-700">
                {formatCurrency(calculateLineTotal(item))}
              </p>
            </div>
          );
        })}

        {orderItems.length > 0 && (
          <div className="pt-2 mt-2 border-t border-gray-100 flex justify-between text-sm">
            <span className="font-medium text-gray-700">Subtotal</span>
            <span className="font-semibold text-gray-900">
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
    <div className="w-screen h-screen bg-gray-50">
      <Navbar />
      <BackButton />
      <OrderDetails orderItems={filteredItems} itemsMap={itemsMap} />
    </div>
  );
};

export default OrderSummary;
