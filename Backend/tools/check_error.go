package tools


func CheckError(err error) error {
	if err != nil {
		return err
	} else {
		return nil
	}
}