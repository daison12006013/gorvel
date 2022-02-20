package paginate

import (
	"math"

	"github.com/daison12006013/gorvel/pkg/facade/logger"
	"github.com/daison12006013/gorvel/pkg/response"
)

var defaultView string = "pkg/pagination/tailwind.go.html"

func Construct(items interface{}, total int, perPage int, currentPage int) *Paginate {
	p := Paginate{}
	p.Reconstruct(items, total, perPage, currentPage)
	return &p
}

func (p *Paginate) Reconstruct(items interface{}, total int, perPage int, currentPage int) *Paginate {
	p.Items = items
	p.PerPage = perPage
	p.CurrentPage = currentPage
	p.Total = total
	p.LastPage = int(math.Ceil(float64(total) / float64(perPage)))

	p.OnEachSide = 3
	p.Fragment = nil

	return p
}

func (p *Paginate) Links() string {
	return p.Render(nil)
}

func (p *Paginate) Render(view *string /*, data array*/) string {
	if view == nil {
		dv := defaultView
		view = &dv
	}

	result, err := response.Render(
		[]string{*view},
		map[string]interface{}{
			"currentPage":     p.CurrentPage,
			"total":           p.Total,
			"firstItem":       p.FirstItem(),
			"hasMorePages":    p.HasMorePages(),
			"hasPages":        p.HasPages(),
			"lastItem":        p.LastItem(),
			"nextPageUrl":     p.NextPageUrl(),
			"onFirstPage":     p.OnFirstPage(),
			"previousPageUrl": p.PreviousPageUrl(),

			// here we provide the $elements
			"elements": p.Elements(),
		},
	)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	return result
}

func (p *Paginate) HasMorePages() bool {
	return p.CurrentPage < p.LastPage
}

func (p *Paginate) NextPageUrl() *string {
	if p.HasMorePages() {
		s := p.Url(p.CurrentPage + 1)
		return &s
	}
	return nil
}

func (p *Paginate) ToArray() map[string]interface{} {
	return map[string]interface{}{
		"current_page":   p.CurrentPage,
		"data":           p.Items,
		"first_page_url": p.Url(1),
		"from":           p.FirstItem(),
		"last_page":      p.LastPage,
		"last_page_url":  p.Url(p.LastPage),
		"next_page_url":  p.NextPageUrl(),
		"per_page":       p.PerPage,
		"prev_page_url":  p.PreviousPageUrl(),
		"to":             p.LastItem(),
		"total":          p.Total,
	}
}

func (p *Paginate) Elements() map[int]string {
	window := UrlWindow(*p).Get()
	elems := window.first
	p.elementsLoop(&elems, window.slider)
	p.elementsLoop(&elems, window.last)
	return elems
}

func (p *Paginate) elementsLoop(elems *map[int]string, m map[int]string) map[int]string {
	if len(m) > 0 {
		(*elems)[len(*elems)+1] = "..."

		for _, value := range m {
			(*elems)[len(*elems)+1] = value
		}
	}
	return *elems
}
