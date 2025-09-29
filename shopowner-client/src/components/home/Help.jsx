import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const Help = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HELP} />
    </div>
  );
};

export default Help;
