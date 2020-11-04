package slot

type SyncNode struct {
	Id                int
	Source            string
	SourcePassword    string
	Target            []string
	TargetPassword    string
	SlotLeftBoundary  int
	SlotRightBoundary int
	Slaves            []string
}
