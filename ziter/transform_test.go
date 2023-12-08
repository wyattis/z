package ziter

var _ Iterator[int] = &transformingIterator[int, int]{} // Ensure interface is implemented
var _ Iterator[int] = &filteringIterator[int]{}         // Ensure interface is implemented
