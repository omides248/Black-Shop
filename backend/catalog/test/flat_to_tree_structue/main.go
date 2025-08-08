package main

import (
	"encoding/json"
	"fmt"
)

type CategoryID string
type Category struct {
	ID       CategoryID
	ParentID *CategoryID
	Name     string
}

var flatCategories = []*Category{
	{ID: "C", ParentID: p("B"), Name: "Skins"},
	{ID: "A", ParentID: nil, Name: "Game"},
	{ID: "B", ParentID: p("A"), Name: "Fortnite"},
	{ID: "D", ParentID: nil, Name: "Fortnite"},
}

func p(s string) *CategoryID {
	c := CategoryID(s)
	return &c
}

type CategoryResponse struct {
	ID            CategoryID         `json:"id"`
	Name          string             `json:"name"`
	Subcategories []CategoryResponse `json:"subcategories,omitempty"`
}

func main() {
	//slowTree := buildTreeSlow(flatCategories)
	//printTree(slowTree)
	fastTree := buildTreeFast(flatCategories)
	printTree(fastTree)
}

func printTree(tree []CategoryResponse) {
	data, _ := json.MarshalIndent(tree, "", "\t")
	fmt.Println(string(data))
}

func buildTreeSlow(categories []*Category) []CategoryResponse {

	var roots []CategoryResponse
	for _, cat := range categories {
		if cat.ParentID == nil {

			rootNode := CategoryResponse{ID: cat.ID, Name: cat.Name}

			rootNode.Subcategories = findChildrenSlow(categories, cat.ID)

			roots = append(roots, rootNode)
		}
	}

	return roots
}

func findChildrenSlow(allCategories []*Category, parentID CategoryID) []CategoryResponse {
	var children []CategoryResponse
	for _, cat := range allCategories {

		// Find child for any parent base on ParentID
		if cat.ParentID != nil && *cat.ParentID == parentID {

			// Save child
			childNode := CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
			}

			childNode.Subcategories = findChildrenSlow(allCategories, cat.ID)

			children = append(children, childNode)
		}
	}
	return children
}

func buildTreeFast(categories []*Category) []CategoryResponse {
	categoryMap := make(map[CategoryID]*CategoryResponse)

	for _, cat := range categories {
		categoryMap[cat.ID] = &CategoryResponse{ID: cat.ID, Name: cat.Name, Subcategories: []CategoryResponse{}}
	}

	//data, _ := json.MarshalIndent(categoryMap, "", "\t")
	//fmt.Println(string(data))

	var rootPointers []*CategoryResponse

	for _, cat := range categories {
		node := categoryMap[cat.ID]

		if cat.ParentID == nil {
			rootPointers = append(rootPointers, node)
		} else {
			parentID := *cat.ParentID
			fmt.Println("child category", cat)
			parent, exists := categoryMap[parentID]

			if exists {
				fmt.Println("exists child category", cat, exists)
				parent.Subcategories = append(parent.Subcategories, *node)
			}
		}

	}

	var finalTree []CategoryResponse
	for _, rootPtr := range rootPointers {
		finalTree = append(finalTree, *rootPtr)
	}

	return finalTree

}
