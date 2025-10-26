import AppRoutes from "./routes/AppRoutes.jsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ToastContainer } from "react-toastify";

const App = () => {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <AppRoutes />
      <ToastContainer />
    </QueryClientProvider>
  );
};

export default App;
