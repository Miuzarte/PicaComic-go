package PicaComic

import (
	"context"
	"fmt"
	"iter"
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

var authorization = ""

func SetToken(token string) {
	authorization = token
}

var threads = 4 // 下载并发数

// SetThreads 设置下载并发数
func SetThreads(n int) {
	if n <= 0 {
		n = 1
	}
	threads = n
}

// SetUseEnvProxy 设置是否使用系统环境变量中的代理
//
// 默认为 true
func SetUseEnvProxy(b bool) {
	ht := httpClient.Transport.(*http.Transport)
	if b {
		ht.Proxy = http.ProxyFromEnvironment
	} else {
		ht.Proxy = nil
	}
}

func SignIn(ctx context.Context, email, password string) (*SignInResp, error) {
	u := toUrl(API_URL)
	u.Path = API_AUTH_SIGNIN
	resp, err := doApiAndDecodeTo[SignInResp](ctx,
		http.MethodPost,
		u.String(),
		map[string]string{
			"email":    email,
			"password": password,
		},
	)
	if err == nil && resp.Token != "" {
		authorization = resp.Token
	}
	return resp, err
}

func Search(ctx context.Context, keyword string, categories []string, sort Sort, page uint) (*SearchResp, error) {
	if sort == "" {
		sort = ORDER_DEFAULT
	}
	p := max(page, 1)

	u := toUrl(API_URL)
	u.Path = API_COMICS_SEARCH

	q := u.Query()
	q.Add("page", strconv.Itoa(int(p)))
	u.RawQuery = q.Encode()

	return doApiAndDecodeTo[SearchResp](ctx,
		http.MethodPost,
		u.String(),
		map[string]any{
			"keyword":    keyword,
			"categories": categories,
			"sort":       sort,
		},
	)
}

func Comics(ctx context.Context, block, tag string, sort Sort, page uint) (*ComicsResp, error) {
	p := max(page, 1)

	u := toUrl(API_URL)
	u.Path = API_COMICS

	q := u.Query()
	q.Add("page", strconv.Itoa(int(p)))
	if block != "" {
		q.Add("c", block)
	}
	if tag != "" {
		q.Add("t", tag)
	}
	if sort != "" {
		q.Add("s", sort)
	}
	u.RawQuery = q.Encode()

	return doApiAndDecodeTo[ComicsResp](ctx, http.MethodGet, u.String(), nil)
}

func ComicInfo(ctx context.Context, bookId string) (*ComicInfoResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId)
	return doApiAndDecodeTo[ComicInfoResp](ctx, http.MethodGet, u.String(), nil)
}

