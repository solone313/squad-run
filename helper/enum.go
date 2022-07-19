package helper

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (s Weekday) String() string {
	return WeekdayToString[s]
}

var WeekdayToString = map[Weekday]string{
	Sunday:    "Sunday",
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
}

var WeekdayToID = map[string]Weekday{
	"Sunday":    Sunday,
	"Monday":    Monday,
	"Tuesday":   Tuesday,
	"Wednesday": Wednesday,
	"Thursday":  Thursday,
	"Friday":    Friday,
	"Saturday":  Saturday,
}
