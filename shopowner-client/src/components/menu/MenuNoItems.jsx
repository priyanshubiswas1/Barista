import React from "react";
import Button from "../Button";
import { Plus } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";

const MenuNoItems = () => {
  const { setOverlay } = useOverlay();
  return (
    <div className="w-full h-[75%] bg-white p-4 flex flex-col items-center justify-center rounded-md border-1 border-gray-200">
      <div className="flex flex-col w-60 gap-5 items-center">
        <div className="flex flex-col w-full gap-1.5 items-center">
          <h2 className="font-semibold">No Items Added</h2>
          <p className="text-gray-400 text-center">
            Your menu is looking a little empty. Add your first item to get
            started
          </p>
        </div>

        <div>
          <Button
            onClick={() => setOverlay(Overlays.ADDITEM)}
            variant="subtle"
            icon={Plus}
          >
            Add Item
          </Button>
        </div>
      </div>
    </div>
  );
};

export default MenuNoItems;
