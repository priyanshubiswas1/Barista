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
import CircleLoading from "../components/CircleLoading";
import ARModelModal from "../components/ARModelModal";

const VegNonVegSelect = ({ value, onChange }) => (
  <div className="flex gap-2 pt-4 pb-4 justify-center">
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
            text-sm font-semibold rounded-full border-2 px-4 py-2 flex items-center gap-2 transition-all duration-200
            ${
              value === type.value
                ? "bg-emerald-600 border-emerald-600 text-white shadow-md scale-105"
                : "bg-white border-gray-300 text-gray-600 hover:border-emerald-500 hover:text-emerald-600 hover:shadow-sm"
            }
          `}
        >
          {type.icon && (
            <div
              className={`
                border rounded-sm p-0.5 flex items-center justify-center
                ${
                  value === type.value
                    ? "text-white border-white"
                    : type.color === "green"
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
  onOpenAR,
}) => (
  <div className="bg-white rounded-2xl border border-gray-200 shadow-md hover:shadow-lg transition-shadow duration-200 overflow-hidden">
    <div
      className="flex items-center justify-between p-4 cursor-pointer select-none bg-gradient-to-r from-gray-50 to-white hover:from-emerald-50 hover:to-white transition-all duration-200"
      onClick={toggleCategory}
    >
      <div className="flex items-center gap-3">
        <h2 className="font-bold text-gray-900 text-base">{category}</h2>
        <span className="bg-emerald-100 text-emerald-700 text-xs font-bold px-2.5 py-1 rounded-full">
          {items.length}
        </span>
      </div>
      <div className={`text-gray-500 transition-transform duration-200 ${isExpanded ? 'rotate-0' : '-rotate-90'}`}>
        <IoMdArrowDropdown size={24} />
      </div>
    </div>
    
    {isExpanded && (
      <div className="p-3 flex flex-col gap-3 border-t border-gray-100 bg-gray-50/30">
        {items.map((item) => (
          <Item
            key={item.itemId}
            item={item}
            orderItems={orderItems}
            setOrderItems={setOrderItems}
            onOpenAR={onOpenAR}
          />
        ))}
      </div>
    )}
  </div>
);

