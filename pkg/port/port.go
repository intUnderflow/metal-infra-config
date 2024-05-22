package port

import "strconv"

type Port int

func ParsePort(portString string) (Port, error) {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return 0, err
	}

	return Port(port), nil
}
