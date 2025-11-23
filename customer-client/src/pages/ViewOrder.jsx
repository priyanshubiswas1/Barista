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
    <div className="w-full px-4 pt-4">
      <div className="p-3 border border-gray-200 rounded-md bg-white">
        <h2 className="text-sm font-medium">Order Details</h2>

        {orderItems.length ? (
          orderItems.map((item) => (
            <div
              key={item.itemId}
              className="flex items-center justify-between py-1"
            >
              <div className="flex flex-1 items-center gap-2">
                <h2 className="text-sm text-gray-600">{item.itemName}</h2>
                <div className="text-green-800 border border-green-800 rounded-sm p-1">
                  <IoEllipseSharp size={6} />
                </div>
              </div>
              <div className="flex flex-1 items-center justify-end gap-2">
                <div className="bg-rose-100 border border-rose-300 rounded-md flex items-center justify-between w-20 px-2 shadow-sm">
                  <button
                    type="button"
                    onClick={() => updateOrderItems(item, item.quantity - 1)}
                    className="text-rose-700 text-xl leading-none"
                  >
                    <IoRemove />
                  </button>
                  <span className="text-rose-700 text-sm font-medium">
                    {item.quantity}
                  </span>
                  <button
                    type="button"
                    onClick={() => updateOrderItems(item, item.quantity + 1)}
                    className="text-rose-700 text-xl leading-none"
                  >
                    <IoAdd />
                  </button>
                </div>
                <h2 className="w-[30px] text-right font-medium">
                  ₹{item.priceAtOrderTime * item.quantity}
                </h2>
              </div>
            </div>
          ))
        ) : (
          <p className="text-gray-500 text-sm py-2">No items in order</p>
        )}

        <div className="pt-2">
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
    <div className="w-full px-4 pt-2">
      <div className="p-3 border border-gray-200 rounded-md bg-white">
        <h2 className="text-sm font-medium">Billing Details</h2>
        <div className="flex flex-col py-2 gap-1 text-sm">
          <div className="flex justify-between">
            <p className="text-gray-600">Subtotal</p>
            <p className="text-gray-900 font-medium">₹{subtotal}</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600">Service Tax (10%)</p>
            <p className="text-gray-900 font-medium">₹{serviceTax}</p>
          </div>
        </div>
        <div className="border-t border-gray-200 mt-2 pt-2 flex justify-between font-medium">
          <h3>Total</h3>
          <h3>₹{total}</h3>
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
    <div className="w-full px-4 pt-2">
      <div className="p-3 border border-gray-200 rounded-md bg-white">
        <h2 className="text-sm font-medium">Basic Details</h2>
        <div className="flex flex-col py-2 gap-1 text-sm">
          <div className="flex justify-between">
            <p className="text-gray-600">Table No.</p>
            <p className="text-gray-900 font-medium">{tableNo}</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600">Date</p>
            <p className="text-gray-900 font-medium">{formattedDate}</p>
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
      <div className="w-full px-4 pt-4 flex justify-start items-center">
        <button
          onClick={() => {
            setOrderItems([]); // reset selections
            navigate(`/${shopId}/${tableId}`);
          }}
          className="bg-white border border-gray-200 rounded-md flex justify-center items-center text-xl text-gray-300 p-2"
        >
          <IoIosArrowBack />
        </button>
      </div>
    );
  };

  return (
    <div className="flex flex-col gap-4 items-center">
      <BackButton />
      <div className="bg-emerald-500 rounded-full w-10 h-10 text-white flex justify-center items-center text-2xl">
        <FaCheck />
      </div>
      <h2 className="font-semibold">Your Order has been placed!</h2>
      <button
        onClick={() => navigate(`/order/${placedOrderId}`)}
        className="bg-emerald-700 rounded-md text-white px-10 py-2"
      >
        View Order Details
      </button>
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

  if (tableData.isLoading)
    return <p className="p-4">Loading...</p>;
  if (tableData.error)
    return <p className="p-4 text-red-500">Error loading data.</p>;

  return (
    <div className="w-screen h-screen">
      <Navbar />

      {isOrderPlaced ? (
        <OrderPlaced placedOrderId={placedOrderId} />
      ) : (
        <>
          <OrderDetails orderItems={orderItems} setOrderItems={setOrderItems} />
          <BillingDetails orderItems={orderItems} />
          <BasicDetails tableNo={tableData.data.table.tableNumber} />

          <div className="fixed bottom-0 left-0 right-0">
            <div className="bg-emerald-700 text-white py-3 flex justify-between px-5 items-center">
              <span className="text-sm">
                {orderItems.reduce((t, i) => t + i.quantity, 0)} items |{" "}
                <span className="font-semibold">₹{finalTotal}</span>
              </span>

              <button
                onClick={handlePlaceOrder}
                disabled={placeOrderMutation.isPending}
                className="px-4 py-2 rounded-md font-semibold text-sm bg-white text-emerald-700"
              >
                {placeOrderMutation.isPending ? "Placing..." : "Place Order"}
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default ViewOrder;
