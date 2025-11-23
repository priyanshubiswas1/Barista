import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";
import ManageCategory from "./ManageCategory.jsx";
import AddItems from "./AddItems.jsx";
import AddTable from "./AddTable.jsx";
import AddNewOrder from "./AddNewOrder.jsx";
import UpdateItem from "./UpdateItem.jsx";
import UpdateTable from "./UpdateTable.jsx";
import OrderDetails from "./OrderDetails.jsx";
import GenerateBill from "./GenerateBill.jsx";

const overlayHeaders = {
  [Overlays.ADDITEM]: "Add New Item",
  [Overlays.UPDATEITEM]: "Update Item",
  [Overlays.MANAGECAT]: "Manage Categories",
  [Overlays.CREATEORDER]: "Add New Order",
  [Overlays.CREATETABLE]: "Add New Table",
  [Overlays.UPDATETABLE]: "Update Table",
  [Overlays.ORDERDETAILS]: "Order Details",
  [Overlays.GENERATEBILL]: "Generate Bill",
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
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
          <AddTable header={header} />
        </div>
      );
    case Overlays.UPDATETABLE:
      return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
          <UpdateTable header={header} table={overlayData} />
        </div>
      );
    case Overlays.CREATEORDER:
      return <AddNewOrder header={header} />;
    case Overlays.ORDERDETAILS:
      return <OrderDetails header={header} />;
    case Overlays.GENERATEBILL:
      return <GenerateBill header={header} />;
    default:
      return null;
  }
};

export const OverlaySidebar = () => {
  const { overlay, setOverlay } = useOverlay();

  if (overlay === Overlays.CREATETABLE || overlay === Overlays.UPDATETABLE) {
    return <OverlayComponent />;
  }

  return (
    <>
      {/* Backdrop */}
      {overlay && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 transition-opacity duration-300"
          onClick={() => setOverlay("")}
        />
      )}
      
      {/* Sidebar */}
      <div
        className={`fixed top-0 right-0 h-screen bg-white shadow-2xl transition-all duration-300 ease-in-out z-50 
          ${overlay ? "translate-x-0 w-[35%] min-w-[500px]" : "translate-x-full w-0"}
        `}
      >
        <OverlayComponent />
      </div>
    </>
  );
};
