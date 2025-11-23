import { useState, useMemo } from "react";
import {
  IoAdd,
  IoRemove,
  IoEllipseSharp,
  IoTriangleSharp,
} from "react-icons/io5";
import { IoMdArrowDropright, IoMdArrowDropdown } from "react-icons/io";
import About from "../components/About";
import Navbar from "../components/Navbar";
import { fetchItemsByShopId, getSession } from "../api/api";
import { useAuthStore } from "../store/authStore";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";

const VegNonVegSelect = ({ value, onChange }) => (
  <div className="flex gap-3 pt-4 pb-4 justify-center">
    {[
      { label: "All", value: "all", icon: null },
      { label: "Veg", value: "veg", icon: <IoEllipseSharp size={8} />, color: "green" },
      { label: "Non Veg", value: "non-veg", icon: <IoTriangleSharp size={8} />, color: "red" },
    ].map((type) => (
      <label key={type.value} className="cursor-pointer group">
        <input
          type="radio"
          name="foodType"
          value={type.value}
          checked={value === type.value}
          onChange={onChange}
          className="hidden"
        />
        <div
          className={`
            text-sm font-medium rounded-full border px-4 py-2 flex items-center gap-2 transition-all duration-200
            ${
              value === type.value
                ? "bg-emerald-50 border-emerald-600 text-emerald-800 shadow-sm"
                : "bg-white border-gray-200 text-gray-600 hover:bg-gray-50"
            }
          `}
        >
          {type.icon && (
            <div
              className={`
                border rounded-sm p-0.5 flex items-center justify-center
                ${
                  type.color === "green"
                    ? "text-green-700 border-green-700"
                    : "text-red-700 border-red-700"
                }
              `}
            >
              {type.icon}
            </div>
          )}
          {type.label}
        </div>
      </label>
    ))}
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
  <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
    <div
      className="flex items-center justify-between p-4 cursor-pointer select-none bg-gray-50/50 hover:bg-gray-50 transition-colors"
      onClick={toggleCategory}
    >
      <div className="flex items-center gap-3">
        <h2 className="font-semibold text-gray-800 text-lg">{category}</h2>
        <span className="bg-gray-200 text-gray-700 text-xs font-bold px-2 py-0.5 rounded-full">
          {items.length}
        </span>
      </div>
      <div className="text-gray-400">
        {isExpanded ? (
          <IoMdArrowDropdown size={24} />
        ) : (
          <IoMdArrowDropright size={24} />
        )}
      </div>
    </div>
    
    {isExpanded && (
      <div className="p-3 flex flex-col gap-3 border-t border-gray-100">
        {items.map((item) => (
          <Item
            key={item.itemId}
            item={item}
            orderItems={orderItems}
            setOrderItems={setOrderItems}
          />
        ))}
      </div>
    )}
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
    <div className="flex justify-between bg-white border border-gray-100 rounded-lg p-3 hover:shadow-md transition-shadow duration-200">
      <div className="flex flex-col justify-between flex-1 pr-4">
        <div className="flex flex-col items-start gap-1.5">
          <div
            className={`
              border rounded-sm w-4 h-4 flex items-center justify-center
              ${
                item.veg
                  ? "text-green-700 border-green-700"
                  : "text-red-700 border-red-700"
              }
            `}
          >
            {item.veg ? (
              <IoEllipseSharp size={8} />
            ) : (
              <IoTriangleSharp size={8} />
            )}
          </div>
          <h3 className="font-medium text-gray-900 leading-tight">{item.itemName}</h3>
          <p className="text-sm text-gray-500 line-clamp-2">{item.description}</p>
        </div>
        <div className="mt-2 font-semibold text-gray-900">₹ {item.price}</div>
      </div>
      
      <div className="relative flex-shrink-0">
        <div
          className="w-24 h-24 rounded-lg bg-gray-100 bg-cover bg-center"
          style={{
            backgroundImage: `url(${item.image})`,
          }}
        />
        <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-20">
          {count === 0 ? (
            <button
              type="button"
              onClick={() => updateOrderItems(count + 1)}
              className="w-full bg-white text-emerald-700 font-bold text-sm border border-gray-200 rounded-md py-1.5 shadow-sm hover:bg-emerald-50 transition-colors uppercase tracking-wide"
            >
              Add
            </button>
          ) : (
            <div className="flex items-center justify-between bg-emerald-700 text-white rounded-md px-2 py-1.5 shadow-md w-full">
              <button
                type="button"
                onClick={() => updateOrderItems(count - 1)}
                className="text-white hover:text-emerald-100 transition-colors"
              >
                <IoRemove size={16} />
              </button>
              <span className="font-bold text-sm">{count}</span>
              <button
                type="button"
                onClick={() => updateOrderItems(count + 1)}
                className="text-white hover:text-emerald-100 transition-colors"
              >
                <IoAdd size={16} />
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const Menu = () => {
  const navigate = useNavigate();
  const { shopid, tableid } = useParams();
  
  // Ideally move this to a useEffect or a more central place, but keeping it here for now as per existing logic
  if (shopid) localStorage.setItem("shopId", shopid);
  if (tableid) localStorage.setItem("tableId", tableid);

  const [expandedCategories, setExpandedCategories] = useState({});
  const [foodType, setFoodType] = useState("all");

  const { orderItems, setOrderItems } = useAuthStore();

  const { data, isLoading, error } = useQuery({
    queryKey: ["menuItems", shopid],
    queryFn: () => fetchItemsByShopId(shopid),
    staleTime: 1000 * 60 * 5, // 5 minutes
  });

  useQuery({
    queryKey: ["session", shopid, tableid],
    queryFn: () => getSession({ shopid, tableid }),
    enabled: !!(shopid && tableid),
    staleTime: 5 * 60 * 1000,
  });

  // Initialize expanded categories when data loads
  // We use useMemo to avoid useEffect, but we need to handle the initial state.
  // A simple way is to default all to true or handle it in the render.
  // Here we'll just check if the category is in the state, if not, we consider it expanded by default (or collapsed).
  // Let's default to expanded for better UX.
  
  const items = useMemo(() => {
    if (!data?.items) return [];
    let filtered = data.items;
    if (foodType === "veg") {
      filtered = filtered.filter((item) => item.veg === true);
    } else if (foodType === "non-veg") {
      filtered = filtered.filter((item) => item.veg === false);
    }
    return filtered;
  }, [data, foodType]);

  const groupedItems = useMemo(() => {
    return items.reduce((acc, item) => {
      acc[item.category] = acc[item.category] || [];
      acc[item.category].push(item);
      return acc;
    }, {});
  }, [items]);

  const toggleCategory = (category) => {
    setExpandedCategories((prev) => ({
      ...prev,
      [category]: prev[category] === undefined ? false : !prev[category], // Default to true if undefined? No, let's explicit.
    }));
  };
  
  // Helper to check if expanded. Default to true if not set.
  const isExpanded = (category) => {
    return expandedCategories[category] !== false;
  };

  const totalPrice = orderItems.reduce(
    (total, item) => total + item.quantity * item.priceAtOrderTime,
    0
  );

  const totalItemsCount = orderItems.reduce((total, item) => total + item.quantity, 0);

  if (isLoading) return (
    <div className="flex items-center justify-center h-screen bg-gray-50">
      <div className="animate-pulse flex flex-col items-center">
        <div className="h-12 w-12 bg-gray-200 rounded-full mb-4"></div>
        <div className="h-4 w-32 bg-gray-200 rounded"></div>
      </div>
    </div>
  );
  
  if (error) return (
    <div className="flex items-center justify-center h-screen bg-gray-50 text-red-600">
      <p>Error loading menu: {error.message}</p>
    </div>
  );

  return (
    <div className="w-full min-h-screen bg-gray-50 pb-24">
      <div className="sticky top-0 z-20 bg-white shadow-sm">
        <Navbar />
      </div>
      
      <div className="max-w-md mx-auto">
        <About />
        
        <div className="px-4 sticky top-[60px] z-10 bg-gray-50/95 backdrop-blur-sm pt-2">
          <VegNonVegSelect
            value={foodType}
            onChange={(e) => setFoodType(e.target.value)}
          />
          <div className="w-full h-px bg-gray-200 mb-4" />
        </div>

        <div className="px-4 flex flex-col gap-4">
          {Object.entries(groupedItems).map(([category, categoryItems]) => (
            <Category
              key={category}
              category={category}
              items={categoryItems}
              isExpanded={isExpanded(category)}
              toggleCategory={() => toggleCategory(category)}
              orderItems={orderItems}
              setOrderItems={setOrderItems}
            />
          ))}
          
          {items.length === 0 && (
            <div className="text-center py-10 text-gray-500">
              No items found for this selection.
            </div>
          )}
        </div>
      </div>

      {orderItems.length > 0 && totalPrice > 0 && (
        <div className="fixed bottom-4 left-4 right-4 z-30 max-w-md mx-auto">
          <div className="bg-emerald-800 text-white rounded-xl shadow-lg p-4 flex justify-between items-center transform transition-transform hover:scale-[1.02]">
            <div className="flex flex-col">
              <span className="font-bold text-lg">
                {totalItemsCount} {totalItemsCount !== 1 ? "Items" : "Item"}
              </span>
              <span className="text-emerald-200 text-sm">
                Total: <span className="font-bold text-white">₹ {totalPrice}</span>
              </span>
            </div>
            <button
              onClick={() => navigate(`/${shopid}/${tableid}/view-order`)}
              className="bg-white text-emerald-800 px-5 py-2.5 rounded-lg font-bold text-sm shadow-sm hover:bg-emerald-50 transition-colors"
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
