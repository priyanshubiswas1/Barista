import { useEffect, useState } from "react";
import { IoClose } from "react-icons/io5";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQueries, useQueryClient } from "@tanstack/react-query";
import { createPayment, fetchItemsById } from "../../api/api";
import Button from "../Button";

const formatCurrency = (amount) => {
  if (typeof amount !== "number") return "₹ 0";
  return `₹ ${amount.toFixed(2)}`;
};

const formatDateTime = (value) => {
  if (!value) return new Date().toLocaleDateString();
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return new Date().toLocaleDateString();
  }
  return date.toLocaleDateString();
};

const calculateLineTotal = (item) =>
  Number(item.priceAtOrderTime || 0) * Number(item.quantity || 0);

const GenerateBill = ({ header }) => {
  const { overlayData, setOverlay } = useOverlay();
  const queryClient = useQueryClient();
  const [phoneNumber, setPhoneNumber] = useState("");

  if (!overlayData) return null;

  const {
    orderIds, // This will be an array of order IDs
    orderNumbers, // This will be an array of order numbers
    tableId,
    tableNumber = "-",
    createdAt,
    orderItems = [],
    sessionId,
  } = overlayData;

  const total = orderItems.reduce(
    (acc, item) => acc + calculateLineTotal(item),
    0
  );

  const uniqueItemIds = Array.from(
    new Set(orderItems.map((item) => item.itemId).filter(Boolean))
  );

  const itemQueries = useQueries({
    queries: uniqueItemIds.map((itemId) => ({
      queryKey: ["item", itemId],
      queryFn: () => fetchItemsById(itemId),
    })),
  });

  const isItemsLoading = itemQueries.some((query) => query.isLoading);
  const itemsError = itemQueries.find((query) => query.error)?.error;
  const itemDetailsMap = uniqueItemIds.reduce((acc, itemId, index) => {
    const query = itemQueries[index];
    if (query?.data) {
      const details = query.data.item || query.data;
      acc[itemId] = details;
    }
    return acc;
  }, {});

  const createPaymentMutation = useMutation({
    mutationFn: createPayment,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      setOverlay("");
    },
    onError: (err) => {
      console.log("Error while closing order", err);
    },
  });

  const handleSendBill = () => {
    console.log("Send Bill clicked");
    // Implement send bill logic here
  };

  const handleCloseOrder = () => {
    if (phoneNumber == "") {
      alert("phone number required");
      return;
    }
    createPaymentMutation.mutate({
      sessionId: sessionId,
      tableId: tableId,
      shopId: localStorage.getItem("shopId"),
      phoneNumber: phoneNumber,
      amount: total,
      paymentMethod: "CASH",
      paymentStatus: "completed",
    });
  };

  return (
    <div className="flex flex-col w-full h-full relative">
      <div className="flex flex-col w-full h-full overflow-auto">
        <div className="px-6 pt-6 pb-4 flex justify-between items-start">
          {header && <h2 className="font-semibold text-2xl">{header}</h2>}
          <button
            onClick={() => setOverlay("")}
            className="text-gray-600 hover:text-gray-900"
            aria-label="Close overlay"
          >
            <IoClose size={24} />
          </button>
        </div>

        <div className="w-full h-px bg-gray-200" />

        <div className="w-full h-full flex flex-col p-4 gap-3">
          <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-2">
            <h2 className="font-semibold">Basic Details</h2>
            <div className="flex justify-between">
              <p>Order ID</p>
              <p>{orderNumbers ? orderNumbers.join(", ") : "-"}</p>
            </div>
            <div className="flex justify-between">
              <p>Table No.</p>
              <p>{tableNumber}</p>
            </div>
            <div className="flex justify-between">
              <p>Date</p> <p>{formatDateTime(createdAt)}</p>
            </div>
          </div>

          <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-2">
            <h2 className="font-semibold">Customer Details</h2>
            <div className="flex justify-between items-center">
              <p>Phone Number</p>
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  value={phoneNumber}
                  onChange={(e) => setPhoneNumber(e.target.value)}
                  placeholder="+ Add Phone Number"
                  className="placeholder-emerald-700 rounded-md px-2 py-1 text-sm w-40"
                />
              </div>
            </div>
          </div>

          <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-3">
            <h2 className="font-semibold">Order Items</h2>
            {isItemsLoading && (
              <p className="text-gray-500 text-sm">
                Fetching latest item details…
              </p>
            )}
            {itemsError && (
              <p className="text-red-500 text-sm">
                {itemsError.message || "Failed to load item details."}
              </p>
            )}
            {orderItems.length === 0 && (
              <p className="text-gray-500">No items in this order.</p>
            )}
            {orderItems.map((item) => (
              <div
                key={`${item.orderItemId || item.itemId}-${item.itemName}`}
                className="flex justify-between items-center"
              >
                <div>
                  <p className="font-medium text-gray-800">
                    {item.itemName ||
                      itemDetailsMap[item.itemId]?.itemName ||
                      itemDetailsMap[item.itemId]?.name ||
                      "Item"}
                  </p>
                  <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                </div>
                <p className="text-gray-700">
                  {formatCurrency(calculateLineTotal(item))}
                </p>
              </div>
            ))}
          </div>

          <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-2">
            <h2 className="font-semibold">Billing Details</h2>
            <div className="flex justify-between">
              <p>Subtotal</p>
              <p>{formatCurrency(total)}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="flex justify-end gap-3 p-4 border-t border-gray-200">
        <Button variant="subtle" onClick={handleSendBill}>
          Send Bill
        </Button>
        <Button
          variant="primary"
          onClick={handleCloseOrder}
          disabled={createPaymentMutation.isPending}
        >
          {createPaymentMutation.isPending ? "Closing..." : "Close Order"}
        </Button>
      </div>
    </div>
  );
};

export default GenerateBill;
