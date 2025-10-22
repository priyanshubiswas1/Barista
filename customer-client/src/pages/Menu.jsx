import { useState } from "react";
import About from "../components/About";
import Navbar from "../components/Navbar";
import Searchbar from "../components/Searchbar";
import { IoEllipseSharp, IoTriangleSharp } from "react-icons/io5";
import { IoMdArrowDropright, IoMdArrowDropdown } from "react-icons/io";

const VegNonVegSelect = ({ value, onChange }) => (
  <div className="flex gap-2 pt-3 pb-2">
    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="veg"
        checked={value === "veg"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        <div className="text-green-800 border border-green-800 rounded-sm p-1">
          <IoEllipseSharp size={6} />
        </div>
        Veg
      </div>
    </label>

    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="non-veg"
        checked={value === "non-veg"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        <div className="text-red-800 border border-red-800 rounded-sm p-1">
          <IoTriangleSharp size={6} />
        </div>
        Non Veg
      </div>
    </label>

    <label className="cursor-pointer">
      <input
        type="radio"
        name="foodType"
        value="best-seller"
        checked={value === "best-seller"}
        onChange={onChange}
        className="hidden"
      />
      <div className="text-xs font-medium text-gray-800 rounded-4xl border border-gray-300 flex items-center justify-around px-3 py-1 gap-2">
        Best Seller
      </div>
    </label>
  </div>
);

const Search = () => {
  const [foodType, setFoodType] = useState("");
  return (
    <div className="px-5 py-2">
      <Searchbar />
      <VegNonVegSelect
        value={foodType}
        onChange={(e) => setFoodType(e.target.value)}
      />
    </div>
  );
};

const Category = () => {
  // const isExpanded = expandedCategories[category] ?? false;
  const isExpanded = true;
  return (
    <div className="p-4">
      <div className="bg-gray-200 border border-gray-300 rounded-md">
        <div className="flex items-center justify-between p-2">
          <div className="flex items-center gap-4">
            <h2 className="font-medium">Shakes</h2>
            <div className="bg-white border border-gray-300 rounded-sm h-5 w-5 flex justify-center items-center text-gray-600 text-sm">
              6
            </div>
          </div>
          <div className="">
            {isExpanded ? (
              <IoMdArrowDropdown size={25} className="text-gray-800" />
            ) : (
              <IoMdArrowDropright size={25} className="text-gray-800" />
            )}
          </div>
        </div>
        <div className="">
          <Item />
        </div>
      </div>
    </div>
  );
};

const Item = () => {
  return (
    <div className="px-3 pb-2">
      <div className="bg-white border border-gray-200 rounded-md">Veg</div>
    </div>
  );
};

const Content = () => {
  return (
    <div className="">
      <Category />
    </div>
  );
};

const Menu = () => {
  return (
    <div className="w-screen h-screen">
      <Navbar />
      <About />
      <Search />
      <div className="w-full h-px bg-gray-300" />
      <Content />
    </div>
  );
};

export default Menu;
