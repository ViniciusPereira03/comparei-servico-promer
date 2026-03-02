package user

type User struct {
	ID          string `json:"id"`
	Status      int    `json:"status"`
	RayDistance int    `json:"ray_distance"`
}
