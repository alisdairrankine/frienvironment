package gfx

import (
	"sort"
	"sync"
)

type Window struct {
	ID    int
	X, Y  int32
	W, H  int32
	Z     int // z-index for draw/click order
	Title string
}

func NewWindowManager() *WindowManager {
	return &WindowManager{
		windows:          make(map[int]*Window),
		nextWinID:        1,
		DraggingWindowID: -1,
	}
}

type WindowManager struct {
	windows   map[int]*Window
	nextWinID int
	topZ      int

	sync.Mutex

	DraggingWindowID int
	DragOffsetX      int32
	DragOffsetY      int32
}

// getSortedWindows returns a slice of windows, sorted by Z-index (bottom to top)
func (wm *WindowManager) SortedWindows() []*Window {
	wm.Lock()
	defer wm.Unlock()

	list := make([]*Window, 0, len(wm.windows))
	for _, w := range wm.windows {
		list = append(list, w)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Z < list[j].Z
	})
	return list
}

func (wm *WindowManager) WindowByID(id int) (*Window, bool) {
	wm.Lock()
	w, ok := wm.windows[id]
	wm.Unlock()
	return w, ok
}

func (wm *WindowManager) MoveWindow(id int, x, y int32) (*Window, bool) {
	wm.Lock()
	w, ok := wm.windows[id]
	w.X = x
	w.Y = y
	wm.Unlock()
	return w, ok
}

// bringToFront gives a window the highest Z-index
func (wm *WindowManager) BringToFront(id int) {
	wm.Lock()
	defer wm.Unlock()
	if w, ok := wm.windows[id]; ok {
		wm.topZ++
		w.Z = wm.topZ
	}
}

func (wm *WindowManager) NewWindow(w *Window) {

	wm.Lock()
	defer wm.Unlock()
	id := wm.nextWinID
	wm.nextWinID++
	wm.topZ++
	w.ID = id
	w.Z = wm.topZ
	wm.windows[id] = w
}

func (wm *WindowManager) DestroyWindow(id int) {
	wm.Lock()
	defer wm.Unlock()
	delete(wm.windows, id)
}
