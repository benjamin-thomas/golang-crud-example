package main

import "fmt"

type pagination struct {
	Per, Page, PrevPage, NextPage, LastPage int
}

func newPagination(per, page string, count int) (pagination, error) {
	p := pagination{}

	p.Per = mustAtoi(per)
	p.Page = mustAtoi(page)

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Per < 1 {
		return p, fmt.Errorf("Bogus per: %d", p.Per)
	}

	p.LastPage = (count / p.Per) + 1

	p.PrevPage = p.Page - 1
	if p.PrevPage < 1 {
		p.PrevPage = 1
	}

	p.NextPage = p.Page + 1
	if p.NextPage > p.LastPage {
		p.NextPage = p.LastPage
	}

	return p, nil
}

func (p pagination) HideFirstLink() bool {
	return p.Page <= 1
}

func (p pagination) HideLastLink() bool {
	return p.Page >= p.LastPage
}
