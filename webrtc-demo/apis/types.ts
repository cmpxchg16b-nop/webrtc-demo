// type ConnectionAttributes map[string]string
export type ConnectionAttributes = Record<string, string>;

// type MessagePayload struct {
//  	Register               *pkgconnreg.RegisterPayload               `json:"register,omitempty"`
//  	Echo                   *pkgconnreg.EchoPayload                   `json:"echo,omitempty"`
//  	AttributesAnnouncement *pkgconnreg.AttributesAnnouncementPayload `json:"attributes_announcement,omitempty"`
// }
export interface MessagePayload {
  node_id?: string;
  register?: RegisterPayload;
  echo?: EchoPayload;
  online?: any;
  rename?: any;
  attributes_announcement?: AttributesAnnouncementPayload;
}

// type RegisterPayload struct {
//  	NodeName string  `json:"node_name"`
//  	Token    *string `json:"token,omitempty"`
// }
export interface RegisterPayload {
  node_name: string;
  token?: string;
}

// type EchoDirection string
// const (
//  	EchoDirectionC2S EchoDirection = "ping"
//  	EchoDirectionS2C EchoDirection = "pong"
// )
export type EchoDirection = "ping" | "pong";

export const EchoDirectionC2S = "ping" as const;
export const EchoDirectionS2C = "pong" as const;

// type EchoPayload struct {
//  	Direction       EchoDirection `json:"direction"`
//  	CorrelationID   string        `json:"correlation_id"`
//  	ServerTimestamp uint64        `json:"server_timestamp"`
//  	Timestamp       uint64        `json:"timestamp"`
//  	SeqID           uint64        `json:"seq_id"`
// }
export interface EchoPayload {
  direction: EchoDirection;
  correlation_id: string;
  server_timestamp: number;
  timestamp: number;
  seq_id: number;
}

// type AttributesAnnouncementPayload struct {
//  	Attributes  ConnectionAttributes `json:"attributes,omitempty"`
//  	Withdrawals []string             `json:"withdrawals,omitempty"`
// }
export interface AttributesAnnouncementPayload {
  attributes?: ConnectionAttributes;
  withdrawals?: string[];
}

// type AuthenticationType string
// const (
// 	AuthenticationTypeJWT  AuthenticationType = "jwt"
// 	AuthenticationTypeMTLS AuthenticationType = "mtls"
// )
export type AuthenticationType = "jwt" | "mtls";

export const AuthenticationTypeJWT = "jwt" as const;
export const AuthenticationTypeMTLS = "mtls" as const;

// type ConnRegistryData struct {
// 	NodeName       *string              `json:"node_name,omitempty"`
// 	ConnectedAt    uint64               `json:"connected_at"`
// 	RegisteredAt   *uint64              `json:"registered_at,omitempty"`
// 	LastHeartbeat  *uint64              `json:"last_heartbeat,omitempty"`
// 	Attributes     ConnectionAttributes `json:"attributes,omitempty"`
// 	QUICConn       *quicGo.Conn         `json:"-"`
// 	Claims         jwt.MapClaims        `json:"-"`
// 	Authentication AuthenticationType   `json:"authentication"`
// }
export interface ConnRegistryData {
  node_name?: string;
  connected_at: number;
  registered_at?: number;
  last_heartbeat?: number;
  attributes?: ConnectionAttributes;
  authentication: AuthenticationType;
}
