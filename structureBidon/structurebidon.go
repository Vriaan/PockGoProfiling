package structurebidon

type Bidon struct {
	someText string
	someID   string
}

func New() *Bidon {
	return &Bidon{}
}
