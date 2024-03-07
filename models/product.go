package models

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	CodeValue    string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseProducts struct {
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type ResponseProduct struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	CodeValue    string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
	Error   bool    `json:"error"`
}
