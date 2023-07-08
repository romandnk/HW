package hw10programoptimization

import (
	"bufio"

	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result, err := countDomains(r, domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	var user User

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, "."+domain) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
