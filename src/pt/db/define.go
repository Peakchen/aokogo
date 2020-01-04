package db

type TQueryLimit struct {
	Min int
	Max int
	GoRoute int
}

const (
	CstMinLimit = int(5000) 
)

const (
	CstPointformat = string(" . ")
)

var (
	CstQueryLimit = []TQueryLimit{
		TQueryLimit{
			Min: 5000,
			Max: 10000,
			GoRoute: 20,
		},
		TQueryLimit{
			Min: 10001,
			Max: 20000,
			GoRoute: 30,
		},
		TQueryLimit{
			Min: 20001,
			Max: 40000,
			GoRoute: 40,
		},
		TQueryLimit{
			Min: 50001,
			Max: 80000,
			GoRoute: 50,
		},
		TQueryLimit{
			Min: 80001,
			Max: 100000,
			GoRoute: 60,
		},
		TQueryLimit{
			Min: 100001,
			Max: 999999999,
			GoRoute: 70,
		},
	}
)

type TQueryField struct {
	Field string
	Data interface{}
}