// Episodes 获取章节列表, page 可为 -1
func Episodes(ctx context.Context, bookId string, page int) (*EpsResp, error) {
	page = max(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "eps")

	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()

	resp, err := doApiAndDecodeTo[EpsResp](ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if page == -1 {
		for page := 2; page <= resp.Eps.Pages; page++ {
			r, err := Episodes(ctx, bookId, page)
			if err != nil {
				return resp, err
			}
			resp.append(r)
		}
		resp.Eps.Limit = len(resp.Eps.Docs)

		if len(resp.Eps.Docs) != resp.Eps.Total {
			return resp, fmt.Errorf("Episodes: len(resp.Eps.Docs)(%d) != resp.Eps.Total(%d)",
				len(resp.Eps.Docs), resp.Eps.Total)
		}
	}

	return resp, nil
}

// Pages 获取图片列表, page 可为 -1
func Pages(ctx context.Context, bookId string, epId, page int) (*PagesResp, error) {
	epId = max(epId, 1)
	page = max(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "order", strconv.Itoa(epId), "pages")

	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()

	resp, err := doApiAndDecodeTo[PagesResp](ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if page == -1 {
		for page := 2; page <= resp.Pages.Pages; page++ {
			r, err := Pages(ctx, bookId, epId, page)
			if err != nil {
				return resp, err
			}
			resp.append(r)
		}
		resp.Pages.Limit = len(resp.Pages.Docs)

		if len(resp.Pages.Docs) != resp.Pages.Total {
			return resp, fmt.Errorf("Pages: len(resp.Pages.Docs)(%d) != resp.Pages.Total(%d)",
				len(resp.Pages.Docs), resp.Pages.Total)
		}
	}

	return resp, nil
}

// Recommendation 看了這本子的人也在看
func Recommendation(ctx context.Context, bookId string) (*RecommendationResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "recommendation")
	return doApiAndDecodeTo[RecommendationResp](ctx, http.MethodGet, u.String(), nil)
}

// Comments 偉論, page 可为 -1
func Comments(ctx context.Context, bookId string, page int) (*CommentsResp, error) {
	page = max(page, 1)

	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "comments")

	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()

	resp, err := doApiAndDecodeTo[CommentsResp](ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if page == -1 {
		for page := 2; page <= resp.Comments.Pages; page++ {
			r, err := Comments(ctx, bookId, page)
			if err != nil {
				return resp, err
			}
			resp.append(r)
		}
		resp.Comments.Limit = len(resp.Comments.Docs)

		if len(resp.Comments.Docs) != resp.Comments.Total {
			return resp, fmt.Errorf("Comments: len(resp.Comments.Docs)(%d) != resp.Comments.Total(%d)",
				len(resp.Comments.Docs), resp.Comments.Total)
		}
	}

	return resp, nil
}

// Like 讚好
func Like(ctx context.Context, bookId string) (*LikeResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "like")
	return doApiAndDecodeTo[LikeResp](ctx, http.MethodPost, u.String(), nil)
}

// Favourite 收藏
func Favourite(ctx context.Context, bookId string) (*FavouriteResp, error) {
	u := toUrl(API_URL)
	u.Path = path.Join(API_COMICS, bookId, "favourite")
	return doApiAndDecodeTo[FavouriteResp](ctx, http.MethodPost, u.String(), nil)
}

// Keywords 紳士都在搜的關鍵字
func Keywords(ctx context.Context) (*KeywordsResp, error) {
	u := toUrl(API_URL)
	u.Path = API_KEYWORDS
	return doApiAndDecodeTo[KeywordsResp](ctx, http.MethodGet, u.String(), nil)
}

// Categories 熱門分類
func Categories(ctx context.Context) (*CategoriesResp, error) {
	u := toUrl(API_URL)
	u.Path = API_CATEGORIES
	return doApiAndDecodeTo[CategoriesResp](ctx, http.MethodGet, u.String(), nil)
}

// UserFavourite 已收藏, page 可为 -1
func UserFavourite(ctx context.Context, sort Sort, page int) (*UserFavouriteResp, error) {
	if sort == "" {
		sort = ORDER_DEFAULT
	}
	p := max(page, 1)

	u := toUrl(API_URL)
	u.Path = API_USERS_FAVOURITE

	q := u.Query()
	q.Add("page", strconv.Itoa(p))
	u.RawQuery = q.Encode()

	resp, err := doApiAndDecodeTo[UserFavouriteResp](ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if page == -1 {
		for page := 2; page <= resp.Comics.Pages; page++ {
			r, err := UserFavourite(ctx, sort, page)
			if err != nil {
				return resp, err
			}
			resp.append(r)
		}
		resp.Comics.Limit = len(resp.Comics.Docs)

		if len(resp.Comics.Docs) != resp.Comics.Total {
			return resp, fmt.Errorf("UserFavourite: len(resp.Comics.Docs)(%d) != resp.Comics.Total(%d)",
				len(resp.Comics.Docs), resp.Comics.Total)
		}
	}

	return resp, nil
}

func DownloadCoversIter(ctx context.Context, search *SearchResp) iter.Seq2[Image, error] {
	return newDownloader(ctx, newCoversDownload(search)).downloadIter()
}

func DownloadPagesIter(ctx context.Context, pages *PagesResp) iter.Seq2[Image, error] {
	return newDownloader(ctx, newPagesDownload(pages)).downloadIter()
}
