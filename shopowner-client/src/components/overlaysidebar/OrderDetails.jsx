import { useEffect, useState } from "react";
import { IoClose } from "react-icons/io5";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQueries, useQueryClient } from "@tanstack/react-query";
import { fetchItemsById, updateOrder, cancelOrder } from "../../api/api";
import { toast } from "react-toastify";

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

const getStatusClasses = (status) => {
  switch (status) {
    case "new":
      return {
        bgColor: "bg-blue-400",
        borderColor: "border-blue-500",
        hoverBgColor: "hover:bg-blue-500",
        focusRingColor: "focus:ring-blue-300",
      };
    case "ready":
      return {
        bgColor: "bg-purple-400",
        borderColor: "border-purple-500",
        hoverBgColor: "hover:bg-purple-500",
        focusRingColor: "focus:ring-purple-300",
      };
    case "completed":
      return {
        bgColor: "bg-teal-400",
        borderColor: "border-teal-500",
        hoverBgColor: "hover:bg-teal-500",
        focusRingColor: "focus:ring-teal-300",
      };
    case "cancelled":
      return {
        bgColor: "bg-red-400",
        borderColor: "border-red-500",
        hoverBgColor: "hover:bg-red-500",
        focusRingColor: "focus:ring-red-300",
      };
    default:
      return {
        bgColor: "bg-gray-400",
        borderColor: "border-gray-500",
        hoverBgColor: "hover:bg-gray-500",
        focusRingColor: "focus:ring-gray-300",
      };
  }
};

const OrderDetails = ({ header }) => {
  const { overlayData, setOverlay } = useOverlay();
  const queryClient = useQueryClient();
  const [isStatusMenuOpen, setIsStatusMenuOpen] = useState(false);
  const [showCancelConfirm, setShowCancelConfirm] = useState(false);
  const statusOptions = ["new", "ready", "completed", "cancelled"];

  if (!overlayData) return null;

  const {
    orderId,
    orderNumber,
    tableId,
    tableNumber = "-",
    status: initialStatus = "new",
    createdAt,
    orderItems = [],
  } = overlayData;

  const [currentStatus, setCurrentStatus] = useState(initialStatus);

  useEffect(() => {
    setCurrentStatus(initialStatus);
  }, [initialStatus]);

  const updateOrderMutation = useMutation({
    mutationFn: updateOrder,
    onSuccess: (data, variables) => {
      console.log("Order update successful:", data);
      toast("Order updated✅");
      setCurrentStatus(variables.orderStatus);
      queryClient.invalidateQueries({ queryKey: ["order-items"] });
      if (tableId) {
        queryClient.invalidateQueries({
          queryKey: ["order-items", tableId],
        });
      }
    },
    onError: (error) => {
      console.error("Update Order failed: ", error);
      toast.error("Failed to update order. Please try again.");
    },
  });

  const cancelOrderMutation = useMutation({
    mutationFn: cancelOrder,
    onSuccess: (data) => {
      console.log("Order cancelled successfully:", data);
      toast("Order cancelled✅");
      setCurrentStatus("cancelled");
      queryClient.invalidateQueries({ queryKey: ["order-items"] });
      if (tableId) {
        queryClient.invalidateQueries({
          queryKey: ["order-items", tableId],
        });
      }
      setShowCancelConfirm(false);
    },
    onError: (error) => {
      console.error("Cancel Order failed: ", error);
      toast.error("Failed to cancel order. Please try again.");
      setShowCancelConfirm(false);
    },
  });

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

  useEffect(() => {
    setIsStatusMenuOpen(false);
  }, [currentStatus]);

  const handleStatusClick = (nextStatus) => {
    if (nextStatus === currentStatus) {
      console.log(`Status already set to ${nextStatus}`);
      return;
    }
    if (!orderId) {
      console.error("Order ID is missing, cannot update order", {
        overlayData,
      });
      toast.error("Error: Order ID is missing. Please refresh and try again.");
      return;
    }

    if (nextStatus === "cancelled") {
      setShowCancelConfirm(true);
      return;
    }

    console.log(`Status changed to ${nextStatus}`, {
      orderId,
      orderStatus: nextStatus,
    });
    setIsStatusMenuOpen(false);
    console.log("Calling updateOrderMutation.mutate with:", {
      orderId,
      orderStatus: nextStatus,
    });
    updateOrderMutation.mutate({
      orderId: orderId,
      orderStatus: nextStatus,
    });
  };

  const confirmCancelOrder = () => {
    if (!orderId) {
      console.error("Order ID is missing, cannot cancel order");
      toast.error("Error: Order ID is missing. Please refresh and try again.");
      return;
    }
    cancelOrderMutation.mutate({ orderId });
  };

  const { bgColor, borderColor, hoverBgColor, focusRingColor } =
    getStatusClasses(currentStatus);

  const isStatusEditable = currentStatus !== "cancelled";

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
              <p>{orderNumber || "-"}</p>
            </div>
            <div className="flex justify-between items-center">
              <p>Order Status</p>
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setIsStatusMenuOpen((prev) => !prev)}
                  disabled={
                    !isStatusEditable ||
                    updateOrderMutation.isPending ||
                    cancelOrderMutation.isPending
                  }
                  className={`flex items-center gap-2 ${bgColor} text-white border ${borderColor} rounded-md px-4 py-1 capitalize shadow-sm ${hoverBgColor} focus:outline-none focus:ring-2 ${focusRingColor} disabled:opacity-50 disabled:cursor-not-allowed`}
                >
                  {updateOrderMutation.isPending ||
                  cancelOrderMutation.isPending
                    ? "Updating..."
                    : currentStatus}
                  {isStatusEditable && (
                    <span className="text-xs">
                      {isStatusMenuOpen ? "▲" : "▼"}
                    </span>
                  )}
                </button>
                {isStatusMenuOpen && isStatusEditable && (
                  <div className="absolute right-0 mt-2 w-32 bg-white border border-gray-200 rounded-md shadow-lg z-10">
                    {statusOptions.map((option) => (
                      <button
                        key={option}
                        type="button"
                        disabled={option === currentStatus}
                        onClick={() => handleStatusClick(option)}
                        className={`w-full text-left px-4 py-2 text-sm capitalize hover:bg-gray-100 ${
                          option === currentStatus
                            ? "text-gray-400 cursor-not-allowed"
                            : "text-gray-700"
                        }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <div className="flex justify-between">
              <p>Table No.</p>
              <p>{tableNumber}</p>
            </div>
            <div className="flex justify-between">
              <p>Date</p>
              <p>{formatDateTime(createdAt)}</p>
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

      {showCancelConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
          <div className="bg-white p-6 rounded-lg shadow-xl text-center mx-4">
            <h3 className="text-lg font-semibold mb-4">Confirm Cancellation</h3>
            <p className="mb-6">Are you sure you want to cancel this order?</p>
            <div className="flex justify-center gap-4">
              <button
                onClick={() => setShowCancelConfirm(false)}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-100"
              >
                No, Go Back
              </button>
              <button
                onClick={confirmCancelOrder}
                disabled={cancelOrderMutation.isPending}
                className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {cancelOrderMutation.isPending
                  ? "Cancelling..."
                  : "Yes, Cancel Order"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrderDetails;
