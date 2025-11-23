import { useEffect, useState } from "react";
import {
  IoAdd,
  IoRemove,
  IoEllipseSharp,
  IoTriangleSharp,
} from "react-icons/io5";
import { IoMdArrowDropright, IoMdArrowDropdown } from "react-icons/io";
import About from "../components/About";
import Navbar from "../components/Navbar";
import { fetchItemsByShopId } from "../api/api";
import { useAuthStore } from "../store/authStore";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";

const VegNonVegSelect = ({ value, onChange }) => (
  <div className="flex gap-2 pt-3 pb-2">
    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="all"
        checked={value === "all"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        {" "}
        All{" "}
      </div>
    </label>
    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="veg"
        checked={value === "veg"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        <div className="text-green-800 border border-green-800 rounded-sm p-1">
          {" "}
          <IoEllipseSharp size={6} />{" "}
        </div>{" "}
        Veg
      </div>
    </label>
    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="non-veg"
        checked={value === "non-veg"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        <div className="text-red-800 border border-red-800 rounded-sm p-1">
          {" "}
          <IoTriangleSharp size={6} />{" "}
        </div>{" "}
        Non Veg
      </div>
    </label>
  </div>
);

const Category = ({
  category,
  items,
  isExpanded,
  toggleCategory,
  orderItems,
  setOrderItems,
}) => (
  <div className="">
    <div className="bg-gray-200 border border-gray-300 rounded-md">
      <div
        className="flex items-center justify-between p-2 cursor-pointer select-none"
        onClick={toggleCategory}
      >
        <div className="flex items-center gap-4">
          <h2 className="font-medium">{category}</h2>
          <div className="bg-white border border-gray-300 rounded-sm h-5 w-5 flex justify-center items-center text-gray-600 text-sm">
            {items.length}
          </div>
        </div>
        {isExpanded ? (
          <IoMdArrowDropdown size={25} />
        ) : (
          <IoMdArrowDropright size={25} />
        )}
      </div>
      <div className="px-4 pb-3 flex flex-col gap-2">
        {isExpanded &&
          items.map((item) => (
            <Item
              key={item.itemId}
              item={item}
              orderItems={orderItems}
              setOrderItems={setOrderItems}
            />
          ))}
      </div>
    </div>
  </div>
);

const Item = ({ item, orderItems, setOrderItems }) => {
  const currentItem = orderItems.find((i) => i.itemId === item.itemId);
  const count = currentItem ? currentItem.quantity : 0;

  const updateOrderItems = (newQuantity) => {
    const exists = orderItems.find((i) => i.itemId === item.itemId);
    if (!exists && newQuantity === 0) return;

    const newOrderItems = exists
      ? orderItems.map((i) =>
          i.itemId === item.itemId
            ? { ...i, quantity: newQuantity, priceAtOrderTime: item.price }
            : i
        )
      : [
          ...orderItems,
          {
            itemId: item.itemId,
            itemName: item.itemName,
            quantity: newQuantity,
            priceAtOrderTime: item.price,
          },
        ];
    setOrderItems(newOrderItems.filter((i) => i.quantity > 0));
  };

  return (
    <div className="flex justify-between bg-white border border-gray-300 rounded-xl p-3">
      <div className="flex flex-col justify-between">
        <div className="flex flex-col items-start gap-1">
          <div
            className={`${
              item.veg
                ? "text-green-800 border-green-800"
                : "text-red-800 border-red-800"
            } border rounded-sm w-4 h-4 flex items-center justify-center`}
          >
            {item.veg ? (
              <IoEllipseSharp size={8} />
            ) : (
              <IoTriangleSharp size={8} />
            )}
          </div>
          <h2 className="font-medium text-gray-800">{item.itemName}</h2>
        </div>
        <h2 className="text-gray-600 text-sm">₹ {item.price}</h2>
      </div>
      <div
        className="relative w-20 h-20 rounded-md flex items-center justify-center"
        style={{
          backgroundImage: `url(${item.image})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
        }}
      >
        {count === 0 ? (
          <button
            type="button"
            onClick={() => {
              updateOrderItems(count + 1);
            }}
            className="absolute bottom-[-8px] bg-rose-700 text-white rounded-md w-8 h-8 flex items-center justify-center shadow-md"
          >
            <IoAdd size={20} />
          </button>
        ) : (
          <div className="absolute bottom-[-8px] bg-rose-100 border border-rose-300 rounded-md flex justify-between items-center px-1 py-1 gap-2 shadow-sm w-full">
            <button
              type="button"
              onClick={() => updateOrderItems(count - 1)}
              className="text-rose-700 text-xl leading-none"
            >
              <IoRemove />
            </button>
            <span className="text-rose-700 text-sm font-medium">{count}</span>
            <button
              type="button"
              onClick={() => updateOrderItems(count + 1)}
              className="text-rose-700 text-xl leading-none"
            >
              <IoAdd />
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

const Content = ({
  items,
  expandedCategories,
  toggleCategory,
  orderItems,
  setOrderItems,
}) => {
  const groupedItems = items.reduce((acc, item) => {
    acc[item.category] = acc[item.category] || [];
    acc[item.category].push(item);
    return acc;
  }, {});

  return (
    <div className="p-4 flex flex-col gap-4">
      {Object.entries(groupedItems).map(([category, categoryItems]) => (
        <Category
          key={category}
          category={category}
          items={categoryItems}
          isExpanded={expandedCategories[category]}
          toggleCategory={() => toggleCategory(category)}
          orderItems={orderItems}
          setOrderItems={setOrderItems}
        />
      ))}
    </div>
  );
};

const Menu = () => {
  const navigate = useNavigate();
  const { shopid, tableid } = useParams();
  localStorage.setItem("shopId", shopid);
  localStorage.setItem("tableId", tableid);

  const [expandedCategories, setExpandedCategories] = useState({});
  const [foodType, setFoodType] = useState("all");

  const { orderItems, setOrderItems } = useAuthStore();

  const { data, isLoading, error } = useQuery({
    queryKey: ["menuItems"],
    queryFn: () => fetchItemsByShopId(shopid),
  });

  useEffect(() => {
    if (data?.items) {
      const categories = data.items.reduce((acc, item) => {
        acc[item.category] = true;
        return acc;
      }, {});
      setExpandedCategories(categories);
    }
  }, [data]);

  if (isLoading) return <p>Loading menu...</p>;
  if (error) return <p>Error: {error.message}</p>;

  let items = data.items;

  if (foodType === "veg") {
    items = items.filter((item) => item.veg === true);
  } else if (foodType === "non-veg") {
    items = items.filter((item) => item.veg === false);
  }

  const toggleCategory = (category) => {
    setExpandedCategories((prev) => ({ ...prev, [category]: !prev[category] }));
  };

  const totalPrice = orderItems.reduce(
    (total, item) => total + item.quantity * item.priceAtOrderTime,
    0
  );

  return (
    <div className="w-screen h-screen">
      <Navbar />
      <About />
      <div className="px-5 py-2">
        <VegNonVegSelect
          value={foodType}
          onChange={(e) => setFoodType(e.target.value)}
        />
      </div>
      <div className="w-full h-px bg-gray-300" />
      <Content
        items={items}
        expandedCategories={expandedCategories}
        toggleCategory={toggleCategory}
        orderItems={orderItems}
        setOrderItems={setOrderItems}
      />
      {orderItems.length > 0 && totalPrice > 0 && (
        <div className="fixed bottom-0 left-0 right-0 bg-emerald-700 text-white py-3">
          <div className="flex justify-between items-center px-5">
            <span className=" text-sm">
              {orderItems.reduce((total, item) => total + item.quantity, 0)}{" "}
              {orderItems.reduce((total, item) => total + item.quantity, 0) !==
              1
                ? "Items"
                : "Item"}{" "}
              added | <span className="font-medium">₹ {totalPrice}</span>
            </span>
            <button
              onClick={() => navigate(`/${shopid}/${tableid}/view-order`)}
              className="px-4 py-2 rounded-md font-semibold text-sm"
            >
              View Order
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Menu;
