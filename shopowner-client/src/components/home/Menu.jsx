import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const Menu = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.MENU} />
    </div>
  );
};

export default Menu;
