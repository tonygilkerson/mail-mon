package msg

type MsgKey string

const (
	MbxTemperature            MsgKey = "MailboxTemperature"
	MbxMuleAlarm                     = "MuleAlarm"
	MbxDoorOpened                    = "MailboxDoorOpened"
	MbxChargerChargeStatusOn         = "ChargerChargeStatusOn"
	MbxChargerChargeStatusOff        = "ChargerChargeStatusOff"
	MbxChargerPowerSourceGood        = "ChargerPowerSourceGood"
	MbxChargerPowerSourceBad         = "ChargerPowerSourceBad"
	MbxRoadMainLoopHeartbeat         = "RoadMainLoopHeartbeat"
)