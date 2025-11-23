import { useMemo } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import Navbar from "../components/Navbar";
import Button from "../components/Button";
import { IoEllipseSharp } from "react-icons/io5";
import { IoIosArrowBack } from "react-icons/io";
import { getOrdersBySessionId } from "../api/api";

const BackButton = () => {
  const shopId = localStorage.getItem("shopId");
  const tableId = localStorage.getItem("tableId");
  const navigate = useNavigate();

  return (
    <div className="w-full px-4 pt-4 flex justify-start items-center">
      <button
        onClick={() => navigate(`/${shopId}/${tableId}`)}
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

const OrderDetails = ({ orderItems }) => {
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
        {orderItems.map((item) => (
          <div key={item.itemId} className="flex justify-between items-center">
            <div className="flex items-center gap-2">
              <h2 className="text-sm text-gray-800 font-medium">
                {item.itemName || "Item"}
              </h2>
              <div className="text-green-800 border border-green-800 rounded-sm p-1">
                <IoEllipseSharp size={6} />
              </div>
              <p className="text-xs text-gray-500">Qty: {item.quantity}</p>
            </div>
            <p className="text-sm text-gray-700">
              {formatCurrency(calculateLineTotal(item))}
            </p>
          </div>
        ))}

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

  if (isLoading) return <p>Loading menu...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const orders = Array.isArray(data.orderItems)
    ? data.orderItems
    : Object.values(data.orderItems || {});

  const filteredItems = orders.filter(
    (item) => String(item.orderId) === String(orderid)
  );

  return (
    <div className="w-screen h-screen bg-gray-50">
      <Navbar />
      <BackButton />
      <OrderDetails orderItems={filteredItems} />
    </div>
  );
};

export default OrderSummary;
