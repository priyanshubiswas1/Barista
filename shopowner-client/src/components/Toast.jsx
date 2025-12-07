const Toast = ({ children, action, onAction }) => {
  return (
    <div className="flex px-3 py-2 gap-8 bg-gray-950 rounded-md shadow-sm items-center">
      <p className="text-white font-medium">{children}</p>
      {action && (
        <button
          onClick={onAction}
          className="text-sm text-indigo-500 font-medium cursor-pointer"
        >
          {action}
        </button>
      )}
    </div>
  );
};

export default Toast;
