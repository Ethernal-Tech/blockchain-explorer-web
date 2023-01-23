package paginationModel

type PaginationData struct {
	NextPage     int
	PreviousPage int
	CurrentPage  int
	TotalPages   int
	TotalRows    int64
	PerPage      int
}
