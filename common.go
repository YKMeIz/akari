package akari

import ()

type Core struct {
	Domain                string
	Port                  string
	CertChain             string
	CertKey               string
	MessageRelativePath   string
	WebsocketRelativePath string
	DatabasePath          string
	Event                 map[string]HandlerFunc
}

type HandlerFunc func(*Message) error

type Message struct {
	// Message sender's token.
	Source string

	// Message receiver's token.
	Destination []string

	// Message content.
	Data map[string]string
}

type PushbulletPush struct {
	pushType    string
	title       string
	body        string
	accessToken string
}

const (
	REGISTEROK string = `{"Status": "ok!"}`
	REGISTERER string = `Wrong user name or token.`
	REGISTEROL string = `Your device is already online.`
	REGISTERTO string = `Id authentication is required.`
	PBPUSHNERR string = `Fail to send Pushbullet notification.`
	HANDLERFER string = `An error occured on running custom handler function.`
	MESSAGEERR string = `Missing source or destination.`
)

var (
	DatabasePath string
)

func formatErrInfo(errInfo string) string {
	return `{"Status": "error! ` + errInfo + `"}`
}
