package value_objects

type UserType int

const (
	memberType   = "member"
	employeeType = "employee"
	invalidType  = "invalid"
)

const (
	Member UserType = iota
	Employee
	Invalid
)

func NewUserType(userType UserType) UserType {
	return userType
}

func NewUserTypeFromString(userType string) UserType {
	switch userType {
	case memberType:
		return Member
	case employeeType:
		return Employee
	default:
		return Invalid
	}
}

func (s UserType) String() string {
	return [...]string{memberType, employeeType, invalidType}[s]
}

func (s UserType) IsValid() bool {
	return s != Invalid
}

func (s UserType) Equals(other UserType) bool {
	return s == other
}
