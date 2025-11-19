import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";
import ManageCategory from "./ManageCategory.jsx";
import AddItems from "./AddItems.jsx";
import AddTable from "./AddTable.jsx";
import AddNewOrder from "./AddNewOrder.jsx";
import UpdateItem from "./UpdateItem.jsx";

const overlayHeaders = {
  [Overlays.ADDITEM]: "Add New Item",
  [Overlays.UPDATEITEM]: "Update Item",
  [Overlays.MANAGECAT]: "Manage Categories",
  [Overlays.CREATEORDER]: "Add New Order",
  [Overlays.CREATETABLE]: "Add New Table",
};

const OverlayComponent = () => {
  const { overlay, overlayData } = useOverlay();
  if (!overlay) return null;

  const header = overlayHeaders[overlay] || null;
  switch (overlay) {
    case Overlays.MANAGECAT:
      return <ManageCategory header={header} />;
    case Overlays.ADDITEM:
      return <AddItems header={header} />;
    case Overlays.UPDATEITEM:
      return <UpdateItem header={header} item={overlayData} />;
    case Overlays.CREATETABLE:
      return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
          <AddTable header={header} />
        </div>
      );
    case Overlays.CREATEORDER:
      return <AddNewOrder header={header} />;
    default:
      return null;
  }
};

export const OverlaySidebar = () => {
  const { overlay } = useOverlay();

  if (overlay === Overlays.CREATETABLE) {
    return <OverlayComponent />;
  }

  return (
    <div
      className={`fixed top-0 right-0 h-screen bg-white shadow-lg transition-transform duration-300 ease-in-out z-50 
        ${overlay ? "translate-x-0 w-[30%]" : "translate-x-full w-0"}
      `}
    >
      <OverlayComponent />
    </div>
  );
};
