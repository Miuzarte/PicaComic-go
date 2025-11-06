package PicaComic

import (
	"context"
	"net/http"
	"path"
	"strconv"
)

type Sort = string

const (
	ORDER_DEFAULT Sort = "ua" // 默认
	ORDER_LATEST  Sort = "dd" // 新到舊
	ORDER_OLDEST  Sort = "da" // 舊到新
	ORDER_LOVED   Sort = "ld" // 最多愛心
	ORDER_POINT   Sort = "vd" // 最多紳士指名
)

const (
	API_URL             = "https://picaapi.picacomic.com"
	API_AUTH_SIGNIN     = "/auth/sign-in"
	API_COMICS_SEARCH   = "/comics/advanced-search"
	API_COMICS          = "/comics"
	API_KEYWORDS        = "/keywords"        // 紳士都在搜的關鍵字
	API_CATEGORIES      = "/categories"      // 熱門分類
	API_USERS_FAVOURITE = "/users/favourite" // 已收藏

	IMG_STATIC = "/static"
)

var (
	OrdersList = map[string]string{
		"ua": "默认",
		"dd": "新到舊",
		"da": "舊到新",
		"ld": "最多愛心",
		"vd": "最多紳士指名",
	}
	CategoriesList = []string{
		"嗶咔漢化",
		"全彩",
		"長篇",
		"同人",
		"短篇",
		"圓神領域",
		"碧藍幻想",

		"CG雜圖",
		"英語 ENG",
		"生肉",

		"純愛",
		"百合花園",
		"耽美花園",
		"偽娘哲學",
		"後宮閃光",
		"扶他樂園",

		"姐姐系",
		"妹妹系",

		"SM",
		"性轉換",
		"足の恋",
		"人妻",
		"NTR",
		"強暴",
		"非人類",

		"艦隊收藏",
		"Love Live",
		"SAO 刀劍神域",
		"Fate",
		"東方",

		"WEBTOON",
		"禁書目錄",
		"歐美",
		"Cosplay",
		"重口地帶",
	}
)

type SignInResp struct {
	Token string `json:"token"`
}

func SignIn(ctx context.Context, email, password string) (*http.Response, *SignInResp, error) {
	u := toUrl(API_URL)
	u.Path = API_AUTH_SIGNIN
	return decodeTo[SignInResp](DoApi(ctx,
		http.MethodPost,
		u.String(),
		map[string]string{
			"email":    email,
			"password": password,
		},
	))
}

type Stats struct {
	Total int `json:"total"` // 总数
	Limit int `json:"limit"` // 返回数 == len(.Docs)
	Pages int `json:"pages"` // 共几页
	Page  int `json:"page"`  // 第几页
}

// maybe[TODO] 合并所有的 XxxComic struct

type SearchComic struct { // 14
	UpdatedAt   string   `json:"updated_at"`
	Thumb       Image    `json:"thumb"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	ChineseTeam string   `json:"chineseTeam"`
	CreatedAt   string   `json:"created_at"`
	Finished    bool     `json:"finished"`
	TotalViews  int      `json:"totalViews"`
	Categories  []string `json:"categories"`
	TotalLikes  int      `json:"totalLikes"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Id          string   `json:"_id"`
	LikesCount  int      `json:"likesCount"`
}

type SearchResp struct {
	Comics struct {
		Docs  []SearchComic `json:"docs"`
		Stats `json:",squash"`
	} `json:"comics"`
}

func Search(ctx context.Context, keyword string, categories []string, sort Sort, page int) (*http.Response, *SearchResp, error) {
	if sort == "" {
		sort = ORDER_DEFAULT
	}
	page = min(page, 1)

	u := toUrl(API_URL)
	u.Path = API_COMICS_SEARCH
	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return decodeTo[SearchResp](DoApi(ctx, http.MethodPost, u.String(), map[string]any{
		"keyword":    keyword,
		"categories": categories,
		"sort":       sort,
	}))
}

