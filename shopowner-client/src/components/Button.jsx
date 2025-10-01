import React from "react";

const BASE_CLASSES =
  "rounded-md shadow-xs disabled:shadow-none text-sm font-medium cursor-pointer disabled:cursor-not-allowed transition-colors duration-120 ease-in-out";

const VARIANTS = {
  primary: "bg-gray-800 text-white hover:bg-gray-950 disabled:bg-gray-100",
  outline:
    "bg-white text-gray-800 border border-gray-200 hover:border-gray-300 disabled:text-gray-200 disabled:border-gray-50",
  subtle:
    "bg-gray-200 text-gray-500 hover:text-gray-700 disabled:bg-gray-100 disabled:text-gray-200",
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
    : "px-4 py-2 flex items-center gap-2";

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
