package db

type DB struct {
}

type Cluster struct {
	Master DB
	Replica DB
}

func (d DB) String() string {
	return "db"
}

func (c Cluster) CleanUp() error {
	return nil
}