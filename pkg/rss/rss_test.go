//Authored by Nikolai Uvarov
//Allowed for free copying
//Package rss implements a simple mechanism for listening of certain rss news feed

package rss

import (
	"NewsAgg/pkg/db/dbmock"
	"reflect"
	"testing"
	"time"

	"github.com/SlyMarbo/rss"
)

func Test_readConfig(t *testing.T) {
	tests := []struct {
		name  string
		want  int
		want1 []string
	}{
		{
			name: "onlyone",
			want: 5,
			want1: []string{
				"http://www.kommersant.ru/RSS/news.xml",
			},
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

func TestCollect(t *testing.T) {
	db := dbmock.New()

	f, _ := rss.Fetch("http://www.kommersant.ru/RSS/news.xml")

	want := len(f.Items)

	Collect(db)

	time.Sleep(10 * time.Second)

	got := db.Len()

	if got != want {
		t.Errorf("Collect() got = %v, want %v ", got, want)
	}
}
