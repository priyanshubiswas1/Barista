import { LuEllipsisVertical } from "react-icons/lu";

const Navbar = () => {
  return (
    <div className="w-full ">
      <div className="py-4 px-5 flex justify-between">
        <h2 className="font-semibold">The Cafe App</h2>
        <div className="">
          <LuEllipsisVertical />
        </div>
      </div>
      <div className="w-full h-px bg-gray-200" />
    </div>
  );
};

export default Navbar;
