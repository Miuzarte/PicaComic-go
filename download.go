package PicaComic

import (
	"context"
	"iter"
	"strconv"
	"sync"
)

type ImageType int

const (
	IMAGE_TYPE_UNKNOWN ImageType = iota

	IMAGE_TYPE_WEBP
	IMAGE_TYPE_JPEG
	IMAGE_TYPE_PNG
)

func (it ImageType) String() string {
	switch it {
	case IMAGE_TYPE_UNKNOWN:
		return "unknown"
	case IMAGE_TYPE_WEBP:
		return "webp"
	case IMAGE_TYPE_JPEG:
		return "jpeg"
	case IMAGE_TYPE_PNG:
		return "png"
	default:
		return ""
	}
}

type Image struct {
	Ii   ImageInfo
	P    int
	Data []byte
	Type ImageType
	// IsFromCache bool // TODO(maybe
}

func (p *Image) String() string {
	return p.Ii.OriginalName + ": " + strconv.Itoa(len(p.Data))
}

type download struct {
	img *Image
	err chan error
	// cache *cacheComic // TODO(maybe
}

func (d *download) start(ctx context.Context) {
	_, body, err := d.img.Ii.Download(ctx)
	if err == nil {
		d.img.Data = body
	}
	d.err <- err
}

func newCoversDownload(search *SearchResp) (dls []*download) {
	c := search.Comics
	pBase := (c.Page - 1) * (c.Limit)
	dls = make([]*download, len(c.Docs))
	for i := range c.Docs {
		dls[i] = &download{
			img: &Image{
				Ii: c.Docs[i].Thumb,
				P:  pBase + i + 1,
			},
			err: make(chan error, 1),
		}
	}
	return dls
}

func newPagesDownload(pages *PagesResp) (dls []*download) {
	p := pages.Pages
	pBase := (p.Page - 1) * (p.Limit)
	dls = make([]*download, len(p.Docs))
	for i := range p.Docs {
		dls[i] = &download{
			img: &Image{
				Ii: p.Docs[i].Media,
				P:  pBase + i + 1,
			},
			err: make(chan error, 1),
		}
	}
	return dls
}

type downloader struct {
	ctx    context.Context
	cancel context.CancelFunc
	items  []*download
}

func newDownloader(ctx context.Context, dls []*download) *downloader {
	ctx, cancel := context.WithCancel(ctx)
	return &downloader{
		ctx:    ctx,
		cancel: cancel,
		items:  dls,
	}
}

func (dl *downloader) startBackground() {
	go func() {
		limiter := newLimiter()
		defer limiter.close()

		for _, item := range dl.items {
			select {
			case <-dl.ctx.Done():
				return
			case limiter.acquire() <- struct{}{}:
			}

			go func() {
				defer limiter.release()
				item.start(dl.ctx)
			}()
		}
	}()
}

func (dl *downloader) downloadIter() iter.Seq2[Image, error] {
	dl.startBackground()
	return func(yield func(Image, error) bool) {
		defer dl.cancel()
		for _, item := range dl.items {
			if !yield(*item.img, <-item.err) {
				return
			}
		}
	}
}

type limiter struct {
	sem  chan struct{}
	once sync.Once
}

func newLimiter() *limiter {
	return &limiter{
		sem: make(chan struct{}, threads),
	}
}

func (l *limiter) acquire() chan<- struct{} {
	return l.sem
}

func (l *limiter) release() {
	<-l.sem
}

func (l *limiter) close() {
	l.once.Do(func() {
		close(l.sem)
	})
}
