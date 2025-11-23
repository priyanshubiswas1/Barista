import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LuEllipsisVertical } from "react-icons/lu";

const Navbar = () => {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const shopId = localStorage.getItem("shopId");
  const tableId = localStorage.getItem("tableId");

  return (
    <div className="w-full ">
      <div className="py-4 px-5 flex justify-between relative">
        <h2 className="font-semibold">The Cafe App</h2>
        <div className="relative">
          <button
            aria-label="menu"
            onClick={() => setOpen((v) => !v)}
            className="p-1 rounded"
          >
            <LuEllipsisVertical />
          </button>

          {open && (
            <div className="absolute right-0 mt-2 w-40 bg-white border border-gray-200 rounded-md shadow-md z-50">
              <div className="p-2">
                <button
                  onClick={() => {
                    setOpen(false);
                    navigate(`/${shopId}/${tableId}/orders`);
                  }}
                  className="w-full text-left px-3 py-2 rounded text-sm hover:bg-gray-100"
                >
                  View Orders
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
      <div className="w-full h-px bg-gray-200" />
    </div>
  );
};

export default Navbar;
