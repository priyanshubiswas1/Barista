import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const History = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.HISTORY} />
    </div>
  );
};

export default History;
