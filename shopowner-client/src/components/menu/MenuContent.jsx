import { CaretDown, DotsThree } from "phosphor-react";

import { useState } from "react";

const Toggle = () => {
  const [isOn, setIsOn] = useState(false);
  const toggleSwitch = () => setIsOn(!isOn);

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
                  <div className="w-9 h-9 bg-gray-300 rounded-full"></div>
                  <div>{item.itemName}</div>
                </div>
              </td>

              {/* Price */}
              <td>₹ {item.price}</td>

              {/* Veg / Non-veg */}
              <td>
                {item.isVeg ? (
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
                  <Toggle />
                </div>
              </td>

              {/* Actions */}
              <td>
                <DotsThree size={16} weight="bold" className="text-gray-500" />
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
