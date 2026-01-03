package models

type RootMeProfile struct {
	Nom         string             `json:"nom"`
	LogoURL     string             `json:"logo_url"`
	Score       string             `json:"score"`
	Position    int                `json:"position"`
	Validations []Validation       `json:"validations"`
}

type Validation struct {
	IDChallenge string `json:"id_challenge"`
	Titre       string `json:"titre"`
	IDRubrique  string `json:"id_rubrique"`
	Date        string `json:"date"`
}
