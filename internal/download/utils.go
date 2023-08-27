package download

import (
	"dlx/internal/download/datafile"
	"sort"
)

func genSortedStreams(streams map[string]*datafile.Stream) []*datafile.Stream {
	sortedStreams := make([]*datafile.Stream, 0, len(streams))
	for _, stream := range streams {
		sortedStreams = append(sortedStreams, stream)
	}
	if len(sortedStreams) > 1 {
		sort.SliceStable(sortedStreams, func(i, j int) bool {
			return sortedStreams[i].Size > sortedStreams[j].Size
		})
	}
	return sortedStreams
}
