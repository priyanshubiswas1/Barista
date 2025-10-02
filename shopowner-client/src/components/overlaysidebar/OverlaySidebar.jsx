import { useOverlay } from "../../store/overlayStore";
import { Overlays } from "../../constants/constants";
import ManageCategory from "./ManageCategory.jsx";
import AddItems from "./AddItems.jsx";

const overlayHeaders = {
  [Overlays.ADDITEM]: "Add Items",
  [Overlays.MANAGECAT]: "Manage Categories",
  [Overlays.CREATEORDER]: "Create Orders",
};

const OverlayComponent = () => {
  const { overlay } = useOverlay(); // <-- Move hook here
  const header = overlayHeaders[overlay] || null;
  switch (overlay) {
    case Overlays.MANAGECAT:
      return <ManageCategory header={header} />;
    case Overlays.ADDITEM:
      return <AddItems header={header} />;
    default:
      return null;
  }
};

export const OverlaySidebar = () => {
  const { overlay } = useOverlay(); // <-- And here
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
