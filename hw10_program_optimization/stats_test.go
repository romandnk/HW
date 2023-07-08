//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func TestGetDomainStat(t *testing.T) {
	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func BenchmarkGetDomainStat(b *testing.B) {
	b.StopTimer()

	users := generateUsers(b.N)
	content := prepareData(users)
	reader := strings.NewReader(content)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(reader, "com")
	}
}

func generateUsers(n int) []User {
	users := make([]User, n)
	for i := 0; i < n; i++ {
		addition := fmt.Sprintf("%d", i)
		users[i] = User{
			ID:       i,
			Name:     "test" + addition,
			Username: "test" + addition,
			Email:    "test" + addition + "@test.com",
			Phone:    "1234567890",
			Password: "test",
			Address:  "test" + addition,
		}
	}
	return users
}

func prepareData(users []User) string {
	var sb strings.Builder
	for _, user := range users {
		jsonStr, _ := json.Marshal(user)
		sb.WriteString(string(jsonStr))
		sb.WriteString("\n")
	}
	return sb.String()
}

//func BenchmarkGetDomainStat(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_, _ = GetDomainStat(bytes.NewBufferString(data), "com")
//	}
//}
