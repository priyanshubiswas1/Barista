import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const Community = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.COMMUNITY} />
    </div>
  );
};
export default Community;
