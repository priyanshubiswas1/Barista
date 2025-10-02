import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import { useMutation } from "@tanstack/react-query";
import { registerUser } from "../api/api";
import { useNavigate } from "react-router-dom";

const Register = () => {
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    phone: "",
    address: "",
  });
  const [err, setErr] = useState({
    name: "",
    email: "",
    password: "",
    phone: "",
    address: "",
  });
  const navigate = useNavigate();

  const { mutate, isPending, isError, isSuccess, error } = useMutation({
    mutationFn: registerUser,
    onSuccess: (data) => {
      // alert("Logged in ✅");
      localStorage.setItem("token", data.token);
      localStorage.setItem("shopId", data.shopId);
      localStorage.setItem("menuId", data.menuId);
      navigate("/menu");
    },
    onError: (error) => {
      console.log("Registration failed: ", error);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newError = { email: "", password: "" };

    if (!formData.email) newError.email = "Email is required";
    else if (!formData.email.includes("@"))
      newError.email = "Enter a valid email";
    if (!formData.password) newError.password = "Password is required";

    setErr(newError);
    if (newError.email || newError.password) return;

    mutate(formData);
  };

  return (
    <div className="flex h-screen">
      <div className="h-screen w-screen flex flex-col gap-6 items-center justify-center bg-white">
        <header className="flex flex-col gap-2 items-center">
          <h1 className="text-2xl font-semibold text-gray-800">
            Register to MyCafeApp
          </h1>
        </header>
        <div className="w-100 flex flex-col gap-2 bg-white p-5 rounded-md">
          <form onSubmit={handleSubmit} className="flex flex-col gap-2">
            <Input
              label="Name"
              type="text"
              name="name"
              placeholder="Enter your name"
              helperText={err.name}
              error={!!err.name}
              value={formData.name}
              onChange={handleChange}
            />
            <Input
              label="Email"
              type="email"
              name="email"
              placeholder="Enter your email"
              helperText={err.email}
              error={!!err.email}
              value={formData.email}
              onChange={handleChange}
            />
            <Input
              label="Phone"
              type="text"
              name="phone"
              placeholder="Enter your phone number"
              prefix={"+91"}
              helperText={err.phone}
              error={!!err.phone}
              value={formData.phone}
              onChange={handleChange}
            />
            <Input
              label="Password"
              type="password"
              name="password"
              placeholder="Enter your password"
              autocomplete="new-password"
              helperText={err.password}
              error={!!err.password}
              value={formData.password}
              onChange={handleChange}
            />
            <Input
              label="Address"
              type="text"
              name="address"
              placeholder="Enter your address"
              helperText={err.address}
              error={!!err.address}
              value={formData.address}
              onChange={handleChange}
            />

            <Button type="submit" variant="primary">
              {isPending ? "Registering..." : "Register"}
            </Button>
            <div className="flex justify-center w-full gap-1 items-center">
              <p className="text-gray-500">Already have an account?</p>
              <button
                type="button"
                className=" text-sm text-gray-800 hover:underline cursor-pointer"
                onClick={() => navigate("/login")}
              >
                Sign In
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Register;
