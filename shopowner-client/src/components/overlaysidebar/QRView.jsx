import { PiPrinter } from "react-icons/pi";
import Button from "../Button";
import { useState } from "react";

const QRView = ({ table, onClose }) => {
  const [qrOnly, setQrOnly] = useState(false);
  return (
    <div className="fixed inset-0 z-50 flex flex-col gap-6 items-center justify-center bg-black/70">
      {qrOnly ? (
        <div className="rounded-md w-[20%] h-[50%] flex flex-col items-center justify-center p-6">
          <div className="bg-white rounded-md p-2">
            <img
              src={`data:image/png;base64,${table.qrUrl}`}
              alt={`QR for Table ${table.tableNumber}`}
              className="w-32 h-32"
            />
          </div>
        </div>
      ) : (
        <div
          className="bg-indigo-800 text-white rounded-md w-[20%] h-[50%] flex flex-col items-center p-6 relative"
          style={{
            backgroundImage: "url(/card-pattern.png)",
            backgroundSize: "cover",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
          }}
        >
          <div className="flex flex-col items-center mb-4">
            <div className="bg-white rounded-full h-10 w-10 mb-2"></div>
            <h2 className="font-medium text-xl">Cafe Name</h2>
          </div>
          <div className="flex flex-col items-center mb-2">
            <h3 className="font-extralight text-sm">
              Table {table.tableNumber}
            </h3>
          </div>
          <div className="bg-white rounded-md p-2">
            <img
              src={`data:image/png;base64,${table.qrUrl}`}
              alt={`QR for Table ${table.tableNumber}`}
              className="w-32 h-32"
            />
          </div>
          <h2 className="font-extralight text-sm mb-4 mt-10">www.cafe.com</h2>
        </div>
      )}
      <div className="flex gap-2">
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
        <Button variant="outline" icon={PiPrinter}>
          Print
        </Button>
        <Button variant="primary" onClick={() => setQrOnly(!qrOnly)}>
          QR only
        </Button>
      </div>
    </div>
  );
};

export default QRView;
