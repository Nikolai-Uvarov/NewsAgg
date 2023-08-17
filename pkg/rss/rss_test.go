//Authored by Nikolai Uvarov
//Allowed for free copying
//Package rss implements a simple mechanism for listening of certain rss news feed

package rss

import (
	"reflect"
	"testing"
)

func Test_readConfig(t *testing.T) {
	tests := []struct {
		name  string
		want  int
		want1 []string
	}{
		// TODO: Add test cases.
		{
			name: "onlyone",
			want: 5,
			want1: []string{
				"http://www.kommersant.ru/RSS/news.xml",
				"https://habr.com/ru/rss/best/daily/?fl=ru",
				"https://tass.ru/rss/v2.xml",
				"http://www.polit.ru/rss/index.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readConfig()
			if got != tt.want {
				t.Errorf("readConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("readConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
