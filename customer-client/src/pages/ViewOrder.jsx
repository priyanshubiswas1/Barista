import { useAuthStore } from "../store/authStore";
import Navbar from "../components/Navbar";
import { useNavigate, useParams } from "react-router-dom";
import { IoAdd, IoEllipseSharp, IoRemove } from "react-icons/io5";
import { FaCheck } from "react-icons/fa6";
import { Plus } from "phosphor-react";
import Button from "../components/Button";
import { fetchTableByTableId, placeOrder } from "../api/api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { IoIosArrowBack } from "react-icons/io";
import CircleLoading from "../components/CircleLoading";

const OrderDetails = ({ orderItems, setOrderItems }) => {
  const navigate = useNavigate();
  const { shopid, tableid } = useParams();

  const updateOrderItems = (item, newQuantity) => {
    setOrderItems((prev) => {
      if (newQuantity <= 0) return prev.filter((i) => i.itemId !== item.itemId);

      const exists = prev.find((i) => i.itemId === item.itemId);
      if (exists) {
        return prev.map((i) =>
          i.itemId === item.itemId ? { ...i, quantity: newQuantity } : i
        );
      } else {
        return [
          ...prev,
          {
            itemId: item.itemId,
            itemName: item.itemName,
            quantity: newQuantity,
            priceAtOrderTime: item.priceAtOrderTime,
          },
        ];
      }
    });
  };

  return (
    <div className="w-full max-w-3xl mx-auto px-4 pt-4">
      <div className="p-4 border-2 border-gray-200 rounded-2xl bg-white shadow-md hover:shadow-lg transition-shadow duration-200">
        <h2 className="text-base font-bold text-gray-900 mb-3">
          Order Details
        </h2>

        {orderItems.length ? (
          orderItems.map((item) => (
            <div
              key={item.itemId}
              className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0"
            >
              <div className="flex flex-1 items-center gap-2">
                <h2 className="text-sm text-gray-800 font-medium">
                  {item.itemName}
                </h2>
                <div className="text-green-700 border-2 border-green-700 rounded p-0.5">
                  <IoEllipseSharp size={6} />
                </div>
              </div>
              <div className="flex flex-1 items-center justify-end gap-3">
                <div className="bg-emerald-50 border-2 border-emerald-600 rounded-lg flex items-center justify-between w-24 px-2 py-1 shadow-sm">
                  <button
                    type="button"
                    onClick={() => updateOrderItems(item, item.quantity - 1)}
                    className="text-emerald-700 text-lg leading-none hover:scale-110 transition-transform active:scale-95"
                  >
                    <IoRemove />
                  </button>
                  <span className="text-emerald-700 text-sm font-bold min-w-[20px] text-center">
                    {item.quantity}
                  </span>
                  <button
                    type="button"
                    onClick={() => updateOrderItems(item, item.quantity + 1)}
                    className="text-emerald-700 text-lg leading-none hover:scale-110 transition-transform active:scale-95"
                  >
                    <IoAdd />
                  </button>
                </div>
                <h2 className="w-[60px] text-right font-bold text-gray-900">
                  ₹{item.priceAtOrderTime * item.quantity}
                </h2>
              </div>
            </div>
          ))
        ) : (
          <p className="text-gray-500 text-sm py-2">No items in order</p>
        )}

        <div className="pt-3">
          <Button
            onClick={() => navigate(`/${shopid}/${tableid}`)}
            variant="outline"
            icon={Plus}
          >
            Add More Items
          </Button>
        </div>
      </div>
    </div>
  );
};

const BillingDetails = ({ orderItems }) => {
  const subtotal = orderItems.reduce(
    (total, item) => total + item.quantity * item.priceAtOrderTime,
    0
  );
  const serviceTax = Math.round(subtotal * 0.1);
  const total = subtotal + serviceTax;

  return (
    <div className="w-full max-w-3xl mx-auto px-4 pt-3">
      <div className="p-4 border-2 border-gray-200 rounded-2xl bg-white shadow-md">
        <h2 className="text-base font-bold text-gray-900 mb-3">
          Billing Details
        </h2>
        <div className="flex flex-col py-2 gap-2 text-sm">
          <div className="flex justify-between">
            <p className="text-gray-600">Subtotal</p>
            <p className="text-gray-900 font-semibold">₹{subtotal}</p>
          </div>
        </div>
      </div>
    </div>
  );
};

const BasicDetails = ({ tableNo }) => {
  const formattedDate = new Date().toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  return (
    <div className="w-full max-w-3xl mx-auto px-4 pt-3">
      <div className="p-4 border-2 border-gray-200 rounded-2xl bg-white shadow-md">
        <h2 className="text-base font-bold text-gray-900 mb-3">
          Basic Details
        </h2>
        <div className="flex flex-col py-2 gap-2 text-sm">
          <div className="flex justify-between">
            <p className="text-gray-600">Table No.</p>
            <p className="text-gray-900 font-semibold">{tableNo}</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600">Date</p>
            <p className="text-gray-900 font-semibold">{formattedDate}</p>
          </div>
        </div>
      </div>
    </div>
  );
};

