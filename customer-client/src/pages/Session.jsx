import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { getSession } from "../api/api";

const Session = () => {
  const { shopid, tableid } = useParams();
  console.log("Params:", shopid, tableid);

  if (!shopid || !tableid) {
    return <div>Loading...</div>;
  }

  const { data, isLoading, isFetched, error } = useQuery({
    queryKey: [shopid, tableid],
    queryFn: () => getSession({ shopid, tableid }),
  });
  if (isFetched) {
    console.log(data);
  }

  return <div>SessionPage</div>;
};

export default Session;
