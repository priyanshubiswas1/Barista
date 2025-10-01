import { IoClose } from "react-icons/io5";
import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";

export const OverlaySidebar = () => {
  const { overlay, setOverlay } = useOverlay();

  // Dynamic content based on overlay type
  let content;
  switch (overlay) {
    case Overlays.ADDITEM:
      content = <h2 className="text-gray-800 text-lg">Add Items</h2>;
      break;
    case Overlays.MANAGECAT:
      content = <h2 className="text-gray-800 text-lg">Manage Categories</h2>;
      break;
    case Overlays.CREATEORDER:
      content = <h2 className="text-gray-800 text-lg">Create Orders</h2>;
      break;
    default:
      content = null;
  }

  return (
    <div
      className={`fixed top-0 right-0 h-screen bg-white shadow-lg transition-transform duration-300 ease-in-out z-50 ${
        overlay ? "translate-x-0 w-64" : "translate-x-full w-0"
      }`}
    >
      {/* Close button */}
      <div className="flex justify-end p-4">
        <button
          onClick={() => setOverlay("")}
          className="text-gray-600 hover:text-gray-900"
        >
          <IoClose size={24} />
        </button>
      </div>

      {/* Dynamic content */}
      <div className="p-4">{content}</div>
    </div>
  );
};
