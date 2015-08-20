package main

import (
	"strconv"
)

type pagination struct {
	Per, Page, PrevPage, NextPage, LastPage, Count int
}

func paginateParams(sPer, sPage string) (per, offset int) {
	per, err := strconv.Atoi(sPer)
	if err != nil {
		per = defaultPer
	}

	page, err := strconv.Atoi(sPage)
	if err != nil {
		page = 1
	}

	offset = per * (page - 1)
	return per, offset
}

// func newPagination(per, page string, count int) (pagination, error) {
//   if page == "" {
//     page = "1"
//   }
//   if per == "" {
//     per = defaultPer
//   }

//   p := pagination{}
//   var err error

//   p.Per, err = strconv.Atoi(per)
//   if err != nil {
//     return p, err
//   }

//   p.Page, err = strconv.Atoi(page)
//   if err != nil {
//     return p, err
//   }

//   p.Count = count

//   if p.Page < 1 {
//     p.Page = 1
//   }

//   if p.Per < 1 {
//     return p, fmt.Errorf("Bogus per: %d", p.Per)
//   }

//   p.LastPage = (count / p.Per) + 1

//   p.PrevPage = p.Page - 1
//   if p.PrevPage < 1 {
//     p.PrevPage = 1
//   }

//   p.NextPage = p.Page + 1
//   if p.NextPage > p.LastPage {
//     p.NextPage = p.LastPage
//   }

//   return p, nil
// }

func (p pagination) HideFirstLink() bool {
	return p.Page <= 1
}

func (p pagination) HideLastLink() bool {
	return p.Page >= p.LastPage
}