type Comic struct { // 13
	UnderlineId string   `json:"_id"` // == .Id
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	TotalViews  int      `json:"totalViews"`
	TotalLikes  int      `json:"totalLikes"`
	PagesCount  int      `json:"pagesCount"`
	EpsCount    int      `json:"epsCount"`
	Finished    bool     `json:"finished"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Thumb       Image    `json:"thumb"`
	Id          string   `json:"id"`
	LikesCount  int      `json:"likesCount"`
}

type ComicsResp struct {
	Comics struct {
		Docs  []Comic `json:"docs"`
		Stats `json:",squash"`
	} `json:"comics"`
}

func Comics(ctx context.Context, block, tag string, order Sort, page int) (*http.Response, *ComicsResp, error) {
	u := toUrl(API_URL)
	u.Path = API_COMICS
	q := u.Query()
	if block != "" {
		q.Add("c", block)
	}
	if tag != "" {
		q.Add("t", tag)
	}
	if order != "" {
		q.Add("s", order)
	}
	if page > 0 {
		q.Add("page", strconv.Itoa(page))
	}
	u.RawQuery = q.Encode()
	return decodeTo[ComicsResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type User struct {
	Id         string   `json:"_id"`
	Gender     string   `json:"gender"` // "m" | "f"?
	Name       string   `json:"name"`
	Title      string   `json:"title"` // "WANDANCE"
	Verified   bool     `json:"verified"`
	Exp        int      `json:"exp"`
	Level      int      `json:"level"`
	Characters []string `json:"characters"` // "knight", "vip", "single_doc", "tool_man"
	Role       string   `json:"role"`       // "knight"
	Avatar     Image    `json:"avatar"`
	Slogan     string   `json:"slogan"` // 签名

	// Character string `json:"character"` // [CommentsResp] constant url, ad?
}

type DetailComic struct { // 24
	Id      string `json:"_id"`
	Creator User   `json:"creator"`

	Title       string   `json:"title"`
	Description string   `json:"description"`
	Thumb       Image    `json:"thumb"`
	Author      string   `json:"author"`
	ChineseTeam string   `json:"chineseTeam"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	PagesCount  int      `json:"pagesCount"`
	EpsCount    int      `json:"epsCount"`
	Finished    bool     `json:"finished"`

	UpdatedAt string `json:"updated_at"` // "2025-11-01T17:30:24.731Z"
	CreatedAt string `json:"created_at"` // "2025-05-19T11:34:55.755Z"

	AllowDownload bool `json:"allowDownload"`
	AllowComment  bool `json:"allowComment"`

	TotalLikes    int `json:"totalLikes"`
	TotalViews    int `json:"totalViews"`
	TotalComments int `json:"totalComments"`
	ViewsCount    int `json:"viewsCount"`
	LikesCount    int `json:"likesCount"`
	CommentsCount int `json:"commentsCount"`

	IsFavourite bool `json:"isFavourite"`
	IsLiked     bool `json:"isLiked"`
}

type ComicInfoResp struct {
	Comic DetailComic `json:"comic"`
}

