# MODULE.md — Communication & WhatsApp Gateway

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Communication        |
| Package          | communication        |
| Tables           | `wa_devices`, `message_templates`, `blast_messages`, `blast_message_logs`, `whatsmeow_device` |
| Bounded Context  | messaging            |
| Module Path      | github.com/epmp/backend/internal/modules/communication |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | `/api/v1/communication` |
| Operations   | CRUD Devices, QR Pairing, CRUD Templates, Blast Messaging, Recipient Preview |

## Endpoints

| Method | Route | Description |
| ------ | ----- | ----------- |
| `GET` | `/devices` | List registered WhatsApp devices |
| `POST` | `/devices` | Register new device (label only, phone auto-detected) |
| `GET` | `/devices/:id/qr` | Get live WhatsApp multi-device pairing QR string |
| `PUT` | `/devices/:id/status` | Disconnect or update device status |
| `DELETE` | `/devices/:id` | Disconnect and remove WhatsApp device & session |
| `GET` | `/templates` | List active message templates |
| `POST` | `/templates` | Create message template with variables |
| `PUT` | `/templates/:id` | Update message template |
| `DELETE` | `/templates/:id` | Soft delete message template |
| `GET` | `/blast` | List blast message campaigns |
| `POST` | `/blast` | Create blast message draft / campaign |
| `GET` | `/blast/:id` | Get blast campaign details |
| `GET` | `/blast/:id/logs` | Get per-recipient delivery logs |
| `POST` | `/blast/:id/send` | Trigger blast dispatch via whatsmeow |
| `GET` | `/recipients` | Preview target recipients by audience filter |

## Behaviors

| Behavior     | Enabled | Details |
| ------------ | ------- | ------- |
| Multi-Device | true    | WhatsApp Web Multi-Device via native `whatsmeow` |
| QR Pairing   | true    | WhatsApp noise handshake raw QR code |
| Auto-JID     | true    | Auto-detects phone number from scan callback |
| Validation   | true    | `IsOnWhatsApp` checks phone registration before send |
| Auto-Prefix  | true    | Converts local `08...` to international `628...` |
| Persistence  | true    | Encrypted session keys stored in PostgreSQL |
| Auto-Reconnect | true  | Reconnects active devices on server boot |

## Structure

```
internal/modules/communication/
  module.go      # Echo HTTP handlers & route registration
  wa_manager.go  # Whatsmeow client manager, session persistence, IsOnWhatsApp
  MODULE.md      # Module specification
```
