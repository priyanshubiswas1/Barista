import { IoMdArrowDropdown } from "react-icons/io";
import { FaEllipsis } from "react-icons/fa6";

import { useState } from "react";

const Toggle = () => {
  const [isOn, setIsOn] = useState(false);
  const toggleSwitch = () => setIsOn(!isOn);

  return (
    <div className="flex justify-center items-center">
      <div
        onClick={toggleSwitch}
        className={`w-12 h-7 rounded-full flex items-center cursor-pointer transition-all duration-300 ease-in-out ${
          isOn ? "bg-gray-800" : "bg-gray-300"
        }`}
      >
        <div
          className={`w-5 h-5 rounded-full bg-white transition-transform duration-300 ease-in-out ${
            isOn ? "transform translate-x-6" : "translate-x-1"
          }`}
        />
      </div>
    </div>
  );
};

const MenuContentHeader = ({ category, count }) => {
  return (
    <div className="flex items-center gap-3">
      <IoMdArrowDropdown className="text-3xl" />
      <h2 className="font-semibold text-lg">{category}</h2>
      <div className="bg-gray-200 border border-gray-50 rounded-sm h-6 w-6 flex items-center justify-center">
        {count}
      </div>
      <FaEllipsis className="text-gray-500" />
    </div>
  );
};

const MenuContentBody = ({ items }) => {
  return (
    <div className="flex flex-col w-full items-center">
      <div className="flex w-full text-gray-400 items-center px-10 pt-4 pb-2">
        <div className="flex flex-1 items-center gap-4">
          <span className="w-9 h-9 ml-5" />
          <h2>Item Name</h2>
        </div>
        <div className="flex flex-1 items-center justify-around gap-12">
          <h2>price</h2>
          <h2>veg/non-veg</h2>
          <h2>availability</h2>
          <h2></h2>
        </div>
      </div>
      <div className="flex w-full text-gray-100 border rounded"></div>
      {items.map((item) => {
        return (
          <div className="flex w-full pt-2 px-10">
            <div className="flex flex-1 items-center gap-4">
              <div className="w-9 h-9 ml-5 bg-gray-800 rounded-full"></div>
              <div>{item.itemName}</div>
            </div>
            <div className="flex flex-1 items-center justify-around gap-12">
              <div>{item.price}</div>
              <div className="w-20 flex justify-center">
                {!item ? (
                  <div className="bg-green-100 text-green-700 text-sm font-semibold rounded-sm flex items-center justify-center px-2">
                    veg
                  </div>
                ) : (
                  <div className="bg-red-100 text-red-700 text-sm font-semibold rounded-sm flex items-center justify-center px-2">
                    non-veg
                  </div>
                )}
              </div>
              <div>
                <Toggle />
              </div>
              <FaEllipsis className="text-gray-500" />
            </div>
          </div>
        );
      })}
    </div>
  );
};

const MenuContent = ({ category, items }) => {
  return (
    <div className="w-full bg-white p-4 rounded-md">
      <MenuContentHeader category={category} count={items.length} />
      <MenuContentBody items={items} />
    </div>
  );
};

export default MenuContent;
