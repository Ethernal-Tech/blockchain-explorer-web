package paginationModel

type PaginationData struct {
	NextPage     int
	PreviousPage int
	CurrentPage  int
	TotalPages   int
	TotalRows    uint64
	PerPage      int
}
