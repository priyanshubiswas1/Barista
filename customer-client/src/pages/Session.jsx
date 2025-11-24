import { useQuery, useMutation } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { getSession, sendOtp } from "../api/api";
import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import { useAuthStore } from "../store/authStore";
import CircleLoading from "../components/CircleLoading";

const SendOTP = () => {
  const [otpSent, setOtpSent] = useState(false);

  const { shopId, tableId } = useAuthStore();

  const [number, setNumber] = useState("");
  const [err, setErr] = useState("");

  const { mutate, isPending, isError, isSuccess, error } = useMutation({
    mutationFn: sendOtp,
    onSuccess: (data) => {
      // alert("Logged in ✅");
      localStorage.setItem("session-token", data.token);
      setOtpSent(true);
    },
    onError: (error) => {
      console.log("Login failed: ", error);
    },
  });

  if (otpSent) {
    return <Verify />;
  }

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!number) {
      setErr("Number is required.");
      return;
    }
    setErr("");

    mutate({
      shopId: shopId,
      tableId: tableId,
      phoneNumber: number,
    });
  };

  return (
    <div className="w-screen h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-md w-full mx-4">
        <div className="bg-white rounded-2xl shadow-2xl p-8 border-2 border-gray-200">
          <div className="text-center mb-6">
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Welcome</h1>
            <p className="text-gray-600 text-sm">Enter your phone number to continue</p>
          </div>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Input
              label="Phone"
              type="number"
              name="number"
              placeholder="Enter your number"
              helperText={err}
              error={!!err}
              value={number}
              onChange={(e) => setNumber(e.target.value)}
            />
            <Button type="submit" variant="primary">
              {isPending ? "Sending..." : "Send OTP"}
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};

const Verify = () => {
  return <div className="">verify</div>;
};

const JoinWithOTP = () => {};

const Session = () => {
  const { shopid, tableid } = useParams();
  const { shopId, tableId, setShopId, setTableId } = useAuthStore();

  if (shopid != null && shopid !== shopId) setShopId(shopid);
  if (tableid != null && tableid !== tableId) setTableId(tableid);

  // const data = null;

  // return (
  //   <div className="w-full h-full">
  //     <SendOTP />
  //   </div>
  // );

  const { data, isLoading, isFetched, error } = useQuery({
    queryKey: [shopid, tableid],
    queryFn: () => getSession({ shopid, tableid }),
  });

  if (isLoading) {
    return <CircleLoading fullScreen text="Loading session..." />;
  }

  if (isFetched) {
    return data == null ? (
      <div className="w-full h-full">
        <SendOTP />
      </div>
    ) : (
      <div className="w-full h-full">
        <JoinWithOTP />
      </div>
    );
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }
};

export default Session;
