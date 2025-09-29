import React from "react";
import { IoFilter } from "react-icons/io5";

const Filter = () => {
  return (
    <div className="flex items-center gap-3 bg-white border border-gray-200 rounded-lg px-3 py-2">
      <IoFilter className="text-xl flex-shrink-0" />
      <h2 className="">Filters</h2>
    </div>
  );
};

export default Filter;
