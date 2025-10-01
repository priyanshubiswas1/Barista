import React from "react";
import Button from "../Button";
import { Plus } from "phosphor-react";

const MenuNoItems = () => {
  return (
    <div className="w-full h-[75%] bg-white p-4">
      <h2 className="font-semibold">No Items Added</h2>
      <p>Your menu is looking a little empty. Add</p>
      <p>your first item to get started</p>
      <Button variant="subtle" icon={Plus}>
        Add Item
      </Button>
    </div>
  );
};

export default MenuNoItems;
