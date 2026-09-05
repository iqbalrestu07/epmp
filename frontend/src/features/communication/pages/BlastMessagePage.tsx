import { useState } from "react";

export function BlastMessagePage() {
  const [message, setMessage] = useState("");

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Blast Message</h1>
        <p className="text-muted-foreground text-sm">
          Send bulk WhatsApp messages to your tenants.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2 bg-white p-6 rounded-xl shadow-sm border border-slate-200 space-y-4">
          
          <div>
            <label className="block text-sm font-medium mb-1">Select Sender Device</label>
            <select className="w-full p-2 border border-slate-300 rounded-lg">
              <option>+62 812-3456-7890 (Finance Dept)</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Target Audience</label>
            <select className="w-full p-2 border border-slate-300 rounded-lg">
              <option>All Active Tenants</option>
              <option>Tenants with Overdue Invoices</option>
              <option>Specific Building (Building A)</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Message Content</label>
            <textarea 
              rows={6}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              className="w-full p-3 border border-slate-300 rounded-lg resize-none"
              placeholder="Hello {{tenant_name}}, this is an announcement regarding..."
            />
            <p className="text-xs text-slate-500 mt-1">Available variables: {'{{tenant_name}}, {{room_name}}, {{invoice_amount}}'}</p>
          </div>

          <div className="pt-4 flex justify-end">
            <button className="bg-orange text-slate-900 px-6 py-2 rounded-lg font-medium hover:bg-orange/90 shadow-sm">
              Send Blast
            </button>
          </div>
        </div>

        {/* Preview Panel */}
        <div className="bg-[#EFEAE2] p-6 rounded-xl border border-slate-200 h-fit">
          <h3 className="font-semibold text-sm mb-4 uppercase tracking-wider text-slate-500">Preview</h3>
          <div className="bg-[#dcf8c6] p-3 rounded-lg rounded-tr-none shadow-sm relative text-sm">
            {message ? message.replace('{{tenant_name}}', 'John Doe') : "Your message preview will appear here..."}
            <span className="block text-[10px] text-slate-500 text-right mt-2">12:00 PM</span>
          </div>
        </div>

      </div>
    </div>
  );
}
