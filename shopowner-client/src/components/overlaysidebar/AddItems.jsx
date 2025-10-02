import { IoClose } from "react-icons/io5";
import { LuCirclePlus } from "react-icons/lu";
import { AiOutlineDelete } from "react-icons/ai";
import { PiPencilSimpleLight } from "react-icons/pi";
import { useOverlay } from "../../store/overlayStore";

const CategoryItem = () => {
  return (
    <div className="w-full">
      <div className="flex justify-between items-center px-6 pb-3">
        <div className="flex text-xl text-gray-900 gap-3">
          <div className="">Shakes</div>
          <div className="bg-gray-100 border border-gray-200 w-7 h-7 rounded flex items-center justify-center">
            <h2>5</h2>
          </div>
        </div>
        <div className="flex text-gray-500 text-xl gap-3">
          <PiPencilSimpleLight />
          <AiOutlineDelete />
        </div>
      </div>
      <div className="w-full h-px bg-gray-200" />
    </div>
  );
};

const AddItems = ({ header }) => {
  const { setOverlay } = useOverlay();
  return (
    <div className="flex flex-col w-full h-full">
      <div className="px-6 pt-6 pb-4 flex justify-between items-start">
        {header && <h2 className="font-semibold text-2xl">{header}</h2>}
        <button
          onClick={() => setOverlay("")}
          className="text-gray-600 hover:text-gray-900"
          aria-label="Close overlay"
        >
          <IoClose size={24} />
        </button>
      </div>
      <div className="w-full h-px bg-gray-200" />
      <div className="flex flex-col w-full h-full">
        <div className="pt-8">
          <CategoryItem />
        </div>
        <div className="pt-3 px-6">
          <button className="flex items-center gap-3">
            <LuCirclePlus />
            <h2 className="text-lg font-medium">Add New Category</h2>
          </button>
        </div>
      </div>
    </div>
  );
};

export default AddItems;
