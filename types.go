package PicaComic

type SignInResp struct {
	Token string `json:"token"`
}

type Stats struct {
	Total int `json:"total"` // 总数
	Limit int `json:"limit"` // 返回数 == len(.Docs)
	Pages int `json:"pages"` // 共几页
	Page  int `json:"page"`  // 第几页
}

// [TODO]? 合并所有的 XxxComic struct

type SearchComic struct { // 14
	UpdatedAt   string    `json:"updated_at"`
	Thumb       ImageInfo `json:"thumb"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	ChineseTeam string    `json:"chineseTeam"`
	CreatedAt   string    `json:"created_at"`
	Finished    bool      `json:"finished"`
	TotalViews  int       `json:"totalViews"`
	Categories  []string  `json:"categories"`
	TotalLikes  int       `json:"totalLikes"` // 有些本没有, 优先 {.LikesCount}
	Title       string    `json:"title"`
	Tags        []string  `json:"tags"`
	Id          string    `json:"_id"`
	LikesCount  int       `json:"likesCount"`
}

type SearchResp struct {
	Comics struct {
		Docs  []SearchComic `json:"docs"`
		Stats `json:",squash"`
	} `json:"comics"`
}

type Comic struct { // 13
	UnderlineId string    `json:"_id"` // == .Id
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	TotalViews  int       `json:"totalViews"`
	TotalLikes  int       `json:"totalLikes"`
	PagesCount  int       `json:"pagesCount"`
	EpsCount    int       `json:"epsCount"`
	Finished    bool      `json:"finished"`
	Categories  []string  `json:"categories"`
	Tags        []string  `json:"tags"`
	Thumb       ImageInfo `json:"thumb"`
	Id          string    `json:"id"`
	LikesCount  int       `json:"likesCount"`
}

type ComicsResp struct {
	Comics struct {
		Docs  []Comic `json:"docs"`
		Stats `json:",squash"`
	} `json:"comics"`
}

type User struct {
	Id         string    `json:"_id"`
	Gender     string    `json:"gender"` // "m" | "f"?
	Name       string    `json:"name"`
	Title      string    `json:"title"` // "WANDANCE"
	Verified   bool      `json:"verified"`
	Exp        int       `json:"exp"`
	Level      int       `json:"level"`
	Characters []string  `json:"characters"` // "knight", "vip", "single_doc", "tool_man"
	Role       string    `json:"role"`       // "knight"
	Avatar     ImageInfo `json:"avatar"`
	Slogan     string    `json:"slogan"` // 签名

	// Character string `json:"character"` // [CommentsResp] constant url, ad?
}

type DetailComic struct { // 24
	Id      string `json:"_id"`
	Creator User   `json:"creator"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumb       ImageInfo `json:"thumb"`
	Author      string    `json:"author"`
	ChineseTeam string    `json:"chineseTeam"`
	Categories  []string  `json:"categories"`
	Tags        []string  `json:"tags"`
	PagesCount  int       `json:"pagesCount"`
	EpsCount    int       `json:"epsCount"`
	Finished    bool      `json:"finished"`

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

type PageDoc struct {
	UnderlineId string    `json:"_id"`
	Media       ImageInfo `json:"media"`
	Id          string    `json:"id"`
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

type RecommendationComic struct { // 9
	Id         string    `json:"_id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	PagesCount int       `json:"pagesCount"`
	EpsCount   int       `json:"epsCount"`
	Finished   bool      `json:"finished"`
	Categories []string  `json:"categories"`
	Thumb      ImageInfo `json:"thumb"`
	LikesCount int       `json:"likesCount"`
}

type RecommendationResp struct {
	Comics []RecommendationComic
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

type LikeResp struct {
	Action string `json:"action"` // "like" | "unlike"
}

type FavouriteResp struct {
	Action string `json:"action"` // "favourite" | "un_favourite"
}

type KeywordsResp struct {
	Keywords []string `json:"keywords"`
}

type CategoryBase struct {
	Title string    `json:"title"`
	Thumb ImageInfo `json:"thumb"`
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

type UserFavouriteComic struct { // 12
	Id         string    `json:"_id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	TotalViews int       `json:"totalViews"`
	TotalLikes int       `json:"totalLikes"`
	PagesCount int       `json:"pagesCount"`
	EpsCount   int       `json:"epsCount"`
	Finished   bool      `json:"finished"`
	Categories []string  `json:"categories"`
	Tags       []string  `json:"tags"`
	Thumb      ImageInfo `json:"thumb"`
	LikesCount int       `json:"likesCount"`
}

type UserFavouriteResp struct {
	Comics struct {
		Docs  []UserFavouriteComic `json:"docs"`
		Stats `json:",squash"`
	}
}
