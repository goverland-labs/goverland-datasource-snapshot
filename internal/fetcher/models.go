package fetcher

type Message struct {
	Space    string `json:"space"`
	Proposal string `json:"proposal"`
}

type Data struct {
	Message Message `json:"message"`
}

type Info struct {
	Data Data `json:"data"`
}
