package dirs

import "os/user"

func GetHomeDir() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.HomeDir
}
