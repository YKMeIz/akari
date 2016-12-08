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

type HandlerFunc func() error

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
	REGISTERER string = `{"Status": "error! Wrong user name or token."}`
	REGISTEROL string = `{"Status": "error! Your device is already online."}`
	REGISTERTO string = `{"Status": "error! Id authentication is required."}`
	PBPUSHNERR string = `{"Status": "error! Fail to send Pushbullet notification."}`
	HANDLERFER string = `{"Status": "error! An error occured on running custom handler function."}`
	MESSAGEERR string = `{"Status": "error! Missing source or destination."}`
)

var (
	DatabasePath string
)
