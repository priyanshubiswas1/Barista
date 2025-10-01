import React from "react";
import Button from "../Button";
import { Plus } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";

const MenuNoItems = () => {
  const { setOverlay } = useOverlay();
  return (
    <div className="w-full h-[75%] bg-white p-4 flex flex-col items-center justify-center">
      <h2 className="font-semibold">No Items Added</h2>
      <div className="flex flex-col items-center justify-center pt-1">
        <p>Your menu is looking a little empty. Add</p>
        <p>your first item to get started</p>
      </div>
      <div className="mt-5">
        <Button
          onClick={() => setOverlay(Overlays.ADDITEM)}
          variant="subtle"
          icon={Plus}
        >
          Add Item
        </Button>
      </div>
    </div>
  );
};

export default MenuNoItems;
