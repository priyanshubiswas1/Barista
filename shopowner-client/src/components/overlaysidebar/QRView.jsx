import { PiPrinter, PiDownloadSimple } from "react-icons/pi";
import { IoClose } from "react-icons/io5";
import Button from "../Button";
import { useState } from "react";

const QRView = ({ table, onClose }) => {
  const [qrOnly, setQrOnly] = useState(false);

  const handleDownload = () => {
    const link = document.createElement("a");
    link.href = `data:image/png;base64,${table.qrUrl}`;
    link.download = `table-${table.tableNumber}-qr.png`;
    link.click();
  };

  const handlePrint = () => {
    const printWindow = window.open("", "", "width=400,height=600");
    printWindow.document.write(`
      <html>
        <head>
          <title>Table ${table.tableNumber} QR Code</title>
          <style>
            body { 
              display: flex; 
              justify-content: center; 
              align-items: center; 
              height: 100vh; 
              margin: 0; 
              font-family: Arial, sans-serif;
            }
            .qr-card {
              text-align: center;
              padding: 40px;
            }
            img { 
              width: 300px; 
              height: 300px; 
            }
            h2 { margin: 20px 0 10px 0; }
          </style>
        </head>
        <body>
          <div class="qr-card">
            <h1>Barista</h1>
            <h2>Table ${table.tableNumber}</h2>
            <img src="data:image/png;base64,${table.qrUrl}" alt="QR Code" />
            <p>Scan to order</p>
          </div>
        </body>
      </html>
    `);
    printWindow.document.close();
    printWindow.print();
  };

  return (
    <div className="fixed inset-0 z-50 flex flex-col gap-6 items-center justify-center bg-black/80 backdrop-blur-sm p-4">
      {qrOnly ? (
        <div className="bg-white rounded-2xl shadow-2xl p-8 flex flex-col items-center justify-center">
          <div className="bg-gray-50 rounded-xl p-4 border-2 border-gray-200">
            <img
              src={`data:image/png;base64,${table.qrUrl}`}
              alt={`QR for Table ${table.tableNumber}`}
              className="w-64 h-64"
            />
          </div>
          <p className="mt-4 text-gray-600 text-sm font-medium">
            Table {table.tableNumber}
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden relative">
          {/* Header with gradient */}
          <div className="bg-gradient-to-br from-indigo-600 to-indigo-700 text-white p-8 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />

            <div className="relative z-10 flex flex-col items-center">
              <div className="bg-white rounded-full h-16 w-16 mb-3 flex items-center justify-center shadow-lg">
                <span className="text-3xl">☕</span>
              </div>
              <h2 className="font-bold text-2xl mb-1">Barista</h2>
              <p className="text-indigo-100 text-sm">Scan & Order</p>
            </div>
          </div>

          {/* QR Code Section */}
          <div className="p-8 flex flex-col items-center bg-gradient-to-b from-gray-50 to-white">
            <div className="mb-4">
              <div className="inline-block bg-indigo-100 border-2 border-indigo-200 rounded-lg px-4 py-2">
                <h3 className="font-bold text-indigo-700 text-lg">
                  Table {table.tableNumber}
                </h3>
              </div>
            </div>

            <div className="bg-white rounded-2xl p-6 shadow-lg border-2 border-gray-200">
              <img
                src={`data:image/png;base64,${table.qrUrl}`}
                alt={`QR for Table ${table.tableNumber}`}
                className="w-48 h-48"
              />
            </div>

            <div className="mt-6 text-center">
              <p className="text-gray-600 text-sm font-medium">
                Scan this code to place your order
              </p>
              <p className="text-gray-400 text-xs mt-1">Powered by Barista</p>
            </div>
          </div>
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex gap-3">
        <button
          onClick={onClose}
          className="bg-white hover:bg-gray-100 text-gray-700 font-medium px-6 py-3 rounded-xl shadow-lg border-2 border-gray-200 transition-all duration-200 hover:shadow-xl flex items-center gap-2"
        >
          <IoClose size={20} />
          Close
        </button>
        <button
          onClick={handleDownload}
          className="bg-white hover:bg-gray-100 text-gray-700 font-medium px-6 py-3 rounded-xl shadow-lg border-2 border-gray-200 transition-all duration-200 hover:shadow-xl flex items-center gap-2"
        >
          <PiDownloadSimple size={20} />
          Download
        </button>
        <button
          onClick={handlePrint}
          className="bg-white hover:bg-gray-100 text-gray-700 font-medium px-6 py-3 rounded-xl shadow-lg border-2 border-gray-200 transition-all duration-200 hover:shadow-xl flex items-center gap-2"
        >
          <PiPrinter size={20} />
          Print
        </button>
        <button
          onClick={() => setQrOnly(!qrOnly)}
          className="bg-indigo-700 hover:bg-indigo-800 text-white font-medium px-6 py-3 rounded-xl shadow-lg transition-all duration-200 hover:shadow-xl"
        >
          {qrOnly ? "Show Card" : "QR Only"}
        </button>
      </div>
    </div>
  );
};

export default QRView;
