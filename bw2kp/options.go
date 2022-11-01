package bw2kp

type Options struct {
	BitwardenSession string `short:"s" long:"bw-session" description:"The session string for Bitwarden" required:"true"`
	DatabasePath     string `short:"l" long:"database-path" description:"The path of the generated Keepass database" required:"true"`
	DatabasePassword string `short:"p" long:"database-password" description:"The password to unlock the generated Keepass database" required:"true"`
}
