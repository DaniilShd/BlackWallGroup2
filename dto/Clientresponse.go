package dto

import "strconv"

type ClientResponse struct {
	ClientId int    `json:"client_id"`
	Sum      string `json:"sum"`
}

func (c *ClientResponse) CheckSum(amount int) bool {
	sum, _ := strconv.Atoi(c.Sum)
	return sum < amount
}
