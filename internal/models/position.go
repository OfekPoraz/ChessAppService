package models

type PossibleMovesPosition struct {
	Row       int  `json:"row"`
	Col       int  `json:"col"`
	IsCapture bool `json:"isCapture"`
}
