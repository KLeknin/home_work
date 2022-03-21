package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

var ErrNilInput = errors.New("nil input")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrNilInput
	}
	reader := bufio.NewReader(r)
	var user User
	var err error
	domainStat := make(DomainStat)
	for {
		line, readError := reader.ReadSlice(byte('\n'))
		if readError != nil && !errors.Is(readError, io.EOF) {
			return nil, err
		}
		if len(line) < 1 {
			line = []byte("{}")
		}

		if err = jsoniter.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		domainName := getDomainName(user, domain)
		if domainName > "" {
			domainStat[domainName]++
		}

		if errors.Is(readError, io.EOF) {
			break
		}
	}
	return domainStat, nil
}

func getDomainName(u User, domain string) (domainName string) {
	domainName = ""
	if strings.HasSuffix(u.Email, "."+domain) {
		dogPlace := strings.Index(u.Email, "@")
		domainName = strings.ToLower(u.Email[dogPlace+1:])
	}
	return
}
