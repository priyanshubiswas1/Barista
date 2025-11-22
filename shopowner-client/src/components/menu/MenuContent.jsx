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
      console.log("Availability update failed: ", error);
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
    <div className="flex items-center gap-3">
      <CaretDown size={14} weight="fill" className="text-gray-800" />
      <h3>{category}</h3>
      <div className="bg-gray-50 border border-gray-100 rounded-sm h-6 w-6 flex items-center justify-center text-xs font-medium">
        {count}
      </div>
      <DotsThree size={16} weight="bold" className="text-gray-500" />
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
      console.log("Delete item failed: ", error);
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
          className="absolute right-0 top-6 z-50 bg-white border border-gray-200 rounded-md shadow-lg min-w-[120px] overflow-hidden"
        >
          <button
            onClick={handleEdit}
            className="w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 text-left transition-colors"
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2 text-sm text-red-700 bg-red-100 hover:bg-red-200 text-left transition-colors"
          >
            Delete
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
          <tr className="text-gray-500 border-b border-gray-200">
            <th className="pl-10 py-3 font-medium text-left w-[40%]">
              Item Name
            </th>
            <th className="py-3 font-medium text-left w-[15%]">Price</th>
            <th className="py-3 font-medium text-left w-[15%]">
              Veg / Non-veg
            </th>
            <th className="py-3 font-medium text-left w-[12%]">Availability</th>
            <th className="py-3 font-medium text-left w-[8%]"></th>
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
              <td className="pl-10">
                <div className="flex items-center gap-4">
                  {item.image ? (
                    <img
                      src={item.image}
                      alt={item.itemName}
                      className="w-9 h-9 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-9 h-9 bg-gray-300 rounded-full" />
                  )}
                  <div>{item.itemName}</div>
                </div>
              </td>

              {/* Price */}
              <td>₹ {item.price}</td>

              {/* Veg / Non-veg */}
              <td>
                {item.veg ? (
                  <div className="bg-green-100 text-green-700 text-xs font-semibold rounded-sm px-2 py-0.5 inline-block">
                    Veg
                  </div>
                ) : (
                  <div className="bg-red-100 text-red-700 text-xs font-semibold rounded-sm px-2 py-0.5 inline-block">
                    Non-veg
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
    <div className="w-full bg-white p-4 rounded-lg border-1 border-gray-200">
      <MenuContentHeader category={category} count={items.length} />
      <MenuContentBody items={items} />
    </div>
  );
};

export default MenuContent;
