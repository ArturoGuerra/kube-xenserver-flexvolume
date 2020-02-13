package xapi


type (
    XClient interface {
        Attach
        Detach
    }

    xClient struct {
       Username string
       Password string
       Host     string
    }
)

func New(username string, password string, host string) *XClient {
    return &xClient{
        username,
        password,
        host,
    }
}
