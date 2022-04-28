package TagPointType

type Enum int

const (
	ADD_POST  Enum = 20
	ADD_FLOOR Enum = 10

	LIKE_POST Enum = 4
	FAV_POST  Enum = 3
	DIS_POST  Enum = 4

	UNLIKE_POST Enum = -4
	UNFAV_POST  Enum = -3
	UNDIS_POST  Enum = -4

	LIKE_FLOOR Enum = 1
	DIS_FLOOR  Enum = 1

	UNLIKE_FLOOR Enum = -1
	UNDIS_FLOOR  Enum = -1
)
