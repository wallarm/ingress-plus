package nginx

func ParseLBMethod(method string) (string, error) {
	if method == "round_robin" {
		return "", nil
	}
	return method, nil
}

func ParseLBMethodForPlus(method string) (string, error) {
	if method == "round_robin" {
		return "", nil
	}
	return method, nil
}
