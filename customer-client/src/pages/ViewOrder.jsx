import { useAuthStore } from "../store/authStore";
import Navbar from "../components/Navbar";
import { useNavigate, useParams } from "react-router-dom";
import { IoAdd, IoEllipseSharp, IoRemove } from "react-icons/io5";
import { Plus } from "phosphor-react";
import Button from "../components/Button";

const OrderDetails = () => {
  const navigate = useNavigate();
  const { shopid, tableid } = useParams();
  const Items = () => {
    return (
      <div className="flex items-center justify-between py-1">
        <div className="flex flex-1 items-center gap-2">
          <h2 className="text-sm text-gray-600">Veg Salad</h2>
          <div className="text-green-800 border border-green-800 rounded-sm p-1">
            <IoEllipseSharp size={6} />
          </div>
        </div>
        <div className="flex flex-1 items-center justify-end gap-2">
          <div className="bg-rose-100 border border-rose-300 rounded-md flex items-center justify-between w-20 px-2 shadow-sm">
            <button
              type="button"
              // onClick={() => updateOrderItems(count - 1)}
              className="text-rose-700 text-xl leading-none"
            >
              <IoRemove />
            </button>
            <span className="text-rose-700 text-sm font-medium">1</span>
            <button
              type="button"
              // onClick={() => updateOrderItems(count + 1)}
              className="text-rose-700 text-xl leading-none"
            >
              <IoAdd />
            </button>
          </div>
          <h2 className="font-medium">₹120</h2>
        </div>
      </div>
    );
  };
  return (
    <div className="w-full px-4 pt-4">
      <div className="p-3 border border-gray-200 hover:border-gray-300 rounded-md bg-white">
        <h2 className="text-sm font-medium">Order Details</h2>
        <Items />
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

const BillingDetails = () => {
  return (
    <div className="w-full px-4 pt-2">
      <div className="p-3 border border-gray-200 hover:border-gray-300 rounded-md bg-white">
        <h2 className="text-sm font-medium">Billing Details</h2>
        <div className="flex flex-col py-2 gap-1">
          <div className="flex justify-between">
            <p className="text-gray-600 font-light text-sm">Subtotal</p>
            <p className=" text-gray-900 font-medium text-sm">₹120</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600 font-light text-sm">Discount</p>
            <p className=" text-gray-900 font-medium text-sm">₹20</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600 font-light text-sm">Service Tax</p>
            <p className=" text-gray-900 font-medium text-sm">₹10</p>
          </div>
        </div>
        <div className="-mx-3 h-px bg-gray-200" />
        <div className="flex justify-between pt-2 items-center">
          <h3>Total</h3>
          <h3>₹120</h3>
        </div>
      </div>
    </div>
  );
};

const BasicDetails = () => {
  return (
    <div className="w-full px-4 pt-2">
      <div className="p-3 border border-gray-200 hover:border-gray-300 rounded-md bg-white">
        <h2 className="text-sm font-medium">Basic Details</h2>
        <div className="flex flex-col py-2 gap-1">
          <div className="flex justify-between">
            <p className="text-gray-600 font-light text-sm">Table No.</p>
            <p className="text-gray-900 font-medium text-sm">2</p>
          </div>
          <div className="flex justify-between">
            <p className="text-gray-600 font-light text-sm">Date</p>
            <p className="text-gray-900 font-medium text-sm">12 June, 2025</p>
          </div>
        </div>
      </div>
    </div>
  );
};

const ViewOrder = () => {
  const navigate = useNavigate();
  const { shopid, tableid } = useParams();

  const { placeOrderItems } = useAuthStore();
  console.log(placeOrderItems);
  return (
    <div className="w-screen h-screen">
      <Navbar />
      <OrderDetails />
      <BillingDetails />
      <BasicDetails />

      <div className="fixed bottom-0 left-0 right-0 flex flex-col gap-2">
        <div className="h-px bg-gray-300 w-full" />
        <div className="px-4 py-1">
          <div className="bg-gray-200 border border-gray-300 w-full rounded-md flex items-center gap-3 p-1">
            <span className="text-gray-600 text-lg">+91 </span>
            <input
              className="text-gray-400 text-lg"
              type="text"
              placeholder="Phone Number"
            />
          </div>
        </div>
        <div className=" bg-emerald-700 text-white py-3">
          <div className="flex justify-between items-center px-5">
            {/* <span className=" text-sm">
            {orderItems.reduce((total, item) => total + item.quantity, 0)}{" "}
            {orderItems.reduce((total, item) => total + item.quantity, 0) !== 1
              ? "Items"
              : "Item"}{" "}
            added | <span className="font-medium">₹ {totalPrice}</span>
          </span> */}
            <button
              onClick={() => navigate(`/${shopid}/${tableid}`)}
              className="px-4 py-2 rounded-md font-semibold text-sm"
            >
              Place Order
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ViewOrder;
