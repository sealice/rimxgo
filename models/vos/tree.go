package vos

type NodeData struct {
	Id       int         `json:"id"`
	ParentId int         `json:"pid"`
	Label    string      `json:"label"`
	Children []*NodeData `json:"children"`
}

type TreeData []*NodeData

func (t TreeData) GenerateTree(rootId int) []*NodeData {
	ls := map[int]*NodeData{}
	treeData := make([]*NodeData, 0)

	for _, node := range t {
		if ls[node.Id] != nil {
			node.Children = append(node.Children, ls[node.Id].Children...)
		}

		ls[node.Id] = node

		if node.ParentId == rootId {
			treeData = append(treeData, node)
			continue
		}

		if ls[node.ParentId] != nil {
			ls[node.ParentId].Children = append(ls[node.ParentId].Children, node)
		} else if node.ParentId != 0 {
			ls[node.ParentId] = &NodeData{Children: []*NodeData{node}}
		}
	}

	return treeData
}
