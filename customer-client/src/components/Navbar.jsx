import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LuEllipsisVertical } from "react-icons/lu";

const Navbar = () => {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const shopId = localStorage.getItem("shopId");
  const tableId = localStorage.getItem("tableId");

  return (
    <div className="w-full bg-white">
      <div className="py-4 px-5 flex justify-between items-center relative max-w-7xl mx-auto">
        <h2 className="font-semibold text-lg text-gray-800">Barista</h2>
        <div className="relative">
          <button
            aria-label="menu"
            onClick={() => setOpen((v) => !v)}
            className={`p-2 rounded-lg hover:bg-gray-100 transition-colors ${
              open ? "bg-gray-100" : ""
            }`}
          >
            <LuEllipsisVertical className="text-gray-700" size={20} />
          </button>

          {open && (
            <>
              <div
                className="fixed inset-0 z-40"
                onClick={() => setOpen(false)}
              />
              <div className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-lg shadow-xl z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
                <div className="py-1">
                  <button
                    onClick={() => {
                      setOpen(false);
                      navigate(`/${shopId}/${tableId}/orders`);
                    }}
                    className="w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:bg-emerald-50 hover:text-emerald-700 transition-colors font-medium"
                  >
                    View Orders
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
      <div className="w-full h-px bg-gradient-to-r from-transparent via-gray-200 to-transparent" />
    </div>
  );
};

export default Navbar;
