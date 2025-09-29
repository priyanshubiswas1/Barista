import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const Settings = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.SETTINGS} />
    </div>
  );
};

export default Settings;
