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
	hResp, resp, err := SignIn(t.Context(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComics(t *testing.T) {
	hResp, resp, err := Comics(t.Context(), "", "", ORDER_LATEST, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComicInfo(t *testing.T) {
	hResp, resp, err := ComicInfo(t.Context(), "682cca5d6f0cd54536b01fd6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestEpisodes(t *testing.T) {
	hResp, resp, err := Episodes(t.Context(), "682cca5d6f0cd54536b01fd6", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestPages(t *testing.T) {
	hResp, resp, err := Pages(t.Context(), "682cca5d6f0cd54536b01fd6", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestRecommendation(t *testing.T) {
	hResp, resp, err := Recommendation(t.Context(), "682cca5d6f0cd54536b01fd6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestKeywords(t *testing.T) {
	hResp, resp, err := Keywords(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestSearch(t *testing.T) {
	hResp, resp, err := Search(t.Context(), "C99", nil, "", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestLike(t *testing.T) {
	hResp, resp, err := Like(t.Context(), "682cca5d6f0cd54536b01fd6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestComments(t *testing.T) {
	hResp, resp, err := Comments(t.Context(), "682cca5d6f0cd54536b01fd6", 1)
	if err != nil {
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestFavourite(t *testing.T) {
	hResp, resp, err := Favourite(t.Context(), "682cca5d6f0cd54536b01fd6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestUserFavourite(t *testing.T) {
	hResp, resp, err := UserFavourite(t.Context(), "", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestCategories(t *testing.T) {
	hResp, resp, err := Categories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", hResp)
	t.Logf("%+v", resp)
	// t.Logf("%s", resp.Todo)
}

func TestImageDownload(t *testing.T) {
	img := Image{
		FileServer:   "https://storage1.picacomic.com",
		OriginalName: "Cosplay.jpg",
		Path:         "24ee03b1-ad3d-4c6b-9f0f-83cc95365006.jpg",
	}
	hResp, body, err := img.Download(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	f, e := os.Create(fmt.Sprintf("/home/miuzarte/git/PicaComic-go/_download/%s", img.OriginalName))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	f.Write(body)
	t.Logf("%+v", hResp)
	t.Logf("%d", len(body))
	if len(body) < 128 {
		t.Logf("%s", body)
	}
}
