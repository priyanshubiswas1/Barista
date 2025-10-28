import { useAuthStore } from "../store/authStore";
import Navbar from "../components/Navbar";
import { useNavigate, useParams } from "react-router-dom";
import { IoEllipseSharp } from "react-icons/io5";
import { IoIosArrowBack } from "react-icons/io";
import Button from "../components/Button";
import { useState } from "react";

const OrderDetails = ({ orderItems }) => {
  const Items = ({ item, orderItems }) => {
    const currentItem = orderItems.find((i) => i.itemId === item.itemId);

    return (
      <div className="flex items-center justify-between py-1">
        <div className="flex flex-1 items-center gap-2">
          <h2 className="text-sm text-gray-600">{item.itemName}</h2>
          <div className="text-green-800 border border-green-800 rounded-sm p-1">
            <IoEllipseSharp size={6} />
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="w-full px-4 pt-4">
      <div className="p-3 border border-gray-200 hover:border-gray-300 rounded-md bg-white">
        <h2 className="text-sm font-medium">Order Details</h2>
        {/* {orderItems.length > 0 ? (
          orderItems.map((item) => (
            <Items
              key={item.itemId}
              item={item}
              orderItems={orderItems}
              setOrderItems={setOrderItems}
            />
          ))
        ) : (
          <p className="text-gray-500 text-sm py-2">No items in order</p>
        )} */}
      </div>
    </div>
  );
};

const BackButton = () => {
  const navigate = useNavigate();

  return (
    <div className="w-full px-4 pt-4 flex justify-start items-center">
      <button
        onClick={() => navigate()}
        className="bg-white border border-gray-200 rounded-md flex justify-center items-center text-xl text-gray-300 p-2"
      >
        <IoIosArrowBack />
      </button>
    </div>
  );
};

const OrderSummary = () => {
  const { orderid } = useParams();

  const { orderItems, setOrderItems } = useState();

  // const sessionData = useQuery({
  //     queryKey: ["session", shopid, tableid],
  //     queryFn: () => getSession({ shopid, tableid }),
  // });

  return (
    <div className="w-screen h-screen">
      <Navbar />
      <BackButton />
      <OrderDetails orderItems={orderItems} />
    </div>
  );
};

export default OrderSummary;
