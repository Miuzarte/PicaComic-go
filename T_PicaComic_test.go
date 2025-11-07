package PicaComic

import (
	"fmt"
	"os"
	"testing"
)

func init() {
	SetToken("")
}

func TestSignIn(t *testing.T) {
	resp, err := SignIn(t.Context(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComics(t *testing.T) {
	resp, err := Comics(t.Context(), "", "", ORDER_LATEST, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComicInfo(t *testing.T) {
	resp, err := ComicInfo(t.Context(), "630f6170c0b3ab7d08f3da8a")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestEpisodes(t *testing.T) {
	resp, err := Episodes(t.Context(), "630f6170c0b3ab7d08f3da8a", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestPages(t *testing.T) {
	resp, err := Pages(t.Context(), "630f6170c0b3ab7d08f3da8a", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestRecommendation(t *testing.T) {
	resp, err := Recommendation(t.Context(), "630f6170c0b3ab7d08f3da8a")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestKeywords(t *testing.T) {
	resp, err := Keywords(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestSearch(t *testing.T) {
	resp, err := Search(t.Context(), "耳で恋した同僚", nil, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestLike(t *testing.T) {
	resp, err := Like(t.Context(), "630f6170c0b3ab7d08f3da8a")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComments(t *testing.T) {
	resp, err := Comments(t.Context(), "630f6170c0b3ab7d08f3da8a", 0)
	if err != nil {
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestFavourite(t *testing.T) {
	resp, err := Favourite(t.Context(), "630f6170c0b3ab7d08f3da8a")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestUserFavourite(t *testing.T) {
	resp, err := UserFavourite(t.Context(), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestCategories(t *testing.T) {
	resp, err := Categories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestImageDownload(t *testing.T) {
	img := ImageInfo{
		FileServer:   "https://storage1.picacomic.com",
		OriginalName: "Cosplay.jpg",
		Path:         "24ee03b1-ad3d-4c6b-9f0f-83cc95365006.jpg",
	}
	_, body, err := img.Download(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	f, e := os.Create(fmt.Sprintf("/home/miuzarte/git/PicaComic-go/_download/%s", img.OriginalName))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	f.Write(body)
	t.Logf("%d", len(body))
	if len(body) < 128 {
		t.Logf("%s", body)
	}
}
