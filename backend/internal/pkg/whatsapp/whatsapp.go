package whatsapp

// WhatsAppClientManager manages multiple whatsmeow clients.
// Each organization or device will have its own session.
type WhatsAppClientManager struct {
	// clients map[string]*whatsmeow.Client
}

func NewWhatsAppClientManager() *WhatsAppClientManager {
	return &WhatsAppClientManager{}
}

// TODO: Implement ConnectDevice(qrCallback func(string))
// TODO: Implement SendMessage(deviceID string, jid string, message string)
// TODO: Implement BlastMessage(deviceID string, jids []string, message string)
