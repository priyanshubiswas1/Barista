import React from "react";
import Button from "../Button";
import { Plus } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";

const MenuNoItems = () => {
  const { setOverlay } = useOverlay();
  return (
    <div className="w-full h-[75%] bg-white p-8 flex flex-col items-center justify-center rounded-xl border-2 border-gray-200 shadow-sm">
      <div className="flex flex-col w-80 gap-6 items-center">
        <div className="w-24 h-24 bg-gradient-to-br from-indigo-100 to-indigo-200 rounded-full flex items-center justify-center mb-2">
          <span className="text-5xl">☕</span>
        </div>
        <div className="flex flex-col w-full gap-2 items-center">
          <h2 className="font-bold text-xl text-gray-900">No Items in Menu</h2>
          <p className="text-gray-500 text-center text-sm leading-relaxed">
            Your menu is empty. Start building your menu by adding your first
            delicious item. Your customers are waiting!
          </p>
        </div>

        <Button
          onClick={() => setOverlay(Overlays.ADDITEM)}
          variant="primary"
          icon={Plus}
        >
          Add Your First Item
        </Button>
      </div>
    </div>
  );
};

export default MenuNoItems;
