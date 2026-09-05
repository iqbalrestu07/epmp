import { useState } from "react";

export function MessagingSettingsPage() {
  const [showQR, setShowQR] = useState(false);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Messaging Settings</h1>
          <p className="text-muted-foreground text-sm">
            Manage your connected WhatsApp devices for sending notifications.
          </p>
        </div>
        <button
          onClick={() => setShowQR(true)}
          className="bg-orange text-slate-900 px-4 py-2 rounded-lg font-medium hover:bg-orange/90"
        >
          Add WhatsApp Device
        </button>
      </div>

      <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
        <h3 className="font-semibold text-lg mb-4">Connected Devices</h3>
        
        <div className="border border-slate-300 rounded-lg p-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
              <span className="text-green-600 font-bold">WA</span>
            </div>
            <div>
              <p className="font-medium">+62 812-3456-7890</p>
              <p className="text-xs text-green-600 flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-green-600 inline-block"></span>
                Connected (Finance Department)
              </p>
            </div>
          </div>
          <button className="text-red-500 text-sm hover:underline">Disconnect</button>
        </div>
      </div>

      {showQR && (
        <div className="fixed inset-0 z-50 bg-slate-1000 flex items-center justify-center">
          <div className="bg-white p-8 rounded-2xl w-full max-w-md text-center shadow-xl">
            <h2 className="text-xl font-bold mb-2">Scan QR Code</h2>
            <p className="text-sm text-slate-500 mb-6">
              Open WhatsApp on your phone, tap Menu or Settings and select Linked Devices. Point your phone to this screen to capture the code.
            </p>
            
            <div className="bg-slate-100 w-64 h-64 mx-auto rounded-lg flex items-center justify-center border-2 border-dashed border-slate-300 mb-6">
              <span className="text-slate-400">Loading QR... (Backend Stub)</span>
            </div>
            
            <button 
              onClick={() => setShowQR(false)}
              className="text-slate-500 hover:text-slate-900 font-medium"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
