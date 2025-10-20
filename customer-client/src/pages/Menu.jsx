import { useState } from "react";
import About from "../components/About";
import Navbar from "../components/Navbar";
import Searchbar from "../components/Searchbar";
import { IoEllipseSharp, IoTriangleSharp } from "react-icons/io5";

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
  return (
    <div className="">
      <div className="">Shakes</div>
      <div className="">
        <Item />
      </div>
    </div>
  );
};

const Item = () => {
  return <div className="">Veg Salad</div>;
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