const Item = ({ item, orderItems, setOrderItems, onOpenAR }) => {
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
    <div className="flex justify-between bg-white border border-gray-200 rounded-xl p-3 hover:shadow-lg hover:border-emerald-200 transition-all duration-200 group">
      <div className="flex flex-col justify-between flex-1 pr-3 min-w-0">
        <div className="flex flex-col items-start gap-1.5">
          <div
            className={`
              border-2 rounded w-4 h-4 flex items-center justify-center flex-shrink-0
              ${
                item.veg
                  ? "text-green-600 border-green-600"
                  : "text-red-600 border-red-600"
              }
            `}
          >
            {item.veg ? (
              <IoEllipseSharp size={8} />
            ) : (
              <IoTriangleSharp size={8} />
            )}
          </div>
          <h3 className="font-semibold text-gray-900 leading-tight text-sm">{item.itemName}</h3>
          <p className="text-xs text-gray-500 line-clamp-2 leading-relaxed">{item.description}</p>
        </div>
        <div className="mt-2 font-bold text-gray-900 text-base">₹{item.price}</div>
      </div>
      
      <div className="relative flex-shrink-0">
        <div
          className="w-28 h-28 rounded-xl bg-gradient-to-br from-gray-100 to-gray-200 bg-cover bg-center shadow-sm group-hover:shadow-md transition-shadow duration-200"
          style={{
            backgroundImage: `url(${item.image})`,
          }}
        />
        {item.scan3dUrl && (
          <button
            type="button"
            onClick={() => onOpenAR?.(item)}
            className="absolute top-2 right-2 bg-black/70 text-white text-[11px] font-bold px-2 py-1 rounded-lg shadow hover:bg-black/80 active:scale-95 transition"
          >
            AR
          </button>
        )}
        <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-24">
          {count === 0 ? (
            <button
              type="button"
              onClick={() => updateOrderItems(count + 1)}
              className="w-full bg-white text-emerald-700 font-bold text-sm border-2 border-emerald-600 rounded-lg py-2 shadow-md hover:bg-emerald-600 hover:text-white transition-all duration-200 uppercase tracking-wide"
            >
              Add
            </button>
          ) : (
            <div className="flex items-center justify-between bg-emerald-600 text-white rounded-lg px-3 py-2 shadow-lg w-full">
              <button
                type="button"
                onClick={() => updateOrderItems(count - 1)}
                className="text-white hover:scale-110 transition-transform duration-150 active:scale-95"
              >
                <IoRemove size={16} />
              </button>
              <span className="font-bold text-base min-w-[20px] text-center">{count}</span>
              <button
                type="button"
                onClick={() => updateOrderItems(count + 1)}
                className="text-white hover:scale-110 transition-transform duration-150 active:scale-95"
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
  const [arItem, setArItem] = useState(null);

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
    <CircleLoading fullScreen text="Loading menu..." />
  );
  
  if (error) return (
    <div className="flex items-center justify-center h-screen bg-gray-50 text-red-600">
      <p>Error loading menu: {error.message}</p>
    </div>
  );

  return (
    <div className="w-full min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 pb-28">
      <div className="sticky top-0 z-20 bg-white/95 backdrop-blur-md shadow-sm">
        <Navbar />
      </div>
      
      <div className="max-w-3xl mx-auto">
        <About />
        
        <div className="px-4 sticky top-[60px] z-10 bg-gradient-to-b from-gray-50 to-transparent backdrop-blur-sm pt-2 pb-3">
          <VegNonVegSelect
            value={foodType}
            onChange={(e) => setFoodType(e.target.value)}
          />
          <div className="w-full h-px bg-gradient-to-r from-transparent via-gray-300 to-transparent mt-4" />
        </div>

        <div className="px-4 flex flex-col gap-5 mt-2">
          {Object.entries(groupedItems).map(([category, categoryItems]) => (
            <Category
              key={category}
              category={category}
              items={categoryItems}
              isExpanded={isExpanded(category)}
              toggleCategory={() => toggleCategory(category)}
              orderItems={orderItems}
              setOrderItems={setOrderItems}
              onOpenAR={(item) => setArItem(item)}
            />
          ))}
          
          {items.length === 0 && (
            <div className="text-center py-16 text-gray-500 bg-white rounded-2xl shadow-sm">
              <p className="text-lg font-medium">No items found</p>
              <p className="text-sm mt-1">Try selecting a different filter</p>
            </div>
          )}
        </div>
      </div>

      {orderItems.length > 0 && totalPrice > 0 && (
        <div className="fixed bottom-4 left-4 right-4 z-30 max-w-3xl mx-auto px-4 sm:px-0">
          <div className="bg-gradient-to-r from-emerald-700 to-emerald-600 text-white rounded-2xl shadow-2xl p-4 flex justify-between items-center transform transition-all hover:shadow-[0_20px_40px_-10px_rgba(5,150,105,0.5)] hover:scale-[1.02] border border-emerald-500/20">
            <div className="flex flex-col">
              <span className="font-bold text-lg">
                {totalItemsCount} {totalItemsCount !== 1 ? "Items" : "Item"}
              </span>
              <span className="text-emerald-100 text-sm">
                Total: <span className="font-bold text-white">₹{totalPrice}</span>
              </span>
            </div>
            <button
              onClick={() => navigate(`/${shopid}/${tableid}/view-order`)}
              className="bg-white text-emerald-700 px-6 py-3 rounded-xl font-bold text-sm shadow-lg hover:bg-emerald-50 transition-all hover:shadow-xl active:scale-95"
            >
              View Order
            </button>
          </div>
        </div>
      )}

      <ARModelModal
        open={!!arItem}
        title={arItem?.itemName}
        src={arItem ? `/api/item/${arItem.itemId}/scan3d` : ""}
        onClose={() => setArItem(null)}
      />
    </div>
  );
};

export default Menu;
