package api

type API struct {
	Client *Client
	Skills *SkillService
}

func NewAPI(client *Client) *API {
	return &API{
		Client: client,
		Skills: NewSkillService(client),
	}
}
