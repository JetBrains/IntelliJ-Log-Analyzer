package analyzer

func (a *Analyzer) GetIndexingFilesList() []string {
	var IndexingFilesList []string
	for _, entity := range a.DynamicEntities {
		if entity.Name == "Indexing diagnostic" {
			for path, _ := range entity.entityInstances {
				IndexingFilesList = append(IndexingFilesList, path)
			}
		}
	}
	return IndexingFilesList
}
