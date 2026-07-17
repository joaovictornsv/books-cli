package models

type EnumValue struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type FieldSpec struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

type SchemaDocument struct {
	Statuses   []EnumValue `json:"statuses"`
	Categories []EnumValue `json:"categories"`
	Fields     []FieldSpec `json:"fields"`
}

func Schema() SchemaDocument {
	return SchemaDocument{
		Statuses:   statusSchema(),
		Categories: categorySchema(),
		Fields:     fieldSchema(),
	}
}

func statusSchema() []EnumValue {
	return []EnumValue{
		{Value: StatusNotStarted.String(), Description: "Owned, not yet reading"},
		{Value: StatusReading.String(), Description: "Currently reading"},
		{Value: StatusRead.String(), Description: "Finished"},
		{Value: StatusToBuy.String(), Description: "Wishlist"},
		{Value: StatusArchived.String(), Description: "Hidden from list and search"},
	}
}

func categorySchema() []EnumValue {
	return []EnumValue{
		{Value: CategoryTheology.String(), Description: "Christian faith, Bible, devotionals, apologetics, pastoral"},
		{Value: CategoryFiction.String(), Description: "Novels, short stories, literary and genre fiction"},
		{Value: CategorySoftware.String(), Description: "Programming, software engineering, CS"},
		{Value: CategoryPhilosophy.String(), Description: "Philosophy, ethics, stoicism, political philosophy"},
		{Value: CategoryHistory.String(), Description: "Historical narrative and historiography"},
		{Value: CategoryPersonalDevelopment.String(), Description: "Self-help, productivity, habits, popular psychology"},
		{Value: CategoryFinanceBusiness.String(), Description: "Money, investing, economics, business"},
		{Value: CategoryScience.String(), Description: "Natural sciences, math popularization"},
		{Value: CategoryPoliticsCulture.String(), Description: "Political/social commentary, cultural criticism"},
		{Value: CategoryBiography.String(), Description: "Biographies, memoirs, autobiographies"},
		{Value: CategoryOther.String(), Description: "Catch-all when no other category fits"},
	}
}

func fieldSchema() []FieldSpec {
	return []FieldSpec{
		{Name: "id", Type: "integer", Description: "Unique book identifier"},
		{Name: "title", Type: "string", Description: "Book title"},
		{Name: "author", Type: "string", Optional: true, Description: "Author name"},
		{Name: "category", Type: "enum", Optional: true, Description: "One of the category enum values"},
		{Name: "status", Type: "enum", Description: "One of the status enum values"},
		{Name: "priority_to_buy", Type: "boolean", Description: "Wishlist priority flag (0 or 1)"},
		{Name: "eligible_to_donate", Type: "boolean", Description: "Eligible to donate flag (0 or 1)"},
		{Name: "donated", Type: "boolean", Description: "Donated flag (0 or 1)"},
		{Name: "notes", Type: "string", Optional: true, Description: "Free-form notes"},
		{Name: "description", Type: "string", Optional: true, Description: "Book synopsis"},
		{Name: "added_at", Type: "timestamp", Description: "When the book was added (RFC3339)"},
		{Name: "started_at", Type: "timestamp", Optional: true, Description: "When reading started (RFC3339)"},
		{Name: "finished_at", Type: "timestamp", Optional: true, Description: "When reading finished (RFC3339)"},
	}
}
