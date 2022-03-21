// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

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

	t.Run("find ''", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("nil input", func(t *testing.T) {
		_, err := GetDomainStat(nil, "domain")
		if !errors.Is(err, ErrNilInput) {
			t.Fatalf("Wrong error, exp:\"%v\" get:\"%v\"", ErrNilInput, err)
		}
	})

	t.Run("empty string input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "domain")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}

{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("separated strings", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "domain")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

}
func BenchmarkGetDomainStat(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()

	if len(r.File) != 1 {
		b.Fatal("len")
	}

	data, err := r.File[0].Open()
	if err != nil {
		b.Fatal(err)
	}
	content, err := ioutil.ReadAll(data)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		data := bytes.NewBuffer(content)
		b.StartTimer()
		stat, err := GetDomainStat(data, "biz")
		_ = stat
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
	}

}