const OrderPlaced = ({ placedOrderId }) => {
  const navigate = useNavigate();
  console.log(placedOrderId);

  const BackButton = () => {
    const shopId = localStorage.getItem("shopId");
    const tableId = localStorage.getItem("tableId");
    const navigate = useNavigate();
    const { setOrderItems } = useAuthStore();

    return (
      <div className="w-full max-w-3xl mx-auto px-4 pt-4 flex justify-start items-center">
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

  return (
    <div className="flex flex-col gap-6 items-center justify-center min-h-[60vh] px-4">
      <BackButton />
      <div className="bg-white rounded-3xl shadow-2xl p-8 max-w-md w-full text-center border-2 border-emerald-100">
        <div className="bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-full w-20 h-20 text-white flex justify-center items-center text-3xl mx-auto mb-6 shadow-lg">
          <FaCheck />
        </div>
        <h2 className="font-bold text-2xl text-gray-900 mb-2">Order Placed!</h2>
        <p className="text-gray-600 mb-6">
          Your order has been successfully placed
        </p>
        <button
          onClick={() => navigate(`/order/${placedOrderId}`)}
          className="w-full bg-gradient-to-r from-emerald-600 to-emerald-700 rounded-xl text-white px-8 py-3 font-bold shadow-lg hover:shadow-xl hover:from-emerald-700 hover:to-emerald-800 transition-all duration-200 active:scale-95"
        >
          View Order Details
        </button>
      </div>
    </div>
  );
};

const ViewOrder = () => {
  const { shopid, tableid } = useParams();
  const { orderItems, setOrderItems } = useAuthStore();

  const [isOrderPlaced, setIsOrderPlaced] = useState(false);
  const [placedOrderId, setPlacedOrderId] = useState("");

  const tableData = useQuery({
    queryKey: ["table", shopid, tableid],
    queryFn: () => fetchTableByTableId({ shopid, tableid }),
  });

  const placeOrderMutation = useMutation({
    mutationFn: placeOrder,
    onSuccess: (data) => {
      const orderId = data?.order?.orderId;
      if (orderId) {
        setPlacedOrderId(orderId);
        setIsOrderPlaced(true);
      } else {
        console.error("Something went wrong.");
      }
    },
  });

  const totalPrice = orderItems.reduce(
    (t, i) => t + i.quantity * i.priceAtOrderTime,
    0
  );
  const serviceTax = Math.round(totalPrice * 0.1);
  const finalTotal = totalPrice + serviceTax;

  const handlePlaceOrder = () => {
    const payloadOrderItems = orderItems.map(({ itemName, ...rest }) => rest);
    placeOrderMutation.mutate({
      shopId: shopid,
      tableId: tableid,
      specialNotes: "",
      orderItems: payloadOrderItems,
    });
  };

  if (tableData.isLoading) return <CircleLoading text="Loading..." />;
  if (tableData.error)
    return <p className="p-4 text-red-500">Error loading data.</p>;

  const totalItemsCount = orderItems.reduce((t, i) => t + i.quantity, 0);

  return (
    <div className="w-screen min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 pb-28">
      <div className="sticky top-0 z-20 bg-white/95 backdrop-blur-md shadow-sm">
        <Navbar />
      </div>

      {isOrderPlaced ? (
        <OrderPlaced placedOrderId={placedOrderId} />
      ) : (
        <>
          <OrderDetails orderItems={orderItems} setOrderItems={setOrderItems} />
          <BillingDetails orderItems={orderItems} />
          <BasicDetails tableNo={tableData.data.table.tableNumber} />

          {orderItems.length > 0 && finalTotal > 0 && (
            <div className="fixed bottom-4 left-4 right-4 z-30 max-w-3xl mx-auto px-4 sm:px-0">
              <div className="bg-gradient-to-r from-emerald-700 to-emerald-600 text-white rounded-2xl shadow-2xl p-4 flex justify-between items-center transform transition-all hover:shadow-[0_20px_40px_-10px_rgba(5,150,105,0.5)] hover:scale-[1.02] border border-emerald-500/20">
                <div className="flex flex-col">
                  <span className="font-bold text-lg">
                    {totalItemsCount} {totalItemsCount !== 1 ? "Items" : "Item"}
                  </span>
                  <span className="text-emerald-100 text-sm">
                    Total:{" "}
                    <span className="font-bold text-white">₹{finalTotal}</span>
                  </span>
                </div>
                <button
                  onClick={handlePlaceOrder}
                  disabled={placeOrderMutation.isPending}
                  className="bg-white text-emerald-700 px-6 py-3 rounded-xl font-bold text-sm shadow-lg hover:bg-emerald-50 transition-all hover:shadow-xl disabled:opacity-60 disabled:cursor-not-allowed active:scale-95"
                >
                  {placeOrderMutation.isPending ? "Placing..." : "Place Order"}
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default ViewOrder;
