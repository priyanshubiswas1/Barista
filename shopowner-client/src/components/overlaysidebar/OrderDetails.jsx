import { IoClose } from "react-icons/io5";
import { useOverlay } from "../../store/overlayStore";

const formatCurrency = (amount) => {
  if (typeof amount !== "number") return "₹ 0";
  return `₹ ${amount.toFixed(2)}`;
};

const formatDateTime = (value) => {
  if (!value) return new Date().toLocaleString();
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return new Date().toLocaleString();
  }
  return date.toLocaleString();
};

const calculateLineTotal = (item) =>
  Number(item.priceAtOrderTime || 0) * Number(item.quantity || 0);

const OrderDetails = ({ header }) => {
  const { overlayData, setOverlay } = useOverlay();

  if (!overlayData) return null;

  const {
    orderId,
    tableNumber = "-",
    status = "new",
    createdAt,
    orderItems = [],
  } = overlayData;

  const total = orderItems.reduce((acc, item) => acc + calculateLineTotal(item), 0);

  return (
    <div className="flex flex-col w-full h-full">
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
              <p>{orderId || "-"}</p>
            </div>
            <div className="flex justify-between">
              <p>Order Status</p>
              <div className="bg-blue-400 text-white border rounded-md px-4 capitalize">
                {status}
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
            {orderItems.length === 0 && (
              <p className="text-gray-500">No items in this order.</p>
            )}
            {orderItems.map((item) => (
              <div
                key={`${item.orderItemId || item.itemId}-${item.itemName}`}
                className="flex justify-between items-center"
              >
                <div>
                  <p className="font-medium text-gray-800">{item.itemName}</p>
                  <p className="text-sm text-gray-500">
                    Qty: {item.quantity}
                  </p>
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
    </div>
  );
};

export default OrderDetails;
