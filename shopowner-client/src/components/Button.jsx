import React from "react";

const BASE_CLASSES =
  "rounded-lg shadow-sm disabled:shadow-none text-sm font-medium cursor-pointer disabled:cursor-not-allowed transition-all duration-200 ease-in-out";

const VARIANTS = {
  primary:
    "bg-indigo-700 text-white hover:bg-indigo-800 hover:shadow-md disabled:bg-gray-100 disabled:text-gray-400",
  outline:
    "bg-white text-gray-700 border border-gray-200 hover:border-indigo-300 hover:text-indigo-700 hover:bg-indigo-50/50 disabled:text-gray-200 disabled:border-gray-50",
  subtle:
    "bg-gray-100 text-gray-600 hover:bg-gray-200 hover:text-gray-800 disabled:bg-gray-50 disabled:text-gray-300",
};

const Button = ({
  children,
  variant = "primary",
  icon: Icon,
  onClick,
  ...props
}) => {
  const variantClasses = VARIANTS[variant] || VARIANTS.primary;

  const isIconOnly = Icon && React.Children.count(children) === 0;

  // Icon-only: square button, no gap
  const paddingClasses = isIconOnly
    ? "w-8 h-8 p-0 flex items-center justify-center"
    : "px-4 h-8 flex items-center gap-2 justify-center";

  const handleClick = (e) => {
    if (onClick) {
      onClick(e);
    }
  };

  return (
    <button
      className={`${BASE_CLASSES} ${paddingClasses} ${variantClasses}`}
      onClick={handleClick}
      {...props}
    >
      {Icon && (
        <span
          className={
            isIconOnly ? "w-4 h-4 flex items-center justify-center" : ""
          }
        >
          <Icon size={16} weight="bold" color="currentColor" />
        </span>
      )}
      {children}
    </button>
  );
};

export default Button;
