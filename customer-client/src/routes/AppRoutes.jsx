import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Menu from "../pages/Menu";
import ViewOrder from "../pages/ViewOrder";
import Orders from "../pages/Orders";
import Session from "../pages/Session";
import OrderSummary from "../pages/OrderSummary";

function Protected({ children }) {
  const token = localStorage.getItem("session-token");
  if (!token) return <Navigate to="/session" />;

  // const { data, isLoading } = useQuery({
  //   queryKey: [token],
  //   queryFn: () => verifyUser(),
  //   onSuccess: (data) => {
  //     useAuthStore.setState({ isAuth: data });
  //     console.log("verified");
  //   },
  // });

  // if (isLoading) return <p>Verifying...</p>;
  // if (!data) {
  //   return <Navigate to="/register" />;
  // }

  return <>{children}</>;
}
function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/:shopid/:tableid" element={<Menu />} />
        <Route path="/:shopid/:tableid/view-order" element={<ViewOrder />} />
        {/* <Route path="/:shopid/:tableid/session" element={<Session />} /> */}
        <Route path="/:shopid/:tableid/orders" element={<Orders />} />
        <Route
          path="/order/:orderid"
          element={
            <Protected>
              <OrderSummary />
            </Protected>
          }
        />
        <Route path="/:shopid/:tableid/*" element={<Navigate to={"/menu"} />} />
      </Routes>
    </BrowserRouter>
  );
}

export default AppRoutes;
