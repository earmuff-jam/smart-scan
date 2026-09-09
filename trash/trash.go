package trash

// MoveToTrash ...
// defines a function that can be used to move files to trash based on the OS
func MoveToTrash(file string) error {
	return moveToTrash(file)
}