func ComicInfo(ctx context.Context, bookId string) (*http.Response, *ComicInfoResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId)
	return decodeTo[ComicInfoResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type Episode struct {
	UnderlineId string `json:"_id"`
	Title       string `json:"title"`
	Order       int    `json:"order"`      // Docs 中从大到小排序
	UpdatedAt   string `json:"updated_at"` // "2025-11-01T11:46:55.870Z"
	Id          string `json:"id"`
}

type EpsResp struct {
	Eps struct {
		Docs  []Episode `json:"docs"`
		Stats `json:",squash"`
	} `json:"eps"`
}

func Episodes(ctx context.Context, bookId string, page int) (*http.Response, *EpsResp, error) {
	page = min(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "eps")
	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return decodeTo[EpsResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type PageDoc struct {
	UnderlineId string `json:"_id"`
	Media       Image  `json:"media"`
	Id          string `json:"id"`
}

type PagesResp struct {
	Pages struct {
		Docs  []PageDoc `json:"docs"`
		Stats `json:",squash"`
	} `json:"pages"`
	Ep struct {
		Id    string `json:"_id"`
		Title string `json:"title"`
	} `json:"ep"`
}

func Pages(ctx context.Context, bookId string, epId, page int) (*http.Response, *PagesResp, error) {
	epId = min(epId, 1)
	page = min(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "order", strconv.Itoa(epId))
	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return decodeTo[PagesResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type RecommendationComic struct { // 9
	Id         string   `json:"_id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	PagesCount int      `json:"pagesCount"`
	EpsCount   int      `json:"epsCount"`
	Finished   bool     `json:"finished"`
	Categories []string `json:"categories"`
	Thumb      Image    `json:"thumb"`
	LikesCount int      `json:"likesCount"`
}

type RecommendationResp struct {
	Comics []RecommendationComic
}

// Recommendation 看了這本子的人也在看
func Recommendation(ctx context.Context, bookId string) (*http.Response, *RecommendationResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "recommendation")
	return decodeTo[RecommendationResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type CommentDoc struct {
	UnderlineId   string `json:"_id"`
	Content       string `json:"content"`
	User          User   `json:"_user"`
	Id            string `json:"id"`
	LikesCount    int    `json:"likes_count"`
	CommentsCount int    `json:"comments_count"`
	IsLiked       bool   `json:"is_liked"`
}

type CommentsResp struct {
	Comments struct {
		Docs  []CommentDoc `json:"docs"`
		Stats `json:",squash"`
		// Page int `json:"page"` // raw string, weakly decode to int
	} `json:"comments"`
}

// Comments 偉論
func Comments(ctx context.Context, bookId string, page int) (*http.Response, *CommentsResp, error) {
	page = min(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "comments")
	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return decodeTo[CommentsResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type LikeResp struct {
	Action string `json:"action"` // "like" | "unlike"
}

// Like 讚好
func Like(ctx context.Context, bookId string) (*http.Response, *LikeResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "like")
	return decodeTo[LikeResp](DoApi(ctx, http.MethodPost, u.String(), nil))
}

type FavouriteResp struct {
	Action string `json:"action"` // "favourite" | "un_favourite"
}

// Favourite 收藏
func Favourite(ctx context.Context, bookId string) (*http.Response, *FavouriteResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "favourite")
	return decodeTo[FavouriteResp](DoApi(ctx, http.MethodPost, u.String(), nil))
}

type KeywordsResp struct {
	Keywords []string `json:"keywords"`
}

// Keywords 紳士都在搜的關鍵字
func Keywords(ctx context.Context) (*http.Response, *KeywordsResp, error) {
	u := toUrl(API_URL)
	u.Path = API_KEYWORDS
	return decodeTo[KeywordsResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type CategoryBase struct {
	Title string `json:"title"`
	Thumb Image  `json:"thumb"`
}

type CategoryLink struct {
	Link   string `json:"link"` // "" if !.IsWeb
	IsWeb  bool   `json:"isWeb"`
	Active bool   `json:"active"` // always true
}

type CategoryTag struct {
	Id          string `json:"_id"`
	Description string `json:"description"`
}

type CategoriesResp struct {
	Categories []struct {
		CategoryBase
		CategoryLink
		CategoryTag
	}
}

// Categories 熱門分類
func Categories(ctx context.Context) (*http.Response, *CategoriesResp, error) {
	u := toUrl(API_URL)
	u.Path = API_CATEGORIES
	return decodeTo[CategoriesResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}

type UserFavouriteComic struct { // 12
	Id         string   `json:"_id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	TotalViews int      `json:"totalViews"`
	TotalLikes int      `json:"totalLikes"`
	PagesCount int      `json:"pagesCount"`
	EpsCount   int      `json:"epsCount"`
	Finished   bool     `json:"finished"`
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
	Thumb      Image    `json:"thumb"`
	LikesCount int      `json:"likesCount"`
}

type UserFavouriteResp struct {
	Comics struct {
		Docs  []UserFavouriteComic `json:"docs"`
		Stats `json:",squash"`
	}
}

// UserFavourite 已收藏
func UserFavourite(ctx context.Context, sort Sort, page int) (*http.Response, *UserFavouriteResp, error) {
	if sort == "" {
		sort = ORDER_DEFAULT
	}
	page = min(page, 1)

	u := toUrl(API_URL)
	u.Path = API_USERS_FAVOURITE
	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return decodeTo[UserFavouriteResp](DoApi(ctx, http.MethodGet, u.String(), nil))
}
