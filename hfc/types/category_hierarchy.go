package types

type CategoryHierarchy struct {
	Name   string
	Parent *CategoryHierarchy
}

func NewCategoryHierarchyItem(name string, parent *CategoryHierarchy) *CategoryHierarchy {
	var item = &CategoryHierarchy{
		Name:   name,
		Parent: parent,
	}
	return item
}

func (catHItem *CategoryHierarchy) ToString() string {
	if catHItem.Parent == nil {
		return catHItem.Name
	}
	return catHItem.Parent.ToString() + "/" + catHItem.Name
}

func (catHItem *CategoryHierarchy) AsSlice() []string {
	if catHItem.Parent == nil {
		return []string{catHItem.Name}
	}
	return append(catHItem.Parent.AsSlice(), catHItem.Name)
}
