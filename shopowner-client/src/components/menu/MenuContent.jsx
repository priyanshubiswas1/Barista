import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CaretDown, DotsThree } from "phosphor-react";

import { useState, useRef, useEffect } from "react";
import { deleteItemById, updateItemAvailable } from "../../api/api";
import { toast } from "react-toastify";
import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";

const Toggle = ({ itemId, isAvailable }) => {
  const queryClient = useQueryClient();
  const updateAvailability = useMutation({
    mutationFn: updateItemAvailable,
    onSuccess: () => {
      toast("Availability has been updated");
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
    },
    onError: (error) => {
      setIsOn(!isOn);
      toast.error("Failed to update availability. Please try again.");
    },
  });

  const [isOn, setIsOn] = useState(isAvailable);
  const toggleSwitch = () => {
    const newValue = !isOn;
    setIsOn(newValue);
    updateAvailability.mutate({
      itemId: itemId,
      shopId: localStorage.getItem("shopId"),
      available: newValue,
    });
  };

  return (
    <div className="flex justify-center items-center">
      <div
        onClick={toggleSwitch}
        className={`relative w-7 h-4 rounded-full cursor-pointer transition-all duration-300 ease-in-out ${
          isOn ? "bg-gray-800" : "bg-gray-300"
        }`}
      >
        <div
          className={`absolute top-[2px] left-[2px] w-3 h-3 bg-white rounded-full transition-all duration-300 ease-in-out ${
            isOn ? "translate-x-[12px]" : "translate-x-0"
          }`}
        />
      </div>
    </div>
  );
};

const MenuContentHeader = ({ category, count }) => {
  return (
    <div className="flex items-center gap-3 pb-3">
      <CaretDown size={16} weight="fill" className="text-gray-700" />
      <h3 className="font-bold text-gray-900">{category}</h3>
      <div className="bg-emerald-100 border border-emerald-200 rounded-lg h-6 px-2.5 flex items-center justify-center text-xs font-bold text-emerald-700">
        {count}
      </div>
    </div>
  );
};

const ItemActions = ({ item }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);
  const buttonRef = useRef(null);
  const queryClient = useQueryClient();

  const { setOverlay } = useOverlay();

  // Delete mutation - hooks must be at top level
  const deleteItemMutation = useMutation({
    mutationFn: deleteItemById,
    onSuccess: () => {
      toast("Item deleted");
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
    },
    onError: (error) => {
      toast.error("Failed to delete item. Please try again.");
    },
  });

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  const handleEdit = () => {
    setIsOpen(false);
    setOverlay(Overlays.UPDATEITEM, item);
  };

  const handleDelete = () => {
    setIsOpen(false);
    deleteItemMutation.mutate({ itemId: item.itemId });
  };

  return (
    <div className="relative">
      <div
        ref={buttonRef}
        onClick={() => setIsOpen(!isOpen)}
        className="cursor-pointer"
      >
        <DotsThree size={16} weight="bold" className="text-gray-500" />
      </div>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute right-0 top-8 z-50 bg-white border border-gray-200 rounded-lg shadow-xl min-w-[140px] overflow-hidden"
        >
          <button
            onClick={handleEdit}
            className="w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 text-left transition-colors font-medium flex items-center gap-2"
          >
            <span>✏️</span> Edit
          </button>
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2.5 text-sm text-red-600 bg-red-50 hover:bg-red-100 text-left transition-colors font-medium flex items-center gap-2"
          >
            <span>🗑️</span> Delete
          </button>
        </div>
      )}
    </div>
  );
};

const MenuContentBody = ({ items }) => {
  return (
    <div className="w-full ">
      <table className="w-full text-sm border-collapse table-fixed">
        {/* Header */}
        <thead>
          <tr className="text-gray-600 bg-gray-50 border-b-2 border-gray-200">
            <th className="pl-10 py-3.5 font-semibold text-left w-[40%]">
              Item Name
            </th>
            <th className="py-3.5 font-semibold text-left w-[15%]">Price</th>
            <th className="py-3.5 font-semibold text-left w-[15%]">
              Type
            </th>
            <th className="py-3.5 font-semibold text-left w-[12%]">Available</th>
            <th className="py-3.5 font-semibold text-left w-[8%]"></th>
          </tr>
        </thead>

        {/* Body */}
        <tbody>
          {items.map((item, index) => (
            <tr
              key={index}
              className={`text-gray-900 ${
                index !== items.length - 1 ? "border-b border-gray-200" : ""
              }`}
              style={{ height: "48px" }}
            >
              {/* Item name */}
              <td className="pl-10 py-2">
                <div className="flex items-center gap-4">
                  {item.image ? (
                    <img
                      src={item.image}
                      alt={item.itemName}
                      className="w-10 h-10 rounded-lg object-cover border border-gray-200"
                    />
                  ) : (
                    <div className="w-10 h-10 bg-gradient-to-br from-gray-200 to-gray-300 rounded-lg flex items-center justify-center">
                      <span className="text-gray-500 text-xs">🍽️</span>
                    </div>
                  )}
                  <div className="font-medium text-gray-900">{item.itemName}</div>
                </div>
              </td>

              {/* Price */}
              <td className="font-semibold text-gray-900">₹{item.price}</td>

              {/* Veg / Non-veg */}
              <td>
                {item.veg ? (
                  <div className="bg-green-100 border border-green-200 text-green-700 text-xs font-bold rounded-lg px-3 py-1 inline-block">
                    🌿 Veg
                  </div>
                ) : (
                  <div className="bg-red-100 border border-red-200 text-red-700 text-xs font-bold rounded-lg px-3 py-1 inline-block">
                    🍗 Non-veg
                  </div>
                )}
              </td>

              {/* Availability */}
              <td>
                <div className="flex items-center">
                  <Toggle itemId={item.itemId} isAvailable={item.available} />
                </div>
              </td>

              {/* Actions */}
              <td>
                <ItemActions item={item} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

const MenuContent = ({ category, items }) => {
  return (
    <div className="w-full bg-white p-5 rounded-xl border-2 border-gray-200 shadow-sm hover:shadow-md transition-all duration-200">
      <MenuContentHeader category={category} count={items.length} />
      <MenuContentBody items={items} />
    </div>
  );
};

export default MenuContent;
