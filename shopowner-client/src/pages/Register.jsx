import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import { useMutation } from "@tanstack/react-query";
import { registerUser } from "../api/api";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/authStore";
import { toast } from "react-toastify";

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
      toast.success("Registration successful!");
      localStorage.setItem("token", data.token);
      localStorage.setItem("shopId", data.shopId);
      localStorage.setItem("menuId", data.menuId);
      localStorage.setItem("name", data.name);
      localStorage.setItem("email", data.email);
      useAuthStore.setState({ isAuth: true });

      navigate("/menu");
    },
    onError: (error) => {
      toast.error("Registration failed: " + error.message);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newError = {
      name: "",
      email: "",
      password: "",
      phone: "",
      address: "",
    };

    if (!formData.name) newError.name = "Please enter your name";
    if (!formData.email) newError.email = "Please enter your email";
    else if (!formData.email.includes("@"))
      newError.email = "Enter a valid email";
    if (!formData.phone) newError.phone = "Please enter your phone number";
    if (!formData.password) newError.password = "Please enter your password";
    if (!formData.address) newError.address = "Please enter your address";

    setErr(newError);
    if (Object.values(newError).some((error) => error !== "")) return;

    mutate(formData);
  };

  return (
    <div className="flex h-screen bg-gradient-to-br from-gray-50 via-indigo-50/30 to-gray-100">
      <div className="h-screen w-screen flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden relative">
          {/* Header with gradient */}
          <div className="bg-gradient-to-br from-indigo-600 to-indigo-700 text-white p-6 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />

            <div className="relative z-10 flex flex-col items-center">
              <div className="bg-white rounded-full h-14 w-14 mb-2 flex items-center justify-center shadow-lg">
                <span className="text-3xl">☕</span>
              </div>
              <h2 className="font-bold text-xl mb-1">Join Barista</h2>
              <p className="text-indigo-100 text-xs">
                Create your café account
              </p>
            </div>
          </div>

          {/* Form Section */}
          <div className="p-6 bg-gradient-to-b from-gray-50 to-white">
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
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
                {isPending ? "Creating..." : "Create Account"}
              </Button>
              <div className="flex justify-center w-full gap-1 items-center text-sm">
                <p className="text-gray-600">Already have an account?</p>
                <button
                  type="button"
                  className="text-indigo-700 hover:text-indigo-800 hover:underline cursor-pointer font-medium transition-colors"
                  onClick={() => navigate("/login")}
                >
                  Sign In
                </button>
              </div>
            </form>

            <div className="mt-4 text-center">
              <p className="text-gray-400 text-xs">Powered by Barista</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register;
